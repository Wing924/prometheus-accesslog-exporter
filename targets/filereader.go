package targets

import (
	"os"

	"github.com/Wing924/prometheus-accesslog-exporter/postions"

	"github.com/pkg/errors"

	"github.com/prometheus/client_golang/prometheus"
)

type FileReader struct {
	file      *os.File
	pos       *positions.Positions
	readBytes prometheus.Counter
}

func (r *FileReader) Seek(offset int64, whence int) (int64, error) {
	return r.file.Seek(offset, whence)
}

func NewFileReader(path string) (*FileReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "os.Open")
	}
	filesActive.Inc()
	return &FileReader{
		file:      file,
		readBytes: readBytes.WithLabelValues(path),
	}, nil
}

func (r *FileReader) Read(p []byte) (n int, err error) {
	n, err = r.file.Read(p)
	r.readBytes.Add(float64(n))
	return
}

func (r *FileReader) Close() error {
	filesActive.Dec()
	return r.file.Close()
}
