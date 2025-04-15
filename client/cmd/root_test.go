package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestExecute_ErrorHandling(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		cmd := &cobra.Command{
			Use: "failcmd",
			RunE: func(cmd *cobra.Command, args []string) error {
				return fmt.Errorf("test error message")
			},
		}

		rootCmd.AddCommand(cmd)
		rootCmd.SetArgs([]string{"failcmd"})
		Execute()
		return
	}

	cmd := exec.Command(os.Args[0], "-test.run=TestExecute_ErrorHandling")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok {
		if !e.Success() {
			assert.Equal(t, 1, e.ExitCode(), "Should exit with code 1")
			assert.Contains(t, stderr.String(), "test error message",
				"Should print error message")
			return
		}
	}

	t.Fatalf("Process ran with err %v, want exit status 1", err)
}
