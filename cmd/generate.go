/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var styleValue = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#7D56F4")).
	PaddingLeft(1)

var styleLabel = lipgloss.NewStyle().
	Bold(true).
	Foreground(lipgloss.Color("#FFF")).
	PaddingTop(1).
	PaddingLeft(1)

var subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
var docStyle = lipgloss.NewStyle().
	Padding(1, 2, 1, 2).
	Border(lipgloss.NormalBorder(), true).
	BorderForeground(lipgloss.Color("#7D56F4"))

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate UUID v4",
	Long:  "Generate UUID v4",
	Run: func(cmd *cobra.Command, args []string) {
		doc := strings.Builder{}
		id := uuid.New()
		desc := lipgloss.JoinVertical(lipgloss.Left,
			styleLabel.Render("UUID String: "),
			styleValue.Render(id.String()),
			styleLabel.Render("UUID Clock ID: "),
			styleValue.Render(strconv.Itoa(id.ClockSequence())),
		)

		doc.WriteString(desc)
		fmt.Println(docStyle.Render(doc.String()))
	},
}

func init() {
	uuidCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
