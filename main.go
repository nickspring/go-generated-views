package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/gommon/color"
	"github.com/nickspring/go-generated-views/generator"
	"github.com/urfave/cli/v2"
)

var (
	version string
	commit  string
	date    string
	builtBy string
)

type Arguments struct {
	FileNames    cli.StringSlice
	BuildTags    cli.StringSlice
	OutputSuffix string
}

func main() {
	var arguments Arguments

	clr := color.New()
	out := func(format string, args ...interface{}) {
		_, _ = fmt.Fprintf(clr.Output(), format, args...)
	}

	app := &cli.App{
		Name:            "go-generated-views",
		Usage:           "A struct views generator for go",
		HideHelpCommand: true,
		Version:         version,
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:        "file",
				Aliases:     []string{"f"},
				EnvVars:     []string{"GOFILE"},
				Usage:       "The file(s) to generate views. Use more than one flag for more files.",
				Required:    true,
				Destination: &arguments.FileNames,
			},
			&cli.StringSliceFlag{
				Name:        "buildtag",
				Aliases:     []string{"b"},
				Usage:       "Adds build tags to a generated view file.",
				Destination: &arguments.BuildTags,
			},
			&cli.StringFlag{
				Name:        "output-suffix",
				Usage:       "Changes the default filename suffix of _view to something else.  `.go` will be appended to the end of the string no matter what, so that `_test.go` cases can be accommodated ",
				Destination: &arguments.OutputSuffix,
			},
		},
		Action: func(ctx *cli.Context) error {
			for _, fileOption := range arguments.FileNames.Value() {

				g := generator.NewGenerator()
				g.Version = version
				g.Revision = commit
				g.BuildDate = date
				g.BuiltBy = builtBy

				g.WithBuildTags(arguments.BuildTags.Value()...)

				var filenames []string
				if fn, err := globFilenames(fileOption); err != nil {
					return err
				} else {
					filenames = fn
				}

				outputSuffix := `_view`
				if arguments.OutputSuffix != "" {
					outputSuffix = arguments.OutputSuffix
				}

				for _, fileName := range filenames {
					originalName := fileName

					out("go-generated-views started. file: %s\n", color.Cyan(originalName))
					fileName, _ = filepath.Abs(fileName)

					outFilePath := fmt.Sprintf("%s%s.go", strings.TrimSuffix(fileName, filepath.Ext(fileName)), outputSuffix)
					if strings.HasSuffix(fileName, "_test.go") {
						outFilePath = strings.Replace(outFilePath, "_test"+outputSuffix+".go", outputSuffix+"_test.go", 1)
					}

					// Parse the file given in arguments
					raw, err := g.GenerateFromFile(fileName)
					if err != nil {
						return fmt.Errorf("failed generating views\nInputFile=%s\nError=%s", color.Cyan(fileName), color.RedBg(err))
					}

					// Nothing was generated, ignore the output and don't create a file.
					if len(raw) < 1 {
						out(color.Yellow("go-generated-views ignored. file: %s\n"), color.Cyan(originalName))
						continue
					}

					mode := 0o644
					err = os.WriteFile(outFilePath, raw, os.FileMode(mode))
					if err != nil {
						return fmt.Errorf("failed writing to file %s: %s", color.Cyan(outFilePath), color.Red(err))
					}
					out("go-generated-views finished. file: %s\n", color.Cyan(originalName))
				}
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

// globFilenames gets a list of filenames matching the provided filename.
// In order to maintain existing capabilities, only glob when a * is in the path.
// Leave execution on par with old method in case there are bad patterns in use that somehow
// work without the Glob method.
func globFilenames(filename string) ([]string, error) {
	if strings.Contains(filename, "*") {
		matches, err := filepath.Glob(filename)
		if err != nil {
			return []string{}, fmt.Errorf("failed parsing glob filepath\nInputFile=%s\nError=%s", color.Cyan(filename), color.RedBg(err))
		}
		return matches, nil
	} else {
		return []string{filename}, nil
	}
}
