package cmd

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/oliverhohn/woodhouse/helpers/organize"
	"github.com/spf13/cobra"
)

// CLI flags
var dryRun bool
var operation string

var organizeCmd = &cobra.Command{
	Use:   "organize",
	Short: "Organizes files into folders for every quarter in the year",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if !isValidOperation(operation) {
			log.Fatalf("unknown operation: %s, valid values are: %s", operation, strings.Join(validOperations, ", "))
		}

		inDir := args[0]
		outDir := args[1]

		fileIndex := uint64(0)

		ctx := context.Background()
		ctx = context.WithValue(ctx, "dryRun", dryRun)

		err := filepath.Walk(inDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				f, err := organize.NewLocalFile(path)
				if err != nil {
					return err
				}

				if shouldIgnoreFile(f) {
					return nil
				}

				switch operation {
				case "copy":
					if err := organize.Copy(ctx, f, outDir, fileIndex); err != nil {
						return err
					}
				case "move":
					if err := organize.Move(ctx, f, outDir, fileIndex); err != nil {
						return err
					}
				default:
					log.Fatalf("unknown operation: %s", operation)
				}

				fileIndex++
			}

			return nil
		})
		if err != nil {
			log.Fatalf("unable to read files in input directory: %v", err)
		}
	},
}

var extensionsToIgnore = []string{".DS_Store"}

func shouldIgnoreFile(f *organize.LocalFile) bool {
	fileExt := f.GetExt()
	for _, ext := range extensionsToIgnore {
		if fileExt == ext {
			return true
		}
	}

	return false
}

var validOperations = []string{"copy", "move"}

func isValidOperation(op string) bool {
	for _, o := range validOperations {
		if o == op {
			return true
		}
	}

	return false
}

func init() {
	rootCmd.AddCommand(organizeCmd)

	organizeCmd.Flags().BoolVar(&dryRun, "dry-run", false, "Does not move/copy the files. Just prints what would happen")
	organizeCmd.Flags().StringVar(&operation, "operation", "copy", "Type of operation to use when organizing files")
}
