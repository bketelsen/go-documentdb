package docdb

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type azureClient struct {
	url    string
	config Config
}

func (c *azureClient) newHttpRequest(method, url, resourceType, resourceId string, buffer *bytes.Buffer) (*documentDbRequest, error) {
	absolute := c.url + url
	req, err := http.NewRequest(method, absolute, buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-ms-date", time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT"))

	r := &documentDbRequest{resourceId, resourceType, req}
	err = r.signRequest(c.config.MasterKey)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (c *azureClient) makeDeleteResourceRequest(url string) (*http.Response, error) {
	r := parseResource(url)
	req, err := c.newHttpRequest("DELETE", url, r.Type, r.Id, &bytes.Buffer{})
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusNoContent {
		err := RequestError{}
		readJson(resp.Body, &err)
		return nil, err
	}

	return resp, nil
}

func (c *azureClient) makeResourceRequest(url string) (*http.Response, error) {
	r := parseResource(url)
	req, err := c.newHttpRequest("GET", url, r.Type, r.Id, &bytes.Buffer{})
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err := RequestError{}
		readJson(resp.Body, &err)
		return nil, err
	}

	return resp, nil
}

func (c *azureClient) makeJsonRequest(method, url, resourceType, resourceId string, body interface{}) (*documentDbRequest, error) {
	var buffer *bytes.Buffer
	if body != nil {
		b, err := json.Marshal(body)

		if err != nil {
			return nil, err
		}
		buffer = bytes.NewBuffer(b)
	}

	r, err := c.newHttpRequest(method, url, resourceType, resourceId, buffer)
	if err != nil {
		return nil, err
	}

	r.Header.Add("Content-Type", "application/json")
	return r, nil
}

func (c *azureClient) Do(req *documentDbRequest) (*http.Response, error) {
	client, err := c.newHttpClient()
	if err != nil {
		return nil, err
	}
	return client.Do(req.Request)
}

func (c *azureClient) newHttpClient() (*http.Client, error) {
	// todo build client based off config
	return &http.Client{}, nil
}
