package targets

import (
	"encoding/csv"
	"io"

	"github.com/prometheus/client_golang/prometheus"
)

type (
	Parser struct {
		reader    *csv.Reader
		readLines prometheus.Counter
	}
)

func (p *Parser) Read() (record []string, err error) {
	record, err = p.reader.Read()
	if err == nil {
		p.readLines.Inc()
	}
	return
}

func NewParser(path string, r io.Reader) *Parser {
	reader := csv.NewReader(r)
	reader.Comma = ' '
	reader.ReuseRecord = true
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1
	return &Parser{
		reader:    reader,
		readLines: readLines.WithLabelValues(path),
	}
}
