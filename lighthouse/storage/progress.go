package storage

import (
	"io"

	"github.com/lighthouse-web3/lighthouse-go-sdk/lighthouse/schema"
)

type progressReader struct {
	r      io.Reader
	n      int64
	total  int64
	onProg schema.ProgressCallback
}

func (p *progressReader) Read(b []byte) (int, error) {
	n, err := p.r.Read(b)
	p.n += int64(n)
	if p.onProg != nil {
		p.onProg(schema.Progress{Uploaded: p.n, Total: p.total})
	}
	return n, err
}

func wrapWithProgress(r io.Reader, total int64, cb schema.ProgressCallback) io.Reader {
	if cb == nil {
		return r
	}
	return &progressReader{r: r, total: total, onProg: cb}
}
