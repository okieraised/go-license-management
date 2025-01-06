package utils

import (
	"bytes"
	"io"
	"mime/multipart"
	"net"
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

func RetrievePhysicalMacAddr() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var as []string
	for _, ifa := range interfaces {
		if ifa.Flags&net.FlagUp != 0 && bytes.Compare(ifa.HardwareAddr, nil) != 0 {
			if ifa.HardwareAddr[0]&2 == 2 {
				continue
			}
			a := ifa.HardwareAddr.String()
			if a != "" {
				as = append(as, a)
			}
		}

	}
	return as, nil
}

func RetrieveMacAddr() ([]string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	var as []string
	for _, ifa := range interfaces {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}
