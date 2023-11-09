package openai

import (
	"os"
	"testing"
)

func TestImages(t *testing.T) {
	_apiKey := os.Getenv("OPENAI_API_KEY")
	_org := os.Getenv("OPENAI_ORGANIZATION")
	_verbose := os.Getenv("VERBOSE")

	client := NewClient(_apiKey, _org)
	client.Verbose = _verbose == "true"

	if len(_apiKey) <= 0 || len(_org) <= 0 {
		t.Errorf("environment variables `OPENAI_API_KEY` and `OPENAI_ORGANIZATION` are needed")
	}

	// === CreateImage ===
	if created, err := client.CreateImage("A cute baby sea otter",
		ImageOptions{}.
			SetModel("dall-e-2").
			SetN(2).
			SetSize(ImageSize1024x1024_DallE2)); err != nil {
		t.Errorf("failed to create image: %s", err)
	} else {
		if len(created.Data) <= 0 {
			t.Errorf("there was no returned item")
		}
	}

	// === CreateImageEdit ===
	if image, err := NewFileParamFromFilepath("./sample/pepe.png"); err != nil {
		t.Errorf("failed to open sample image: %s", err)
	} else {
		if edited, err := client.CreateImageEdit(image, "A cute baby sea otter wearing a beret",
			ImageEditOptions{}.
				SetModel("dall-e-2").
				SetN(2).
				SetSize(ImageSize1024x1024_DallE2)); err != nil {
			t.Errorf("failed to create edited image: %s", err)
		} else {
			if len(edited.Data) <= 0 {
				t.Errorf("there was no returned item")
			}
		}
	}

	// === CreateImageVariation ===
	if image, err := NewFileParamFromFilepath("./sample/pepe.png"); err != nil {
		t.Errorf("failed to open sample image: %s", err)
	} else {
		if variation, err := client.CreateImageVariation(image,
			ImageVariationOptions{}.
				SetModel("dall-e-2").
				SetN(2).
				SetSize(ImageSize1024x1024_DallE2)); err != nil {
			t.Errorf("failed to create image variation: %s", err)
		} else {
			if len(variation.Data) <= 0 {
				t.Errorf("there was no returned item")
			}
		}
	}
}
