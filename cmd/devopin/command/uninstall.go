package command

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall Devopin CLI and remove systemd service",
	Long:  `Remove Devopin CLI binary, configuration, and systemd service (if installed).`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Uninstalling Devopin CLI...")
		
		// Stop and disable systemd service
		fmt.Println("Stopping systemd service (if running)...")
		if _, err := runCommand("systemctl", "stop", "devopin-resource-alert"); err != nil {
			fmt.Printf("Note: Could not stop service: %v\n", err)
		}
		
		if _, err := runCommand("systemctl", "disable", "devopin-resource-alert"); err != nil {
			fmt.Printf("Note: Could not disable service: %v\n", err)
		}
		
		// Remove systemd service file
		serviceFile := "/etc/systemd/system/devopin-resource-alert.service"
		if _, err := os.Stat(serviceFile); err == nil {
			fmt.Printf("Removing systemd service: %s\n", serviceFile)
			if err := os.Remove(serviceFile); err != nil {
				fmt.Printf("Error removing service file: %v\n", err)
			}
		}
		
		// Reload systemd daemon
		if _, err := runCommand("systemctl", "daemon-reload"); err != nil {
			fmt.Printf("Note: Could not reload systemd daemon: %v\n", err)
		}
		
		// Remove binary
		binaryPath := "/usr/local/bin/devopin"
		if _, err := os.Stat(binaryPath); err == nil {
			fmt.Printf("Removing binary: %s\n", binaryPath)
			if err := os.Remove(binaryPath); err != nil {
				fmt.Printf("Error removing binary: %v\n", err)
			}
		}
		
		// Remove config directory (optional - ask user)
		configDir := "/etc/devopin"
		if _, err := os.Stat(configDir); err == nil {
			fmt.Printf("Config directory exists: %s\n", configDir)
			fmt.Println("To remove configuration, run: sudo rm -rf /etc/devopin")
		}
		
		fmt.Println("\nUninstallation complete!")
		fmt.Println("Note: You may need to manually remove /etc/devopin if you want to delete all configurations.")
	},
}

// Helper function to run shell commands
func runCommand(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
