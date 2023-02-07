/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/elewis787/boa"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "dtc",
	Short:   "A collection of tools a developer would need on a daily basis.",
	Version: "0.5.1",
	Long: `
	A collection of tools a developer would need on a daily basis. These are
  tools I would normally reach to a browser for.

	* Run dtc help for checking out inline help

	Made with ❤ by https://github.com/jim-at-jibba
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	styles := boa.DefaultStyles()
	styles.Title.BorderForeground(lipgloss.Color("6")).Foreground(lipgloss.Color("4"))
	styles.Border.BorderForeground(lipgloss.Color("4"))
	styles.SelectedItem.Foreground(lipgloss.Color("#3C3C3C")).
		Background(lipgloss.Color("4"))

	b := boa.New(boa.WithStyles(styles))
	rootCmd.SetHelpFunc(b.HelpFunc)
	rootCmd.SetUsageFunc(b.UsageFunc)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.devtools.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
