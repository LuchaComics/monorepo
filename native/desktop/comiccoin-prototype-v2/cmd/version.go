package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

const Major = "1"
const Minor = "0"
const Fix = "0"
const ReleaseType = "alpha"

// Configured via -ldflags during build
var GitCommit string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Describes version.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fmt.Sprintf("Version: %s.%s.%s-%s", Major, Minor, Fix, ReleaseType))
	},
}
