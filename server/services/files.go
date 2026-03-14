package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"github.com/google/uuid"
	"github.com/vague2k/blkhell/config"
	"github.com/vague2k/blkhell/server/database"
	serverErrors "github.com/vague2k/blkhell/server/errors"
)

var (
	MimeJpeg      = "image/jpeg"
	MimePng       = "image/png"
	MimePhotoshop = "image/vnd.adobe.photoshop"
	MimeMp3       = "audio/mp3"
	MimeWav       = "audio/wav"
	MimeFlac      = "audio/flac"
)

type FileMetadata struct {
	Filename string
	Ext      string
	Path     string
	Mimetype string
	Size     int64

	// fields are private (for all intents and purposes)
	UserID    string
	OwnerID   string
	OwnerType string
}

type FilesService struct {
	config *config.Config
}

func NewFilesService(config *config.Config) *FilesService {
	return &FilesService{config: config}
}

func (s *FilesService) DownloadFile(w http.ResponseWriter, ctx context.Context, id string) error {
	file, err := s.config.Database.GetFileByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("Could not get file to prepare download")
		}
		return serverErrors.ErrDb
	}

	osFile, err := os.Open(s.config.UploadsDir + file.Path)
	if err != nil {
		return serverErrors.ErrInternal
	}
	defer osFile.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, file.FullFilename()))
	io.Copy(w, osFile)

	return nil
}

func (s *FilesService) DeleteFile(ctx context.Context, id string) (*database.File, error) {
	file, err := s.config.Database.DeleteFile(ctx, id)
	if err != nil {
		return nil, serverErrors.ErrDb
	}

	err = os.Remove(s.config.UploadsDir + file.Path)
	if err != nil {
		return nil, serverErrors.ErrInternal
	}

	return &file, nil
}

func (s *FilesService) Upload(w http.ResponseWriter, r *http.Request, userID, ownerID, ownerType string) (*database.File, error) {
	// max file size: 100MB
	r.Body = http.MaxBytesReader(w, r.Body, 100<<20+1)

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, serverErrors.ErrInternal
	}
	defer file.Close()

	if fileHeader.Size > 100<<20 {
		return nil, errors.New("The uploaded file is too big. Please choose a file that's less than 100MB in size")
	}

	metadata, err := s.WriteToDisk(file, fileHeader)
	if err != nil {
		return nil, err
	}

	metadata.UserID = userID
	metadata.OwnerID = ownerID
	metadata.OwnerType = ownerType

	asset, err := s.WriteToDb(r.Context(), metadata)
	if err != nil {
		return nil, err
	}

	return asset, nil
}

func (s *FilesService) WriteToDisk(file multipart.File, fileHeader *multipart.FileHeader) (*FileMetadata, error) {
	buf := make([]byte, 512)
	_, err := file.Read(buf)
	if err != nil {
		return nil, serverErrors.ErrInternal
	}

	// create the uploads dir if not exist
	filetype := mimetype.Detect(buf).String()
	fmt.Println(filetype)
	dir, err := s.mimetypeDir(filetype)
	if err != nil {
		return nil, err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, serverErrors.ErrInternal
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
		return nil, errors.New("The path doesn't exist")
	}
	defer dst.Close()
	size, err := io.Copy(dst, file)
	if err != nil {
		return nil, serverErrors.ErrInternal
	}

	_, relativeDir, ok := strings.Cut(filePath, "uploads")
	if !ok {
		return nil, errors.New("what the fuck?")
	}

	return &FileMetadata{
		Filename: fileName,
		Path:     relativeDir,
		Ext:      fileExt,
		Mimetype: filetype,
		Size:     size,
	}, nil
}

func (s *FilesService) WriteToDb(ctx context.Context, metadata *FileMetadata) (*database.File, error) {
	file, err := s.config.Database.CreateFile(ctx, database.CreateFileParams{
		ID:        uuid.NewString(),
		Path:      metadata.Path,
		Filename:  metadata.Filename,
		Ext:       metadata.Ext,
		Mimetype:  metadata.Mimetype,
		Size:      metadata.Size,
		UserID:    metadata.UserID,
		OwnerType: metadata.OwnerType,
		OwnerID:   metadata.OwnerID,
	})
	if err != nil {
		return nil, serverErrors.ErrDb
	}
	return &file, nil
}

func (s *FilesService) mimetypeDir(mimetype string) (string, error) {
	// TODO: move to main func somehow
	var subDir string

	switch mimetype {
	case MimeJpeg, MimePng:
		subDir = filepath.Join(s.config.UploadsDir, "images")
	// case MimePhotoshop:
	// 	subDir = filepath.Join(s.config.UploadsDir, "photoshop")
	case MimeMp3, MimeFlac, MimeWav:
		subDir = filepath.Join(s.config.UploadsDir, "audio")
	default:
		return "", errors.New("The file format is not supported")
	}

	if err := os.MkdirAll(subDir, 0o755); err != nil {
		return "", serverErrors.ErrInternal
	}

	return subDir, nil
}
