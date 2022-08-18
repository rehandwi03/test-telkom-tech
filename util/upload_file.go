package util

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"
)

func GetFileName(filename, data string) (string, []byte, error) {
	base64Decode, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Println(err)
		return "", nil, err
	}
	mimeType := http.DetectContentType(base64Decode)

	split := strings.Split(mimeType, "/")

	newFileName := filename + "." + split[1]

	return newFileName, base64Decode, nil
}
