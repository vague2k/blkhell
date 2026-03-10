package services

import (
	"context"
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
	"github.com/vague2k/blkhell/server/database"
)

var (
	MimeJpeg      = "image/jpeg"
	MimePng       = "image/png"
	MimePhotoshop = "image/vnd.adobe.photoshop"
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
	db *database.Queries
}

func NewFilesService(db *database.Queries) *FilesService {
	return &FilesService{db: db}
}

func (s *FilesService) Upload(r *http.Request, userID, ownerID, ownerType string) (*database.File, error) {
	// max file size: 100MB
	// r.Body = http.MaxBytesReader(w, r.Body, 100<<20) // aggresively strict max size
	if err := r.ParseMultipartForm(100 << 20); err != nil {
		return nil, errors.New("The uploaded file is too big. Please choose an file that's less than 100MB in size")
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, ErrInternal
	}
	defer file.Close()

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
		return nil, ErrInternal
	}

	// create the uploads dir if not exist
	filetype := mimetype.Detect(buf).String()
	fmt.Println(filetype)
	dir, err := s.createUploadDirectories(filetype)
	if err != nil {
		return nil, err
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return nil, ErrInternal
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
		return nil, ErrInternal
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
	file, err := s.db.CreateFile(ctx, database.CreateFileParams{
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
		return nil, ErrDb
	}
	return &file, nil
}

func (s *FilesService) createUploadDirectories(mimetype string) (string, error) {
	// TODO: move to main func somehow
	uploadsDir := os.Getenv("UPLOADS_DIR")
	var uploadsWithSubDir string

	switch mimetype {
	case MimeJpeg, MimePng:
		uploadsWithSubDir = filepath.Join(uploadsDir, "images")
	// case MimePhotoshop:
	// 	uploadsWithSubDir = filepath.Join(uploadsDir, "photoshop")
	default:
		return "", errors.New("The file format is not supported")
	}

	if err := os.MkdirAll(uploadsWithSubDir, 0o755); err != nil {
		return "", ErrInternal
	}

	return uploadsWithSubDir, nil
}
