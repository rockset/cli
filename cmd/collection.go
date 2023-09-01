package cmd

import (
	"fmt"
	"github.com/rockset/cli/format"
	"github.com/spf13/cobra"
	"strings"

	"github.com/rockset/rockset-go-client/dataset"
	"github.com/rockset/rockset-go-client/openapi"
	"github.com/rockset/rockset-go-client/option"
)

func newDeleteCollectionCmd() *cobra.Command {
	c := cobra.Command{
		Use:     "collection",
		Aliases: []string{"coll"},
		Short:   "delete collection",
		Long:    "delete Rockset collection",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			ws, _ := cmd.Flags().GetString(WorkspaceFlag)
			name := args[0]

			err = rs.DeleteCollection(ctx, ws, name)
			if err != nil {
				return err
			}

			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "collection '%s.%s' deleted\n", ws, name)
			return nil
		},
	}

	c.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "workspace for the collection")

	return &c
}

func newGetCollectionCmd() *cobra.Command {
	c := cobra.Command{
		Use:     "collection",
		Aliases: []string{"coll"},
		Short:   "get collection",
		Long:    "get Rockset collection",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			wide, _ := cmd.Flags().GetBool(WideFlag)
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)
			name := args[0]

			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			c, err := rs.GetCollection(ctx, ws, name)
			if err != nil {
				return err
			}

			f := format.FormatterFor(cmd.OutOrStdout(), "table", true)

			return f.Format(wide, c)
		},
	}

	c.Flags().Bool(WideFlag, false, "display more information")
	c.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "workspace for the collection")

	return &c
}

func newListCollectionsCmd() *cobra.Command {
	c := cobra.Command{
		Use:     "collections",
		Aliases: []string{"collection", "coll"},
		Short:   "list collections",
		Long:    "list Rockset collections",
		Args:    cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			wide, _ := cmd.Flags().GetBool(WideFlag)
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)

			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			var list []openapi.Collection
			if ws == "" {
				list, err = rs.ListCollections(ctx)
			} else {
				list, err = rs.ListCollections(ctx, option.WithWorkspace(ws))
			}
			if err != nil {
				return err
			}

			f := format.FormatterFor(cmd.OutOrStdout(), "table", true)

			return f.FormatList(wide, format.ToInterfaceArray(list))
		},
	}

	c.Flags().Bool(WideFlag, false, "display more information")
	c.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "workspace for the collection")

	return &c
}

func newCreateCollectionCmd() *cobra.Command {
	c := cobra.Command{
		Use:     "collection NAME",
		Aliases: []string{"coll"},
		Short:   "create collection for use with the write API",
		Long:    "create collection for use with the write API",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			name := args[0]
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)

			options := getCommonCollectionFlags(cmd)

			rs, err := rockClient(cmd)
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

	addCommonCollectionFlags(&c)

	return &c
}

func newCreateS3CollectionCmd() *cobra.Command {
	c := cobra.Command{
		Use:     "collection NAME",
		Aliases: []string{"coll"},
		Short:   "create S3 collection",
		Long:    "create S3 collection",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			name := args[0]
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)

			options := getCommonCollectionFlags(cmd)

			integration, _ := cmd.Flags().GetString(IntegrationFlag)
			bucket, _ := cmd.Flags().GetString(BucketFlag)

			var s3Opts []option.S3SourceOption
			if region, _ := cmd.Flags().GetString(RegionFlag); region != "" {
				s3Opts = append(s3Opts, option.WithS3Region(region))
			}
			if pattern, _ := cmd.Flags().GetString(PatternFlag); pattern != "" {
				s3Opts = append(s3Opts, option.WithS3Pattern(pattern))
			}

			options = append(options, option.WithS3Source(integration, bucket, option.WithJSONFormat(), s3Opts...))

			rs, err := rockClient(cmd)
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

	c.Flags().String(IntegrationFlag, "", "integration name")
	c.Flags().String(BucketFlag, "", "S3 bucket")
	c.Flags().String(PatternFlag, "", "S3 pattern")
	c.Flags().String(RegionFlag, "", "AWS region of the S3 bucket")
	c.Flags().String("source-format", "json", "data source format")

	_ = cobra.MarkFlagRequired(c.Flags(), IntegrationFlag)
	_ = cobra.MarkFlagRequired(c.Flags(), BucketFlag)

	addCommonCollectionFlags(&c)

	return &c
}

func newCreateSampleCollectionCmd() *cobra.Command {
	c := cobra.Command{
		Use:     "collection NAME",
		Aliases: []string{"coll"},
		Args:    cobra.ExactArgs(1),
		Short:   "create sample collection",
		Long:    "create collection with sample data",
		Example: `	## create a sample collection using the movies dataset and wait for the collection to be ready 
	rockset create sample collection --wait --dataset movies movies

	## create a sample collection using the movies dataset with an ingest transformation
	rockset create sample collection --ingest-transformation ingest.sql --dataset movies movies
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			name := args[0]
			ws, _ := cmd.Flags().GetString(WorkspaceFlag)
			from, _ := cmd.Flags().GetString(DatasetFlag)
			wait, _ := cmd.Flags().GetBool(WaitFlag)
			ds := dataset.Sample(from)

			options := getCommonCollectionFlags(cmd)
			pattern := dataset.Lookup(ds)
			if pattern == "" {
				datasets := []string{string(dataset.Movies), string(dataset.MovieRatings)}
				return fmt.Errorf("public dataset %s not found, valid options are: %s", from, strings.Join(datasets, ", "))
			}

			options = append(options, option.WithSampleDataset(ds), option.WithSampleDatasetPattern(pattern))

			rs, err := rockClient(cmd)
			if err != nil {
				return err
			}

			result, err := rs.CreateCollection(ctx, ws, name, options...)
			if err != nil {
				return err
			}

			if wait {
				if err = rs.WaitUntilCollectionReady(ctx, ws, name); err != nil {
					return fmt.Errorf("failed to wait for %s.%s to be ready: %v", ws, name, err)
				}
			}
			_, _ = fmt.Fprintf(cmd.OutOrStdout(), "collection '%s.%s' is %s\n", ws, name, result.GetStatus())

			return nil
		},
	}

	c.Flags().String(DatasetFlag, "", "create sample collection from this dataset")
	_ = cobra.MarkFlagRequired(c.Flags(), DatasetFlag)

	addCommonCollectionFlags(&c)

	return &c
}

func addCommonCollectionFlags(cmd *cobra.Command) {
	cmd.Flags().StringP(WorkspaceFlag, WorkspaceShortFlag, DefaultWorkspace, "workspace for the collection")
	cmd.Flags().String(DescriptionFlag, "", "collection description")
	cmd.Flags().Duration(RetentionFlag, 0, "collection retention")
	cmd.Flags().Bool(WaitFlag, false, "wait until collection is ready")
}

func getCommonCollectionFlags(cmd *cobra.Command) []option.CollectionOption {
	var options []option.CollectionOption

	if retention, _ := cmd.Flags().GetDuration(RetentionFlag); retention != 0 {
		options = append(options, option.WithCollectionRetention(retention))
	}
	if description, _ := cmd.Flags().GetString(DescriptionFlag); description != "" {
		options = append(options, option.WithCollectionDescription(description))
	}
	if compression, _ := cmd.Flags().GetString(CompressionFlag); compression != "" {
		// TODO validate compression
		options = append(options, option.WithStorageCompressionType(option.StorageCompressionType(compression)))
	}

	return options
}
