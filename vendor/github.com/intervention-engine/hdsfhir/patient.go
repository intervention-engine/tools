package hdsfhir

import (
	"encoding/json"
	"fmt"
	"reflect"

	fhir "github.com/intervention-engine/fhir/models"
)

type Patient struct {
	TemporallyIdentified
	MedicalRecordNumber string          `json:"medical_record_number"`
	FirstName           string          `json:"first"`
	LastName            string          `json:"last"`
	BirthTime           UnixTime        `json:"birthdate"`
	Gender              string          `json:"gender"`
	Encounters          []*Encounter    `json:"encounters"`
	Conditions          []*Condition    `json:"conditions"`
	VitalSigns          []*VitalSign    `json:"vital_signs"`
	Procedures          []*Procedure    `json:"procedures"`
	Medications         []*Medication   `json:"medications"`
	Immunizations       []*Immunization `json:"immunizations"`
	Allergies           []*Allergy      `json:"allergies"`
}

// TODO: :care_goals, :medical_equipment, :results, :social_history, :support, :advance_directives, :insurance_providers, :functional_statuses

func (p *Patient) MatchingEncounterReference(entry Entry) *fhir.Reference {
	for _, encounter := range p.Encounters {
		// TODO: Overlaps may not be the right thing here... maybe closest?
		if encounter.StartTime <= entry.EndTime && encounter.EndTime >= entry.StartTime {
			return encounter.FHIRReference()
		}
	}
	return nil
}

func (p *Patient) FHIRModel() *fhir.Patient {
	fhirPatient := &fhir.Patient{}
	fhirPatient.Id = p.GetTempID()
	if p.MedicalRecordNumber != "" {
		fhirPatient.Identifier = []fhir.Identifier{
			{
				Type: &fhir.CodeableConcept{
					Coding: []fhir.Coding{
						{
							System:  "http://hl7.org/fhir/v2/0203",
							Code:    "MR",
							Display: "Medical Record Number",
						},
					},
					Text: "Medical Record Number",
				},
				Value: p.MedicalRecordNumber,
			},
		}
	}
	fhirPatient.Name = []fhir.HumanName{fhir.HumanName{Given: []string{p.FirstName}, Family: []string{p.LastName}}}
	switch p.Gender {
	case "M":
		fhirPatient.Gender = "male"
	case "F":
		fhirPatient.Gender = "female"
	default:
		fhirPatient.Gender = "unknown"
	}
	fhirPatient.BirthDate = p.BirthTime.FHIRDate()
	return fhirPatient
}

func (p *Patient) FHIRModels() []interface{} {
	var models []interface{}
	models = append(models, p.FHIRModel())
	for _, encounter := range p.Encounters {
		models = append(models, encounter.FHIRModels()...)
	}
	for _, condition := range p.Conditions {
		models = append(models, condition.FHIRModels()...)
	}
	for _, observation := range p.VitalSigns {
		models = append(models, observation.FHIRModels()...)
	}
	for _, procedure := range p.Procedures {
		models = append(models, procedure.FHIRModels()...)
	}
	for _, medication := range p.Medications {
		models = append(models, medication.FHIRModels()...)
	}
	for _, immunization := range p.Immunizations {
		models = append(models, immunization.FHIRModels()...)
	}
	for _, allergy := range p.Allergies {
		models = append(models, allergy.FHIRModels()...)
	}

	return models
}

// FHIRTransactionBundle returns a FHIR bundle representing a transaction to post all patient data to a server
func (p *Patient) FHIRTransactionBundle(conditionalUpdate bool) *fhir.Bundle {
	bundle := new(fhir.Bundle)
	bundle.Type = "transaction"
	fhirModels := p.FHIRModels()
	bundle.Entry = make([]fhir.BundleEntryComponent, len(fhirModels))
	for i := range fhirModels {
		bundle.Entry[i].FullUrl = "urn:uuid:" + reflect.ValueOf(fhirModels[i]).Elem().FieldByName("Id").String()
		bundle.Entry[i].Resource = fhirModels[i]
		name := reflect.TypeOf(fhirModels[i]).Elem().Name()
		// Only use conditional put if requested, is a patient, and has an alternate ID (i.e., MRN)
		if conditionalUpdate && name == "Patient" && len(fhirModels[i].(*fhir.Patient).Identifier) > 0 {
			pt := fhirModels[i].(*fhir.Patient)
			bundle.Entry[i].Request = &fhir.BundleEntryRequestComponent{
				Method: "PUT",
				Url:    fmt.Sprintf("%s?identifier=%s", name, pt.Identifier[0].Value),
			}
		} else {
			bundle.Entry[i].Request = &fhir.BundleEntryRequestComponent{
				Method: "POST",
				Url:    name,
			}
		}
	}
	return bundle
}

// The "patient" sub-type is needed to avoid infinite recursion in UnmarshalJSON
type patient Patient

func (p *Patient) UnmarshalJSON(data []byte) (err error) {
	p2 := patient{}
	if err = json.Unmarshal(data, &p2); err == nil {
		*p = Patient(p2)
		for _, encounter := range p.Encounters {
			encounter.Patient = p
		}
		for _, condition := range p.Conditions {
			condition.Patient = p
		}
		for _, observation := range p.VitalSigns {
			observation.Patient = p
		}
		for _, procedure := range p.Procedures {
			procedure.Patient = p
		}
		for _, medication := range p.Medications {
			medication.Patient = p
		}
		for _, immunization := range p.Immunizations {
			immunization.Patient = p
		}
		for _, allergy := range p.Allergies {
			allergy.Patient = p
		}

	}
	return
}
