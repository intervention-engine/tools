package main

import (
	"flag"

	"github.com/intervention-engine/fhir/upload"
	"github.com/intervention-engine/ptgen"
)

func main() {
	registerURL := flag.String("fhirURL", "", "URL for the FHIR server")
	num := flag.Int("n", 100, "Number of patients to generate")
	flag.Parse()
	count := *num
	for i := 0; i < count; i++ {
		resources := ptgen.GeneratePatient()
		upload.UploadResources(resources, *registerURL)
	}
}
