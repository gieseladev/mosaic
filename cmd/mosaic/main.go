package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/gieseladev/mosaic/internal/app/mosaicc"
	"gopkg.in/urfave/cli.v2"
	"os"
)

func main() {
	app := &cli.App{
		Name:                  "mosaic",
		Usage:                 "generate image collages",
		Version:               "0.0.1",
		EnableShellCompletion: true,
		Action: func(c *cli.Context) error {
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "showcase-gen",
				Usage: "generate the example image showing all composers",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "output",
						Aliases: []string{"o"},
						Usage:   "path to write output image to",
					},
				},
				Action: func(c *cli.Context) error {
					if c.NArg() == 0 {
						return cli.Exit("at least one input image required", 1)
					}

					outputPath := c.String("output")
					if outputPath == "" {
						return cli.Exit("output path required", 1)
					}

					images, err := mosaicc.LoadImages(c.Args().Slice())
					if err != nil {
						return err
					}

					generated, err := mosaicc.GenerateComposerShowcase(images)
					if err != nil {
						return err
					}

					return gg.SavePNG(outputPath, generated)
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
