package interactor

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/golang/mock/gomock"
	"io/ioutil"
	"morrah77.com/client_emp/common/mock_common"
	"testing"
)

// TODO generalise mock client setup
func TestNewInteractor(t *testing.T) {
	in := bufio.NewReader(&bytes.Buffer{})
	out := bufio.NewWriter(&bytes.Buffer{})
	er := bufio.NewWriter(&bytes.Buffer{})
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mock_common.NewMockClient(ctrl)
	client.
		EXPECT().
		FetchList().
		Return(nil, nil).
		Times(0)
	client.
		EXPECT().
		FetchItem(gomock.Any()).
		Return(nil, nil).
		Times(0)
	interacror, err := NewInteractor(client, in, out, er)
	if err != nil {
		t.Fatalf("Could not create interacror: %s", err.Error())
	}
	if interacror.client != client {
		t.Errorf("client is wrong:\nexpected: %v\ngot: %v", client, interacror.client)
	}
	if interacror.inputStream != in {
		t.Errorf("inputStream is wrong:\nexpected: %v\ngot: %v", in, interacror.inputStream)
	}
	if interacror.outputStream != out {
		t.Errorf("outputStream is wrong:\nexpected: %v\ngot: %v", out, interacror.outputStream)
	}
	if interacror.errorStream != er {
		t.Errorf("errorStream is wrong:\nexpected: %v\ngot: %v", er, interacror.errorStream)
	}
}

// TODO split to separate cases, generalise mock client setup
func TestInteractor_Run(t *testing.T) {
	listJson := `{"foo":"bar"}`
	listItemOK := `{"id":1}`
	listItemErr := `{"status":"error"}`
	errFetchItemError := errors.New("Status 400")
	errInvalidId := errors.New("Invalid ID")
	expectedErrorContents, err := ioutil.ReadFile("./test_data/interactor_error.txt")
	if err != nil {
		t.Fatalf("Could not read from file: %s", err.Error())
	}
	expectedOutContents, err := ioutil.ReadFile("./test_data/interactor_output.txt")
	if err != nil {
		t.Fatalf("Could not read from file: %s", err.Error())
	}
	bufIn := &bytes.Buffer{}
	bufOut := &bytes.Buffer{}
	bufErr := &bytes.Buffer{}
	in := bufio.NewReader(bufIn)
	out := bufio.NewWriter(bufOut)
	er := bufio.NewWriter(bufErr)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := mock_common.NewMockClient(ctrl)
	client.
		EXPECT().
		FetchList().
		Return([]byte(listJson), nil).
		Times(2)
	client.
		EXPECT().
		FetchItem("1").
		Return([]byte(listItemOK), nil).
		Times(1)
	client.
		EXPECT().
		FetchItem("2").
		Return([]byte(listItemErr), errFetchItemError).
		Times(1)
	client.
		EXPECT().
		FetchItem("foo").
		Return([]byte(listItemErr), errInvalidId).
		Times(1)

	interacror, err := NewInteractor(client, in, out, er)
	if err != nil {
		t.Fatalf("Could not create interacror: %s", err.Error())
	}
	bufIn.Write([]byte("\n1\n2\nl\nfoo\nq\n"))
	interacror.Run()

	_ = out.Flush()
	if string(expectedOutContents) != bufOut.String() {
		t.Errorf("outputStream is wrong:\nexpected: %v\ngot: %v", string(expectedOutContents), bufOut.String())
	}
	_ = er.Flush()
	if string(expectedErrorContents) != bufErr.String() {
		t.Errorf("outputStream is wrong:\nexpected: %v\ngot: %v", string(expectedErrorContents), bufErr.String())
	}
}
