package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
)

func coalesce(a, b string) string {
	if a != "" {
		return a
	}

	return b
}

func isSupportedFormat(formatString string) bool {
	formatString = strings.ToLower(formatString)

	if formatString == "" {
		return true
	}

	supportedFormats := map[string]bool{
		"png":  true,
		"jpeg": true,
		"jpg":  true,
	}

	if _, ok := supportedFormats[formatString]; !ok {
		return false
	}

	return true
}

func (a *app) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	a.Logger.Error(err.Error(), "method", method, "uri", uri)
	fmt.Println(string(debug.Stack()))

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (a *app) clientError(w http.ResponseWriter, message string, code int) {
	http.Error(w, message, code)
}

func (a *app) render(w http.ResponseWriter, r *http.Request, status int, page string) {
	ts, ok := a.TemplateCache[page]
	if !ok {
		a.serverError(w, r, fmt.Errorf("The template %s does not exist", page))
		return
	}

	buff := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buff, "base", nil)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")

	w.WriteHeader(status)
	buff.WriteTo(w)
}
