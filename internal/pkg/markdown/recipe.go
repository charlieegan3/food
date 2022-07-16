package markdown

import (
	"fmt"
	"github.com/charlieegan3/food/internal/pkg/recipes"
	"gopkg.in/yaml.v2"
	"strings"
)

func RecipeYAML(r recipes.Recipe) (string, error) {
	data := map[string]interface{}{
		"title":        r.Title,
		"description":  r.Description,
		"ingredients":  ingredientsMarkdown(r),
		"instructions": instructionsMarkdown(r),
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
		return "", fmt.Errorf("failed to generate yaml data: %w", err)
	}
	return string(str), nil
}

func ingredientsMarkdown(r recipes.Recipe) string {
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

func instructionsMarkdown(r recipes.Recipe) string {
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
