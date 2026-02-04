package tools

import (
	"fmt"
	"net/http"
	"strings"
)

func AsciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		RenderErrorTemplate(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	text := r.FormValue("text")
	banner := r.FormValue("banner")
	color := SanitizeColor(r.FormValue("color"))

	if text == "" || banner == "" {
		RenderErrorTemplate(w, http.StatusBadRequest, "Missing text or banner")
		return
	}

	fontFile := fmt.Sprintf("banners/%s.txt", banner)
	result, err := PrintAsciiArt(text, fontFile)
	if err != nil {
		RenderErrorTemplate(w, http.StatusInternalServerError, fmt.Sprintf("Error generating ASCII art: %v", err))
		return
	}

	backgroundColor := GetContrastingBackground(color)

	data := struct {
		AsciiArt        string
		Color           string
		BackgroundColor string
	}{
		AsciiArt:        strings.Join(result, "\n"),
		Color:           color,
		BackgroundColor: backgroundColor,
	}

	if err := templates.ExecuteTemplate(w, "result.html", data); err != nil {
		RenderErrorTemplate(w, http.StatusInternalServerError, "Error rendering result")
	}
}
