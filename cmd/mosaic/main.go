package mosaic

import (
	"gopkg.in/urfave/cli.v2"
	"log"
	"os"
)

func action(c *cli.Context) error {
	return nil
}

func main() {
	app := &cli.App{
		Name:    "mosaic",
		Usage:   "generate image collages",
		Version: "0.0.1",
		Action:  action,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
