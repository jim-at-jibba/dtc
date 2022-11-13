/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	jwt "github.com/dgrijalva/jwt-go"
	"golang.design/x/clipboard"

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
		if len(os.Getenv("DEBUG")) > 0 {
			f, err := tea.LogToFile("debug.log", "debug")
			if err != nil {
				fmt.Println("fatal:", err)
				os.Exit(1)
			}
			defer f.Close()
		}
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
	rawString     textinput.Model
	decodedHeader string
	decodedClaims string
	decoded       string
	err           error
}

type decodeJWTStr struct {
	decoded       string
	decodedHeader string
	decodedClaims string
}

func DecodeJWT(raw string) (decodeJWTStr, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(raw, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return decodeJWTStr{}, errMsgJWT{err: err}
	}

	fmt.Println(token)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println(err)
		return decodeJWTStr{}, errMsgJWT{err: err}
	}

	jsonClaim, err := json.Marshal(claims)
	if err != nil {
		fmt.Println(err)
		return decodeJWTStr{}, errMsgJWT{err: err}
	}

	header := token.Header

	fmt.Sprint(header)
	var str = "{\n"
	for x, y := range header {
		str += fmt.Sprintf("  %v, %v\n", x, y)
	}
	str += "}"

	clipboard.Write(clipboard.FmtText, []byte(jsonClaim))
	return decodeJWTStr{decodedHeader: str, decodedClaims: string(jsonClaim)}, nil
}

func (m jwtModel) decodeJWTMsg() tea.Msg {
	decoded, err := DecodeJWT(m.rawString.Value())
	if err != nil {
		return errMsgJWT{err: err}
	}
	return decoded
}

func initialJWTModel() jwtModel {
	ti := textinput.New()
	ti.Placeholder = "JWT to decode"
	ti.Focus()

	return jwtModel{
		rawString:     ti,
		decodedHeader: "",
		decodedClaims: "",
		decoded:       "",
		err:           nil,
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
		case "enter":
			return m, m.decodeJWTMsg
		}

	case decodeJWTStr:
		decoded := msg
		m.decodedHeader = decoded.decodedHeader
		m.decodedClaims = decoded.decodedClaims
		return m, tea.Quit

	// We handle errors just like any other message
	case errMsgJWT:
		m.err = msg.err
		return m, nil

	}

	m.rawString, cmd = m.rawString.Update(msg)
	return m, cmd
}

func (m jwtModel) View() string {
	if len(m.decodedClaims) > 0 {
		return tui.ContainerStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				tui.LabelStyle.Render("Decoded JWT (No singature due to Unverified JWT):"),
				tui.Spacer.Render(""),
				lipgloss.NewStyle().Foreground(lipgloss.Color("1")).PaddingLeft(1).Render(m.decodedHeader),
				tui.Spacer.Render(""),
				lipgloss.NewStyle().Foreground(lipgloss.Color("2")).PaddingLeft(1).Render(m.decodedClaims),
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
		text := "Enter the JWT you want to debug."
		return lipgloss.JoinVertical(lipgloss.Left,
			tui.LabelStyle.Render(text),
			tui.Spacer.Render(""),
			lipgloss.NewStyle().PaddingLeft(1).Render(m.rawString.View()),
			tui.Spacer.Render(""),
			tui.ValueStyle.Render("(q to quit)"),
		)
	}
}
