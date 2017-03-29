package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/facebookgo/httpdown"
	"github.com/porthos-rpc/porthos-playground/collector"
	"github.com/porthos-rpc/porthos-playground/handlers"
	"github.com/porthos-rpc/porthos-playground/storage"
)

func defaultValue(a, b string) string {
	if len(a) == 0 {
		return b
	}

	return a
}

func main() {
	bindAddress := flag.String("bind", defaultValue(os.Getenv("BIND_ADDRESS"), ":3000"), "Bind Address.")
	brokerURL := flag.String("broker", defaultValue(os.Getenv("BROKER_URL"), "amqp://"), "Broker URL.")
	db := flag.String("db", defaultValue(os.Getenv("DB_PATH"), ":memory:"), "DB Path / Memory")

	flag.Parse()

	s := storage.NewStorage(storage.NewDb("sqlite3", *db))

	c := collector.NewCollector(*brokerURL, s)
	go c.Start()

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", handlers.IndexHandler)
	serveMux.HandleFunc("/api/services", handlers.NewServicesHandler(s))
	serveMux.HandleFunc("/api/rpc", handlers.NewRPCHandler(*brokerURL))
	serveMux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("client/public"))))

	server := &http.Server{
		Addr:           *bindAddress,
		Handler:        serveMux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	hd := &httpdown.HTTP{
		StopTimeout: 1 * time.Second,
		KillTimeout: 1 * time.Second,
	}

	fmt.Printf("Listening to %s\n", *bindAddress)
	fmt.Println("Hit CTRL-C to exit...")

	if err := httpdown.ListenAndServe(server, hd); err != nil {
		panic(err)
	}
}
