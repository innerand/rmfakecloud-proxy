package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"
)

const Homepage = "https://github.com/yi-jiayu/secure"

var tmpl = template.Must(template.New("version.go").Parse(`//go:generate go run generate/versioninfo.go
// Code generated by versioninfo.go; DO NOT EDIT.

package main

const Version = "secure {{ .Version }} ({{ .GOOS }}-{{ .GOARCH }}) {{ .GoVersion }}\n{{ .Homepage }}"
`))

type VersionInfo struct {
	Version   string
	GOOS      string
	GOARCH    string
	GoVersion string
	Homepage  string
}

func main() {
	describeCmd := exec.Command("git", "describe", "--tags", "--always")
	version, err := describeCmd.Output()
	if err != nil {
		log.Fatal(fmt.Errorf("error describing commit: %v", err))
	}

	goos := os.Getenv("GOOS")
	if goos == "" {
		goos = runtime.GOOS
	}

	goarch := os.Getenv("GOARCH")
	if goarch == "" {
		goarch = runtime.GOARCH
	}

	goVersion := runtime.Version()

	vi := VersionInfo{
		Version:   strings.TrimSpace(string(version)),
		GOARCH:    goarch,
		GOOS:      goos,
		GoVersion: goVersion,
		Homepage:  Homepage,
	}

	f, err := os.Create("version.go")
	if err != nil {
		log.Fatal(fmt.Errorf("error creating version.go: %v", err))
	}

	err = tmpl.Execute(f, vi)
	if err != nil {
		log.Fatal(fmt.Errorf("error executing template: %v", err))
	}
}
