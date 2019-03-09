package targets

import (
	"path/filepath"

	"github.com/Wing924/prometheus-accesslog-exporter/pattern"

	"github.com/mitchellh/go-homedir"

	"github.com/Wing924/prometheus-accesslog-exporter/config"

	"github.com/pkg/errors"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	readBytes = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "accesslog",
		Name:      "read_bytes_total",
		Help:      "Number of bytes read.",
	}, []string{"path"})

	readLines = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "accesslog",
		Name:      "read_lines_total",
		Help:      "Number of lines read.",
	}, []string{"path"})

	filesActive = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "accesslog",
		Name:      "files_active_total",
		Help:      "Number of active files.",
	})
)

func init() {
	prometheus.Register(readBytes)
	prometheus.Register(readLines)
	prometheus.Register(filesActive)
}

// FileTarget describes a particular set of logs.
type FileTarget struct {
	path   string
	format pattern.Format
}

// NewFileTarget create a new FileTarget.
func NewFileTarget(cfg config.TargetConfig) (*FileTarget, error) {
	path, err := homedir.Expand(cfg.Filepath)
	if err != nil {
		return nil, errors.Wrap(err, "homedir.Expand")
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return nil, errors.Wrap(err, "filepath.Abs")
	}

	return &FileTarget{
		path:   path,
		format: pattern.Parse(cfg.TimeScale, cfg.LogFormat),
	}, nil
}

func (t *FileTarget) OpenFiles() {

}

func (t *FileTarget) listFiles() ([]string, error) {
	matches, err := filepath.Glob(t.path)
	return matches, errors.Wrap(err, "filepath.Glob")
}
