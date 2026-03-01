package handlers

import (
	"fmt"
	"net/http"

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

	user, err := h.AuthService.Authenticate(r.Context(), username, password)
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	sessionToken, expires, err := h.AuthService.CreateSession(r.Context(), user.ID)
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

func (h *Handler) CreateBand(w http.ResponseWriter, r *http.Request) {
	bandName := r.FormValue("band-name")
	bandCountry := r.FormValue("band-country")
	releaseName := r.FormValue("release-name")
	releaseType := r.FormValue("release-type")
	releaseNum := r.FormValue("release-number")
	projectName := r.FormValue("project-name")
	projectType := r.FormValue("project-type")

	switch true {
	case bandName == "":
		toastError(w, r, "'Band name' is required")
		return
	case releaseName == "":
		toastError(w, r, "'Release Name' is required")
		return
	case releaseType == "":
		toastError(w, r, "'Release Type' is required")
		return
	case releaseNum == "":
		toastError(w, r, "'Release No.' is required")
		return
	case bandCountry == "":
		toastError(w, r, "'Band country' is required")
		return
	}

	if (projectName == "" && projectType != "") || (projectName != "" && projectType == "") {
		toastError(w, r, "'Project Name' and 'Project Type' must both be filled or both left empty")
		return
	}

	band, err := h.DB.CreateBand(r.Context(), database.CreateBandParams{
		ID:      uuid.NewString(),
		Name:    bandName,
		Country: bandCountry,
	})
	if err != nil {
		toastError(w, r, "Database error, try again. Contact admin if issue occurs")
	}

	release, err := h.DB.CreateRelease(r.Context(), database.CreateReleaseParams{
		ID:     uuid.NewString(),
		BandID: band.ID,
		Name:   releaseName,
		Type:   releaseType,
		Number: releaseNum,
	})
	if err != nil {
		toastError(w, r, "Database error, try again. Contact admin if issue occurs")
	}

	if projectName != "" && projectType != "" {
		_, err := h.DB.CreateProject(r.Context(), database.CreateProjectParams{
			ID:        uuid.NewString(),
			BandID:    band.ID,
			ReleaseID: release.ID,
			Name:      projectName,
			Type:      projectType,
		})
		if err != nil {
			toastError(w, r, "Database error, try again. Contact admin if issue occurs")
		}
	}

	toastSuccess(w, r, "Endpoint hit check console")
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	user, ok := h.AuthService.UserFromContext(r.Context())
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

	metadata, err := h.FilesService.WriteToDisk(file, fileHeader)
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	metadata.UserID = user.ID
	metadata.OwnerType = "label"
	metadata.OwnerID = "label"

	image, err := h.FilesService.WriteToDb(r.Context(), metadata)
	if err != nil {
		toastError(w, r, err.Error())
		return
	}

	toastSuccess(w, r, fmt.Sprintf("'%s.%s' was uploaded successfully!", image.Filename, image.Ext))
}
