/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	b64 "encoding/base64"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jim-at-jibba/dtc/tui"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
)

// encodeBase64Cmd represents the encodeBase64 command
var decodeCmd = &cobra.Command{
	Use:   "decode",
	Short: "Decode base64 string",
	Long:  "Decode base64 string",
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetBool("url")
		p := tea.NewProgram(initialDecodeModel(url))

		if err := p.Start(); err != nil {
			fmt.Println("WHat", err)
			os.Exit(1)
		}

	},
}

func init() {
	base64Cmd.AddCommand(decodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	decodeCmd.Flags().BoolP("url", "u", false, "URL-compatible base64 format")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encodeBase64Cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type decodeModel struct {
	rawString textinput.Model
	decoded   string
	err       error
	url       bool
}

type decodeStr struct {
	decoded string
}

func Decode(raw string, url bool) (decodeStr, error) {

	if url {
		sDec, err := b64.URLEncoding.DecodeString(strings.TrimSpace(raw))
		if err != nil {
			fmt.Println("sDec", string(sDec), err)
			return decodeStr{}, err
		}
		clipboard.Write(clipboard.FmtText, []byte(sDec))
		return decodeStr{decoded: string(sDec)}, nil
	} else {
		sDec, err := b64.StdEncoding.DecodeString(strings.TrimSpace(raw))
		if err != nil {
			fmt.Println("sDec", string(sDec), err)
			return decodeStr{}, err
		}
		clipboard.Write(clipboard.FmtText, []byte(sDec))
		return decodeStr{decoded: string(sDec)}, nil
	}
}

func (m decodeModel) decodeMsg() tea.Msg {
	decoded, err := Decode(m.rawString.Value(), m.url)
	if err != nil {
		return errMsg{err: err}
	}
	return decoded
}

func initialDecodeModel(url bool) decodeModel {
	ti := textinput.New()
	ti.Placeholder = "String to decode"
	ti.Focus()

	return decodeModel{
		rawString: ti,
		decoded:   "",
		err:       nil,
		url:       url,
	}
}

func (m decodeModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m decodeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			return m, m.decodeMsg
		}

	case decodeStr:
		decoded := msg
		m.decoded = decoded.decoded
		return m, tea.Quit

	// We handle errors just like any other message
	case errMsg:
		m.err = msg.err
		return m, tea.Quit

	}

	m.rawString, cmd = m.rawString.Update(msg)
	return m, cmd
}

func (m decodeModel) View() string {
	if len(m.decoded) > 0 {
		return tui.ContainerStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				tui.LabelStyle.Render("Decoded string:"),
				tui.Spacer.Render(""),
				tui.ValueStyle.Render(m.decoded),
				tui.LabelStyle.Render("(Copied to clipboard)"),
			),
		)
	} else if m.err != nil {
		return tui.ErrorContainerStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				tui.LabelStyle.Render("Decoding error:"),
				tui.Spacer.Render(""),
				tui.ValueStyle.Render(m.err.Error()),
			),
		)
	} else {
		var text string
		if m.url {
			text = "Enter the string you want to URL-compatible decode"
		} else {
			text = "Enter the string you want to deccode."
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
