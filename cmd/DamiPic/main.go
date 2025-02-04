package main

import (
	"crypto/tls"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

type app struct {
	Logger        *slog.Logger
	TemplateCache map[string]*template.Template
}

func loadTemplates() (map[string]*template.Template, error) {
	cache := make(map[string]*template.Template)

	pages, err := filepath.Glob("./ui/html/*.tmpl.html")
	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		tmpl, err := template.ParseFiles(page)
		if err != nil {
			return cache, err
		}

		cache[name] = tmpl
	}

	return cache, nil
}

func main() {
	addr := flag.String("addr", ":8080", "The address to listen on")

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	templates, err := loadTemplates()
	if err != nil {
		log.Fatal(err)
	}

	app := &app{
		Logger:        logger,
		TemplateCache: templates,
	}

	// serve static files
	mux.Handle("GET /static/", http.StripPrefix("/static", neuter(fileServer)))

	// routes
	mux.HandleFunc("GET /img/{size}", app.placeholderImgHandler)
	mux.HandleFunc("GET /{$}", app.homehandler)

	server := http.Server{
		Addr:         *addr,
		Handler:      app.recoverPanic(app.logger(commonHeaders(gzipCompression(mux)))),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 15 * time.Second,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}

	logger.Info("Starting server", "addr", *addr)
	log.Fatal(server.ListenAndServeTLS("./cert.pem", "./key.pem"))
}
