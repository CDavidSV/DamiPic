package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type app struct {
	Logger *slog.Logger
}

func main() {
	addr := flag.String("addr", ":8080", "The address to listen on")

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	app := &app{
		Logger: logger,
	}

	mux.Handle("GET /static/", http.StripPrefix("/static", neuter(fileServer)))
	mux.HandleFunc("GET /img/{size}", app.placeholderImgHandler)

	server := http.Server{
		Addr:         *addr,
		Handler:      app.recoverPanic(app.logger(commonHeaders(gzipCompression(mux)))),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
	}

	log.Fatal(server.ListenAndServe())
}
