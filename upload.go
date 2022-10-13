package chirano

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type UploadedFile struct {
	NewFileName, OriginalFileName string
	FileSize                      int64
}

func (h *Helpers) UploadFiles(r *http.Request, uploadDir string, rename ...bool) ([]*UploadedFile, error) {

	renameFile := true
	if len(rename) > 0 {
		renameFile = rename[0]
	}

	var uploadedFiles []*UploadedFile

	if h.MaxFileSize == 0 {
		h.MaxFileSize = 1024 * 1024 * 1024
	}

	err := r.ParseMultipartForm(int64(h.MaxFileSize))
	r.

	if err != nil {
		return nil, errors.New("The uploaded file is too big!")
	}

	for _, fHeaders := range r.MultipartForm.File {
		for _, header := range fHeaders {

			uploadedFiles, err = func(uploadedFiles []*UploadedFile) ([]*UploadedFile, error) {

				var uploadedFile UploadedFile
				infile, err := header.Open()

				if err != nil {
					return nil, err
				}

				defer infile.Close()

				buff := make([]byte, 512)

				_, err = infile.Read(buff)

				if err != nil {
					return nil, err
				}

				//check to see if file type is permitted
				allowed := false
				fileType := http.DetectContentType(buff)

				if len(h.AllowedFileTypes) > 0 {

					for _, x := range h.AllowedFileTypes {
						if strings.EqualFold(fileType, x) {
							allowed = true
						}
					}
				} else {
					allowed = true
				}

				if !allowed {
					return nil, errors.New("the uploaded file type is not allowed")
				}

				_, err = infile.Seek(0, 0)

				if err != nil {
					return nil, err
				}

				if renameFile {
					uploadedFile.NewFileName = fmt.Sprintf("%s%s", h.RandomString(25), filepath.Ext(header.Filename))
				} else {
					uploadedFile.NewFileName = header.Filename
				}

				var outfile *os.File

				defer outfile.Close()

				if outfile, err = os.Create(filepath.Join(uploadDir, uploadedFile.NewFileName)); err != nil {
					return nil, err
				} else {
					fileSize, err := io.Copy(outfile, infile)

					if err != nil {
						return nil, err
					}
					uploadedFile.FileSize = fileSize
				}

				uploadedFiles = append(uploadedFiles, &uploadedFile)

				return uploadedFiles, nil

			}(uploadedFiles)

			if err != nil {
				return uploadedFiles, err
			}
		}
	}
	return uploadedFiles, nil
}
