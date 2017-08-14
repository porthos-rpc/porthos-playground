package main

import (
	"crypto/subtle"
	"flag"
	"fmt"
	"github.com/facebookgo/httpdown"
	_ "github.com/mattn/go-sqlite3"
	"github.com/porthos-rpc/porthos-playground/collector"
	"github.com/porthos-rpc/porthos-playground/handlers"
	"github.com/porthos-rpc/porthos-playground/storage"
	"net/http"
	"os"
	"time"
)

func defaultValue(a, b string) string {
	if len(a) == 0 {
		return b
	}

	return a
}

func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {
	if username != "" && password != "" {
		return func(w http.ResponseWriter, r *http.Request) {
			user, pass, ok := r.BasicAuth()

			if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
				w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
				w.WriteHeader(401)
				w.Write([]byte("Unauthorised.\n"))
				return
			}

			handler(w, r)
		}
	}

	return handler
}

func main() {
	bindAddress := flag.String("bind", defaultValue(os.Getenv("BIND_ADDRESS"), ":3000"), "Bind Address.")
	brokerURL := flag.String("broker", defaultValue(os.Getenv("BROKER_URL"), "amqp://"), "Broker URL.")
	db := flag.String("db", defaultValue(os.Getenv("DB_PATH"), ":memory:"), "DB Path / Memory")
	username := *flag.String("username", defaultValue(os.Getenv("USERNAME"), ""), "HTTP Username")
	password := *flag.String("password", defaultValue(os.Getenv("PASSWORD"), ""), "HTTP Password")
	realm := "eventials"

	flag.Parse()

	s := storage.NewStorage(storage.NewDb("sqlite3", *db))

	c := collector.NewCollector(*brokerURL, s)
	go c.Start()

	serveMux := http.NewServeMux()

	serveMux.HandleFunc("/", BasicAuth(handlers.IndexHandler, username, password, realm))
	serveMux.HandleFunc("/api/services", BasicAuth(handlers.NewServicesHandler(s), username, password, realm))
	serveMux.HandleFunc("/api/rpc", BasicAuth(handlers.NewRPCHandler(*brokerURL), username, password, realm))
	serveMux.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

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
