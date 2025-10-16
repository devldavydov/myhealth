package main

import (
	"bytes"
	"flag"
	"html/template"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var tmplConstantsGo = template.Must(template.New("").Parse(`package constants

var TotalConstants = map[string]string {
	{{- range $key, $value := .Constants }}
	"{{ $key }}": "{{ $value }}",
	{{- end }}
}
`))

var tmplCommonJS = template.Must(template.New("").Parse(`
var Constants = {
	{{- range $key, $value := .Constants }}
	"{{ $key }}": "{{ $value }}",
	{{- end }}
};
`))

type ServiceConfig struct {
	Constants map[string]string `yaml:"constants"`
}

func main() {
	cfgFile := flag.String("in", "", "")
	flag.Parse()

	// Read config file
	cfgData, err := os.ReadFile(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	cfg := ServiceConfig{}
	if err := yaml.Unmarshal(cfgData, &cfg); err != nil {
		log.Fatal(err)
	}

	// Generate constants Go
	var buf bytes.Buffer
	if err := tmplConstantsGo.Execute(&buf, cfg); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("constants/constants_generated.go", buf.Bytes(), os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// Generate common SJ
	buf.Reset()
	if err := tmplCommonJS.Execute(&buf, cfg); err != nil {
		log.Fatal(err)
	}

	if err := os.WriteFile("static/myhealth/js/common_generated.js", buf.Bytes(), os.ModePerm); err != nil {
		log.Fatal(err)
	}
}
