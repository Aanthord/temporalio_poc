package watson

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
)

func WatsonPostWalletCreateGP(nextUserId string) {
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
	body, err := io.ReadAll(resp.Body)
	// Log the request body
	bodyString := string(body)
	log.Print(bodyString)

	if err != nil {
		log.Printf("Reading body failed: %s", err)
		return
	}

	log.Printf("Posted to CreateGPNFTWatston")
}
