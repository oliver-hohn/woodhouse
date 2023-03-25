package organize

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/djherbis/times"
)

const Undefined = "UNDEFINED"

type LocalFile struct {
	Path         string
	HasCreatedAt bool
	CreatedAt    time.Time
}

func (f *LocalFile) GetYear() string {
	if !f.HasCreatedAt {
		return Undefined
	}

	return strconv.Itoa(f.CreatedAt.Year())
}

func (f *LocalFile) GetQuarter() string {
	if !f.HasCreatedAt {
		return Undefined
	}

	// Month() starts at 1 for January, and ends at 12 for December.
	// By subtracting by -1 we can ensure that:
	// Jan -> March is 0
	// April -> June is 1
	// July -> September is 2
	// October -> December is 3
	switch quarterIndex := (f.CreatedAt.Month() - 1) / 3; quarterIndex {
	case 0:
		return "00_jan_to_mar"
	case 1:
		return "01_apr_to_jun"
	case 2:
		return "02_jul_to_sep"
	case 3:
		return "03_oct_to_dec"
	default:
		panic(fmt.Errorf("invalid quarter index: %d", quarterIndex))
	}
}

func (f *LocalFile) GetName() string {
	nameWithExt := filepath.Base(f.Path)
	ext := f.GetExt()

	// For example:
	// nameWithExt = "foo.bar"
	// ext == ".bar"
	// name == "foo"
	return nameWithExt[0 : len(nameWithExt)-len(ext)]
}

func (f *LocalFile) GetExt() string {
	return filepath.Ext(f.Path)
}

func NewLocalFile(path string) (*LocalFile, error) {
	stat, err := times.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("unable to create new file: %w", err)
	}

	file := LocalFile{Path: path}

	if stat.HasBirthTime() {
		file.HasCreatedAt = true
		file.CreatedAt = stat.BirthTime()
	}

	return &file, nil
}
