package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

const Undefined = "UNDEFINED"

func (f *File) GetYear() string {
	if !f.HasCreatedAt {
		return Undefined
	}

	return strconv.Itoa(f.CreatedAt.Year())
}

func (f *File) GetQuarter() string {
	if !f.HasCreatedAt {
		return Undefined
	}

	switch quarterIndex := f.CreatedAt.Month() / 3; quarterIndex {
	case 0:
		return "jan_to_march"
	case 1:
		return "april_to_june"
	case 2:
		return "july_to_september"
	case 3:
		return "october_to_december"
	default:
		panic(fmt.Errorf("invalid quarter index: %d", quarterIndex))
	}
}

func main() {
	flag.Parse()

	if *inDir == "" || *outDir == "" {
		log.Fatal("missing required \"in\" and \"out\" parameters")
	}

	if *inDir == *outDir {
		log.Fatalf("\"in\" and \"out\" parameters have the same value: %s", *inDir)
	}

	err := filepath.Walk(*inDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			f, err := NewFile(path)
			if err != nil {
				return err
			}

			err = copy(f)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Fatalf("unable to read files in input directory: %v", err)
	}
}

func copy(f *File) error {
	source, err := os.Open(f.Path)
	if err != nil {
		return fmt.Errorf("unable to open input file: %s, due to: %w", f.Path, err)
	}
	defer source.Close()

	outPrefix := filepath.Join(*outDir, f.GetYear(), f.GetQuarter())
	// 0777 = Read, Write, Execute for Users, Groups, and Other
	err = os.MkdirAll(outPrefix, 0777)
	if err != nil {
		return fmt.Errorf("unable to create output directory: %s, due to: %w", outPrefix, err)
	}

	out := filepath.Join(outPrefix, filepath.Base(f.Path))
	dest, err := os.Create(out)
	if err != nil {
		return fmt.Errorf("unable to create output file: %s, due to: %w", out, err)
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	if err != nil {
		return fmt.Errorf("unable to copy: %s, to: %s, due to: %w", f.Path, out, err)
	}
	fmt.Printf("Copied %s to %s\n", f.Path, out)

	return nil
}
