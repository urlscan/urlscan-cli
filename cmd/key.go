package cmd

import (
	"errors"
	"fmt"

	"io"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/urlscan/urlscan-cli/pkg/utils"

	"golang.org/x/term"
)

type PasswordReader struct {
	stdout io.Writer
	stdin  io.Reader
}

func NewPasswordReader(stdout io.Writer, stdin io.Reader) *PasswordReader {
	return &PasswordReader{
		stdout: stdout,
		stdin:  stdin,
	}
}

func (p *PasswordReader) ReadPassword(args []string) ([]byte, error) {
	// check stdin first
	reader := utils.StringReaderFromCmdArgs(args)
	token, err := reader.ReadString()
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}
	if token != "" {
		return []byte(token), nil
	}

	// read from terminal next
	_, err = fmt.Fprint(p.stdout, "Enter a urlscan.io API key: ")
	if err != nil {
		return nil, err
	}

	b, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(p.stdout, "") //nolint:errcheck
	if err != nil {
		return nil, fmt.Errorf("read a urlscan.io API key from terminal: %w", err)
	}
	return b, nil
}

var setKeyCmd = &cobra.Command{
	Use:   "set",
	Short: "Set urlscan.io API key",
	RunE: func(cmd *cobra.Command, args []string) error {
		reader := NewPasswordReader(cmd.OutOrStdout(), cmd.InOrStdin())
		b, err := reader.ReadPassword(args)
		if err != nil {
			return err
		}

		key := strings.TrimSpace(string(b))
		if key == "" {
			return fmt.Errorf("API key cannot be empty")
		}

		if err := utils.NewKeyManager().SetKey(key); err != nil {
			return err
		}

		return nil
	},
}

var removeKeyCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove urlscan.io API key",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return cmd.Usage()
		}

		return utils.NewKeyManager().RemoveKey()
	},
}

var keyCmd = &cobra.Command{
	Use:   "key",
	Short: "Manage API key",
}

func init() {
	keyCmd.AddCommand(setKeyCmd)
	keyCmd.AddCommand(removeKeyCmd)

	RootCmd.AddCommand(keyCmd)
}
