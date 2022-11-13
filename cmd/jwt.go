/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jim-at-jibba/dtc/tui"
	"github.com/spf13/cobra"
)

// jwtCmd represents the jwt command
var jwtCmd = &cobra.Command{
	Use:   "jwt-debugger",
	Short: "Decode and inspect JWTs",
	Long:  "Decode and inspect JWTs",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialJWTModel())

		if err := p.Start(); err != nil {
			fmt.Println("WHat", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(jwtCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jwtCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jwtCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type errMsgJWT struct {
	err error
}

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsgJWT) Error() string { return e.err.Error() }

type jwtModel struct {
	rawString textinput.Model
	decoded   string
	err       error
}

func initialJWTModel() jwtModel {
	ti := textinput.New()
	ti.Placeholder = "JWT to decode"
	ti.Focus()

	return jwtModel{
		rawString: ti,
		decoded:   "",
		err:       nil,
	}
}

func (m jwtModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m jwtModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}

	// We handle errors just like any other message
	case errMsgJWT:
		m.err = msg
		return m, nil

	}

	m.rawString, cmd = m.rawString.Update(msg)
	return m, cmd
}

func (m jwtModel) View() string {
	text := "Enter the JWT you want to decode."
	return lipgloss.JoinVertical(lipgloss.Left,
		tui.LabelStyle.Render(text),
		tui.Spacer.Render(""),
		lipgloss.NewStyle().PaddingLeft(1).Render(m.rawString.View()),
		tui.Spacer.Render(""),
		tui.ValueStyle.Render("(q to quit)"),
	)
}
