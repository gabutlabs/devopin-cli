package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd adalah command dasar saat aplikasi dijalankan tanpa subcommand
var rootCmd = &cobra.Command{
	Use:   "Devopin CLI", // Ganti 'nama-aplikasi' sesuai keinginan
	Short: "Aplikasi sederhana untuk monitoring resource dan worker",
	Long: `Devopin CLI adalah aplikasi command-line sederhana
untuk memonitor resource seperti CPU, disk, dan memory,
serta mengelola worker untuk tugas-tugas tertentu.`,
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
}
