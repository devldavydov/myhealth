package main

import (
	"bytes"
	"flag"
	"log"
	"os"
	"text/template"

	"github.com/stretchr/testify/assert/yaml"
)

var tmpl = template.Must(template.New("").Parse(`package cmdproc

// Code generated by "go generate". DO NOT EDIT!
{{ $cfg := . }}
import (
	"fmt"
	"time"
	"strconv"
	"strings"	

	"github.com/devldavydov/myhealth/internal/storage"

	"go.uber.org/zap"
	tele "gopkg.in/telebot.v4"
)	

func (r *CmdProcessor) process(c tele.Context, cmd string, userID int64) error {
	cmdParts := []string{}
	for _, part := range strings.Split(cmd, ",") {
		cmdParts = append(cmdParts, strings.Trim(part, " "))
	}

	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid command",
			zap.String("command", cmd),
			zap.Int64("userID", userID),
		)
		return c.Send(MsgErrInvalidCommand)
	}

	var resp []CmdResponse

	switch cmdParts[0] {
	{{ range .Config.Commands -}}
	case "{{ .Name }}":
		resp = r.process_{{ .Name }}("{{ .Name }}", cmdParts[1:], userID)
	{{ end -}}
	case "h":
		resp = r.processHelp()
	default:
		r.logger.Error(
			"unknown command",
			zap.String("command", cmd),
			zap.Int64("userID", userID),
		)
		resp = NewSingleCmdResponse(MsgErrInvalidCommand)
	}	

	if r.debugMode {
		if err := c.Send("!!! ОТЛАДОЧНЫЙ РЕЖИМ !!!"); err != nil {
			return err
		}
	}

	for _, rItem := range resp {
		if err := c.Send(rItem.what, rItem.opts...); err != nil {
			return err
		}
	}

	return nil	
}
{{ range $cfg.Config.Commands }}
{{ with $cmd := . -}}
func (r *CmdProcessor) process_{{ $cmd.Name }}(baseCmd string, cmdParts []string, userID int64) []CmdResponse {
	if len(cmdParts) == 0 {
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		return NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	var resp []CmdResponse

	switch cmdParts[0] {
	{{ range $cmd.SubCommands -}}
	case "{{ .Name }}":
		{{- if (ne (len .Args) 0) }}
		if len(cmdParts[1:]) != {{ len .Args }} {
			return NewSingleCmdResponse(MsgErrInvalidArgsCount)
		}
		
		cmdParts = cmdParts[1:]
		{{ range $index, $arg := .Args }}
		{{- if (eq $arg.Type "timestamp") }}
		val{{ $index }}, err := parseTimestamp(r.tz, cmdParts[{{ $index }}])
		{{ end -}}
		{{- if (eq $arg.Type "floatG0") }}
		val{{ $index }}, err := parseFloatG0(cmdParts[{{ $index }}])
		{{ end -}}
		{{- if (eq $arg.Type "floatGE0") }}
		val{{ $index }}, err := parseFloatGE0(cmdParts[{{ $index }}])
		{{ end -}}
		{{- if (eq $arg.Type "stringG0") }}
		val{{ $index }}, err := parseStringG0(cmdParts[{{ $index }}])
		{{ end -}}
		{{- if (eq $arg.Type "stringGE0") }}
		val{{ $index }}, err := parseStringGE0(cmdParts[{{ $index }}])
		{{ end -}}
		{{- if (eq $arg.Type "gender") }}
		val{{ $index }}, err := parseGender(cmdParts[{{ $index }}])
		{{ end -}}
		{{- if (eq $arg.Type "meal") }}
		val{{ $index }}, err := parseMeal(cmdParts[{{ $index }}])
		{{ end -}}
		{{- if (eq $arg.Type "stringArr") }}
		val{{ $index }}, err := parseStringArr(cmdParts[{{ $index }}])
		{{ end -}}
		{{- if (eq $arg.Type "intArr") }}
		val{{ $index }}, err := parseIntArr(cmdParts[{{ $index }}])
		{{ end -}}			 
		if err != nil {
			return argError("{{ $arg.Name }}")
		}
		{{ end }}
		resp = r.{{ .Func }}(
			userID,
			{{ range $index, $arg := .Args -}}
			val{{ $index }},
			{{ end -}}
		)
		{{ else }}
		resp = r.{{ .Func }}(userID)
		{{ end }}		
	{{ end -}}
	case "h":
		return NewSingleCmdResponse(
			newCmdHelpBuilder(baseCmd, "{{ $cmd.Description }}").
			{{ range $cmd.SubCommands -}}
			{{ if (eq .Comment "") -}}
			addCmd(
				"{{ .Description }}",
				"{{ .Name }}",
				{{ range .Args -}}
				"{{ .Name }} [{{ (index $cfg.TypesMap .Type).DescriptionShort }}]",
				{{ end -}}
			).
			{{ else -}}
			addCmdWithComment(
				"{{ .Description }}",
				"{{ .Name }}",
				"{{ .Comment }}",
				{{ range .Args -}}
				"{{ .Name }} [{{ (index $cfg.TypesMap .Type).DescriptionShort }}]",
				{{ end -}}
			).	
			{{ end -}}		 
			{{ end -}}
			build(),
		optsHTML)

	default:
		r.logger.Error(
			"invalid command",
			zap.Strings("cmdParts", cmdParts),
			zap.Int64("userID", userID),
		)
		resp = NewSingleCmdResponse(MsgErrInvalidCommand)
	}

	return resp
}
{{ end -}}
{{ end }}
func (r *CmdProcessor) processHelp() []CmdResponse {
	var sb strings.Builder
	sb.WriteString("<b>Команды помощи по разделам:</b>\n")
	{{- range $cfg.Config.Commands }}
	sb.WriteString("<b>\u2022 {{ .Name }},h</b> - {{ .DescriptionShort }}\n")
	{{- end }}
	sb.WriteString("\n<b>Типы данных:</b>\n")
	{{- range $cfg.Config.Types }}
	sb.WriteString("<b>\u2022 {{ .DescriptionShort }}</b> - {{ .Description }}\n")
	{{- end }}
	return NewSingleCmdResponse(sb.String(), optsHTML)
}

func parseTimestamp(tz *time.Location, arg string) (time.Time, error) {
	var t time.Time
	var err error

	if arg == "" {
		t = time.Now().In(tz)
	} else {
		t, err = time.Parse("02.01.2006", arg)
		if err != nil {
			return time.Time{}, err
		}
	}

	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, tz), nil
}

func parseFloatG0(arg string) (float64, error) {
	val, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return 0, err
	}

	if val <= 0 {
		return 0, fmt.Errorf("not above zero")
	}

	return val, nil
}

func parseFloatGE0(arg string) (float64, error) {
	val, err := strconv.ParseFloat(arg, 64)
	if err != nil {
		return 0, err
	}

	if val < 0 {
		return 0, fmt.Errorf("not above or equal zero")
	}

	return val, nil
}

func parseStringG0(arg string) (string, error) {
	if len(arg) == 0 {
		return "", fmt.Errorf("empty string")
	}
	
	return arg, nil
}

func parseStringGE0(arg string) (string, error) {
	return arg, nil
}

func parseGender(arg string) (string, error) {
	switch arg {
	case "m", "f":
		return arg, nil
	default:
		return "", fmt.Errorf("wrong gender")
	}
}

func parseMeal(arg string) (storage.Meal, error) {
	return storage.NewMealFromString(arg)
}

func parseStringArr(arg string) ([]string, error) {
	parts := []string{}
	for _, part := range strings.Split(arg, "|") {
		parts = append(parts, strings.Trim(part, " "))
	}

	if len(parts) == 0 {
		return nil, fmt.Errorf("empty array")
	}

	return parts, nil
}

func parseIntArr(arg string) ([]int64, error) {
	parts := []int64{}
	for _, part := range strings.Split(arg, "|") {
		part = strings.Trim(part, " ")
		val, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			return nil, err
		}

		parts = append(parts, val)
	}

	if len(parts) == 0 {
		return nil, fmt.Errorf("empty array")
	}

	return parts, nil
}

func argError(argName string) []CmdResponse {
	return NewSingleCmdResponse(fmt.Sprintf("%s: %s", MsgErrInvalidArg, argName))
}

func formatTimestamp(ts time.Time) string {
	return ts.Format("02.01.2006")
}

type cmdHelpItem struct {
	label   string
	cmd     string
	comment string
	args    []string
}

type cmdHelpBuilder struct {
	baseCmd string
	label   string
	items   []cmdHelpItem
}

func newCmdHelpBuilder(baseCmd, label string) *cmdHelpBuilder {
	return &cmdHelpBuilder{baseCmd: baseCmd, label: label}
}

func (r *cmdHelpBuilder) addCmd(label, cmd string, args ...string) *cmdHelpBuilder {
	r.items = append(r.items, cmdHelpItem{
		label: label,
		cmd:   cmd,
		args:  args,
	})
	return r
}

func (r *cmdHelpBuilder) addCmdWithComment(label, cmd, comment string, args ...string) *cmdHelpBuilder {
	r.items = append(r.items, cmdHelpItem{
		label:   label,
		cmd:     cmd,
		comment: comment,
		args:    args,
	})
	return r
}

func (r *cmdHelpBuilder) build() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("<b>%s</b>\n", r.label))
	for i, item := range r.items {
		sb.WriteString(fmt.Sprintf("<b>\u2022 %s</b>\n", item.label))
		sb.WriteString(fmt.Sprintf("%s,%s", r.baseCmd, item.cmd))

		if len(item.args) > 0 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}

		for j, arg := range item.args {
			sArg := arg
			if strings.Contains(sArg, "|") {
				parts := strings.Split(sArg, "|")
				sArg = fmt.Sprintf("%s\n ИЛИ\n %s", parts[0], parts[1])
			}

			if j == len(item.args)-1 {
				sb.WriteString(fmt.Sprintf(" %s\n", sArg))
			} else {
				sb.WriteString(fmt.Sprintf(" %s,\n", sArg))
			}
		}

		if item.comment != "" {
			sb.WriteString(fmt.Sprintf("\n<i>Примечание</i>: %s\n", item.comment))
		}

		if i != len(r.items)-1 {
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
`))

