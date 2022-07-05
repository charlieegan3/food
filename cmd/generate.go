package cmd

import (
	"archive/zip"
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type recipe struct {
	ID string `json:"recipe_id"`

	SourceID string  `json:"id"`
	Date     float64 `json:"date"`

	Title        string `json:"title"`
	Ingredients  string `json:"ingredients"`
	Instructions string `json:"instructions"`
	Description  string `json:"text"`
	Notes        string `json:"notes"`

	Images []string `json:"images"`

	Link string `json:"link"`

	Categories []string `json:"categories"`

	Yield string `json:"yield"`

	Favorite   bool `json:"favorite"`
	WantToCook bool `json:"wantToCook"`
}

type contentSection struct {
	Title    string
	Lines    []string
	Numbered bool
}

func (s *contentSection) Markdown() string {
	lines := []string{
		fmt.Sprintf("### %s", s.Title),
	}

	for _, l := range s.Lines {
		if s.Numbered && len(s.Lines) > 1 {
			lines = append(lines, fmt.Sprintf("1. %s", l))
		} else {
			lines = append(lines, fmt.Sprintf("- %s", l))
		}
	}

	return strings.Join(lines, "\n")
}

func (r *recipe) JSON() string {
	cp := *r
	cp.Images = []string{}
	json, err := json.MarshalIndent(cp, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	return string(json)
}

func (r *recipe) IngredientsMarkdown() string {
	var sections []contentSection

	currentSection := contentSection{Title: "", Lines: []string{}}
	for i, e := range strings.Split(r.Ingredients, "\n") {
		if i == 0 {
			if strings.HasPrefix(e, "#") {
				currentSection.Title = strings.TrimPrefix(e, "#")
				continue
			}
		} else {
			if strings.HasPrefix(e, "#") {
				completeSection := currentSection
				sections = append(sections, completeSection)
				currentSection = contentSection{
					Title: strings.TrimPrefix(e, "#"),
					Lines: []string{},
				}
				continue
			}
		}

		currentSection.Lines = append(currentSection.Lines, e)
	}
	sections = append(sections, currentSection)

	markdown := ""
	for _, s := range sections {
		markdown += s.Markdown()
		markdown += "\n"
	}

	return markdown
}

func (r *recipe) InstructionsMarkdown() string {
	var sections []contentSection

	currentSection := contentSection{Title: "", Lines: []string{}, Numbered: true}
	for i, e := range strings.Split(r.Instructions, "\n") {
		if i == 0 {
			if strings.HasPrefix(e, "#") {
				currentSection.Title = strings.TrimPrefix(e, "#")
				continue
			}
		} else {
			if strings.HasPrefix(e, "#") {
				completeSection := currentSection
				sections = append(sections, completeSection)
				currentSection = contentSection{
					Title:    strings.TrimPrefix(e, "#"),
					Lines:    []string{},
					Numbered: true,
				}
				continue
			}
		}

		currentSection.Lines = append(currentSection.Lines, e)
	}
	sections = append(sections, currentSection)

	markdown := ""
	for _, s := range sections {
		markdown += s.Markdown()
		markdown += "\n"
	}

	return markdown
}

func (r *recipe) YAML() string {
	data := map[string]interface{}{
		"title":        r.Title,
		"description":  r.Description,
		"ingredients":  r.IngredientsMarkdown(),
		"instructions": r.InstructionsMarkdown(),
		"images":       len(r.Images),
		"link":         r.Link,
		"categories":   r.Categories,
		"yield":        r.Yield,
		"favorite":     r.Favorite,
		"want_to_cook": r.WantToCook,
		"notes":        strings.ReplaceAll(r.Notes, "\n", "\n\n"),
	}

	str, err := yaml.Marshal(data)
	if err != nil {
		log.Fatal(err)
	}
	return string(str)
}

var generateCmd = &cobra.Command{
	Use: "generate",
	Run: func(cmd *cobra.Command, args []string) {

		data, err := ioutil.ReadFile("Recipes.melarecipes")
		if err != nil {
			log.Fatal(err)
		}

		archive, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
		if err != nil {
			log.Fatal(err)
		}

		var recipes []recipe
		for _, f := range archive.File {
			if !strings.HasSuffix(f.Name, ".melarecipe") {
				continue
			}

			recipeFile, err := f.Open()
			if err != nil {
				log.Fatal(err)
			}

			data, err := ioutil.ReadAll(recipeFile)
			if err != nil {
				log.Fatal(err)
			}

			var r recipe
			err = json.Unmarshal(data, &r)
			if err != nil {
				log.Fatal(err)
			}

			// create a new standard ID from inconsistent generated IDs
			h := sha1.New()
			h.Write([]byte(r.SourceID))
			r.ID = hex.EncodeToString(h.Sum(nil))

			recipes = append(recipes, r)
		}

		for _, r := range recipes {
			mdContent := fmt.Sprintf(`
---
%s
---
`, r.YAML())

			path := fmt.Sprintf("site/content/recipes/%s/", r.ID)

			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			ioutil.WriteFile(fmt.Sprintf("%s/index.md", path), []byte(mdContent), 0644)

			imagePath := fmt.Sprintf("site/content/recipes/%s/images", r.ID)
			err = os.MkdirAll(imagePath, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}

			for index, i := range r.Images {
				imageData, err := base64.StdEncoding.DecodeString(i)
				if err != nil {
					log.Fatal(err)
				}
				ioutil.WriteFile(fmt.Sprintf("%s/%d.jpg", imagePath, index), []byte(imageData), 0644)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
