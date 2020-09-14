package console

import (
	"context"
	"html/template"
	"os"

	"github.com/micheam/imgcontent"
)

// UploadResultHandler ...
func UploadResultHandler() imgcontent.UploadResultHandler {
	return func(ctx context.Context, data imgcontent.UploadOutput) error {
		model := struct {
			Scheme, Host, Path, Alt string
		}{
			Scheme: data.URL.Scheme,
			Host:   data.URL.Host,
			Path:   data.URL.EscapedPath(),
			Alt:    data.Filename.Value,
		}

		// TODO(michema): spec template from file
		resultTemplate := `![{{.Alt}}]({{.Scheme}}://{{.Host}}/{{.Path}})`
		t := template.Must(template.New("result-template").Parse(resultTemplate))
		t.Execute(os.Stdout, model)
		return nil
	}
}
