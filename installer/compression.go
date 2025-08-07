package installer

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ulikunitz/xz"
)

type compression struct{}

// Extracts the given file according to the file ending
func (c *compression) Extract(filePath string, dstPath string, withoutRootFolder bool) error {
	switch {
	case strings.HasSuffix(filePath, ".zip"):
		return c.ExtractZip(filePath, dstPath, withoutRootFolder)
	case strings.HasSuffix(filePath, ".tar.gz"):
		return c.ExtractTarGz(filePath, dstPath, withoutRootFolder)
	case strings.HasSuffix(filePath, ".tar.xz"):
		return c.ExtractTarXz(filePath, dstPath, withoutRootFolder)
	case strings.HasSuffix(filePath, ".tar.bz2"):
		return c.ExtractTarBz2(filePath, dstPath, withoutRootFolder)
	default:
		return fmt.Errorf("unknown file type for file: %s", filePath)
	}
}

func (c *compression) ExtractZip(filePath string, dstPath string, withoutRootFolder bool) error {
	// Make sure the destination path exits and is clean
	if err := os.RemoveAll(dstPath); err != nil {
		return fmt.Errorf("failed cleaning destination folder: %w", err)
	}
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		return fmt.Errorf("failed creating destination folder: %w", err)
	}

	// Open the zip file
	archive, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer archive.Close()

	// Read the zip entries
	for _, zipEntry := range archive.File {
		entryPath := zipEntry.Name
		// Remove the root folder if wanted
		if withoutRootFolder {
			parts := strings.Split(entryPath, string(filepath.Separator))
			entryPath = strings.Join(parts[1:], string(filepath.Separator))
			if entryPath == "" {
				// This is the initial folder which is now empty
				continue
			}
		}

		currentItemPath := filepath.Join(dstPath, entryPath)
		if zipEntry.FileInfo().IsDir() {
			// Folder
			if err := os.MkdirAll(currentItemPath, os.ModePerm); err != nil {
				return fmt.Errorf("failed creating folder '%s': %w", currentItemPath, err)
			}
			continue
		} else {
			// File
			// Fallback to create the required folders as a zip is not required to include the folders as separate entries.
			if err := os.MkdirAll(filepath.Dir(currentItemPath), os.ModePerm); err != nil {
				return fmt.Errorf("failed creating folder for file '%s': %w", currentItemPath, err)
			}
			// Create the file
			outFile, err := os.OpenFile(currentItemPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fs.FileMode(zipEntry.Mode()))
			if err != nil {
				return fmt.Errorf("failed creating file '%s': %w", currentItemPath, err)
			}
			fileInArchive, err := zipEntry.Open()
			if err != nil {
				return fmt.Errorf("failed opening file from archive '%s': %w", currentItemPath, err)
			}
			if _, err := io.Copy(outFile, fileInArchive); err != nil {
				outFile.Close() // This error is omitted as the "Copy" error is more interesting
				return fmt.Errorf("failed copying content to '%s': %w", currentItemPath, err)
			}
			if err := outFile.Close(); err != nil {
				return fmt.Errorf("close failed on '%s': %w", currentItemPath, err)
			}
			fileInArchive.Close()
		}
	}
	return nil
}

// Extracts a tar.gz file into a folder.
func (c *compression) ExtractTarGz(filePath string, dstPath string, withoutRootFolder bool) error {
	// Open the file
	gzipStream, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed opening file '%s': %w", filePath, err)
	}
	defer gzipStream.Close()

	// Create an uncompressed reader
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}
	defer uncompressedStream.Close()

	// Extract the tar
	return c.extractTarStream(uncompressedStream, dstPath, withoutRootFolder)
}

// Extracts a tar.xz file into a folder.
func (c *compression) ExtractTarXz(filePath string, dstPath string, withoutRootFolder bool) error {
	// Open the file
	xzFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer xzFile.Close()

	// Create an uncompressed reader
	r := bufio.NewReader(xzFile)
	uncompressedStream, err := xz.NewReader(r)
	if err != nil {
		return err
	}

	// Extract the tar
	return c.extractTarStream(uncompressedStream, dstPath, withoutRootFolder)
}

// Extracts a tar.bz2 file into a folder.
func (c *compression) ExtractTarBz2(filePath string, dstPath string, withoutRootFolder bool) error {
	// Open the file
	bzip2Stream, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed opening file '%s': %w", filePath, err)
	}
	defer bzip2Stream.Close()

	// Create an uncompressed reader
	uncompressedStream := bzip2.NewReader(bzip2Stream)

	// Extract the tar
	return c.extractTarStream(uncompressedStream, dstPath, withoutRootFolder)
}

func (c *compression) extractTarStream(input io.Reader, dstPath string, withoutRootFolder bool) error {
	// Make sure the destination path exits
	if err := os.MkdirAll(dstPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed creating destination folder: %w", err)
	}

	// Create the tar reader and read the files
	tarReader := tar.NewReader(input)
	var header *tar.Header
	var err error
	for header, err = tarReader.Next(); err == nil; header, err = tarReader.Next() {
		entryPath := header.Name
		// Remove the root folder if wanted
		if withoutRootFolder {
			parts := strings.Split(entryPath, string(filepath.Separator))
			entryPath = strings.Join(parts[1:], string(filepath.Separator))
			if entryPath == "" {
				// This is the initial folder which is now empty
				continue
			}
		}
		currentItemPath := filepath.Join(dstPath, entryPath)
		switch header.Typeflag {
		case tar.TypeDir:
			// Folder
			if err := os.MkdirAll(currentItemPath, os.ModePerm); err != nil {
				return fmt.Errorf("failed creating folder '%s': %w", currentItemPath, err)
			}
		case tar.TypeReg:
			// File
			// Not all tar files include the folders, so make sure the folder path exists
			if err := os.MkdirAll(filepath.Dir(currentItemPath), os.ModePerm); err != nil {
				return fmt.Errorf("failed creating folder for file '%s': %w", currentItemPath, err)
			}
			// Create the file
			outFile, err := os.OpenFile(currentItemPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, fs.FileMode(header.Mode))
			if err != nil {
				return fmt.Errorf("failed creating file '%s': %w", currentItemPath, err)
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close() // This error is omitted as the "Copy" error is more interesting
				return fmt.Errorf("failed copying content to '%s': %w", currentItemPath, err)
			}
			if err := outFile.Close(); err != nil {
				return fmt.Errorf("close failed on '%s': %w", currentItemPath, err)
			}
		case tar.TypeSymlink:
			// Symlink
			if err := os.Symlink(header.Linkname, header.Name); err != nil {
				if errors.Is(err, fs.ErrNotExist) {
					fmt.Printf("warning: created symlink to non-existing file: '%s' to '%s'\n", currentItemPath, header.Linkname)
					continue
				}
				return fmt.Errorf("failed creating symlink '%s' to '%s': %w", currentItemPath, header.Linkname, err)
			}
		default:
			return fmt.Errorf("unknown type: %b in '%s'", header.Typeflag, header.Name)
		}
	}
	if err != io.EOF {
		return fmt.Errorf("getting Next() failed: %w", err)
	}
	return nil
}
