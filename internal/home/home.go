package home

import (
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/foureyez/requestr/internal/http"
	"github.com/foureyez/requestr/internal/validator"
)

const (
	httpMethod = iota
	url
)

type Home struct {
	client  http.Client
	resView viewport.Model
	inputs  []textinput.Model
}

func NewHome() Home {
	vp := viewport.New(30, 5)
	inputs := make([]textinput.Model, 2)
	inputs[httpMethod] = textinput.New()
	inputs[httpMethod].CharLimit = 5
	inputs[httpMethod].Width = 5
	inputs[httpMethod].Prompt = ""
	inputs[httpMethod].SetValue(http.Get.String())
	inputs[httpMethod].SetSuggestions([]string{http.Get.String(), http.Post.String(), http.Put.String()})

	inputs[url] = textinput.New()
	inputs[url].Placeholder = "Enter Url"
	inputs[url].CharLimit = 200
	inputs[url].Width = 30
	inputs[url].Prompt = ""
	// inputs[url].Validate = validator.Url
	inputs[url].Focus()

	return Home{
		client:  http.NewClient(),
		inputs:  inputs,
		resView: vp,
	}
}

func (h Home) Init() tea.Cmd {
	return textinput.Blink
}

func (h Home) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(h.inputs))
	var cm tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return h, tea.Quit
		case tea.KeyEnter:
			return h, h.httpCmd
		}
	case http.Response:
		h.resView.SetContent(fmt.Sprint(msg))
	case error:
		h.resView.SetContent(fmt.Sprint(msg))
	}

	for i := range h.inputs {
		h.inputs[i], cmds[i] = h.inputs[i].Update(msg)
	}
	cmds = append(cmds, cm)

	return h, tea.Batch(cmds...)
}

func (h Home) View() string {
	return appStyle.Render(lipgloss.JoinVertical(
		lipgloss.Center,
		inputStyle.Render(h.inputs[httpMethod].View()),
		inputStyle.Render(h.inputs[url].View()),
		inputStyle.Render(h.resView.View()),
	))
}

func (h Home) httpCmd() tea.Msg {
	if err := h.validate(); err != nil {
		return err
	}
	r := h.createRequest()
	h.resView.SetContent(fmt.Sprintf("%v", r))
	res, err := h.client.Execute(r)
	if err != nil {
		return err
	}
	return res
}

func (h Home) createRequest() http.Request {
	url := h.inputs[url].Value()
	return http.Request{
		Url: url,
	}
}

func (h Home) validate() error {
	url := h.inputs[url].Value()
	if err := validator.Url(url); err != nil {
		return err
	}
	return nil
}
