package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/rockset/cli/completion"
	"github.com/rockset/cli/config"
	"github.com/rockset/cli/flag"
	"github.com/rockset/cli/format"
	"github.com/rockset/cli/sort"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"time"

	"github.com/rockset/rockset-go-client"
	"github.com/rockset/rockset-go-client/dataset"
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/rockset/rockset-go-client/option"
)

func newDeleteCollectionCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "collection",
		Aliases:           []string{"coll", "c"},
		Short:             "delete collection",
		Long:              "delete Rockset collection",
		Annotations:       group("collection"),
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completion.Collection,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			ws, _ := cmd.Flags().GetString(flag.Workspace)
			name := args[0]

			err = rs.DeleteCollection(ctx, ws, name)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "collection '%s.%s' deleted\n", ws, name)
			return nil
		},
	}

	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.DefaultWorkspace, "workspace for the collection")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace)

	return &cmd
}

func newGetCollectionCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:               "collection",
		Aliases:           []string{"coll", "c"},
		Short:             "get collection",
		Long:              "get Rockset collection",
		Annotations:       group("collection"),
		Args:              cobra.ExactArgs(1),
		ValidArgsFunction: completion.Collection,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ws, _ := cmd.Flags().GetString(flag.Workspace)
			output, _ := cmd.Flags().GetString("output")
			name := args[0]

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			collection, err := rs.GetCollection(ctx, ws, name)
			if err != nil {
				return err
			}

			if output != "" {
				var out = cmd.OutOrStdout()

				if output != "-" {
					out, err = os.OpenFile(output, os.O_CREATE|os.O_WRONLY, 0600)
					if err != nil {
						return err
					}
				}

				return json.NewEncoder(out).Encode(translate(collection))
			}

			return formatOne(cmd, collection)
		},
	}

	cmd.Flags().Bool(flag.Wide, false, "display more information")
	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.DefaultWorkspace, "workspace for the collection")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace)

	cmd.Flags().String("output", "", "save json for create collection request to output file, use `-` for stdout")
	_ = cobra.MarkFlagFilename(cmd.Flags(), "output")

	return &cmd
}

func newListCollectionsCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "collections",
		Aliases:     []string{"collection", "coll", "c"},
		Short:       "list collections",
		Long:        "list Rockset collections",
		Annotations: group("collection"),
		Args:        cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			ws, _ := cmd.Flags().GetString(flag.Workspace)

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			var list []openapi.Collection
			if ws == "" || ws == flag.AllWorkspaces {
				list, err = rs.ListCollections(ctx)
			} else {
				list, err = rs.ListCollections(ctx, option.WithWorkspace(ws))
			}
			if err != nil {
				return err
			}

			ms := sort.Multi[openapi.Collection]{
				LessFuncs: []func(p1 *openapi.Collection, p2 *openapi.Collection) bool{
					sort.ByWorkspace[*openapi.Collection],
					sort.ByName[*openapi.Collection],
				},
			}
			ms.Sort(list)

			return formatList(cmd, format.ToInterfaceArray(list))
		},
	}

	cmd.Flags().Bool(flag.Wide, false, "display more information")
	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.AllWorkspaces, "workspace for the collection")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace)

	return &cmd
}

func newCreateCollectionCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "collection NAME",
		Aliases:     []string{"coll", "c"},
		Short:       "create collection for use with the write API",
		Long:        "create collection for use with the write API",
		Annotations: group("collection"),
		Args:        cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			name := args[0]
			ws, _ := cmd.Flags().GetString(flag.Workspace)
			retention, _ := cmd.Flags().GetDuration(flag.Retention)
			transform, _ := cmd.Flags().GetString(flag.IngestTransformation)

			input, _ := cmd.Flags().GetString("input")

			options := getCommonCollectionFlags(cmd)

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			if input != "" {
				var in = cmd.InOrStdin()
				if input != "-" {
					in, err = os.Open(input)
					if err != nil {
						return err
					}
				}

				dec := json.NewDecoder(in)
				var request openapi.CreateCollectionRequest
				if err = dec.Decode(&request); err != nil {
					return err
				}

				request.Name = &name
				options = append(options, option.WithCollectionRequest(request))
			}

			if retention != 0 {
				options = append(options, option.WithCollectionRetention(retention))
			}
			if transform != "" {
				options = append(options, option.WithIngestTransformation(transform))
			}

			result, err := rs.CreateCollection(ctx, ws, name, options...)
			if err != nil {
				return fmt.Errorf("failed to create collection: %w", err)
			}

			if err = waitForCollection(ctx, cmd, rs, ws, name); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "collection '%s.%s' is %s\n", ws, name, result.GetStatus())

			return nil
		},
	}

	cmd.Flags().String("input", "", "input file for create collection request, use `-` to read from stdin")
	_ = cobra.MarkFlagFilename(cmd.Flags(), "input")
	addCommonCollectionFlags(&cmd)

	return &cmd
}

