package watson

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

type post struct {
	Userid string `json:"user_id"`
}

func WatsonPostWalletCreateJobPosting(nextUserId string) {
	watsonURL := (os.Getenv("watsonURL"))
	params := url.Values{}
	params.Add("user_id", nextUserId)
	resp, err := http.PostForm(watsonURL,
		params)
	if err != nil {
		log.Printf("Request Failed: %s", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	// Log the request body
	bodyString := string(body)
	log.Print(bodyString)
	// Unmarshal result
	post := Post{}
	err = json.Unmarshal(body, &post)
	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}

	log.Printf("Post added")
}
