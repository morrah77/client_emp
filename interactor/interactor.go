package interactor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"morrah77.com/client_emp/common"
)

type Interactor struct {
	client common.Client
	inputStream io.Reader
	outputStream io.Writer
	errorStream io.Writer
}

func NewInteractor(client common.Client, inputStream io.Reader, outputStream io.Writer,
	errorStream io.Writer) (*Interactor, error) {

	return &Interactor{
		client:       client,
		inputStream:  inputStream,
		outputStream: outputStream,
		errorStream:  errorStream,
	}, nil
}

func(it *Interactor) Run() int {
	res, err := it.client.FetchList()
	if err != nil {
		it.outputError(err, "An error occured: %s\nShutting down...\n")
		return 1
	}
	it.outputResult(res, "List")

	var input string
LOOP:
	for {
		it.outputHint()
		_, err := fmt.Fscan(it.inputStream, &input)
		if err != nil {
			it.outputError(err, "An error occured: %s\nShutting down...")
			return 1
		}
		switch input {
		case "q":
			break LOOP
		case "l":
			res, err := it.client.FetchList()
			if err != nil {
				it.outputError(err, "An error occured: %s\nPlease try again")
			}
			it.outputResult(res, "List")
		default:
			res, err := it.client.FetchItem(input)
			if err != nil {
				it.outputError(err, "An error occured: %s\nPlease try again")
			}
			if res != nil {
				it.outputResult(res, "Item")
			}
		}
	}
	return 0
}

func(it *Interactor) outputHint() {
	_, _ = fmt.Fprintln(it.outputStream, "\n" +
		"Please type ID to view item details,\n"+
		"`l` to view list,\n"+
		"`q` to quit")
}

func(it *Interactor) outputResult(res []byte, message string) {
	b := &bytes.Buffer{}
	err := json.Indent(b, res, "", "\t")
	if err != nil {
		_, _ = fmt.Fprintln(it.outputStream, fmt.Sprintf("%s:\n%s", message, string(res)))
	} else {
		_, _ = fmt.Fprintln(it.outputStream, fmt.Sprintf("%s:\n%s", message, string(b.Bytes())))
	}
}

func(it *Interactor) outputError(err error, message string)  {
	_, _ = fmt.Fprintln(it.errorStream, fmt.Sprintf(message, err.Error()))
}

