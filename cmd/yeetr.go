/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// yeetrCmd represents the yeetr command
var yeetrCmd = &cobra.Command{
	Use:   "yeetr",
	Short: "For sharing files.",
	Long:  "Sharing is, by design, ephemeral, so, the link that yeetr provides will expire, after a given time or when the file is downloaded. The idea for this tool was lifted from https://www.npmjs.com/package/yeetr",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("yeetr called")
	},
}

func init() {
	rootCmd.AddCommand(yeetrCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// yeetrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// yeetrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
