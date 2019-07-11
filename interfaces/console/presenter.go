package console

import (
	"context"
	"github.com/micheam/imgcontent/usecases"
	"os"
	"text/template"
)

var _ usecases.UploadPresenter = (*ConsoleUploadPresenter)(nil)
var _ usecases.UploadPresenter = ConsoleUploadPresenter{}

// TODO: Presenter はそもそもこれで良いのか？
type ConsoleUploadPresenter struct{}

func (c ConsoleUploadPresenter) Complete(ctx context.Context, data usecases.UploadOutput) error {

	model := struct {
		Scheme, Host, Path, Alt string
	}{
		Scheme: data.URL.Scheme,
		Host: data.URL.Host,
		Path: data.URL.EscapedPath(),
		Alt: data.Filename.Value,
	}

    const tpl = `<img src="{{.Scheme}}://{{.Host}}/{{.Path}}" alt="{{.Alt}}" />`
	t := template.Must(template.New("view").Parse(tpl))
	t.Execute(os.Stdout, model)

    return nil
}
