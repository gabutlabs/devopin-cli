package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number",
	Long:  `Print the version number of Devopin CLI`,
	Run: func(cmd *cobra.Command, args []string) {
		if Version == "" {
			fmt.Println("Devopin CLI v0.0.0-dev (build from source)")
		} else {
			fmt.Printf("Devopin CLI %s\n", Version)
		}

	},
}

// Version will be set by ldflags during build
var Version string
