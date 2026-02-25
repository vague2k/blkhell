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

var (
	MimeJpeg = "image/jpeg"
	MimePng  = "image/png"
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
	user, ok := h.Auth.UserFromContext(r.Context())
	if !ok {
		toastError(w, r, "Could not get user.")
		return
	}
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

	// create the uploads dir if not exist
	filetype := http.DetectContentType(buf)
	dir, err := createUploadDirectories(filetype)
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fileExt := strings.ToLower(filepath.Ext(fileHeader.Filename))
	fileName := strings.TrimSuffix(fileHeader.Filename, fileExt)
	filePath := fmt.Sprintf(
		"%s/%s-%d%s",
		dir,
		strings.ReplaceAll(fileName, " ", "-"), // replaces spaces with dashes due to browser nonsense
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

	_, parentDir, ok := strings.Cut(filePath, "uploads")
	if !ok {
		toastError(w, r, "500 internal error: what the fuck?")
		fmt.Println(parentDir)
		return
	}
	image, err := h.DB.CreateFile(r.Context(), database.CreateFileParams{
		ID:       uuid.NewString(),
		UserID:   user.ID,
		Path:     parentDir,
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
func createUploadDirectories(mimetype string) (string, error) {
	// TODO: move to main func somehow
	uploadsDir := os.Getenv("UPLOADS_DIR")
	if uploadsDir == "" {
		panic("UPLOADS_DIR env var is not set")
	}
	var uploadsWithSubDir string

	switch mimetype {
	case MimeJpeg, MimePng:
		uploadsWithSubDir = filepath.Join(uploadsDir, "images")
	default:
		return "", fmt.Errorf("The provided file format is not supported yet.")
	}

	if err := os.MkdirAll(uploadsWithSubDir, 0o755); err != nil {
		return "", fmt.Errorf("could not create uploads dir: %w", err)
	}

	return uploadsWithSubDir, nil
}

func getDirectoryPathFromEnv(key string) string {
	dir := os.Getenv(key)
	if dir == "" {
		panic(fmt.Sprintf("%s env var is not set", key))
	}
	return dir
}
