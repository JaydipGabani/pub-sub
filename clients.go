//
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type KustoIngestionClient struct{}

type RawMsgs struct {
	Records []PubsubMsg `json:"records,omitempty"`
}

func (client *KustoIngestionClient) SendAsync(msg PubsubMsg) error {
	targetUri := os.Getenv("TARGETURI")
	if targetUri == "" {
		log.Fatalf("Failed to set target uri env var")
	}

	functionKey := os.Getenv("FUNCTIONKEY")
	if functionKey == "" {
		log.Fatalf("Failed to set function key")
	}

	requestUri := targetUri + "?code=" + functionKey

	rawMsgs := RawMsgs{
		Records: []PubsubMsg{msg},
	}

	buf, err := json.Marshal(rawMsgs)
	if err != nil {
		log.Fatalf(err.Error())
	}

	var resp *http.Response
	for i := 0; i < 3; i++ {
		resp, err = http.Post(requestUri, "application/json", bytes.NewBuffer(buf))
		if err == nil {
			break
		}
		fmt.Println(err)
	}

	if err != nil {
		return err
	}

	fmt.Println("Response Status:", resp.Status)

	return nil
}
