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

func loadIGpackage(url_entry string, fhir_server string, params *Params) {

	// debug.SetGCPercent(200)
	// debug.SetMaxStack(2000000000)
	trimmed := strings.Replace(url_entry, "index.html", "", -1)
	u, err := url.Parse(trimmed)
	if err != nil {
		fmt.Println("invalid url")
	}

	// clean url
	u.Path = path.Join(u.Path, "package.tgz")

	client := resty.New()
	resp, _ := client.R().Get(u.String())

	fmt.Println("Reached Published IG with Status Code:", resp.StatusCode())
	if resp.StatusCode() != 200 {
		fmt.Println("Check that the URL for the IG is correct.")
	}

	y := bytes.NewReader(resp.Body())
	archive, err := gzip.NewReader(y)

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

			// Order mostly from: https://github.com/nmdp-bioinformatics/igloader/blob/main/igloader/igloader.py#L33
			// and https://github.com/hapifhir/hapi-fhir/blob/75c74a22dbd1f0dde3631b540d1898eef2a2666f/hapi-fhir-jpaserver-base/src/main/java/ca/uhn/fhir/jpa/packages/PackageInstallerSvcImpl.java#L85-L93

			stuff := []string{"NamingSystem", "CodeSystem", "ValueSet", "OperationDefinition", "StructureDefinition", "ConceptMap", "SearchParameter", "Subscription", "CapabilityStatement"}
			color.Set(color.FgBlue)
			fmt.Printf("Loading %s\n", stuff)
			color.Unset()
			for _, y := range stuff {
				for _, x := range msg.Filesrep {
					if x.ResourceType == y && x.Filename != "ig-r4.json" {
						getpushJSON(fhir_server, url_entry, x.Filename, x.ResourceType, x.Id, params)
					}
				}
			}

			stuff2 := []string{"Patient", "Practitioner", "Organization", "Location", "Library", "Measure", "MeasureReport", "Questionnaire", "QuestionnaireResponse", "Procedure", "ImplementationGuide"}
			color.Set(color.FgBlue)
			fmt.Printf("Loading %s\n", stuff2)
			color.Unset()
			for _, b := range stuff2 {
				for _, a := range msg.Filesrep {
					if a.ResourceType == b && a.Filename != "ig-r4.json" {
						getpushJSON(fhir_server, url_entry, a.Filename, a.ResourceType, a.Id, params)
					}
				}
			}

			color.Blue("2nd pass: Load resources again (except ig-r4.json) to address customized dependencies in IGs.")
			for _, dog := range msg.Filesrep {
				// get and push ths stuff
				// if dog.Filename != "ig-r4.json" {
				// 	getpushJSON(fhir_server, url_entry, dog.Filename, dog.ResourceType, dog.Id)
				// getpushJSON(fhir_server, url_entry, dog.Filename, dog.ResourceType, dog.Id)
				// if dog.Filename != "ig-r4.json" && dog.ResourceType != "ImplementationGuide" {
				if dog.Filename != "ig-r4.json" {
					getpushJSON(fhir_server, url_entry, dog.Filename, dog.ResourceType, dog.Id, params)
				}
			}
			color.Green("If there are still errors, you may choose to run the tool again.")
		}
	}
}

func getpushJSON(fhir_server string, ig string, filename string, resourcetype string, id string, params *Params) {

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

	switch params.TypeAuth {
	// TODO: On some IGs this panics: "panic: runtime error: invalid memory address or nil pointer dereference"
	case "None":
		put, err := client.R().SetBody(resp.Body()).
			SetHeader("Content-Type", "application/fhir+json").
			Put(p.String())
		if err != nil {
			fmt.Println("error with put, is it the fhir url?")
			fmt.Println(ig, filename, resourcetype, id)
		}

		if put.StatusCode() != 200 && put.StatusCode() != 201 {
			color.Set(color.FgYellow)
			fmt.Println(put.RawResponse.Status) // this causes the panic
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
	case "Basic":
		put, err := client.R().SetBody(resp.Body()).
			SetHeader("Content-Type", "application/fhir+json").
			SetBasicAuth(params.BasicUser, params.BasicPass).
			Put(p.String())
		if err != nil {
			fmt.Println("error with put, is it the fhir url?")
			fmt.Println(ig, filename, resourcetype, id)
		}

		if put.StatusCode() != 200 && put.StatusCode() != 201 {
			color.Set(color.FgYellow)
			fmt.Println(put.RawResponse.Status) // this causes the panic
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
	case "Token":
		put, err := client.R().SetBody(resp.Body()).
			SetHeader("Content-Type", "application/fhir+json").
			SetAuthToken(params.Token).
			Put(p.String())
		if err != nil {
			fmt.Println("error with put, is it the fhir url?")
			fmt.Println(ig, filename, resourcetype, id)
		}

		if put.StatusCode() != 200 && put.StatusCode() != 201 {
			color.Set(color.FgYellow)
			fmt.Println(put.RawResponse.Status) // this causes the panic
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

	case "Custom":
		custom := "Custom" + " " + params.Token
		put, err := client.R().SetBody(resp.Body()).
			SetHeader("Content-Type", "application/fhir+json").
			SetHeader("Authorization", custom).
			Put(p.String())
		if err != nil {
			fmt.Println("error with put, is it the fhir url?")
			fmt.Println(ig, filename, resourcetype, id)
		}

		if put.StatusCode() != 200 && put.StatusCode() != 201 {
			color.Set(color.FgYellow)
			fmt.Println(put.RawResponse.Status) // this causes the panic
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

}
