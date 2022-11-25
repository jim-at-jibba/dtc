/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
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

		state, err := NewModel()
		if err != nil {
			fmt.Printf(fmt.Sprintf("Error starting init command: %s\n", err))
			os.Exit(1)
		}
		// tea.NewProgram starts the Bubbletea framework which will render our
		// application using our state.
		if err := tea.NewProgram(state).Start(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(loremIpsumCmd)
}

// List Model

// END List Model

// Amount Model

// END Amount Model

// Main Model
type Words struct {
	title       string
	description string
}

func (w Words) FilterValue() string {
	return w.title
}

func (w Words) Title() string {
	return w.title
}

func (w Words) Description() string {
	return w.description
}

func NewModel() (*loremIpsumModel, error) {
	return &loremIpsumModel{}, nil
}

type loremIpsumModel struct {
	list   list.Model
	loaded bool
}

var _ tea.Model = (*loremIpsumModel)(nil)

func (m loremIpsumModel) Init() tea.Cmd {
	return nil
}

func (m *loremIpsumModel) initList(width, height int) {
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), width, height)
	m.list.Title = "What type?"
	m.list.SetItems([]list.Item{
		Words{title: "Word", description: "Word type"},
		Words{title: "Sentence", description: "I want a sentence"},
		Words{title: "Paragraph", description: "Give me paragraphs"},
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
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyCtrlBackslash:
			return m, tea.Quit
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
		return "loading..."
	}
}

// END Main Model
