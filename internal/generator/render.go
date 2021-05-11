package generator

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"github.com/isacikgoz/mattermost-suite-utilities/internal/model"
)

const templatesDir = "templates"

// Render generates the file for the struct
func Render(st *model.Struct, outputPath string) error {
	tmpl, err := initTemplate()
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, st)
	if err != nil {
		return err
	}

	os.MkdirAll(outputPath, 0700)

	dstFile := filepath.Join(outputPath, "client.go")
	return ioutil.WriteFile(dstFile, buf.Bytes(), 0664)
}

func initTemplate() (*template.Template, error) {
	data, err := ioutil.ReadFile(filepath.Join(templatesDir, "client.go.tmpl"))
	if err != nil {
		return nil, err
	}

	tmpl, err := template.New("client").Funcs(funcMap).Parse(string(data))
	if err != nil {
		return nil, err
	}

	return tmpl, nil
}
