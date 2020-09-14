package imgcontent

import (
	"bytes"
	"context"
	"text/template"
	"time"

	"github.com/pkg/errors"
)

// TimeBaseContentPathBuilder ...
// TODO(micheam): add factory func
type TimeBaseContentPathBuilder struct {
	BaseTime *time.Time
	Tmpl     *string
}

// DefaultContentPathTemplate ...
const DefaultContentPathTemplate = `{{.BaseTime.Format "2006/01/02/030405"}}.{{.FileName}}`

// Build ...
func (c TimeBaseContentPathBuilder) Build(
	ctx context.Context, fname Filename) (path ContentPath, err error) {

	basetime := time.Now()
	if c.BaseTime != nil {
		basetime = *c.BaseTime
	}

	model := struct {
		BaseTime time.Time
		FileName string
	}{
		BaseTime: basetime,
		FileName: fname.Value,
	}

	tmpl := DefaultContentPathTemplate
	if c.Tmpl != nil {
		tmpl = *c.Tmpl
	}

	t := template.Must(template.New("content-path-template").Parse(tmpl))
	var buf bytes.Buffer
	if err = t.Execute(&buf, model); err != nil {
		return path, errors.Wrap(err, "failed to build content path")
	}

	path = ContentPath(buf.String())
	return
}
