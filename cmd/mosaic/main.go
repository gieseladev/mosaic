package main

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/gieseladev/mosaic"
	"github.com/gieseladev/mosaic/internal/app/mosaicc"
	"gopkg.in/urfave/cli.v2"
	"image"
	"os"
)

func loadImages(c *cli.Context) ([]image.Image, error) {
	if c.NArg() == 0 {
		return nil, cli.Exit("at least one input image required", 1)
	}

	outputPath := c.String("output")
	if outputPath == "" {
		return nil, cli.Exit("output path required", 1)
	}

	return mosaicc.LoadImages(c.Args().Slice())
}

func getComposer(c *cli.Context, count int) (composer mosaic.ComposerInfo, err error) {
	id := c.String("composer")

	switch id {
	case "":
		fallthrough
	case "random":
		composer = mosaic.RecommendComposers(count)[0]

	default:
		var ok bool
		composer, ok = mosaic.GetComposer(id)
		if !ok {
			err = cli.Exit(fmt.Sprintf("no composer %q found", id), 1)
			return
		}
	}

	return
}

func getDimensions(c *cli.Context) (int, int) {
	width := c.Int("width")
	height := c.Int("height")

	if width == 0 && height == 0 {
		return 512, 512
	} else if width == 0 {
		width = height
	} else if height == 0 {
		height = width
	}

	return width, height
}

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:    "output",
			Aliases: []string{"o"},
			Usage:   "path to write output image to",
		},
		&cli.IntFlag{
			Name:  "width",
			Usage: "width of composition",

			DefaultText: "512, or same as height if set",
		},
		&cli.IntFlag{
			Name:  "height",
			Usage: "height of composition",

			DefaultText: "512, or same as width if set",
		},
	}

	app := &cli.App{
		Name:    "mosaic",
		Usage:   "generate image collages",
		Version: "0.0.2",

		Commands: []*cli.Command{
			{
				Name:      "generate",
				Aliases:   []string{"gen"},
				Usage:     "generate a composition",
				ArgsUsage: "<image>...",

				Flags: append([]cli.Flag{
					&cli.StringFlag{
						Name:        "composer",
						Aliases:     []string{"c"},
						Usage:       "use specific composer",
						DefaultText: "random",
					},
				}, flags...),

				Action: func(c *cli.Context) error {
					outputPath := c.String("output")
					if outputPath == "" {
						return cli.Exit("output path required", 1)
					}

					composer, err := getComposer(c, c.NArg())
					if err != nil {
						return err
					}

					dc := gg.NewContext(getDimensions(c))

					images, err := loadImages(c)
					if err != nil {
						return err
					}

					imgCount := composer.RecommendImageCount(len(images))
					err = composer.Compose(dc, images[:imgCount]...)
					if err != nil {
						return err
					}

					return dc.SavePNG(outputPath)
				},
			},
			{
				Name:  "showcase-gen",
				Usage: "generate the example image showing all composers",

				Flags: flags,

				Action: func(c *cli.Context) error {
					outputPath := c.String("output")
					if outputPath == "" {
						return cli.Exit("output path required", 1)
					}

					images, err := loadImages(c)
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
