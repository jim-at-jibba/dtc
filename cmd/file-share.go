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
	"github.com/jim-at-jibba/dev-tools-cli/tui"
	"github.com/spf13/cobra"
)

// yeetrCmd represents the yeetr command
var fileShareCmd = &cobra.Command{
	Use:   "file-share",
	Short: "For sharing files.",
	Long:  "Sharing is, by design, ephemeral, so, the link that file-share provides will expire, after a given time or when the file is downloaded. The idea for this tool was lifted from https://www.npmjs.com/package/yeetr",
	Run: func(cmd *cobra.Command, args []string) {
		expires, _ := cmd.Flags().GetString("expires")

		p := tea.NewProgram(initialModel(expires))

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
	fileShareCmd.Flags().String("expires", "14d", "Expire time for the link. In the format d, w, m. e.g. 2w = 2 weeks")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// yeetrCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type errFileShareMsg struct {
	err error
}

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errFileShareMsg) Error() string { return e.err.Error() }

type yeetResponse struct {
	Link string
}

func GetFileShareUrl(fileName string, expires string) (string, error) {
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
	url := fmt.Sprintf("https://file.io/?expires=%s", expires)
	req, err := http.NewRequest("POST", url, bytes.NewReader(body.Bytes()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer rsp.Body.Close()
	decoder := json.NewDecoder(rsp.Body)
	var y yeetResponse
	err = decoder.Decode(&y)
	fmt.Println(y)
	if err != nil {
		return "", err
	}
	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	}
	return y.Link, nil
}

func (m model) fileShareMsg() tea.Msg {
	link, err := GetFileShareUrl(m.choice, m.expires)
	fmt.Println("WHAT", err)
	if err != nil {
		return errFileShareMsg{err: err}
	}
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
	expires string
	err     error
}

func initialModel(expires string) model {
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
	return model{spinner: s, list: list, expires: expires}
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
			return m, m.fileShareMsg
		}
	case tea.WindowSizeMsg:
		h, v := tui.ContainerNoBorderStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	case uploadUrl:
		m.link = msg.link
		m.yeeting = false
		return m, tea.Quit

	case errFileShareMsg:
		fmt.Println("In err mess")
		m.err = msg.err
		m.yeeting = false
		return m, tea.Quit
	}

	var cmd tea.Cmd
	m.list, _ = m.list.Update(msg)
	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if !m.yeeting && m.link == "" && m.err == nil {
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
	} else if m.err != nil {
		return tui.ErrorContainerStyle.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				tui.LabelStyle.Render("File sharing error:"),
				tui.Spacer.Render(""),
				tui.ValueStyle.Render(m.err.Error()),
			),
		)

	} else {
		return "Readt"
	}
}
