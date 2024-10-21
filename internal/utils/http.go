package utils

import (
	"io"
	"mime/multipart"
)

func MultipartToBytes(in *multipart.FileHeader) ([]byte, error) {
	fInfo, err := in.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		cErr := fInfo.Close()
		if cErr != nil && err == nil {
			err = cErr
		}
	}()

	content, err := io.ReadAll(fInfo)
	if err != nil {
		return nil, err
	}

	return content, err
}
