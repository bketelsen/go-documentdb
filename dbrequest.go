package docdb

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type documentDbRequest struct {
	resourceId   string
	resourceType string
	*http.Request
}

func (r *documentDbRequest) signRequest(key string) error {
	masterToken := "master"
	tokenVersion := "1.0"

	parts := []string{r.Method, r.resourceType, r.resourceId, r.Header.Get("x-ms-date"), r.Header.Get("Date"), ""}
	stringToSign := strings.ToLower(strings.Join(parts, "\n"))
	sig, err := signString(stringToSign, key)
	if err != nil {
		return err
	}

	header := url.QueryEscape("type=" + masterToken + "&ver=" + tokenVersion + "&sig=" + sig)
	r.Header.Add("Authorization", header)
	return nil
}

type RequestError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (d RequestError) Error() string {
	return fmt.Sprintf("%v, %v", d.Code, d.Message)
}
