package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

func main() {
	srv := http.Server{
		Addr:              ":4444",
		ReadTimeout:       30 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		Handler:           handler(),
	}

	fmt.Println("Listening on :4444")
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println("ERROR: listen and serve:", err)
		os.Exit(1)
	}
}

func handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/test", handlerTest)

	return mux
}

type Payload struct {
	Action      string      `json:"action"`
	PullRequest PullRequest `json:"pull_request"`
}

type PullRequest struct {
	State  string `json:"state"`
	User   User   `json:"user"`
	Merged bool   `json:"merged"`
}

type User struct {
	Login string `json:"login"`
}

var (
	approvedUsers = []string{"rms1000watt"}
)

func handlerTest(w http.ResponseWriter, r *http.Request) {
	payload, err := getPayload(r)
	if err != nil {
		errStr := "Failed getting payload: " + err.Error()
		fmt.Println(errStr)
		http.Error(w, errStr, http.StatusInternalServerError)
		return
	}

	if err = validatePayload(payload); err != nil {
		errStr := "Invalid payload: " + err.Error()
		fmt.Println(errStr)
		http.Error(w, errStr, http.StatusUnauthorized)
		return
	}

	fmt.Println(payload)
}

func getPayload(r *http.Request) (payload Payload, err error) {
	payload = Payload{}
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		fmt.Println("Failed unmarshalling payload:", err)
		return
	}

	return
}

func validatePayload(payload Payload) (err error) {
	if !approvedUser(payload.PullRequest.User.Login) {
		err = fmt.Errorf("Not an approved user: %s", payload.PullRequest.User.Login)
		return
	}

	if payload.Action != "closed" {
		err = fmt.Errorf("Not an approved action: %s", payload.Action)
		return
	}

	if !payload.PullRequest.Merged {
		err = fmt.Errorf("Not merged PR")
		return
	}

	if payload.PullRequest.State != "closed" {
		err = fmt.Errorf("Not an approved PR state: %s", payload.PullRequest.State)
		return
	}

	return
}

func approvedUser(in string) (approved bool) {
	for _, user := range approvedUsers {
		if user == in {
			approved = true
			return
		}
	}

	return
}
