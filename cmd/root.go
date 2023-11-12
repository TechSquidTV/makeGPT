/*
Copyright © 2023 Kyle Tryon (a.k.a TechSquidTV) makeGPT@TechSquidTV.com
*/
package cmd

import (
	"os"

	"github.com/TechSquidTV/makeGPT/packages/api"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "makeGPT",
	Short: "Create, manage, and update your OpenAI custom GPTs",
	Long: `Easily create, update, manage, and publish your OpenAI custom GPTs from the command line.
makeGPT enables you to manage your custom GPTs with Git and utilize CI/CD pipelines to automatically publish updates.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var nonInteractive bool

func init() {
	api.CheckBearerToken()

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.makeGPT.yaml)")
	// Define a --non-interactive flag
	rootCmd.PersistentFlags().BoolVar(&nonInteractive, "non-interactive", false, "Run in non-interactive mode")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
