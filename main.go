package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()

	window := myApp.NewWindow("Simple HTTP Request Tool")
	WindowAdjustment(window)

	selectedProtocol := "http"
	protocolRadio := widget.NewRadioGroup([]string{"http", "https"}, func(selected string) {
		selectedProtocol = selected
	})

	protocolRadio.SetSelected(selectedProtocol)

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("URL")

	methodSelector := widget.NewSelect([]string{"GET", "POST", "PUT", "DELETE"}, nil)
	methodSelector.SetSelected("GET")

	defaultHeader := binding.NewString()
	defaultHeader.Set("{}")
	headersEntry := widget.NewEntryWithData(defaultHeader)
	headersEntry.SetPlaceHolder("Headers (JSON)")

	bodyEntry := widget.NewMultiLineEntry()
	bodyEntry.SetPlaceHolder("Request Body")

	resultLabel := widget.NewLabel("Response")
	resultLabel.Wrapping = fyne.TextWrapWord

	resultContent := container.NewVBox(
		resultLabel,
	)

	sendButton := widget.NewButton("Send Request", func() {
		url := urlEntry.Text
		method := methodSelector.Selected
		headersJSON := headersEntry.Text
		if headersJSON == "" {
			headersJSON = "{}"
		}
		body := bodyEntry.Text

		headers := make(map[string]string)
		if err := json.Unmarshal([]byte(headersJSON), &headers); err != nil {
			resultLabel.SetText("Invalid Headers JSON")
			return
		}

		request := &CustomRequest{
			URL:     CombineUrl(selectedProtocol, url),
			Method:  method,
			Headers: headers,
			Body:    body,
		}

		resp, err := request.Do()
		if err != nil {
			resultLabel.SetText(fmt.Sprintf("Error: %s", err))
			resultLabel.Show()
			return
		}
		responseBody := "Response Body:\n" + readResponseBody(resp)
		responseStatus := fmt.Sprintf("Response Status: %s", resp.Status)

		resultLabel.SetText(responseStatus + "\n" + responseBody)

		if resp.StatusCode >= 400 {
			responseMessage := http.StatusText(resp.StatusCode)
			resultLabel.SetText(resultLabel.Text + "\n" + responseMessage)
			resultLabel.Show()
		}
	})

	content := container.NewVBox(
		protocolRadio,
		urlEntry,
		methodSelector,
		headersEntry,
		bodyEntry,
		sendButton,
		resultContent,
	)
	// content.Add(resultContent)
	window.SetContent(content)
	window.ShowAndRun()
}

func WindowAdjustment(window fyne.Window) {
	window.Resize(fyne.NewSize(800, 500))
	window.CenterOnScreen()
}
