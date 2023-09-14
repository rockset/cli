package cmd

import (
	"github.com/rockset/cli/format"
	"github.com/spf13/cobra"
)

const (
	ClusterFLag        = "cluster"
	DescriptionFlag    = "description"
	WorkspaceFlag      = "workspace"
	WorkspaceShortFlag = "W"
	RetentionFlag      = "retention"
	CompressionFlag    = "compression"
	IntegrationFlag    = "integration"
	RegionFlag         = "region"
	BucketFlag         = "bucket"
	PatternFlag        = "pattern"
	RoleARNFlag        = "role-arn"
	DatasetFlag        = "dataset"
	WaitFlag           = "wait"
	ValidateFlag       = "validate"
	FileFlag           = "file"
	ContextFLag        = "context"
	DebugFlag          = "debug"
	SizeFlag           = "size"

	FormatFlag = "format"
	WideFlag   = "wide"
	HeaderFlag = "header"
)

const (
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
