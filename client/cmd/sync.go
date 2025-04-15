package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/GlebRadaev/password-manager/client/services"
)

// syncService is the shared sync service instance
var syncService SyncServiceInterface = services.NewSyncService()

// syncCmd represents the sync command for synchronizing data with the server
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync data with server",
	Long: `Synchronizes local password entries with the remote server.
Detects and reports any conflicts that need resolution.`,
	Run: func(cmd *cobra.Command, args []string) {
		resp, err := syncService.Sync()
		if err != nil {
			log.Fatalf("Sync failed: %v", err)
		}

		if len(resp.Conflicts) > 0 {
			fmt.Printf("Found %d conflicts:\n", len(resp.Conflicts))
			for _, c := range resp.Conflicts {
				fmt.Printf("- %s (ID: %s)\n", c.DataID, c.ConflictID)
			}
		} else {
			fmt.Println("Sync completed successfully")
		}
	},
}

// resolveCmd represents the command for resolving synchronization conflicts
var resolveCmd = &cobra.Command{
	Use:   "resolve <conflict-id>",
	Short: "Resolve sync conflict",
	Long: `Resolves a synchronization conflict using specified strategy.
Available strategies: client (keep local), server (keep remote), merge (combine).`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		strategy, _ := cmd.Flags().GetString("strategy")

		resp, err := syncService.Resolve(args[0], strategy)
		if err != nil {
			log.Fatalf("Resolve failed: %v", err)
		}
		fmt.Println("Conflict resolved:", resp.Message)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd, resolveCmd)

	resolveCmd.Flags().StringP("strategy", "s", "", "Resolution strategy (client, server, merge)")
	resolveCmd.MarkFlagRequired("strategy")
}
