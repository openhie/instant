package main

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"path"
	"strings"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	pb "github.com/openhie/instant/goinstant/fhirnpmproto"
	"google.golang.org/protobuf/encoding/protojson"
)

func loadIGpackage(url_entry string) {

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

		if hdr.Name == "package/.index.json" {
			// convert the []byte to a string
			// s := string(bs)
			// fmt.Println(s)

			msg := &pb.IndexJson{}
			if err := protojson.Unmarshal(bs, msg); err != nil {
				fmt.Println(err)
				return
			}

			for _, dog := range msg.Filesrep {
				// fmt.Printf("%s\n", dog.Filename)
				// fmt.Printf("%s\n", dog.ResourceType)
				// fmt.Printf("%s\n", dog.Id)

				// get and push ths stuff
				if dog.Filename != "ig-r4.json" {
					getpushJSON(url_entry, dog.Filename, dog.ResourceType, dog.Id)
				}

			}

		}
	}

}

func getpushJSON(ig string, filename string, resourcetype string, id string) {

	trimmed := strings.Replace(ig, "index.html", "", -1)
	u, err := url.Parse(trimmed)
	if err != nil {
		fmt.Println("invalid url")
	}
	// clean url
	u.Path = path.Join(u.Path, filename)
	// fmt.Printf("%s\n", u.String())
	client := resty.New()
	resp, _ := client.R().Get(u.String())
	// fmt.Println(resp.String())

	p, err := url.Parse("http://hapi.fhir.org/baseR4")
	if err != nil {
		fmt.Println("invalid url")
	}
	p.Path = path.Join(p.Path, resourcetype, id)
	// fmt.Printf("%s\n", p.String())

	put, err := client.R().SetBody(resp.Body()).SetHeader("Content-Type", "application/fhir+json").Put(p.String())
	if err != nil {
		fmt.Println("error with put, is it the fhir url?")
	}

	if put.StatusCode() != 200 && put.StatusCode() != 201 {
		color.Set(color.FgYellow)
		fmt.Println(put.RawResponse.Status)
		fmt.Println(u.String())
		// color.Yellow(put.Status())
		fmt.Println(put.String())
		fmt.Println("")
		color.Unset()
	} else {
		color.Set(color.FgGreen)
		fmt.Println(put.RawResponse.Status, u.String())
		color.Unset()
	}
}
