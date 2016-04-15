package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/codegangsta/cli"
	"github.com/intervention-engine/fhir/models"
	"github.com/intervention-engine/hdsfhir"
)

func main() {
	app := cli.NewApp()
	app.Name = "upload"
	app.Usage = "Convert health-data-standards JSON to FHIR JSON and upload it to a FHIR Server"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "offset, o",
			Usage: "How many years to offset dates by",
		},
		cli.StringFlag{
			Name:  "fhir, f",
			Usage: "URL for the FHIR server",
		},
		cli.StringFlag{
			Name:  "json, j",
			Usage: "Path to the directory of JSON files",
		},
		cli.StringFlag{
			Name:  "single, s",
			Usage: "Path to the a single JSON file",
		},
		cli.BoolFlag{
			Name:  "conditional, c",
			Usage: "Indicates that data should be uploaded using conditional updates",
		},
	}
	app.Action = func(c *cli.Context) {
		offset := c.Int("offset")
		fhirUrl := c.String("fhir")
		path := c.String("json")
		singlePath := c.String("single")
		if fhirUrl == "" || (path == "" && singlePath == "") {
			fmt.Println("You must provide a FHIR URL and path to JSON files")
		} else {
			var fileNames []string
			if path != "" {
				files, err := ioutil.ReadDir(path)
				if err != nil {
					panic("Couldn't read the directory" + err.Error())
				}
				for _, file := range files {
					fileNames = append(fileNames, path+"/"+file.Name())
				}
			} else {
				fileNames = []string{singlePath}
			}
			for _, file := range fileNames {
				patient := &hdsfhir.Patient{}
				jsonBlob, err := ioutil.ReadFile(file)
				if err != nil {
					panic("Couldn't read the JSON file" + err.Error())
				}
				json.Unmarshal(jsonBlob, patient)
				patient.BirthTime = hdsfhir.UnixTime(patient.BirthTime.Time().AddDate(offset, 0, 0).Unix())

				for _, cond := range patient.Conditions {
					cond.StartTime = hdsfhir.UnixTime(cond.StartTime.Time().AddDate(offset, 0, 0).Unix())
				}
				for _, enc := range patient.Encounters {
					enc.StartTime = hdsfhir.UnixTime(enc.StartTime.Time().AddDate(offset, 0, 0).Unix())
				}
				for _, med := range patient.Medications {
					med.StartTime = hdsfhir.UnixTime(med.StartTime.Time().AddDate(offset, 0, 0).Unix())
				}
				for _, vit := range patient.VitalSigns {
					vit.StartTime = hdsfhir.UnixTime(vit.StartTime.Time().AddDate(offset, 0, 0).Unix())
				}
				for _, proc := range patient.Procedures {
					proc.StartTime = hdsfhir.UnixTime(proc.StartTime.Time().AddDate(offset, 0, 0).Unix())
				}

				// Convert the patient bundle to JSON data
				data, err := json.Marshal(patient.FHIRTransactionBundle(c.Bool("conditional")))
				if err != nil {
					panic("Couldn't convert FHIR patient bundle to JSON data: " + err.Error())
				}

				// Post the data
				res, err := http.Post(fhirUrl+"/", "application/json", bytes.NewBuffer(data))
				if err != nil {
					panic("Couldn't upload FHIR patient bundle: " + err.Error())
				}

				// Decode the response
				decoder := json.NewDecoder(res.Body)
				responseBundle := &models.Bundle{}
				err = decoder.Decode(responseBundle)
				if err != nil {
					panic("Uploaded FHIR patient bundle, but couldn't process response: " + err.Error())
				}
				fmt.Printf("Successfully uploaded patient bundle with %d resources\n", *responseBundle.Total)
			}

		}
	}

	app.Run(os.Args)
}
