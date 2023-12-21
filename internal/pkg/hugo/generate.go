package hugo

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/charlieegan3/food/internal/pkg/markdown"
	"github.com/charlieegan3/food/internal/pkg/recipes"
)

func Generate(siteBasePath string, rs []recipes.Recipe) error {
	for _, r := range rs {
		recipeYAML, err := markdown.RecipeYAML(r)
		if err != nil {
			return fmt.Errorf("failed to generate yaml frontmatter for recipe %q: %w", r.ID, err)
		}

		mdContent := fmt.Sprintf(`
---
%s
---
`, recipeYAML)

		path := fmt.Sprintf("%s/%s/", siteBasePath, r.ID)

		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to mkdir for recipe %q: %w", path, err)
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s/index.md", path), []byte(mdContent), 0644)
		if err != nil {
			return fmt.Errorf("failed to write file for recipe %q: %w", path, err)
		}

		imagePath := fmt.Sprintf("%s/%s/images", siteBasePath, r.ID)
		err = os.MkdirAll(imagePath, os.ModePerm)
		if err != nil {
			if err != nil {
				return fmt.Errorf("failed to make dir for recipe images for recipe %q: %w", r.ID, err)
			}
		}

		for index, i := range r.Images {
			imageData, err := base64.StdEncoding.DecodeString(i)
			if err != nil {
				return fmt.Errorf("failed to encode image for recipe %q: %w", r.ID, err)
			}
			err = ioutil.WriteFile(fmt.Sprintf("%s/%d.jpg", imagePath, index), []byte(imageData), 0644)
			if err != nil {
				return fmt.Errorf("failed to write image file for recipe %q: %w", path, err)
			}
		}
	}
	return nil
}
