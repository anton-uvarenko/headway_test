package cmd

import (
	"github.com/anton-uvarenko/headway_test/internal/pkg"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "nasa",
	Short: "",
	Long:  "",
	RunE:  nasaCommnad,
}

func Execute() error {
	return rootCmd.Execute()
}

func nasaCommnad(cmd *cobra.Command, args []string) error {
	err := pkg.FetchLastSevenDays()
	return err
}
