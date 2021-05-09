package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
)

func loadIGpackage(fhir_server string, url_entry string) {

	// trim index.html if there
	trimmed := strings.Replace(url_entry, "index.html", "", -1)
	u, err := url.Parse(trimmed)
	if err != nil {
		fmt.Println("invalid url")
	}

	// clean url
	u.Path = path.Join(u.Path, "package.tgz")
	// fmt.Printf("%s\n", u.String())

	client := resty.New()
	resp, _ := client.R().Get(u.String())
	// if err != nil {
	// 	fmt.Println("Please check the URL again for the IG. It should not end in index.html")
	// }
	fmt.Println("Status Code:", resp.StatusCode())
	if resp.StatusCode() != 200 {
		fmt.Println("Check that the URL for the IG is correct.")
	}

	// x := resp.Body()
	// fmt.Println("type:", reflect.TypeOf(x))

	y := bytes.NewReader(resp.Body())

	archive, err := gzip.NewReader(y)
	// archive, err := gzip.NewReader(resp.RawBody())

	if err != nil {
		fmt.Println("There is a problem - is this a tgz?")
	}
	tr := tar.NewReader(archive)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Printf("Contents of %s:\n", hdr.Name)

		// read the complete content of the file h.Name into the bs []byte
		bs, _ := ioutil.ReadAll(tr)

		type IndexJSON struct {
			IndexVersion int32 `json:"index-version"`
			Filesrep     []struct {
				Filename     string `json:"filename"`
				ResourceType string `json:"resourceType"`
				Id           string `json:"id"`
				Url          string `json:"url"`
				Version      string `json:"version"`
				Kind         string `json:"kind"`
				Type         string `json:"type"`
			} `json:"files"`
		}

		if hdr.Name == "package/.index.json" {
			// convert the []byte to a string
			// s := string(bs)
			// fmt.Println(s)

			var msg IndexJSON
			err := json.Unmarshal(bs, &msg)
			if err != nil {
				panic(err)
			}

			// msg := &pb.IndexJson{}
			// if err := protojson.Unmarshal(bs, msg); err != nil {
			// 	fmt.Println(err)
			// 	return
			// }

			for _, CodeSystem := range msg.Filesrep {
				if CodeSystem.ResourceType == "CodeSystem" {
					getpushJSON(fhir_server, url_entry, CodeSystem.Filename, CodeSystem.ResourceType, CodeSystem.Id)
				}
			}
			for _, ValueSet := range msg.Filesrep {
				if ValueSet.ResourceType == "ValueSet" {
					getpushJSON(fhir_server, url_entry, ValueSet.Filename, ValueSet.ResourceType, ValueSet.Id)
				}
			}
			for _, ConceptMap := range msg.Filesrep {
				if ConceptMap.ResourceType == "ConceptMap" {
					getpushJSON(fhir_server, url_entry, ConceptMap.Filename, ConceptMap.ResourceType, ConceptMap.Id)
				}
			}
			for _, SearchParameter := range msg.Filesrep {
				if SearchParameter.ResourceType == "SearchParameter" {
					getpushJSON(fhir_server, url_entry, SearchParameter.Filename, SearchParameter.ResourceType, SearchParameter.Id)
				}
			}
			for _, OperationDefinition := range msg.Filesrep {
				if OperationDefinition.ResourceType == "OperationDefinition" {
					getpushJSON(fhir_server, url_entry, OperationDefinition.Filename, OperationDefinition.ResourceType, OperationDefinition.Id)
				}
			}
			for _, StructureDefinition := range msg.Filesrep {
				if StructureDefinition.ResourceType == "StructureDefinition" {
					getpushJSON(fhir_server, url_entry, StructureDefinition.Filename, StructureDefinition.ResourceType, StructureDefinition.Id)
				}
			}
			for _, CapabilityStatement := range msg.Filesrep {
				if CapabilityStatement.ResourceType == "CapabilityStatement" {
					getpushJSON(fhir_server, url_entry, CapabilityStatement.Filename, CapabilityStatement.ResourceType, CapabilityStatement.Id)
				}
			}
			for _, Organization := range msg.Filesrep {
				if Organization.ResourceType == "Organization" {
					getpushJSON(fhir_server, url_entry, Organization.Filename, Organization.ResourceType, Organization.Id)
				}
			}
			for _, Location := range msg.Filesrep {
				if Location.ResourceType == "Location" {
					getpushJSON(fhir_server, url_entry, Location.Filename, Location.ResourceType, Location.Id)
				}
			}
			for _, Library := range msg.Filesrep {
				if Library.ResourceType == "Library" {
					getpushJSON(fhir_server, url_entry, Library.Filename, Library.ResourceType, Library.Id)
				}
			}
			for _, Questionnaire := range msg.Filesrep {
				if Questionnaire.ResourceType == "Questionnaire" {
					getpushJSON(fhir_server, url_entry, Questionnaire.Filename, Questionnaire.ResourceType, Questionnaire.Id)
				}
			}
			for _, QuestionnaireResponse := range msg.Filesrep {
				if QuestionnaireResponse.ResourceType == "QuestionnaireResponse" {
					getpushJSON(fhir_server, url_entry, QuestionnaireResponse.Filename, QuestionnaireResponse.ResourceType, QuestionnaireResponse.Id)
				}
			}

			for _, dog := range msg.Filesrep {
				// get and push ths stuff
				// if dog.Filename != "ig-r4.json" {
				// 	getpushJSON(fhir_server, url_entry, dog.Filename, dog.ResourceType, dog.Id)
				// getpushJSON(fhir_server, url_entry, dog.Filename, dog.ResourceType, dog.Id)
				// if dog.Filename != "ig-r4.json" && dog.ResourceType != "ImplementationGuide" {
				if dog.Filename != "ig-r4.json" {
					getpushJSON(fhir_server, url_entry, dog.Filename, dog.ResourceType, dog.Id)
				}
			}
		}
	}
}

func getpushJSON(fhir_server string, ig string, filename string, resourcetype string, id string) {

	trimmed := strings.Replace(ig, "index.html", "", -1)
	u, err := url.Parse(trimmed)
	if err != nil {
		fmt.Println("invalid url")
	}
	// clean url
	u.Path = path.Join(u.Path, filename)
	client := resty.New()
	// client.SetDebug(true)
	resp, _ := client.R().Get(u.String())

	p, err := url.Parse(fhir_server)
	if err != nil {
		fmt.Println("invalid url")
	}
	p.Path = path.Join(p.Path, resourcetype, id)

	// TODO: On some IGs this panics: "panic: runtime error: invalid memory address or nil pointer dereference"
	put, err := client.R().SetBody(resp.Body()).SetHeader("Content-Type", "application/fhir+json").Put(p.String())
	if err != nil {
		fmt.Println("error with put, is it the fhir url?")
		fmt.Println(ig, filename, resourcetype, id)
	}

	if put.StatusCode() != 200 && put.StatusCode() != 201 {
		color.Set(color.FgYellow)
		// fmt.Println(put.RawResponse.Status) // this causes the panic
		fmt.Println(u.String())
		// color.Yellow(put.Status())
		fmt.Println(put.String())
		fmt.Println("")
		color.Unset()
	} else {
		color.Set(color.FgGreen)
		fmt.Println(put.RawResponse.Status, filename)
		color.Unset()
	}
}
