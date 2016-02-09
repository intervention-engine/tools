package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/intervention-engine/fhir/models"
)

func main() {
	fhirURL := flag.String("fhirURL", "", "URL for the FHIR server")
	single := flag.String("bundle", "", "Path to the a single JSON file containing a bundle to upload")
	flag.Parse()

	if *fhirURL == "" || *single == "" {
		panic("Must provide parameter values for fhirURL and bundle")
	}

	// Read in the data in FHIR bundle format
	data, err := os.Open(*single)
	if err != nil {
		panic("Couldn't read the JSON file: " + err.Error())
	}

	// Post the data
	res, err := http.Post(*fhirURL+"/", "application/json", data)
	if err != nil {
		panic("Couldn't upload FHIR bundle: " + err.Error())
	}

	decoder := json.NewDecoder(res.Body)
	responseBundle := &models.Bundle{}
	err = decoder.Decode(responseBundle)
	if err != nil {
		panic("Uploaded bundle, but couldn't process response: " + err.Error())
	}
	fmt.Printf("Successfully uploaded %d resources\n", *responseBundle.Total)
}
