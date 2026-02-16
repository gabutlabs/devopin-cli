package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd adalah command dasar saat aplikasi dijalankan tanpa subcommand
var rootCmd = &cobra.Command{
	Use:   "devopin",
	Short: "Devopin CLI - Resource monitoring and alerting tool",
	Long: `Devopin CLI is a command-line application
to monitor resources like CPU, disk, and memory,
and send alerts via Telegram when thresholds are exceeded.`,
}

// Execute adalah fungsi yang dipanggil oleh main.go.
// Fungsi ini akan menjalankan command yang sesuai.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// init() akan dipanggil saat package ini di-load.
// Di sini kita mendaftarkan semua subcommand.
func init() {
	rootCmd.AddCommand(resourceAlertCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(versionCmd)
}
