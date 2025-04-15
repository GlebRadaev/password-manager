package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cobra"

	"github.com/GlebRadaev/password-manager/client/models"
	"github.com/GlebRadaev/password-manager/client/services"
)

// dataService is the shared data service instance
var dataService DataServiceInterface = services.NewDataService()

// addCmd represents the command to add a new data entry
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add new data entry",
	Long: `Add a new entry to the password manager.
Supported entry types: login, note, card, binary.
Example: pm add -t login -d '{"username":"user","password":"pass"}'`,
	Run: func(cmd *cobra.Command, args []string) {
		dataType, _ := cmd.Flags().GetString("type")
		content, _ := cmd.Flags().GetString("data")

		entry := &models.DataEntry{
			ID:        uuid.New().String(),
			Type:      models.DataTypeFromString(dataType),
			Data:      []byte(content),
			CreatedAt: time.Now().Unix(),
			UpdatedAt: time.Now().Unix(),
		}

		if err := dataService.Add(entry); err != nil {
			log.Fatalf("Add failed: %v", err)
		}
		fmt.Printf("Added entry with ID: %s\n", entry.ID)
	},
}

// listCmd represents the command to list all entries
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all entries",
	Long:  `Display a summary list of all stored entries including ID, type and last update time.`,
	Run: func(cmd *cobra.Command, args []string) {
		entries, err := dataService.List()
		if err != nil {
			log.Fatalf("List failed: %v", err)
		}

		for i, e := range entries {
			fmt.Printf("%d. %s [%s] %s\n", i+1, e.ID, e.Type.String(),
				time.Unix(e.UpdatedAt, 0).Format("2006-01-02"))
		}
	},
}

// viewCmd represents the command to view entry details
var viewCmd = &cobra.Command{
	Use:   "view <id>",
	Short: "View entry details",
	Long: `Display complete details for a specific entry including:
- Full metadata (creation/update times)
- Complete data content
Example: pm view 3a9b8c7d-6e5f-4a3b-2c1d-0e9f8a7b6c5d`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		entry, err := dataService.Get(args[0])
		if err != nil {
			log.Fatalf("View failed: %v", err)
		}

		fmt.Printf("ID: %s\n", entry.ID)
		fmt.Printf("Type: %s\n", entry.Type.String())
		fmt.Printf("Created: %s\n", time.Unix(entry.CreatedAt, 0).Format(time.RFC822))
		fmt.Printf("Updated: %s\n", time.Unix(entry.UpdatedAt, 0).Format(time.RFC822))
		fmt.Printf("Data: %s\n", string(entry.Data))
	},
}

// deleteCmd represents the command to delete an entry
var deleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete entry",
	Long: `Permanently remove an entry from the password manager.
Example: pm delete 3a9b8c7d-6e5f-4a3b-2c1d-0e9f8a7b6c5d`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := dataService.Delete(args[0]); err != nil {
			log.Fatalf("Delete failed: %v", err)
		}
		fmt.Println("Entry deleted")
	},
}

func init() {
	rootCmd.AddCommand(addCmd, listCmd, viewCmd, deleteCmd)

	// Add command flags with descriptions
	addCmd.Flags().StringP("type", "t", "", "Entry type (login|note|card|binary)")
	addCmd.Flags().StringP("data", "d", "", "Entry content (JSON format for structured types)")
	addCmd.MarkFlagRequired("type")
	addCmd.MarkFlagRequired("data")
}
