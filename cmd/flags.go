package cmd

import (
	"github.com/rockset/cli/format"
	"github.com/spf13/cobra"
)

// global flags
const (
	ClusterFLag  = "cluster"
	ContextFLag  = "context"
	DebugFlag    = "debug"
	FormatFlag   = "format"
	HeaderFlag   = "header"
	SelectorFlag = "selector"
	WideFlag     = "wide"
)

// command specific flags
const (
	BucketFlag           = "bucket"
	CompressionFlag      = "compression"
	CollectionFlag       = "collection"
	DatasetFlag          = "dataset"
	DescriptionFlag      = "description"
	FileFlag             = "file"
	ForceFlag            = "force"
	IngestTransformation = "ingest-transformation"
	IntegrationFlag      = "integration"
	PatternFlag          = "pattern"
	RegionFlag           = "region"
	RetentionFlag        = "retention"
	RoleARNFlag          = "role-arn"
	SizeFlag             = "size"
	SQLFlag              = "sql"
	ValidateFlag         = "validate"
	VersionFlag          = "version"
	WaitFlag             = "wait"
	WorkspaceFlag        = "workspace"
	WorkspaceShortFlag   = "W"
)

const (
	AllWorkspaces    = "all"
	DefaultFormat    = "table"
	DefaultWorkspace = "commons"
)

func FormatFromCommand(cmd *cobra.Command) format.Format {
	f, err := cmd.Flags().GetString(FormatFlag)
	if err != nil {
		if cmd.Parent() == nil {
			return "unknown"
		}
		return FormatFromCommand(cmd.Parent())
	}

	return format.Format(f)
}
