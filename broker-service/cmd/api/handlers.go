package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

// Create a predicatble JSON Format to adhere to

type RequestPayload struct {
	Action string      `json:"action"`
	Auth   AuthPayload `json:"auth, omitempty"`
}

type AuthPayload struct {
	Email    string `json:"string"`
	Password string `json:"string"`
}

func (app *Config) Broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonResponse{
		Error:   false,
		Message: "Hit the broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload)
}

func (app *Config) HandleSubmissions(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := app.readJSON(w, r, &requestPayload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	switch requestPayload.Action {
	case "auth":
		app.authenticate(w, requestPayload.Auth)

	default:
		app.errorJSON(w, errors.New("unknown action"))
	}
}

func (app *Config) authenticate(w http.ResponseWriter, a AuthPayload) {
	// create some json to send to auth microservce
	jsonData, _ := json.MarshalIndent(a, "", "\t")

	// call the service
	request, err := http.NewRequest(
		"POST",
		"http;//authentication-service/authenticate",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errorJSON(w, err)
	}
	defer response.Body.Close()
	// make sure we get nack the correct status code
	// be sure to check for unauthorized status
	// and check fro status accepted
	// otherwise youd have some other successful unauthorized code
	if response.StatusCode == http.StatusUnauthorized {
		app.errorJson(w, errors.New("Invalid Credentials"))
		return
	} else if response.StatusCode != http.StatusAccepted {
		app.errorJson(w, errors.New("error calling auth service"))
		return
	}

	// create a variable we'll read response.Bpdy into
	var jsonFromService jsonResponse
	// decode json from auth service

	err = json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errorJson(w, errors.New("Invalid Credentials"))
		return
	}

	if jsonFromService.Error {
		app.errorJSON(w, err, http.StatusUnauthorized)
		return
	}

	var payload jsonResponse
	payload.Error = false
	payload.Message = "Authenticated!"
	payload.Data = jsonFromService.Data

	app.writeJSON(w, http.StatusAccepted, payload)
}
