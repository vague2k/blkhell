package database

import (
	"fmt"

	"github.com/vague2k/blkhell/server/data"
)

func (b Band) FullCountry() string {
	return data.Countries[b.Country]
}

func (f File) ReadableSize() string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	fileSize := f.Size
	switch {
	case fileSize >= GB:
		return fmt.Sprintf("%.2fGB", float64(fileSize)/float64(GB))
	case fileSize >= MB:
		return fmt.Sprintf("%.2fMB", float64(fileSize)/float64(MB))
	case fileSize >= KB:
		return fmt.Sprintf("%.2fKB", float64(fileSize)/float64(KB))
	default:
		return fmt.Sprintf("%dB", fileSize)
	}
}