func newCreateS3CollectionCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "collection NAME",
		Aliases:     []string{"coll", "c"},
		Short:       "create S3 collection",
		Long:        "create S3 collection",
		Annotations: group("collection"),
		Args:        cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			name := args[0]
			ws, _ := cmd.Flags().GetString(flag.Workspace)

			options := getCommonCollectionFlags(cmd)

			integration, _ := cmd.Flags().GetString(flag.Integration)
			bucket, _ := cmd.Flags().GetString(flag.Bucket)

			var s3Opts []option.S3SourceOption
			if region, _ := cmd.Flags().GetString(flag.Region); region != "" {
				s3Opts = append(s3Opts, option.WithS3Region(region))
			}
			if pattern, _ := cmd.Flags().GetString(flag.Pattern); pattern != "" {
				s3Opts = append(s3Opts, option.WithS3Pattern(pattern))
			}

			options = append(options, option.WithS3Source(integration, bucket, option.WithJSONFormat(), s3Opts...))

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			result, err := rs.CreateCollection(ctx, ws, name, options...)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "collection '%s.%s' is %s\n", ws, name, result.GetStatus())

			return nil
		},
	}

	cmd.Flags().String(flag.Integration, "", "integration name")
	cmd.Flags().String(flag.Bucket, "", "S3 bucket")
	cmd.Flags().String(flag.Pattern, "", "S3 pattern")
	cmd.Flags().String(flag.Region, "", "AWS region of the S3 bucket")
	cmd.Flags().String("source-format", "json", "data source format")

	_ = cobra.MarkFlagRequired(cmd.Flags(), flag.Integration)
	_ = cmd.RegisterFlagCompletionFunc(flag.Integration, completion.Integration)

	_ = cobra.MarkFlagRequired(cmd.Flags(), flag.Bucket)

	addCommonCollectionFlags(&cmd)

	return &cmd
}

func newCreateSampleCollectionCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "collection NAME",
		Aliases:     []string{"coll", "c"},
		Annotations: group("collection"),
		Args:        cobra.ExactArgs(1),
		Short:       "create sample collection",
		Long:        "create collection with sample data",
		Example: `	## create a sample collection using the movies dataset and wait for the collection to be ready 
	rockset create sample collection --wait --dataset movies movies

	## create a sample collection using the movies dataset with an ingest transformation
	rockset create sample collection --ingest-transformation ingest.sql --dataset movies movies
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			name := args[0]
			ws, _ := cmd.Flags().GetString(flag.Workspace)
			from, _ := cmd.Flags().GetString(flag.Dataset)
			ds := dataset.Sample(from)

			options := getCommonCollectionFlags(cmd)
			pattern := dataset.Lookup(ds)
			if pattern == "" {
				datasets := []string{string(dataset.Movies), string(dataset.MovieRatings)}
				return fmt.Errorf("public dataset %s not found, valid options are: %s", from, strings.Join(datasets, ", "))
			}

			options = append(options, option.WithSampleDataset(ds), option.WithSampleDatasetPattern(pattern))

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			result, err := rs.CreateCollection(ctx, ws, name, options...)
			if err != nil {
				return err
			}

			if err = waitForCollection(ctx, cmd, rs, ws, name); err != nil {
				return err
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "collection '%s.%s' is %s\n", ws, name, result.GetStatus())

			return nil
		},
	}

	cmd.Flags().String(flag.Dataset, "", "create sample collection from this dataset")
	_ = cobra.MarkFlagRequired(cmd.Flags(), flag.Dataset)
	// TODO add completion

	addCommonCollectionFlags(&cmd)

	return &cmd
}

func newCreateTailCollectionCmd() *cobra.Command {
	cmd := cobra.Command{
		Use:         "collection NAME",
		Aliases:     []string{"t"},
		Annotations: group("collection"),
		Args:        cobra.ExactArgs(1),
		Short:       "tail a collection",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			name := args[0]
			ws, _ := cmd.Flags().GetString(flag.Workspace)
			frequency, _ := cmd.Flags().GetDuration("frequency")
			timeField, _ := cmd.Flags().GetString("time-field")

			rs, err := config.Client(cmd)
			if err != nil {
				return err
			}

			// use a large window in the beginning
			when := fmt.Sprintf("CURRENT_TIMESTAMP - SECONDS(10)")
			timer := time.NewTimer(frequency)
			defer timer.Stop()

			for {
				sql := fmt.Sprintf(`SELECT *
