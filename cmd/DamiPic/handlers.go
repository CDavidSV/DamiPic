package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"
	"strings"
)

func (a *app) homehandler(w http.ResponseWriter, r *http.Request) {
	a.render(w, r, 200, "index.tmpl.html")
}

func (a *app) placeholderImgHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	imgSizeInput := r.PathValue("size")
	text := params.Get("text")
	format := coalesce(params.Get("format"), "png")
	// bgColor := params.Get("bg-color")
	// textColor := params.Get("text-color")

	// We validate that the size is in the correct format (e.g. 300x200)
	size := strings.Split(strings.ToLower(imgSizeInput), "x")
	if len(size) != 2 {
		a.clientError(w, "Invalid size path parameter provided", http.StatusBadRequest)
		return
	}

	width, err := strconv.Atoi(size[0])
	if err != nil {
		a.clientError(w, "Width parameter must be a number. Got: "+size[0], http.StatusBadRequest)
		return
	}

	height, err := strconv.Atoi(size[1])
	if err != nil {
		a.clientError(w, "Height parameter must be a number. Got: "+size[1], http.StatusBadRequest)
		return
	}

	if height > 4000 || width > 4000 || height < 10 || width < 10 {
		a.clientError(w, "width and height values must be between 10px and 4000px", http.StatusBadRequest)
		return
	}

	if len(text) > 120 {
		a.clientError(w, "Text must be less than 120 characters long", http.StatusBadRequest)
		return
	}

	if text == "" {
		text = fmt.Sprintf("%d×%d", width, height)
	}

	img, err := generatePlaceholderImg(width, height, text)
	if err != nil {
		a.serverError(w, r, err)
		return
	}

	switch format {
	case "jpg":
	case "jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
		jpeg.Encode(w, img, &jpeg.Options{
			Quality: 80,
		})
		break
	case "png":
		w.Header().Set("Content-Type", "image/png")
		png.Encode(w, img)
		break
	default:
		a.clientError(w, format+" is not a valid or supported image format: Try using png, jpeg or jpg", http.StatusBadRequest)
		break
	}
}
