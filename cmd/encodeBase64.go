/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	b64 "encoding/base64"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jim-at-jibba/devtools/tui"
	"github.com/spf13/cobra"
)

// encodeBase64Cmd represents the encodeBase64 command
var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encode base64 string",
	Long:  "Encode base64 string",
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(initialModel())

		if err := p.Start(); err != nil {
			fmt.Println("WHat", err)
			os.Exit(1)
		}

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

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type (
	errMsg error
)

type model struct {
	rawString textinput.Model
	encoded   string
	err       error
}

type EncodeStr struct {
	encoded string
}

func Encode(raw string) EncodeStr {
	sEnc := b64.StdEncoding.EncodeToString([]byte(raw))
	return EncodeStr{encoded: sEnc}
}

func (m model) encodeMsg() tea.Msg {
	encoded := Encode(m.rawString.Value())
	return encoded
}

func initialModel() model {
	ti := textinput.New()
	ti.Placeholder = "String to encode"
	ti.Focus()

	return model{
		rawString: ti,
		encoded:   "",
		err:       nil,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			return m, m.encodeMsg
		}

	case EncodeStr:
		encoded := msg
		m.encoded = encoded.encoded
		return m, tea.Quit

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil

	}

	m.rawString, cmd = m.rawString.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if len(m.encoded) > 0 {
		return tui.DocStyle.Render(fmt.Sprintf(
			m.encoded,
		))
	} else {
		return fmt.Sprintf(
			"Enter the string you want to encode.\n\n%s\n\n%s",
			m.rawString.View(),
			"(esc to quit)",
		) + "\n"
	}
}
