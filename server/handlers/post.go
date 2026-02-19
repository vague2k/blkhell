package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/vague2k/blkhell/server/database"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, err := h.Auth.Authenticate(r.Context(), username, password)
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	sessionToken, expires, err := h.Auth.CreateSession(r.Context(), user.ID)
	if err != nil {
		toastError(w, r, "500 Internal error: Could not create session.")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
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
	dir, err := createUploadDirectories()
	if err != nil {
		toastError(w, r, "500 Internal error: "+err.Error())
		return
	}

	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
	fileName := strings.TrimSuffix(fileHeader.Filename, fileExt)
	filePath := fmt.Sprintf(
		"%s/%s-%d%s",
		dir,
		fileName,
		time.Now().UnixNano(),
		fileExt,
	)

	dst, err := os.Create(filePath)
	if err != nil {
		toastError(w, r, "Could not create file.")
		return
	}
	defer dst.Close()
	size, err := io.Copy(dst, file)
	if err != nil {
		toastError(w, r, "Could not save file.")
		return
	}

	user, err := h.Auth.GetUserFromRequest(r)
	if err != nil {
		toastError(w, r, "Could not get user.")
		return
	}

	image, err := h.DB.CreateImage(r.Context(), database.CreateImageParams{
		ID:       uuid.NewString(),
		UserID:   user.ID,
		Path:     filePath,
		Filename: fileName,
		Ext:      strings.TrimPrefix(fileExt, "."),
		Size:     size,
	})
	if err != nil {
		toastError(w, r, "500 Internal error: Could not create image in database.")
		return
	}

	toastSuccess(w, r, fmt.Sprintf("'%s.%s' was uploaded successfully!", image.Filename, image.Ext))
}

// TODO: move this method somewhere else
func createUploadDirectories() (string, error) {
	uploadsDir := getDirectoryPathFromEnv("UPLOADS_DIR")

	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		return "", fmt.Errorf("could not create uploads dir: %w", err)
	}

	return uploadsDir, nil
}

func getDirectoryPathFromEnv(key string) string {
	dir := os.Getenv(key)
	if dir == "" {
		panic(fmt.Sprintf("%s env var is not set", key))
	}
	return dir
}
