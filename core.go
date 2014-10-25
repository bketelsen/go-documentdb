package docdb

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"strings"
)

func signString(value, key string) (string, error) {
	var ret string
	enc := base64.StdEncoding

	salt, err := enc.DecodeString(key)
	if err != nil {
		return ret, err
	}

	hmac := hmac.New(sha256.New, salt)
	hmac.Write([]byte(value))
	b := hmac.Sum(nil)

	ret = enc.EncodeToString(b)
	return ret, nil
}

func readJson(reader io.Reader, data interface{}) error {
	return json.NewDecoder(reader).Decode(&data)
}

type ResourceInfo struct {
	Id   string
	Type string
}

func parseResource(resourcePath string) *ResourceInfo {
	if resourcePath[len(resourcePath)-1] != '/' {
		resourcePath = resourcePath + "/"
	}
	if resourcePath[0] != '/' {
		resourcePath = "/" + resourcePath
	}

	parts := strings.Split(resourcePath, "/")
	partsLength := len(parts)

	var resourceId string
	var resourceType string
	if partsLength%2 == 0 {
		resourceId = parts[partsLength-2]
		resourceType = parts[partsLength-3]
	} else {
		resourceId = parts[partsLength-3]
		resourceType = parts[partsLength-2]
	}

	if resourceId[len(resourceId)-1] == '\\' {
		resourceId = resourceId[:len(resourceId)-1]
	}
	if resourceType[len(resourceType)-1] == '\\' {
		resourceType = resourceType[:len(resourceType)-1]
	}

	return &ResourceInfo{resourceId, resourceType}
}
