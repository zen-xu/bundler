package cmd

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	bundle "github.com/zen-xu/bundler/pkg"
	"github.com/zen-xu/bundler/pkg/utils"
)

type Option struct {
	Verbose    bool
	OutputPath string
}

var option = Option{}
var rootCmd = &cobra.Command{
	Use:   "bundler [flags] config",
	Short: "Bundler is tool for bundling resources into a single executable",
	Long:  "Bundler is tool for bundling resources into a single executable",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires config")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		config := args[0]
		bundler, err := bundle.NewBundler(config)
		utils.CheckError(err, "Fail to init bundler")

		ignorePaths := bundler.Bundle(option.OutputPath, option.Verbose)
		if option.Verbose && len(ignorePaths) > 0 {
			fmt.Println(utils.Yellow("Warning: ignored bundle files"))
			for _, path := range ignorePaths {
				fmt.Println(utils.Blue(path))
			}
			fmt.Println()
		}
		fmt.Println(utils.Green("bundle success"))
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if option.OutputPath == "" {
			config := args[0]
			configName := strings.TrimSuffix(config, filepath.Ext(config))
			option.OutputPath = fmt.Sprintf("%s.bundle", configName)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize()

	rootCmd.Flags().StringVarP(&option.OutputPath, "output", "o", "", "bundle output path")
	rootCmd.Flags().BoolVarP(&option.Verbose, "verbose", "v", false, "increase verbosity")
}
