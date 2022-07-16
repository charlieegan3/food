package cmd

import (
	"github.com/charlieegan3/food/internal/pkg/hugo"
	"github.com/charlieegan3/food/internal/pkg/recipes"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
)

var generateCmd = &cobra.Command{
	Use: "generate",
	Run: func(cmd *cobra.Command, args []string) {
		siteBasePath := "site/content/recipes"
		sourceFilePath := "Recipes.melarecipes"

		data, err := ioutil.ReadFile(sourceFilePath)
		if err != nil {
			log.Fatal(err)
		}

		rs, err := recipes.Parse(data)
		if err != nil {
			log.Fatal(err)
		}

		err = hugo.Generate(siteBasePath, rs)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
