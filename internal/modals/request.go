package modals

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
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
	headers
	request
	response
)

type RequestModal struct {
	client        http.Client
	headers       map[string]string
	responseModel viewport.Model
	inputModels   []textinput.Model
	headersModel  table.Model
	requestTime   time.Duration
	currFocus     int
}

func NewRequestModal() RequestModal {
	vp := viewport.New(30, 5)
	inputs := make([]textinput.Model, 2)
	inputs[httpMethod] = textinput.New()
	inputs[httpMethod].Focus()
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

	headers := table.New(table.WithColumns([]table.Column{{
		Title: "Key",
		Width: 20,
	}, {
		Title: "Value",
		Width: 50,
	}}), table.WithHeight(1))

	return RequestModal{
		client:        http.NewClient(),
		inputModels:   inputs,
		headersModel:  headers,
		responseModel: vp,
	}
}

func (h *RequestModal) Init() tea.Cmd {
	return textinput.Blink
}

func (h *RequestModal) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(h.inputModels))
	var cm tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return h, tea.Quit
		case tea.KeyEnter:
			return h, h.httpCmd
		case tea.KeyTab:
			h.currFocus = (h.currFocus + 1) % 2
			h.changeFocus()
		}
	case http.Response:
		h.responseModel.SetContent(fmt.Sprint(msg))
	case error:
		h.responseModel.SetContent(fmt.Sprint(msg))
	}

	for i := range h.inputModels {
		h.inputModels[i], cmds[i] = h.inputModels[i].Update(msg)
	}
	cmds = append(cmds, cm)

	return h, tea.Batch(cmds...)
}

func (h *RequestModal) View() string {
	return appStyle.Render(lipgloss.JoinVertical(
		lipgloss.Center,
		lipgloss.JoinHorizontal(
			lipgloss.Left,
			httpMethodStyle.Render(h.inputModels[httpMethod].View()),
			inputStyle.Render(h.inputModels[url].View()),
		),
		headerStyle.Render(h.headersModel.View()),
		inputStyle.Render(h.responseModel.View()),
		h.renderRequestStat(),
	))
}

func (h *RequestModal) renderRequestStat() string {
	return fmt.Sprintf("Last request took: %d ms", h.requestTime/time.Millisecond)
}

func (h *RequestModal) httpCmd() tea.Msg {
	if err := h.validate(); err != nil {
		return err
	}
	startTime := time.Now()
	defer func() {
		h.requestTime = time.Since(startTime)
	}()
	r := h.createRequest()
	res, err := h.client.Execute(r)
	if err != nil {
		return err
	}
	return res
}

func (h *RequestModal) createRequest() http.Request {
	url := h.inputModels[url].Value()
	return http.Request{
		Url: url,
	}
}

func (h *RequestModal) validate() error {
	url := h.inputModels[url].Value()
	if err := validator.Url(url); err != nil {
		return err
	}
	return nil
}

func (h *RequestModal) changeFocus() {
	h.blurAll()
	switch h.currFocus {
	case httpMethod:
		h.inputModels[httpMethod].Focus()
	case url:
		h.inputModels[url].Focus()
	case request:
	case response:
	}
}

func (h *RequestModal) blurAll() {
	for i := range h.inputModels {
		h.inputModels[i].Blur()
	}
}
