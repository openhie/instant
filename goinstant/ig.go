package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/pkg/errors"
)

type indexJSON struct {
	IndexVersion int32      `json:"index-version"`
	Files        []filesRep `json:"files"`
}

type filesRep struct {
	Filename     string `json:"filename"`
	ResourceType string `json:"resourceType"`
	Id           string `json:"id"`
	Url          string `json:"url"`
	Version      string `json:"version"`
	Kind         string `json:"kind"`
	Type         string `json:"type"`
}

func loadIGpackage(url_entry string, fhir_server string, params *Params) error {
	trimmed := strings.Replace(url_entry, "index.html", "", -1)
	u, err := url.Parse(trimmed)
	if err != nil {
		return errors.Wrap(err, "Invalid url")
	}

	// clean url
	u.Path = path.Join(u.Path, "package.tgz")

	client := resty.New()
	resp, err := client.R().Get(u.String())
	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		return errors.Errorf("Status code: %v. Check that the URL for the IG is correct", resp.StatusCode())
	}
	fmt.Printf("Reached Published IG with Status Code: %v\n", resp.StatusCode())

	y := bytes.NewReader(resp.Body())
	archive, err := gzip.NewReader(y)
	if err != nil {
		return errors.New("There is a problem - is this a tgz?")
	}
	tr := tar.NewReader(archive)

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		bs, err := ioutil.ReadAll(tr)
		if err != nil {
			return err
		}

		if hdr.Name == "package/.index.json" {
			var msg indexJSON
			err := json.Unmarshal(bs, &msg)
			if err != nil {
				gracefulPanic(err, "")
			}

			// Order mostly from: https://github.com/nmdp-bioinformatics/igloader/blob/main/igloader/igloader.py#L33
			// and https://github.com/hapifhir/hapi-fhir/blob/75c74a22dbd1f0dde3631b540d1898eef2a2666f/hapi-fhir-jpaserver-base/src/main/java/ca/uhn/fhir/jpa/packages/PackageInstallerSvcImpl.java#L85-L93

			stuff := []string{"NamingSystem", "CodeSystem", "ValueSet", "OperationDefinition", "StructureDefinition", "ConceptMap", "SearchParameter", "Subscription", "CapabilityStatement"}
			color.Set(color.FgBlue)
			fmt.Printf("Loading primary conformance resources %s\n", stuff)
			color.Unset()
			for _, y := range stuff {
				for _, x := range msg.Files {
					if x.ResourceType == y {
						err = getpushJSON(fhir_server, url_entry, x.Filename, x.ResourceType, false, x.Id, params)
						if err != nil {
							return err
						}
					}
				}
			}

			stuff2 := []string{"Patient", "Practitioner", "Organization", "Location", "Library", "Measure", "MeasureReport", "Questionnaire", "QuestionnaireResponse", "Procedure"}
			color.Set(color.FgBlue)
			fmt.Printf("Loading other singular resources %s\n", stuff2)
			color.Unset()
			for _, b := range stuff2 {
				for _, a := range msg.Files {
					if a.ResourceType == b {
						err = getpushJSON(fhir_server, url_entry, a.Filename, a.ResourceType, false, a.Id, params)
						if err != nil {
							return err
						}
					}
				}
			}

			color.Blue("2nd pass: Load resources again (except ig-r4.json or bundles) to address customized dependencies in IGs.")
			for _, dog := range msg.Files {
				if dog.ResourceType != "Bundle" && dog.ResourceType != "ImplementationGuide" {
					err = getpushJSON(fhir_server, url_entry, dog.Filename, dog.ResourceType, false, dog.Id, params)
					if err != nil {
						return err
					}
				}
			}

			color.Blue("3rd pass - Explicit Bundles (not Structure Definitions)")
			for _, cat := range msg.Files {
				if cat.ResourceType == "Bundle" && cat.Type == "transaction" {
					err = getpushJSON(fhir_server, url_entry, cat.Filename, cat.ResourceType, true, cat.Id, params)
					if err != nil {
						return err
					}
				}
			}

			color.Blue("3rd pass - Implementation Guide")
			for _, mouse := range msg.Files {
				if mouse.Filename != "ig-r4.json" && mouse.ResourceType == "ImplementationGuide" {
					err = getpushJSON(fhir_server, url_entry, mouse.Filename, mouse.ResourceType, false, mouse.Id, params)
					if err != nil {
						return err
					}
				}
			}

			color.Green("If there are still errors, you may choose to run the tool again.")
		}
	}

	return nil
}

