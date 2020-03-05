package loaders

import (
	"fmt"
	"io"
	"os"

	"errors"

	"github.com/gosuri/uiprogress"
	"github.com/mitchellh/mapstructure"
	// "github.com/romnnn/mongoimport"
)

type InputType int32

const (
	SingleInput InputType = iota
	MultipleInput
)

// ImportLoader ...
type ImportLoader interface {
	Load() (map[string]interface{}, error)
	Start() error
	Finish() error
	Describe() string
	Create(readers []io.Reader, sanitize bool) (ImportLoader, error)
}

// Loader ...
type Loader struct {
	Files            []string
	SpecificLoader   ImportLoader
	file             *os.File
	read             int64
	total            int64
	reader           io.Reader
	Bar              *uiprogress.Bar
	SkipSanitization bool
	ready            bool
}

// Describe ..
func (l *Loader) Describe() string {
	return l.SpecificLoader.Describe()
}

// GetProgress ..
func (l *Loader) GetProgress() (int64, int64) {
	return l.read, l.total
}

// SetProgressPercent ..
func (l *Loader) SetProgressPercent(percent float32) {
	if l.Bar != nil {
		l.Bar.Set(int(percent * float32(l.Bar.Total)))
	}
}

// UpdateProgress ...
func (l *Loader) UpdateProgress() {
	done, total := l.GetProgress()
	if l.Bar != nil {
		l.Bar.Total = int(total)
		l.Bar.Set(int(done))
	}

}

// Load ...
func (l *Loader) Load() (map[string]interface{}, error) {
	if !l.ready {
		return nil, fmt.Errorf("Attempt to call Load() without calling Start()")
	}
	return l.SpecificLoader.Load()
}

type test struct {
	file string
	read int64
}

func (l *Loader) Write(p []byte) (n int, err error) {
	l.read = l.read + int64(len(p))
	return n, nil
}

// Start ...
func (l *Loader) Start(bar *uiprogress.Bar) error {
	err := l.SpecificLoader.Start()
	if err != nil {
		return err
	}
	l.Bar = bar
	l.ready = true
	return nil
}

// Create ...
func (l *Loader) Create(files []string) (*Loader, error) {
	readers := []io.Reader{} // make([]io.Reader, len(files))

	loader := &Loader{
		Files:            files,
		SkipSanitization: l.SkipSanitization,
		ready:            false,
		total:            l.total,
	}

	for _, f := range files {
		// Open the file
		f, fileErr := openFile(f)
		if fileErr != nil {
			return nil, fileErr
		}

		// Create the reader
		stats, statErr := f.Stat()
		if statErr == nil {
			l.total = stats.Size()
		}

		readers = append(readers, io.TeeReader(f, loader))
	}

	loader.total = l.total

	spec, err := l.SpecificLoader.Create(readers, l.SkipSanitization)
	if err != nil {
		return nil, err
	}
	loader.SpecificLoader = spec
	return loader, nil
}

// Finish ...
func (l *Loader) Finish() error {
	err := l.SpecificLoader.Finish()
	l.file.Close()
	return err
}

func openFile(file string) (*os.File, error) {
	if file == "" {
		return nil, errors.New("Got invalid empty file path")
	}
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (l Loader) createStruct(values map[string]interface{}, result interface{}) error {
	return mapstructure.Decode(values, result)
}
