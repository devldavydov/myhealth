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

function createPage(tmpl) {
	$('#page').html(tmpl);
}

function hideElement(sel) {
	el = $(sel);
	el.addClass('d-none');
}

function showElement(sel) {
	el = $(sel);
	el.removeClass('d-none');
}

function getQueryParams() {
	return new URLSearchParams(window.location.search);
}

function tmplLoader() {
	return ` + "`" + `
	<div id="loader" class="spinner-border" role="status">
		<span class="visually-hidden">Loading...</span>
	</div>	
	` + "`" + `;
}

function tmplSearch() {
	return ` + "`" + `
	<div class="input-group">
    	<span class="input-group-text"><i class="bi bi-search"></i></span>
    	<input id="search" type="text" class="form-control" placeholder="{{ .Constants.Common_Search }}">
    	<button class="btn btn-outline-secondary" type="button" id="btnSearchClear"><i class="bi bi-x-lg"></i></button>
	</div>
	` + "`" + `;
}

function tmplToast() {
	return ` + "`" + `
	<div class="toast-container position-fixed top-0 end-0 p-3">
	<div id="liveToast" class="toast bg-light" role="alert" aria-live="assertive" aria-atomic="true">
		<div class="d-flex">
		<div id="toastBody" class="toast-body">
		</div>
		<button type="button" class="btn-close me-2 m-auto" data-bs-dismiss="toast" aria-label="Close"></button>
		</div>
	</div>
	</div>	
	` + "`" + `;
}

function tmplAlert(alClass, alID, urlBack) {
	return ` + "`" + `
	<div id="${alID}" class="d-none">
		<div class="msg alert ${alClass}" role="alert">
		</div>
		<div class="mb-3">
    		<a href="${urlBack}" class="btn btn-primary"><i class="bi bi-arrow-90deg-left"></i></a>
		</div>
	</div>
	` + "`" + `;
}
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