FROM %s.%s c
WHERE c.%s > %s
ORDER BY c.%s ASC`,
					ws, name, timeField, when, timeField)
				logger.Debug("getting records", "sql", sql)
				result, err := rs.Query(ctx, sql)
				if err != nil {
					return err
				}

				for _, r := range result.Results {
					// TODO make sure _event_time exists
					t := r[timeField].(string)
					// TODO error handling
					c := r["collections"].([]any)

					fmt.Printf("%s: %v\n", t, strings.Join(Map(c, func(a any) string { return a.(string) }), ", "))
					when = fmt.Sprintf("PARSE_TIMESTAMP_ISO8601('%s')", t)
				}
				select {
				case <-timer.C:
					timer.Reset(frequency)
				case <-ctx.Done():
					// return nil to avoid triggering the error handling in main
					return nil
				}
			}
		},
	}

	cmd.Flags().Duration("frequency", time.Second, "polling frequency to get new documents")
	cmd.Flags().String("time-field", "_event_time", "field name for the time")

	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.DefaultWorkspace, "workspace for the collection")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace)

	return &cmd
}

func Map[IN, OUT any](array []IN, fn func(IN) OUT) []OUT {
	var result = make([]OUT, len(array))

	for i, a := range array {
		result[i] = fn(a)
	}

	return result
}

func waitForCollection(ctx context.Context, cmd *cobra.Command, rs *rockset.RockClient, ws, name string) error {
	wait, _ := cmd.Flags().GetBool(flag.Wait)

	if wait {
		_, _ = fmt.Fprintf(cmd.OutOrStdout(), "waiting for collection '%s.%s' to be READY\n", ws, name)
		if err := rs.Wait.UntilCollectionReady(ctx, ws, name); err != nil {
			return fmt.Errorf("failed to wait for %s.%s to be ready: %v", ws, name, err)
		}
	}

	return nil
}

func addCommonCollectionFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(flag.Workspace, flag.WorkspaceShort, flag.DefaultWorkspace, "workspace for the collection")
	_ = cmd.RegisterFlagCompletionFunc(flag.Workspace, completion.Workspace)

	cmd.Flags().String(flag.Description, "", "collection description")
	cmd.Flags().Duration(flag.Retention, 0, "collection retention")

	cmd.Flags().String(flag.IngestTransformation, "", "ingest transformation SQL")
	cmd.Flags().StringP("ingest-transformation-file", "I", "", "read ingest transformation SQL from file")
	cmd.Flags().Bool(flag.Wait, false, "wait until collection is ready")
}

func getCommonCollectionFlags(cmd *cobra.Command) []option.CollectionOption {
	var options []option.CollectionOption

	if retention, _ := cmd.Flags().GetDuration(flag.Retention); retention != 0 {
		options = append(options, option.WithCollectionRetention(retention))
	}
	if description, _ := cmd.Flags().GetString(flag.Description); description != "" {
		options = append(options, option.WithCollectionDescription(description))
	}
	if compression, _ := cmd.Flags().GetString(flag.Compression); compression != "" {
		// TODO validate compression
		options = append(options, option.WithStorageCompressionType(option.StorageCompressionType(compression)))
	}

	return options
}

func translate(in openapi.Collection) openapi.CreateCollectionRequest {
	out := openapi.CreateCollectionRequest{
		Description:       in.Description,
		FieldMappingQuery: in.FieldMappingQuery,
		RetentionSecs:     in.RetentionSecs,
		Sources:           nil,
	}

	for _, s := range in.Sources {
		s.Status = nil
		out.Sources = append(out.Sources, s)
	}

	return out
}
