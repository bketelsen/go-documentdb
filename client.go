package docdb

import (
	"bytes"
	"errors"
	"net/http"
)

var (
	ErrorInvalidSqlSyntax = errors.New("The sql syntax provided is invalid")
	ErrorUnauthorized     = errors.New("The request could not be authenticated")
)

type Client struct {
	azureClient
}

type Database struct {
	Resource
}

type Collection struct {
	Resource
	IndexingPolicy IndexingPolicy `json:"indexingPolicy,omitempty"`
}

type Document struct {
	Id string `json:"id,omitempty"`
	Resource
}

type IndexingPolicy struct {
	IsAutomatic  bool   `json:"automatic"`
	IndexingMode string `json:"indexingMode"`
}

type Config struct {
	MasterKey string
}

func NewClient(url string, config Config) *Client {
	return &Client{
		azureClient{
			url:    url,
			config: config,
		},
	}
}

func (c *Client) CreateDatabase(databaseName string) (*Database, error) {
	data := struct {
		Id string `json:"id"`
	}{databaseName}

	req, err := c.makeJsonRequest("POST", "/dbs", "dbs", "", &data)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		err := &RequestError{}
		readJson(resp.Body, err)
		return nil, err
	}

	db := &Database{}
	readJson(resp.Body, db)
	return db, nil
}

func (c *Client) CreateCollection(databaseLink, collectionName string, indexingPolicy *IndexingPolicy) (*Collection, error) {
	data := struct {
		Id             string          `json:"id"`
		IndexingPolicy *IndexingPolicy `json:"indexingPolicy,omitempty"`
	}{collectionName, nil}

	path := "/" + databaseLink + "colls/"
	r := parseResource(path)

	req, err := c.makeJsonRequest("POST", path, "colls", r.Id, &data)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		err := &RequestError{}
		readJson(resp.Body, err)
		return nil, err
	}

	collection := &Collection{}
	readJson(resp.Body, collection)
	return collection, nil
}

func (c *Client) CreateDocument(collectionLink string, document interface{}) error {
	path := "/" + collectionLink + "docs/"
	r := parseResource(path)

	req, err := c.makeJsonRequest("POST", path, "docs", r.Id, document)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		err := &RequestError{}
		readJson(resp.Body, err)
		return err
	}

	return readJson(resp.Body, document)
}

func (c *Client) ReadDocuments(collectionLink string, body interface{}) error {
	path := "/" + collectionLink + "docs"

	data := struct {
		Documents interface{} `json:"Documents,omitempty"`
		Count     int         `json:"_count,omitempty"`
	}{Documents: body}

	err := c.readResource(path, &data)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ReadDocument(documentLink string, body interface{}) error {
	path := "/" + documentLink
	err := c.readResource(path, body)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) ReadCollections(databaseLink string) ([]Collection, error) {
	path := "/" + databaseLink + "colls/"

	data := struct {
		Collections []Collection `json:"DocumentCollections,omitempty"`
	}{}

	err := c.readResource(path, &data)
	if err != nil {
		return nil, err
	}
	return data.Collections, nil
}

func (c *Client) ReadCollection(collectionLink string) (*Collection, error) {
	path := "/" + collectionLink

	collection := &Collection{}
	err := c.readResource(path, collection)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

func (c *Client) ReadDatabase(databaseLink string) (*Database, error) {
	path := "/" + databaseLink

	db := &Database{}
	err := c.readResource(path, db)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (c *Client) DeleteDocument(documentLink string) error {
	path := "/" + documentLink
	return c.deleteResource(path)
}

func (c *Client) DeleteCollection(collectionLink string) error {
	path := "/" + collectionLink
	return c.deleteResource(path)
}

func (c *Client) DeleteDatabase(databaseLink string) error {
	path := "/" + databaseLink
	return c.deleteResource(path)
}

func (c *Client) QueryDocuments(collectionLink, query string, body interface{}) error {
	path := "/" + collectionLink + "docs"

	data := struct {
		Documents interface{} `json:"Documents,omitempty"`
		Count     int         `json:"_count,omitempty"`
	}{Documents: body}

	return c.queryResource(path, query, &data)
}

func (c *Client) queryResource(uri, query string, body interface{}) error {
	r := parseResource(uri)
	buffer := bytes.NewBufferString(query)
	req, err := c.newHttpRequest("POST", uri, r.Type, r.Id, buffer)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/sql")
	req.Header.Add("x-ms-documentdb-isquery", "True")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return ErrorInvalidSqlSyntax
	} else if resp.StatusCode != http.StatusOK {
		err := RequestError{}
		return readJson(resp.Body, &err)
	}

	// count := resp.Header.Get("x-ms-item-count")
	// continuation := resp.Header.Get("x-ms-continuation")

	return readJson(resp.Body, body)
}

func (c *Client) readResource(uri string, body interface{}) error {
	resp, err := c.makeResourceRequest(uri)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return readJson(resp.Body, body)
}

func (c *Client) deleteResource(uri string) error {
	_, err := c.makeDeleteResourceRequest(uri)
	if err != nil {
		return err
	}

	return nil
}
