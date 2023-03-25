package organize

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/fatih/color"
)

func Copy(ctx context.Context, f *LocalFile, outDir string, fileIndex uint64) error {
	dryRun := ctx.Value("dryRun").(bool)

	outPrefix := getOutputDir(f, outDir)
	// Suffix filename with an index to avoid clashes for similarly named files
	outFilename := getOutputFilename(f, strconv.FormatUint(fileIndex, 10))
	out := filepath.Join(outPrefix, outFilename)

	// Prevent overwrites by raising an error if the output file already exists
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		if dryRun {
			logError("(dryrun) unable to copy %s to: %s, as the output file already exists\n", f.Path, out)
			return nil
		}
		return fmt.Errorf("unable to copy %s to: %s, as the output file already exists", f.Path, out)
	}

	if dryRun {
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

func Move(ctx context.Context, f *LocalFile, outDir string, fileIndex uint64) error {
	dryRun := ctx.Value("dryRun").(bool)

	outPrefix := getOutputDir(f, outDir)
	// Suffix filename with an index to avoid clashes for similarly named files
	outFilename := getOutputFilename(f, strconv.FormatUint(fileIndex, 10))
	out := filepath.Join(outPrefix, outFilename)

	// Prevent overwrites by raising an error if the output file already exists
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		if dryRun {
			logError("(dryrun) unable to move %s to: %s, as the output file already exists\n", f.Path, out)
			return nil
		}
		return fmt.Errorf("unable to move %s to: %s, as the output file already exists", f.Path, out)
	}

	if dryRun {
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

func getOutputDir(f *LocalFile, outDir string) string {
	return filepath.Join(outDir, f.GetYear(), f.GetQuarter())
}

func getOutputFilename(f *LocalFile, suffix string) string {
	if suffix != "" {
		return fmt.Sprintf("%s_%s%s", f.GetName(), suffix, f.GetExt())
	}

	return fmt.Sprintf("%s%s", f.GetName(), f.GetExt())
}
