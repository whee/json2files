package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"text/template"
)

type JSONData map[string]interface{}

func tmplFilename(t *template.Template, jsd JSONData) (string, error) {
	var fnBuf bytes.Buffer
	err := t.Execute(&fnBuf, jsd)
	if err != nil {
		return "", err
	}
	return fnBuf.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: json2files [output filename template]\n")
	}

	fileTmpl := os.Args[1]

	tmpl, err := template.New("filename").Parse(fileTmpl)
	if err != nil {
		log.Fatal(err)
	}

	dec := json.NewDecoder(os.Stdin)

	for {
		var jsd JSONData
		if err := dec.Decode(&jsd); err != nil {
			log.Fatal(err)
		}

		fn, err := tmplFilename(tmpl, jsd)
		if err != nil {
			log.Fatal(err)
		}

		basepath := path.Dir(fn)
		if err := os.MkdirAll(basepath, 0770); err != nil {
			log.Fatal(err)
		}

		outJSON, err := json.MarshalIndent(jsd, "", " ")
		if err != nil {
			log.Fatal(err)
		}

		if err := ioutil.WriteFile(fn, outJSON, 0640); err != nil {
			log.Fatal(err)
		}
	}
}
