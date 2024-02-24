package home

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type Home struct {
	err             error
	httpRes         *httpRes
	resView         viewport.Model
	urlInput        textinput.Model
	httpMethodInput textinput.Model
}

func NewHome() Home {
	vp := viewport.New(30, 5)
	url := textinput.New()
	url.Placeholder = "Enter URL"
	url.Focus()

	return Home{
		urlInput: url,
		resView:  vp,
	}
}

func (h Home) Init() tea.Cmd {
	return textinput.Blink
}

func (h Home) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	h.urlInput, tiCmd = h.urlInput.Update(msg)
	h.resView, vpCmd = h.resView.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return h, tea.Quit
		case tea.KeyEnter:
			return h, makeHttpCall(h.urlInput.Value())
		}
	case httpRes:
		res := httpRes(msg)
		h.resView.SetContent(fmt.Sprintf("%+v", res))
	case error:
		err := error(msg)
		h.resView.SetContent(fmt.Sprintf("%+v", err))
	}

	return h, tea.Batch(tiCmd, vpCmd)
}

func (h Home) View() string {
	return fmt.Sprintf(
		"%s\n\n%s",
		h.urlInput.View(),
		h.resView.View(),
	) + "\n\n"
}

func makeHttpCall(url string) tea.Cmd {
	return func() tea.Msg {
		c := &http.Client{Timeout: 10 * time.Second}
		res, err := c.Get(url)
		if err != nil {
			return err
		}
		body, _ := io.ReadAll(res.Body)
		return httpRes{
			statusCode: res.StatusCode,
			status:     res.Status,
			body:       body,
			headers:    res.Header,
		}
	}
}

type httpRes struct {
	headers    map[string][]string
	status     string
	body       []byte
	statusCode int
}
