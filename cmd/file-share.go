/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jim-at-jibba/devtools/tui"
	"github.com/spf13/cobra"
)

// yeetrCmd represents the yeetr command
var fileShareCmd = &cobra.Command{
	Use:   "file-share",
	Short: "For sharing files.",
	Long:  "Sharing is, by design, ephemeral, so, the link that file-share provides will expire, after a given time or when the file is downloaded. The idea for this tool was lifted from https://www.npmjs.com/package/yeetr",
	Run: func(cmd *cobra.Command, args []string) {

		p := tea.NewProgram(initialModel())

		if err := p.Start(); err != nil {
			fmt.Println("Error running program:", err)
			os.Exit(1)
		}

	},
}

func init() {
	rootCmd.AddCommand(fileShareCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// yeetrCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// yeetrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type yeetResponse struct {
	Link string
}

func Yeet(fileName string) (string, error) {
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	fw, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", err
	}
	file, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	_, err = io.Copy(fw, file)
	if err != nil {
		return "", err
	}
	writer.Close()
	req, err := http.NewRequest("POST", "https://file.io/?expires=2m", bytes.NewReader(body.Bytes()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, err := client.Do(req)
	check(err)
	defer rsp.Body.Close()
	// b, err := io.ReadAll(rsp.Body)
	// check(err)
	decoder := json.NewDecoder(rsp.Body)
	var y yeetResponse
	err = decoder.Decode(&y)
	check(err)
	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	}
	return y.Link, nil
}

func (m model) yeetMsg() tea.Msg {
	fmt.Println("Yeeting...", m.choice)
	link, _ := Yeet(m.choice)
	return uploadUrl{link: link}
}

type item struct {
	name, size string
}

func (i item) Title() string       { return i.name }
func (i item) Description() string { return i.size }
func (i item) FilterValue() string { return i.name }

type model struct {
	list    list.Model
	spinner spinner.Model
	choice  string
	yeeting bool
	link    string
}

func initialModel() model {
	// Set up spinner
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = tui.Spinner

	// Set up initial list of files
	items := []list.Item{}
	files, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if !f.IsDir() {
			items = append(items, item{name: f.Name(), size: strconv.FormatInt(f.Size()/1000, 10) + "KB"})
		}
	}

	list := list.New(items, list.NewDefaultDelegate(), 0, 0)
	list.Title = "What file do you want to share?"
	return model{spinner: s, list: list}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

type uploadUrl struct {
	link string
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c":
			return m, tea.Quit

		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				m.choice = string(i.name)
				m.yeeting = true
			}
			return m, m.yeetMsg
		}
	case tea.WindowSizeMsg:
		h, v := tui.ContainerNoBorderStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case uploadUrl:
		m.link = msg.link
		m.yeeting = false
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.list, _ = m.list.Update(msg)
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if !m.yeeting && m.link == "" {
		return tui.ContainerNoBorderStyle.Render(m.list.View())
	} else if m.yeeting {
		return tui.ContainerNoBorderStyle.
			Render(
				lipgloss.JoinHorizontal(lipgloss.Left,
					m.spinner.View(),
					tui.LabelStyle.Render("Generating link...press q to quit"),
				),
			)
	} else if len(m.link) > 0 {
		return tui.ContainerStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				tui.LabelStyle.Render("Ephemeral file link:"),
				tui.Spacer.Render(""),
				tui.ValueStyle.Render(m.link),
			),
		)
	} else {
		return "Readt"
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
