package main

import (
	"encoding/json"
	"github.com/PaesslerAG/jsonpath"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

const (
	unsplashSearchPhotosEndpoint = "https://api.unsplash.com/photos/random"
)

var clientId = os.Getenv("WB_UNSPLASH_API_KEY")

func SearchImage(query string) (*string, error) {
	u, err := url.Parse(unsplashSearchPhotosEndpoint)
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("count", "1")
	params.Add("query", query)
	params.Add("client_id", clientId)

	u.RawQuery = params.Encode()
	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	rbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsedBody interface{}
	err = json.Unmarshal(rbody, &parsedBody)
	if err != nil {
		return nil, err
	}

	photoUrl, err := jsonpath.Get("$[0].urls.regular", parsedBody)
	if err != nil {
		return nil, err
	}
	photoUrlString := photoUrl.(string)
	return &photoUrlString, nil
}

func SearchFrog() (*string, error) {
	return SearchImage("frogs,toad")
}

func DownloadFrogPhoto(writer io.Writer) error {
	frog, err := SearchFrog()
	if err != nil {
		return err
	}
	responce, err := http.Get(*frog)
	if err != nil {
		return err
	}
	defer responce.Body.Close()

	n, err := io.Copy(writer, responce.Body)
	log.Printf("DownloadFrogPhoto got %d bytes", n)
	return err
}
