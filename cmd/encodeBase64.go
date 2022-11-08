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
	"github.com/jim-at-jibba/dev-tools-cli/tui"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
)

// encodeBase64Cmd represents the encodeBase64 command
var encodeCmd = &cobra.Command{
	Use:   "encode",
	Short: "Encode base64 string",
	Long:  "Encode base64 string",
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetBool("url")
		p := tea.NewProgram(initialEncodeModel(url))

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
	encodeCmd.Flags().BoolP("url", "u", false, "URL-compatible base64 format")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encodeBase64Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type errMsg struct {
	err error
}

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

type encodeModel struct {
	rawString textinput.Model
	encoded   string
	err       error
	url       bool
}

type encodeStr struct {
	encoded string
}

func Encode(raw string, url bool) encodeStr {
	if url {
		sEnc := b64.URLEncoding.EncodeToString([]byte(raw))
		clipboard.Write(clipboard.FmtText, []byte(sEnc))
		return encodeStr{encoded: sEnc}
	} else {
		sEnc := b64.StdEncoding.EncodeToString([]byte(raw))
		clipboard.Write(clipboard.FmtText, []byte(sEnc))
		return encodeStr{encoded: sEnc}
	}
}

func (m encodeModel) encodeMsg() tea.Msg {
	encoded := Encode(m.rawString.Value(), m.url)
	return encoded
}

func initialEncodeModel(url bool) encodeModel {
	ti := textinput.New()
	ti.Placeholder = "String to encode"
	ti.Focus()

	return encodeModel{
		rawString: ti,
		encoded:   "",
		err:       nil,
		url:       url,
	}
}

func (m encodeModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m encodeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			return m, m.encodeMsg
		}

	case encodeStr:
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

func (m encodeModel) View() string {
	if len(m.encoded) > 0 {
		return tui.ContainerStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				tui.LabelStyle.Render("Encoded string:"),
				tui.Spacer.Render(""),
				tui.ValueStyle.Render(m.encoded),
				tui.LabelStyle.Render("(Copied to clipboard)"),
			),
		)
	} else {
		var text string
		if m.url {
			text = "Enter the string you want to URL-compatible encode"
		} else {
			text = "Enter the string you want to encode."
		}
		return lipgloss.JoinVertical(lipgloss.Left,
			tui.LabelStyle.Render(text),
			tui.Spacer.Render(""),
			lipgloss.NewStyle().PaddingLeft(1).Render(m.rawString.View()),
			tui.Spacer.Render(""),
			tui.ValueStyle.Render("(q to quit)"),
		)
	}
}
