package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/intervention-engine/fhir/models"
	"math/rand"
	"net/http"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "get-risky"
	app.Usage = "Upload risk assessments for patients retrieved from a FHIR URL"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "fhir, f",
			Usage: "URL for the FHIR server",
		},
	}
	app.Action = func(c *cli.Context) {
		fhirUrl := c.String("fhir")
		if fhirUrl == "" {
			fmt.Println("You must provide a FHIR URL")
		} else {
			patientIds := GetPatientIds(fhirUrl + "/Patient")
			r := rand.New(rand.NewSource(99))
			for _, id := range patientIds {
				UploadRiskAssessment(fhirUrl+"/Observation", fhirUrl+"/Patient", id, r)
			}
		}
	}

	app.Run(os.Args)
}

func GetPatientIds(patientUrl string) []string {
	resp, err := http.Get(patientUrl)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(resp.Body)
	patientBundle := &models.PatientBundle{}
	err = decoder.Decode(patientBundle)
	if err != nil {
		panic(err)
	}
	var patientIds []string
	for _, p := range patientBundle.Entry {
		patientIds = append(patientIds, p.Id)
	}
	return patientIds
}

func UploadRiskAssessment(observationUrl, patientUrl, patientId string, r *rand.Rand) {
	observation := models.Observation{Reliability: "ok", Status: "final"}
	randomRisk := r.Intn(4)
	if randomRisk == 0 {
		randomRisk = 1
	}
	observation.ValueQuantity = models.Quantity{Value: float64(randomRisk)}
	cc := models.CodeableConcept{}
	cc.Coding = []models.Coding{models.Coding{System: "http://loinc.org", Code: "75492-9"}}
	observation.Name = cc
	observation.AppliesDateTime = models.FHIRDateTime{Precision: models.Timestamp, Time: time.Now()}
	observation.Subject = models.Reference{Reference: patientUrl + "/" + patientId}
	json, _ := json.Marshal(observation)
	body := bytes.NewReader(json)
	_, err := http.Post(observationUrl, "application/json+fhir", body)
	if err != nil {
		panic("HTTP request failed")
	}
}
