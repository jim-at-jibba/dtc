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
	"github.com/jim-at-jibba/devtools/tui"
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

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate UUID v4",
	Long:  "Generate UUID v4",
	Run: func(cmd *cobra.Command, args []string) {
		count, _ := cmd.Flags().GetString("count")
		doc := strings.Builder{}

		if count != "" {
			desc := lipgloss.JoinVertical(lipgloss.Left,
				styleLabel.Render("UUIDs: "),
			)

			int, err := strconv.Atoi(count)
			if err != nil {
				fmt.Println("err")
			}

			ids := ""
			for i := 0; i < int; i++ {
				id := uuid.New()
				ids += styleValue.Render(id.String()) + "\n"

			}

			doc.WriteString(desc + "\n" + ids)
		} else {
			id := uuid.New()
			desc := lipgloss.JoinVertical(lipgloss.Left,
				styleLabel.Render("UUID: "),
				styleValue.Render(id.String()),
				styleLabel.Render("UUID Clock ID: "),
				styleValue.Render(strconv.Itoa(id.ClockSequence())),
			)
			doc.WriteString(desc)
		}

		fmt.Println(tui.DocStyle.Render(doc.String()))
	},
}

func init() {
	uuidCmd.AddCommand(generateCmd)

	uuidCmd.PersistentFlags().String("count", "", "Number of UUIDs to generate")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
