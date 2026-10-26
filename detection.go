package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type packet struct {
	Filename       string
	StoragePath    string
	MIMEType       string
	FileSize       int64
	ext            string
	detectedFormat string
	Prefix         []byte
}

func fileSniff(filePath string) (packet, error) {
	pack := packet{}

	file, err := os.Open(filePath)
	if err != nil {
		return pack, fmt.Errorf("failed opening %s: %w", filePath, err)
	}

	defer file.Close()

	buffer := make([]byte, 512)
	n, err := io.ReadFull(file, buffer)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		return pack, fmt.Errorf("failed buff reading %s: %w", filePath, err)
	}

	info, err := file.Stat()
	if err != nil {
		return pack, fmt.Errorf("Failed reading stat of file %s: %w", filePath, err)
	}

	pack.Filename = file.Name()
	pack.ext = filepath.Ext(filePath)
	pack.StoragePath = filePath
	pack.FileSize = info.Size()
	pack.Prefix = buffer[:n]
	pack.MIMEType = http.DetectContentType(buffer)

	return pack, nil

}

func detectFormat(pack *packet) error {
	prefix := pack.Prefix

	switch {
	case bytes.HasPrefix(prefix, []byte{137, 80, 78, 71, 13, 10, 26, 10}):
		pack.detectedFormat = "png"

	case bytes.HasPrefix(prefix, []byte{255, 216, 255}):
		pack.detectedFormat = "jpeg"

	case bytes.HasPrefix(prefix, []byte{80, 75, 3, 4}):
		pack.detectedFormat = "zip"

	case bytes.HasPrefix(prefix, []byte("%PDF")):
		pack.detectedFormat = "pdf"

	case bytes.HasPrefix(prefix, []byte("MZ")):
		pack.detectedFormat = "pe"

	default:
		pack.detectedFormat = "unknown"
	}

	return nil
}

func packVerifier(pack *packet) string {
	ext := pack.ext

	switch ext {
	case ".png":
		if pack.detectedFormat == "png" {
			return "MATCH"
		}

	case ".jpg", ".jpeg":
		if pack.detectedFormat == "jpeg" {
			return "MATCH"
		}

	case ".zip":
		if pack.detectedFormat == "zip" {
			return "MATCH"
		}

	case ".docx":
		if pack.detectedFormat == "zip" {
			return "POSSIBLE MATCH (DOCX structure not checked)"
		}

	case ".pdf":
		if pack.detectedFormat == "pdf" {
			return "MATCH"
		}

	case ".exe", ".dll":
		if pack.detectedFormat == "pe" {
			return "MATCH"
		}

	case ".txt":
		if pack.detectedFormat == "unknown" {
			return "INCONCLUSIVE (no universal TXT signature)"
		}
		return "MISMATCH"

	default:
		return "UNSUPPORTED EXTENSION"
	}

	if pack.detectedFormat == "unknown" {
		return "INCONCLUSIVE"
	}

	return "MISMATCH"
}
