package main

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
)

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
