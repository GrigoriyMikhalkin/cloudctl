package cmd

import (
	"fmt"

	"github.com/metal-pod/updater"
	"github.com/spf13/cobra"
)

const (
	downloadURLPrefix = "https://blobstore.fi-ts.io/cloud/" + programName + "/"
)

var (
	updateCmd = &cobra.Command{
		Use:   "update",
		Short: "update the program",
	}
	updateCheckCmd = &cobra.Command{
		Use:   "check",
		Short: "check for update of the program",
		RunE: func(cmd *cobra.Command, args []string) error {
			u := updater.New(downloadURLPrefix, programName)
			return u.Check()
		},
	}
	updateDoCmd = &cobra.Command{
		Use:   "do",
		Short: "do the update of the program",
		RunE: func(cmd *cobra.Command, args []string) error {
			u := updater.New(downloadURLPrefix, programName)
			return u.Do()
		},
	}
	updateDumpCmd = &cobra.Command{
		Use:   "dump <binary>",
		Short: "dump the version update file",
		RunE: func(cmd *cobra.Command, args []string) error {
			u := updater.New(downloadURLPrefix, programName)
			if len(args) < 1 {
				return fmt.Errorf("full path to program required")
			}
			return u.Dump(args[0])
		},
	}
)

func init() {
	updateCmd.AddCommand(updateCheckCmd)
	updateCmd.AddCommand(updateDoCmd)
	updateCmd.AddCommand(updateDumpCmd)
}
