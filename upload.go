package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/intervention-engine/hdsfhir"
	"io/ioutil"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "upload"
	app.Usage = "Convert health-data-standards JSON to FHIR JSON and upload it to a FHIR Server"
	app.Flags = []cli.Flag{
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
	}
	app.Action = func(c *cli.Context) {
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
				patient.PostToFHIRServer(fhirUrl)
			}

		}
	}

	app.Run(os.Args)
}
