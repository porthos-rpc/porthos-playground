package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/porthos-rpc/porthos-go"
	"github.com/porthos-rpc/porthos-playground/storage"
	"github.com/sirupsen/logrus"
)

type responseToClient struct {
	StatusCode  int32       `json:"statusCode"`
	ContentType string      `json:"contentType"`
	Body        interface{} `json:"body"`
}

type Handler struct {
	broker  *porthos.Broker
	clients map[string]*porthos.Client
	closed  bool
	m       sync.RWMutex
}

func NewHandler(url string) *Handler {
	broker, err := porthos.NewBroker(url)

	if err != nil {
		logrus.WithError(err).Fatal("Failed to create porthos broker.")
	}

	return &Handler{
		broker:  broker,
		clients: make(map[string]*porthos.Client, 15),
	}
}

func (h *Handler) Close() {
	h.m.Lock()
	defer h.m.Unlock()

	for _, client := range h.clients {
		client.Close()
	}

	h.broker.Close()

	h.closed = true
}

func (h *Handler) getClient(serviceName string) (*porthos.Client, error) {
	h.m.Lock()
	defer h.m.Unlock()

	if h.closed {
		return nil, errors.New("Handler is closed.")
	}

	client, ok := h.clients[serviceName]
	if ok {
		return client, nil
	}

	client, err := porthos.NewClient(h.broker, serviceName, 120)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create client.")
		return nil, err
	}

	h.clients[serviceName] = client

	return client, nil
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
			http.Error(w, fmt.Sprintf("Error getting the specs from the storage %s", err), http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(specs)

		if err != nil {
			http.Error(w, fmt.Sprintf("Error converting the specs to json %s", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}

// NewRPCHandler creates a new handler to return all specs.
func (h *Handler) NewRPCHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceName := r.PostFormValue("service")
		procedure := r.PostFormValue("procedure")
		contentType := r.PostFormValue("contentType")
		timeout, _ := strconv.Atoi(r.PostFormValue("timeout"))
		body := r.PostFormValue("body")

		client, err := h.getClient(serviceName)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error creating porthos client %s", err), http.StatusInternalServerError)
			return
		}

		fmt.Printf("RPC Call to Service %s, Procedure: %s, ContentType: %s", serviceName, procedure, contentType)

		// call the remote method.
		response, err := client.Call(procedure).
			WithTimeout(time.Duration(timeout)*time.Second).
			WithBodyContentType([]byte(body), contentType).
			Sync()

		if err != nil {
			http.Error(w, fmt.Sprintf("Error performing rpc request %s", err), http.StatusInternalServerError)
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
			http.Error(w, fmt.Sprintf("Error converting the response to json %s", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(json)
	}
}