func getpushJSON(fhir_server string, ig string, filename string, resourcetype string, bundle bool, id string, params *Params) error {
	trimmed := strings.Replace(ig, "index.html", "", -1)
	u, err := url.Parse(trimmed)
	if err != nil {
		fmt.Println("invalid url, must end in index.html or be a directory")
	}
	// clean url
	u.Path = path.Join(u.Path, filename)

	client := resty.New()
	resp, err := client.R().Get(u.String())
	if err != nil {
		return err
	}

	p, err := url.Parse(fhir_server)
	if err != nil {
		return errors.Wrap(err, "invalid url")
	}

	if !bundle {
		p.Path = path.Join(p.Path, resourcetype, id)
	}

	if params == nil {
		return errors.New("Nil pointer... variable 'params' not initialised")
	}
	switch params.TypeAuth {
	case "None":
		if bundle {
			put, err := client.R().SetBody(resp.Body()).
				SetHeader("Content-Type", "application/fhir+json").Post(p.String())
			if err != nil {
				return errors.Wrapf(err, "Status code: %v. Tracing... IG %v, Filename %v, Resource type %v, ID %v", put.StatusCode(), ig, filename, resourcetype, id)
			}
			code := put.StatusCode()
			status := put.RawResponse.Status
			url := u.String()
			file := put.String()
			printStatus(code, status, url, file, filename)

		} else {
			put, err := client.R().SetBody(resp.Body()).
				SetHeader("Content-Type", "application/fhir+json").Put(p.String())
			if err != nil {
				return errors.Wrapf(err, "Status code: %v. Tracing... IG %v, Filename %v, Resource type %v, ID %v", put.StatusCode(), ig, filename, resourcetype, id)
			}
			code := put.StatusCode()
			status := put.RawResponse.Status
			url := u.String()
			file := put.String()
			printStatus(code, status, url, file, filename)

		}

	case "Basic":
		if bundle {
			put, err := client.R().SetBody(resp.Body()).
				SetHeader("Content-Type", "application/fhir+json").
				SetBasicAuth(params.BasicUser, params.BasicPass).Post(p.String())
			if err != nil {
				return errors.Wrapf(err, "Status code: %v. Tracing... IG %v, Filename %v, Resource type %v, ID %v", put.StatusCode(), ig, filename, resourcetype, id)
			}
			code := put.StatusCode()
			status := put.RawResponse.Status
			url := u.String()
			file := put.String()
			printStatus(code, status, url, file, filename)

		} else {
			put, err := client.R().SetBody(resp.Body()).
				SetHeader("Content-Type", "application/fhir+json").
				SetBasicAuth(params.BasicUser, params.BasicPass).Put(p.String())
			if err != nil {
				return errors.Wrapf(err, "Status code: %v. Tracing... IG %v, Filename %v, Resource type %v, ID %v", put.StatusCode(), ig, filename, resourcetype, id)
			}
			code := put.StatusCode()
			status := put.RawResponse.Status
			url := u.String()
			file := put.String()
			printStatus(code, status, url, file, filename)

		}
	case "Token":
		if bundle {
			put, err := client.R().SetBody(resp.Body()).
				SetHeader("Content-Type", "application/fhir+json").
				SetAuthToken(params.Token).Post(p.String())
			if err != nil {
				return errors.Wrapf(err, "Status code: %v. Tracing... IG %v, Filename %v, Resource type %v, ID %v", put.StatusCode(), ig, filename, resourcetype, id)
			}
			code := put.StatusCode()
			status := put.RawResponse.Status
			url := u.String()
			file := put.String()
			printStatus(code, status, url, file, filename)

		} else {
			put, err := client.R().SetBody(resp.Body()).
				SetHeader("Content-Type", "application/fhir+json").
				SetAuthToken(params.Token).Put(p.String())
			if err != nil {
				return errors.Wrapf(err, "Status code: %v. Tracing... IG %v, Filename %v, Resource type %v, ID %v", put.StatusCode(), ig, filename, resourcetype, id)
			}
			code := put.StatusCode()
			status := put.RawResponse.Status
			url := u.String()
			file := put.String()
			printStatus(code, status, url, file, filename)

		}
	case "Custom":
		custom := "Custom" + " " + params.Token
		if bundle {
			put, err := client.R().SetBody(resp.Body()).
				SetHeader("Content-Type", "application/fhir+json").
				SetHeader("Authorization", custom).Put(p.String())
			if err != nil {
				return errors.Wrapf(err, "Status code: %v. Tracing... IG %v, Filename %v, Resource type %v, ID %v", put.StatusCode(), ig, filename, resourcetype, id)
			}
			code := put.StatusCode()
			status := put.RawResponse.Status
			url := u.String()
			file := put.String()
			printStatus(code, status, url, file, filename)

		} else {
			put, err := client.R().SetBody(resp.Body()).
				SetHeader("Content-Type", "application/fhir+json").
				SetHeader("Authorization", custom).Put(p.String())
			if err != nil {
				return errors.Wrapf(err, "Status code: %v. Tracing... IG %v, Filename %v, Resource type %v, ID %v", put.StatusCode(), ig, filename, resourcetype, id)
			}
			code := put.StatusCode()
			status := put.RawResponse.Status
			url := u.String()
			file := put.String()
			printStatus(code, status, url, file, filename)
		}
	}
	return nil
}

func printStatus(code int, status string, url string, file, filename string) {
	if code != 200 && code != 201 {
		color.Set(color.FgYellow)
		fmt.Println(status)
		fmt.Println(url)
		fmt.Println(file)
		fmt.Println("")
		color.Unset()
	} else {
		color.Set(color.FgGreen)
		fmt.Println(status, filename)
		color.Unset()
	}
}
