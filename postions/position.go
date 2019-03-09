package positions

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"gopkg.in/yaml.v2"
)

const positionFileMode = 0700

// Config describes where to get postition information from.
type Config struct {
	SyncPeriod    time.Duration `yaml:"sync_period"`
	PositionsFile string        `yaml:"filename"`
}

// Positions tracks how far through each file we've read.
type Positions struct {
	log       *logrus.Logger
	cfg       Config
	mtx       sync.Mutex
	positions map[string]int64
	quit      chan struct{}
	done      chan struct{}
}

// File format for the positions data.
type File struct {
	Positions map[string]int64 `yaml:"positions"`
}

// New makes a new Positions.
func New(log *logrus.Logger, cfg Config) (*Positions, error) {
	positions, err := readPositionsFile(cfg.PositionsFile)
	if err != nil {
		return nil, err
	}

	p := &Positions{
		log:       log,
		cfg:       cfg,
		positions: positions,
		quit:      make(chan struct{}),
		done:      make(chan struct{}),
	}

	go p.run()
	return p, nil
}

// Stop the Position tracker.
func (p *Positions) Stop() {
	close(p.quit)
	<-p.done
}

// Put records (asynchronously) how far we've read through a file.
func (p *Positions) Put(path string, pos int64) {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	p.positions[path] = pos
}

// Get returns how far we've read through a file.
func (p *Positions) Get(path string) int64 {
	p.mtx.Lock()
	defer p.mtx.Unlock()
	return p.positions[path]
}

// Remove removes the position tracking for a filepath
func (p *Positions) Remove(path string) {
	p.mtx.Lock()
	delete(p.positions, path)
	p.mtx.Unlock()
}

// SyncPeriod returns how often the positions file gets resynced
func (p *Positions) SyncPeriod() time.Duration {
	return p.cfg.SyncPeriod
}

func (p *Positions) run() {
	defer func() {
		p.save()
		p.log.Debug("positions saved")
		close(p.done)
	}()

	ticker := time.NewTicker(p.cfg.SyncPeriod)
	for {
		select {
		case <-p.quit:
			return
		case <-ticker.C:
			p.save()
		}
	}
}

func (p *Positions) save() {
	p.mtx.Lock()
	positions := make(map[string]int64, len(p.positions))
	for k, v := range p.positions {
		positions[k] = v
	}
	p.mtx.Unlock()

	if err := writePositionFile(p.cfg.PositionsFile, positions); err != nil {
		p.log.WithError(err).Error("error writing positions file")
	}
}

func readPositionsFile(filename string) (map[string]int64, error) {
	buf, err := ioutil.ReadFile(filepath.Clean(filename))
	if err != nil {
		if os.IsNotExist(err) {
			return map[string]int64{}, nil
		}
		return nil, err
	}

	var p File
	if err := yaml.UnmarshalStrict(buf, &p); err != nil {
		return nil, err
	}

	return p.Positions, nil
}

func writePositionFile(filename string, positions map[string]int64) error {
	buf, err := yaml.Marshal(File{
		Positions: positions,
	})
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Clean(filename), buf, os.FileMode(positionFileMode))
}
