package auth

import (
	"os"

	"github.com/spf13/cobra"
)

func init() {
	os.UserHomeDir()
}

func Auth(cmd *cobra.Command, args []string) {}