type CommandProcessorConfig struct {
	Commands []Command  `yaml:"commands"`
	Types    []DataType `yaml:"types"`
}

type Command struct {
	Name             string       `yaml:"name"`
	Description      string       `yaml:"description"`
	DescriptionShort string       `yaml:"description_short"`
	SubCommands      []SubCommand `yaml:"subcommands"`
}

type SubCommand struct {
	Name        string `yaml:"name"`
	Func        string `yaml:"func"`
	Description string `yaml:"description"`
	Comment     string `yaml:"comment"`
	Args        []Arg  `yaml:"args"`
}

type Arg struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
}

type DataType struct {
	Name             string `yaml:"name"`
	Description      string `yaml:"description"`
	DescriptionShort string `yaml:"description_short"`
}

func main() {
	cfgFile := flag.String("in", "", "")
	outFileName := flag.String("out", "", "")
	flag.Parse()

	// Read config file
	cfgData, err := os.ReadFile(*cfgFile)
	if err != nil {
		log.Fatal(err)
	}

	cfg := CommandProcessorConfig{}
	if err := yaml.Unmarshal(cfgData, &cfg); err != nil {
		log.Fatal(err)
	}

	// Generate template
	type tmplData struct {
		Config   *CommandProcessorConfig
		TypesMap map[string]DataType
	}

	typesMap := map[string]DataType{}
	for _, t := range cfg.Types {
		typesMap[t.Name] = t
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, tmplData{
		Config:   &cfg,
		TypesMap: typesMap,
	}); err != nil {
		log.Fatal(err)
	}

	// Save to out file
	if err := os.WriteFile(*outFileName, buf.Bytes(), os.ModePerm); err != nil {
		log.Fatal(err)
	}
}
