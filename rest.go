package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	requestsFile        *string = flag.String("requests_file", "requests.http", "File that lists all possible requests that can be run.")
	requestName         *string = flag.String("name", "", "Name of the request to execute.")
	createRequestSchema *bool   = flag.Bool("gen_schema", false, "Creates new requests schema file.")
)

type RequestPayload struct {
	Name    string         `json:"name"`
	Url     string         `json:"url"`
	Method  string         `json:"method"`
	Headers httpHeaders    `json:"headers"`
	Body    map[string]any `json:"body"`
}

type httpHeaders = map[string]string

func main() {
	flag.Parse()

	if flag.Parsed() != true {
		log.Fatal("failed to parse cli args")
	}

	if *createRequestSchema {
		if err := generateSchemaFile(); err != nil {
			log.Fatalf("error to create schema: %v", err)
		}
		return
	}

	payloads, err := parseRequestsFile(*requestsFile)
	if err != nil {
		log.Fatal(err)
	}

	requestPayload, err := selectRequest(payloads, *requestName)
	if err != nil {
		log.Fatal(err)
	}

	if err = makeRequest(requestPayload, os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func selectRequest(payloads []RequestPayload, name string) (RequestPayload, error) {
	for _, req := range payloads {
		if req.Name == name {
			return req, nil
		}
	}

	return RequestPayload{}, fmt.Errorf("request with name %s was not found", name)
}

func makeRequest(req RequestPayload, w io.Writer) error {
	var body io.Reader
	if req.Body != nil {
		buffer, _ := json.Marshal(req.Body)
		body = bytes.NewBuffer(buffer)
	}

	request, _ := http.NewRequest(strings.ToUpper(req.Method), req.Url, body)
	for h, v := range req.Headers {
		request.Header.Add(h, v)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}

	defer response.Body.Close()
	io.WriteString(w, response.Status+"\r\n")
	for k, v := range response.Header {
		fmt.Fprintf(w, "%s: %s\r\n", k, strings.Join(v, ""))
	}
	w.Write([]byte("\r\n"))

	_, err = io.Copy(w, response.Body)
	return err
}

func parseRequestsFile(filePath string) ([]RequestPayload, error) {
	file, err := os.OpenFile(filePath, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var requests []RequestPayload
	return requests, json.NewDecoder(file).Decode(&requests)
}

func generateSchemaFile() error {
	const filename string = "requests_schema.json"
	buf, _ := json.MarshalIndent([]RequestPayload{{
		Name:    "",
		Method:  "",
		Headers: httpHeaders{},
		Body:    nil,
	}}, "", strings.Repeat(" ", 2))

	return os.WriteFile(filename, buf, 0755)
}
