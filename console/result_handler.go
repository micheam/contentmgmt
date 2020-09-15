package console

import (
	"context"
	"fmt"
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
			Alt:    data.Name.Value,
		}

		// TODO(michema): spec template from file
		resultTemplate := `![{{.Alt}}]({{.Scheme}}://{{.Host}}/{{.Path}})`
		t := template.Must(template.New("result-template").Parse(resultTemplate))
		t.Execute(os.Stdout, model)
		return nil
	}
}

// ListResultHandler return imgcontent.ListResultHandler which output result
// into terminal.
func ListResultHandler() imgcontent.ListResultHandler {
	return func(ctx context.Context, data imgcontent.ListOutput) error {
		for c := range data.Contents {
			fmt.Fprintln(os.Stdout, c.Path)
		}
		return nil
	}
}
