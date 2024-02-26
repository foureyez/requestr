package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
	"github.com/rs/zerolog/log"
)

func main() {
	app := tview.NewApplication()
	newPrimitive := func(text string) *tview.TextView {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}
	menu := newPrimitive("Right Side")
	main := newPrimitive("Main content")
	sideBar := newPrimitive("Side Bar")

	urlInputField := tview.NewInputField().
		SetPlaceholder("Enter url")

	urlInputField.SetDoneFunc(func(key tcell.Key) {
		url := urlInputField.GetText()
		res, err := makeHttpCall(url)
		if err != nil {
			main.SetText(err.Error())
		} else {
			main.SetText(fmt.Sprintf("%+v", res))
		}
	})

	grid := tview.NewGrid().
		SetRows(3, 0).
		SetColumns(30, 0, 30).
		SetBorders(true).
		AddItem(urlInputField, 0, 0, 1, 3, 1, 0, true)

	grid.SetBordersColor(tcell.Color44)

	// Layout for screens narrower than 100 cells (menu and side bar are hidden).
	grid.AddItem(menu, 0, 0, 0, 0, 0, 0, false).
		AddItem(main, 1, 0, 1, 3, 0, 0, false).
		AddItem(sideBar, 0, 0, 0, 0, 0, 0, false)

	// Layout for screens wider than 100 cells.
	grid.AddItem(menu, 1, 0, 1, 1, 0, 100, false).
		AddItem(main, 1, 1, 1, 1, 0, 100, false).
		AddItem(sideBar, 1, 2, 1, 1, 0, 100, false)

	if err := app.SetRoot(grid, true).Run(); err != nil {
		log.Fatal().Msgf("Unable to start application: %s", err.Error())
	}
}

func makeHttpCall(url string) (*httpRes, error) {
	c := &http.Client{Timeout: 10 * time.Second}
	res, err := c.Get(url)
	if err != nil {
		return nil, err
	}
	body, _ := io.ReadAll(res.Body)
	return &httpRes{
		statusCode: res.StatusCode,
		status:     res.Status,
		body:       body,
		headers:    res.Header,
	}, nil
}

type httpRes struct {
	headers    map[string][]string
	status     string
	body       []byte
	statusCode int
}
