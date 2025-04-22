package cmd

import (
	"fmt"

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
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := syncService.Sync()
		if err != nil {
			cmd.PrintErrf("Sync failed: %v\n", err)
			return fmt.Errorf("sync operation failed: %w", err)
		}

		if len(resp.Conflicts) > 0 {
			cmd.Printf("Found %d conflicts:\n", len(resp.Conflicts))
			for _, c := range resp.Conflicts {
				cmd.Printf("- %s (ID: %s)\n", c.DataID, c.ConflictID)
			}
		} else {
			cmd.Println("Sync completed successfully")
		}
		return nil
	},
}

// resolveCmd represents the command for resolving synchronization conflicts
var resolveCmd = &cobra.Command{
	Use:   "resolve <conflict-id>",
	Short: "Resolve sync conflict",
	Long: `Resolves a synchronization conflict using specified strategy.
Available strategies: client (keep local), server (keep remote), merge (combine).`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		strategy, _ := cmd.Flags().GetString("strategy")

		resp, err := syncService.Resolve(args[0], strategy)
		if err != nil {
			cmd.PrintErrf("Resolve failed: %v\n", err)
			return fmt.Errorf("resolve operation failed for conflict %s: %w", args[0], err)
		}
		cmd.Println("Conflict resolved:", resp.Message)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(syncCmd, resolveCmd)

	resolveCmd.Flags().StringP("strategy", "s", "", "Resolution strategy (client, server, merge)")
	resolveCmd.MarkFlagRequired("strategy")
}
