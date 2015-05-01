package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/intervention-engine/hdsfhir"
	"io/ioutil"
	"time"
	"os"
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
				patient.UnixBirthTime = time.Unix(patient.UnixBirthTime,0).AddDate(offset,0,0).Unix()

        for _,cond := range patient.Conditions {
          cond.StartTime = time.Unix(cond.StartTime, 0).AddDate(offset,0,0).Unix()
        }
        for _,enc := range patient.Encounters {
          enc.StartTime = time.Unix(enc.StartTime, 0).AddDate(offset,0,0).Unix()
        }
        for _,med := range patient.Medications {
          med.StartTime = time.Unix(med.StartTime, 0).AddDate(offset,0,0).Unix()
        }
        for _,vit := range patient.VitalSigns {
          vit.StartTime = time.Unix(vit.StartTime, 0).AddDate(offset,0,0).Unix()
        }
        for _,proc := range patient.Procedures {
          proc.StartTime = time.Unix(proc.StartTime, 0).AddDate(offset,0,0).Unix()
        }
				patient.PostToFHIRServer(fhirUrl)
			}

		}
	}

	app.Run(os.Args)
}
