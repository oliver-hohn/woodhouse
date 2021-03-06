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
	"github.com/fatih/color"
)

var inDir = flag.String("in", "", "Path to directory with files to organize")
var outDir = flag.String("out", "", "Path to directory where to place the organized files")
var dryrun = flag.Bool("dryrun", false, "Enable to only print how the files would be organised, instead of actually organizing them in the out directory.")
var opType = flag.String("op_type", "copy", "Type of operation to use when organizing. Allowed values are: copy, or move.")

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

func (f *File) GetName() string {
	nameWithExt := filepath.Base(f.Path)
	ext := f.GetExt()

	// For example:
	// nameWithExt = "foo.bar"
	// ext == ".bar"
	// name == "foo"
	return nameWithExt[0 : len(nameWithExt)-len(ext)]
}

func (f *File) GetExt() string {
	return filepath.Ext(f.Path)
}

func (f *File) GetOutputDir() string {
	return filepath.Join(*outDir, f.GetYear(), f.GetQuarter())
}

func (f *File) GetOutputFilename(suffix string) string {
	if suffix != "" {
		return fmt.Sprintf("%s_%s%s", f.GetName(), suffix, f.GetExt())
	}

	return fmt.Sprintf("%s%s", f.GetName(), f.GetExt())
}

func main() {
	flag.Parse()

	if *inDir == "" || *outDir == "" {
		log.Fatal("missing required \"in\" and \"out\" parameters")
	}

	if *inDir == *outDir {
		log.Fatalf("\"in\" and \"out\" parameters have the same value: %s", *inDir)
	}

	fileIndex := uint64(0)

	err := filepath.Walk(*inDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			f, err := NewFile(path)
			if err != nil {
				return err
			}

			if shouldIgnoreFile(f) {
				return nil
			}

			switch *opType {
			case "copy":
				if err := copy(f, fileIndex); err != nil {
					return err
				}
			case "move":
				if err := move(f, fileIndex); err != nil {
					return err
				}
			default:
				log.Fatalf("unknown op_type: %s", *opType)
			}

			fileIndex++
		}

		return nil
	})
	if err != nil {
		log.Fatalf("unable to read files in input directory: %v", err)
	}
}

var extensionsToIgnore = []string{".DS_Store"}

func shouldIgnoreFile(f *File) bool {
	fileExt := f.GetExt()
	for _, ext := range extensionsToIgnore {
		if fileExt == ext {
			return true
		}
	}

	return false
}

func copy(f *File, fileIndex uint64) error {
	outPrefix := f.GetOutputDir()
	// Suffix filename with an index to avoid clashes for similarly named files
	outFilename := f.GetOutputFilename(strconv.FormatUint(fileIndex, 10))
	out := filepath.Join(outPrefix, outFilename)

	// Prevent overwrites by raising an error if the output file already exists
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		if *dryrun {
			logError("(dryrun) unable to copy %s to: %s, as the output file already exists\n", f.Path, out)
			return nil
		}
		return fmt.Errorf("unable to copy %s to: %s, as the output file already exists", f.Path, out)
	}

	if *dryrun {
		fmt.Printf("(dryrun) Copy %s to %s\n", f.Path, out)
		return nil
	}

	source, err := os.Open(f.Path)
	if err != nil {
		return fmt.Errorf("unable to open input file: %s, due to: %w", f.Path, err)
	}
	defer source.Close()

	// 0777 = Read, Write, Execute for Users, Groups, and Other
	err = os.MkdirAll(outPrefix, 0777)
	if err != nil {
		return fmt.Errorf("unable to create output directory: %s, due to: %w", outPrefix, err)
	}

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

func move(f *File, fileIndex uint64) error {
	outPrefix := f.GetOutputDir()
	// Suffix filename with an index to avoid clashes for similarly named files
	outFilename := f.GetOutputFilename(strconv.FormatUint(fileIndex, 10))
	out := filepath.Join(outPrefix, outFilename)

	// Prevent overwrites by raising an error if the output file already exists
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		if *dryrun {
			logError("(dryrun) unable to move %s to: %s, as the output file already exists\n", f.Path, out)
			return nil
		}
		return fmt.Errorf("unable to move %s to: %s, as the output file already exists", f.Path, out)
	}

	if *dryrun {
		fmt.Printf("(dryrun) Move %s to %s\n", f.Path, out)
		return nil
	}

	// 0777 = Read, Write, Execute for Users, Groups, and Other
	if err := os.MkdirAll(outPrefix, 0777); err != nil {
		return fmt.Errorf("unable to create output directory: %s, due to: %w", outPrefix, err)
	}

	if err := os.Rename(f.Path, out); err != nil {
		return fmt.Errorf("unable to move: %s, to: %s, due to: %w", f.Path, out, err)
	}
	fmt.Printf("Moved %s to %s\n", f.Path, out)

	return nil
}

func logError(msg string, args ...interface{}) {
	color.Red(fmt.Sprintf(msg, args...))
}
