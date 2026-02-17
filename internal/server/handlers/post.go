package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.Auth.Authenticate(r.Context(), username, password)
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	sessionID, expires, err := h.Auth.CreateSession(r.Context(), user.ID)
	if err != nil {
		toastError(w, r, "500 Internal error: Could not create session.")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  expires,
	})

	w.Header().Set("HX-Redirect", "/dashboard")
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	// max file size: 100MB
	// r.Body = http.MaxBytesReader(w, r.Body, 100<<20) // aggresively strict max size
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		toastError(w, r, "The uploaded file is too big. Please choose an file that's less than 100MB in size")
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		toastError(w, r, "500 Internal error: Could not return file from form.")
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	_, err = file.Read(buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filetype := http.DetectContentType(buf)
	if filetype != "image/jpeg" && filetype != "image/png" {
		toastError(w, r, "The provided file format is not supported yet. Please upload a JPG or PNG image.")
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create the uploads dir if not exist
	dir, err := createUploadDirectory()
	if err != nil {
		toastError(w, r, "500 Internal error: "+err.Error())
		return
	}

	dst, err := os.Create(fmt.Sprintf("%s/%d%s", dir, time.Now().UnixNano(), filepath.Ext(fileHeader.Filename)))
	if err != nil {
		toastError(w, r, "500 Internal error: Could not create file.")
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	toastSuccess(w, r, "Your upload was Successful!")
}

// TODO: move this method somewhere else
func createUploadDirectory() (string, error) {
	uploadsDir := os.Getenv("UPLOADS_DIR")
	if uploadsDir == "" {
		panic("UPLOADS_DIR env var is not set")
	}
	dir := filepath.Join(uploadsDir, "uploads")

	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("could not create uploads dir: %w", err)
	}

	return dir, nil
}
