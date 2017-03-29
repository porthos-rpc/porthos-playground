package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/porthos-rpc/porthos-go"
	"github.com/porthos-rpc/porthos-playground/storage"
)

type responseToClient struct {
	StatusCode  int16       `json:"statusCode"`
	ContentType string      `json:"contentType"`
	Body        interface{} `json:"body"`
}

// IndexHandler will display the dashboard index page.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("templates/index.html")
	t.Execute(w, nil)
}

// NewServicesHandler creates a new handler to return all specs.
func NewServicesHandler(storage storage.Storage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		specs, err := storage.GetSpecs()

		if err != nil {
			fmt.Errorf("Error getting the specs from the storage %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(specs)

		if err != nil {
			fmt.Errorf("Error converting the specs to json %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}

// NewRPCHandler creates a new handler to return all specs.
func NewRPCHandler(amqpURL string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceName := r.PostFormValue("service")
		procedure := r.PostFormValue("procedure")
		contentType := r.PostFormValue("contentType")
		timeout, _ := strconv.Atoi(r.PostFormValue("timeout"))
		body := r.PostFormValue("body")

		b, err := porthos.NewBroker(amqpURL)

		if err != nil {
			fmt.Errorf("Error creating porthos broker %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer b.Close()

		service, err := porthos.NewClient(b, serviceName, time.Duration(timeout)*time.Second)

		if err != nil {
			fmt.Errorf("Error creating porthos client %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		defer service.Close()

		fmt.Printf("RPC Call to Service %s, Procedure: %s, ContentType: %s", serviceName, procedure, contentType)

		// call a remote method that is "void".
		response, err := service.Call(procedure).WithBodyContentType([]byte(body), contentType).Sync()

		if err != nil {
			fmt.Errorf("Error performing rpc request %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var responseBody interface{}

		if response.ContentType == "application/json" {
			response.UnmarshalJSONTo(&responseBody)
		} else {
			responseBody = string(response.Content)
		}

		json, err := json.Marshal(responseToClient{
			StatusCode:  response.StatusCode,
			ContentType: response.ContentType,
			Body:        responseBody,
		})

		if err != nil {
			fmt.Errorf("Error converting the response to json %s", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}
