package example

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/okeefem2/ok-httpclient/okhttp"
)

func TestMain(m *testing.M) {
	fmt.Println("Starting tests for ok http client examples")
		// Tell library to mock requests
		okhttp.StartMockServer()
		os.Exit(m.Run())
}

func TestGet(t* testing.T) {
	t.Run("TestErrorFetching", func(t *testing.T) {
		key := "0fb9c447-ae36-41ba-bfa5-57f89d08160a"

		okhttp.AddMock(okhttp.Mock{
			Method: http.MethodGet,
			Url: "http://localhost:5005/games/" + key,
			Error: errors.New("Timeout fetching game"),
		})

		game, err := GetGame(key)

		if game != nil {
			t.Error("No game expected")
		}

		if err == nil {
			t.Error("Expected error")
		}

		if err.Error() != "Timeout fetching game" {
			t.Error("Invalid error message")
		}

	});

	t.Run("TestErrorUnmarshalResponseBody", func(t *testing.T) {
		key := "0fb9c447-ae36-41ba-bfa5-57f89d08160a"

		okhttp.AddMock(okhttp.Mock{
			Method: http.MethodGet,
			Url: "http://localhost:5005/games/" + key,
			ResponseBody: `{ "game": true }`,
		})

		game, err := GetGame(key)

		if game != nil {
			t.Error("No game expected")
		}

		if err == nil {
			t.Error("Expected error")
		}

		if !strings.Contains(err.Error(), "cannot unmarshal") {
			t.Error("Invalid error message")
		}

	});

	t.Run("TestNoError", func(t *testing.T) {
		key := "0fb9c447-ae36-41ba-bfa5-57f89d08160a"

		okhttp.AddMock(okhttp.Mock{
			Method: http.MethodGet,
			Url: "http://localhost:5005/games/" + key,
			ResponseBody: `{ "game": { "key": "0fb9c447-ae36-41ba-bfa5-57f89d08160a" } }`,
		})

		game, err := GetGame(key)

		if game == nil {
				t.Error("Expected game to be defined")
				return
		}

		if game.Game["key"] != key {
			t.Error("Expected key to be " + key)
		}

		if err != nil {
			t.Error("Expected no error")
		}
	});

}
