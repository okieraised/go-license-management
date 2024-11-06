package utils

import (
	"bytes"
	"compress/gzip"
	"io"
)

func Compress(s []byte) ([]byte, error) {
	var err error

	buf := bytes.Buffer{}
	zipped := gzip.NewWriter(&buf)

	_, err = zipped.Write(s)
	if err != nil {
		return nil, err
	}
	err = zipped.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decompress(s []byte) ([]byte, error) {
	rdr, err := gzip.NewReader(bytes.NewReader(s))
	if err != nil {
		return nil, err
	}
	defer func() {
		cErr := rdr.Close()
		if cErr != nil && err == nil {
			err = cErr
		}
	}()

	data, err := io.ReadAll(rdr)
	if err != nil {
		return nil, err
	}
	return data, nil
}
