/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// encodeBase64Cmd represents the encodeBase64 command
var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encode base64 string",
	Long:  "Encode base64 string",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("encodeBase64 called")
	},
}

func init() {
	base64Cmd.AddCommand(encodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encodeBase64Cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encodeBase64Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
