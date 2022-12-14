/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jim-at-jibba/dtc/internal/utils"
	"github.com/jim-at-jibba/dtc/tui"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
	loremipsum "gopkg.in/loremipsum.v1"
)

type status int

/* MODEL MANAGMENT */
var models []tea.Model

const (
	selection status = iota
	form
)

// loremIpsumCmd represents the loremIpsum command
var loremIpsumCmd = &cobra.Command{
	Use:   "lorem-ipsum",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if f, err := tea.LogToFile("debug.log", "help"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			defer func() {
				err = f.Close()
				if err != nil {
					log.Fatal(err)
				}
			}()
		}

		models = []tea.Model{NewModel(), NewAmountModel("Word")}
		m := models[selection]
		// m := NewModel()
		p := tea.NewProgram(m)
		if err := p.Start(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loremIpsumCmd)
}

type errMsgLoremIpsum struct {
	err error
}

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsgLoremIpsum) Error() string { return e.err.Error() }

// Amount Model
type amountModel struct {
	amount        textinput.Model
	generatedText string
	textType      string
	err           error
}

type GeneratedText struct {
	text string
}

func GenerateLoremIpsom(textType, amount string) (GeneratedText, error) {
	loremIpsumGenerator := loremipsum.New()
	int, err := strconv.Atoi(amount)
	words := ""
	if err != nil {
		return GeneratedText{}, err
	}
	switch textType {
	case "Word":
		words = loremIpsumGenerator.Words(int)
	case "Sentence":
		words = loremIpsumGenerator.Sentences(int)
	case "Paragraph":
		words = loremIpsumGenerator.Paragraphs(int)

	}

	clipboard.Write(clipboard.FmtText, []byte(words))
	return GeneratedText{text: utils.TruncateString(words, 250)}, nil
}

func (m amountModel) generateLoremMsg() tea.Msg {
	text, err := GenerateLoremIpsom(m.textType, m.amount.Value())
	if err != nil {
		return errMsgLoremIpsum{err: err}
	}
	return text
}

func NewAmountModel(textType string) *amountModel {
	form := textinput.New()
	form.Focus()
	return &amountModel{
		amount:   form,
		textType: textType,
	}
}

func (m amountModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m amountModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		// case "esc":
		// 	models[selection] = NewModel()
		// 	return models[selection].Update(nil)
		case "enter":
			return m, m.generateLoremMsg
		}
	case GeneratedText:
		text := msg
		m.generatedText = text.text
		return m, tea.Quit
	case errMsg:
		m.err = msg.err
		return m, tea.Quit
	}

	m.amount, cmd = m.amount.Update(msg)
	return m, cmd
}

func (m amountModel) View() string {
	if m.generatedText == "" {
		return lipgloss.JoinVertical(lipgloss.Left,
			tui.LabelStyle.Render(fmt.Sprintf("How many %ss", m.textType)),
			tui.Spacer.Render(""),
			lipgloss.NewStyle().PaddingLeft(1).Render(m.amount.View()),
			tui.Spacer.Render(""),
			tui.ValueStyle.Render("(q to quit)"),
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
		return tui.ContainerStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				tui.LabelStyle.Render("Generate lorem ipsum (truncated for brevity):"),
				tui.Spacer.Render(""),
				tui.ValueStyle.Render(fmt.Sprintf("%s...", m.generatedText)),
				tui.LabelStyle.Render("(Copied to clipboard)"),
			),
		)
	}
}

// END Amount Model

// Main Model
type words struct {
	title       string
	description string
}

func (w words) FilterValue() string {
	return w.title
}

func (w words) Title() string {
	return w.title
}

func (w words) Description() string {
	return w.description
}

func NewModel() *loremIpsumModel {
	items := []list.Item{
		words{title: "Word", description: "Word type"},
		words{title: "Sentence", description: "I want a sentence"},
		words{title: "Paragraph", description: "Give me paragraphs"},
	}
	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "What type of text?"
	return &loremIpsumModel{
		list: list,
	}
}

type loremIpsumModel struct {
	list    list.Model
	choosen string
}

var _ tea.Model = (*loremIpsumModel)(nil)

func (m loremIpsumModel) Init() tea.Cmd {
	return nil
}

func (m loremIpsumModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(words)
			if ok {
				m.choosen = i.title
				models[form] = NewAmountModel(m.choosen)
				return models[form].Update(nil)
			}
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

// View renders output to the CLI.
func (m loremIpsumModel) View() string {
	return m.list.View()
}

// END Main Model
