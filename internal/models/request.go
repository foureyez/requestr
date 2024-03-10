package models

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

type RequestModel struct {
	client        http.Client
	responseModel viewport.Model
	inputModels   []textinput.Model
	headersModel  table.Model
	requestTime   time.Duration
	currFocus     int
}

func NewRequestModal() RequestModel {
	vp := viewport.New(30, 5)
	inputs := make([]textinput.Model, 2)
	inputs[httpMethod] = textinput.New()
	inputs[httpMethod].Focus()
	inputs[httpMethod].CharLimit = 6
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

	headers := make(map[string]string)
	headers["User-Agent"] = "requestr-alpha"

	headersTable := table.New(table.WithColumns([]table.Column{{
		Title: "Key",
		Width: 20,
	}, {
		Title: "Value",
		Width: 50,
	}}), table.WithHeight(1))
	headersTable.SetRows([]table.Row{{
		"User-Agent", "requestr-alpha",
	}})

	return RequestModel{
		client:        http.NewClient(),
		inputModels:   inputs,
		headersModel:  headersTable,
		responseModel: vp,
	}
}

func (h *RequestModel) Init() tea.Cmd {
	return textinput.Blink
}

func (h *RequestModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			h.currFocus = (h.currFocus + 1) % 3
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

func (h *RequestModel) View() string {
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

func (h *RequestModel) renderRequestStat() string {
	return fmt.Sprintf("Last request took: %d ms", h.requestTime/time.Millisecond)
}

func (h *RequestModel) httpCmd() tea.Msg {
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

func (h *RequestModel) createRequest() http.Request {
	url := h.inputModels[url].Value()
	return http.Request{
		Url: url,
	}
}

func (h *RequestModel) validate() error {
	url := h.inputModels[url].Value()
	if err := validator.Url(url); err != nil {
		return err
	}
	return nil
}

func (h *RequestModel) changeFocus() {
	h.blurAll()
	switch h.currFocus {
	case httpMethod:
		h.inputModels[httpMethod].Focus()
	case url:
		h.inputModels[url].Focus()
	case headers:
		h.headersModel.Focus()
	case request:
	case response:
	}
}

func (h *RequestModel) blurAll() {
	for i := range h.inputModels {
		h.inputModels[i].Blur()
	}
}
