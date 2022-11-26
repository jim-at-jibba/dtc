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
	"github.com/jim-at-jibba/dtc/tui"
	"github.com/spf13/cobra"
	"golang.design/x/clipboard"
	"gopkg.in/loremipsum.v1"
)

type status int

/* MODEL MANAGMENT */
var models []tea.Model

const (
	selection status = iota
	form
)

func TruncateString(str string, length int) string {
	if length <= 0 {
		return ""
	}

	// This code cannot support Japanese
	// orgLen := len(str)
	// if orgLen <= length {
	//     return str
	// }
	// return str[:length]

	// Support Japanese
	// Ref: Range loops https://blog.golang.org/strings
	truncated := ""
	count := 0
	for _, char := range str {
		truncated += string(char)
		count++
		if count >= length {
			break
		}
	}
	return truncated
}

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

type GeneratedText struct {
	text string
}

func GenerateLoremIpsom(textType, amount string) string {
	loremIpsumGenerator := loremipsum.New()
	int, err := strconv.Atoi(amount)
	words := ""
	if err != nil {
		return ""
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
	return TruncateString(words, 500)
}

func (f amountModel) generateLoremMsg() tea.Msg {
	text := GenerateLoremIpsom(f.textType, f.amount.Value())
	return GeneratedText{text: text}
}

// Amount Model
type amountModel struct {
	amount        textinput.Model
	generatedText string
	textType      string
}

func NewAmountModel(textType string) *amountModel {
	form := textinput.New()
	form.Focus()
	return &amountModel{
		amount:   form,
		textType: textType,
	}
}

func (f amountModel) Init() tea.Cmd {
	return textinput.Blink
}

func (f amountModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return f, tea.Quit
		case "enter":
			return f, f.generateLoremMsg
		}
	case GeneratedText:
		text := msg
		f.generatedText = text.text
		return f, nil
	}

	f.amount, cmd = f.amount.Update(msg)
	return f, cmd
}

func (f amountModel) View() string {
	if f.generatedText == "" {
		return lipgloss.JoinVertical(lipgloss.Left,
			tui.LabelStyle.Render(fmt.Sprintf("How many %ss", f.textType)),
			tui.Spacer.Render(""),
			lipgloss.NewStyle().PaddingLeft(1).Render(f.amount.View()),
			tui.Spacer.Render(""),
			tui.ValueStyle.Render("(q to quit)"),
		)
	} else {
		return tui.ContainerStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				tui.LabelStyle.Render("Generate lorem ipsum (truncated for brevity):"),
				tui.Spacer.Render(""),
				tui.ValueStyle.Render(fmt.Sprintf("%s...", f.generatedText)),
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
	return &loremIpsumModel{}
}

type loremIpsumModel struct {
	list    list.Model
	loaded  bool
	choosen string
}

var _ tea.Model = (*loremIpsumModel)(nil)

func (m loremIpsumModel) Init() tea.Cmd {
	return nil
}

func (m *loremIpsumModel) initList(width, height int) {
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list.Title = "What type of text?"
	m.list.SetItems([]list.Item{
		words{title: "Word", description: "Word type"},
		words{title: "Sentence", description: "I want a sentence"},
		words{title: "Paragraph", description: "Give me paragraphs"},
	})
}

func (m loremIpsumModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !m.loaded {
			m.initList(msg.Width, msg.Height)
			m.loaded = true
		}
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
	if m.loaded {
		return m.list.View()
	} else {
		return tui.ContainerStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				tui.LabelStyle.Render("Loading..."),
			),
		)
	}
}

// END Main Model
