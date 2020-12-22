package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/djherbis/times"
)

var inDir = flag.String("in", "", "Path to directory with files to organize")
var outDir = flag.String("out", "", "Path to directory where to place the organized files")

type File struct {
	Path         string
	HasCreatedAt bool
	CreatedAt    time.Time
}

func NewFile(path string) (*File, error) {
	stat, err := times.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("unable to create new file: %w", err)
	}

	file := File{Path: path}

	if stat.HasBirthTime() {
		file.HasCreatedAt = true
		file.CreatedAt = stat.BirthTime()
	}

	return &file, nil
}

func (f *File) GetYear() int32 {
	if f.HasCreatedAt {
		return int32(f.CreatedAt.Year())
	}

	return -1
}

func (f *File) GetQuarterIndex() int32 {
	if f.HasCreatedAt {
		return int32(f.CreatedAt.Month() / 3)
	}

	return -1
}

func main() {
	flag.Parse()

	if *inDir == "" || *outDir == "" {
		log.Fatal("missing required \"in\" and \"out\" parameters")
	}

	if *inDir == *outDir {
		log.Fatalf("\"in\" and \"out\" parameters have the same value: %s", *inDir)
	}

	files := []*File{}
	err := filepath.Walk(*inDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			f, err := NewFile(path)
			if err != nil {
				return err
			}

			files = append(files, f)
		}

		return nil
	})
	if err != nil {
		log.Fatalf("unable to read files in input directory: %v", err)
	}
}
