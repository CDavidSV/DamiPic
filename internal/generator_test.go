package internal

import (
	"image"
	"image/color"
	"math/rand"
	"sync"
	"testing"
)

type hexColorParserTest struct {
	name      string
	input     string
	wantColor color.RGBA
	wantErr   bool
}

type generatorTest struct {
	name      string
	width     int
	height    int
	text      string
	bgColor   string
	textColor string
	wantErr   bool
}

func TestHexColorParser(t *testing.T) {
	tests := []hexColorParserTest{
		{"Valid Black", "000000", color.RGBA{0, 0, 0, 255}, false},
		{"Valid White", "FFFFFF", color.RGBA{255, 255, 255, 255}, false},
		{"Valid Red", "FF0000", color.RGBA{255, 0, 0, 255}, false},
		{"Invalid Length", "FFF", color.RGBA{}, true},
		{"Invalid Characters", "ZZZZZZ", color.RGBA{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseHexColorString(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseHexColorString(%s) = %v, want error", tt.input, got)
				}
			} else {
				if err != nil {
					t.Errorf("ParseHexColorString(%s) = %v, want %v", tt.input, got, tt.wantColor)
				}
				if got != tt.wantColor {
					t.Errorf("ParseHexColorString(%s) = %v, want %v", tt.input, got, tt.wantColor)
				}
			}
		})
	}
}

func TestGenerator(t *testing.T) {
	tests := []generatorTest{
		{"Default Colors", 200, 100, "Hello", "", "", false},
		{"Custom Colors", 200, 100, "Test", "FF5733", "33FF57", false},
		{"Invalid Background Color", 200, 100, "Test", "ZZZZZZ", "000000", true},
		{"Invalid Text Color", 200, 100, "Test", "000000", "XXXXXX", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := GeneratePlaceholderImg(tt.width, tt.height, tt.text, tt.bgColor, tt.textColor)
			if tt.wantErr {
				if err == nil {
					t.Errorf("expected an error for input (bgColor: %s, textColor: %s), but got none", tt.bgColor, tt.textColor)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if img == nil {
					t.Errorf("expected a valid image but got nil")
				}
				if img.Bounds() != image.Rect(0, 0, tt.width, tt.height) {
					t.Errorf("expected image bounds %v, got %v", image.Rect(0, 0, tt.width, tt.height), img.Bounds())
				}
			}
		})
	}
}

func TestGeneratorConcurrent(t *testing.T) {
	const goroutineCount = 500
	var wg sync.WaitGroup
	wg.Add(goroutineCount)

	errors := make(chan error, goroutineCount)

	for i := 0; i < goroutineCount; i++ {
		go func(i int) {
			defer wg.Done()
			width := rand.Intn(4000-10) + 10
			height := rand.Intn(4000-10) + 10

			img, err := GeneratePlaceholderImg(width, height, "Concurrent Test", "", "")
			if err != nil {
				errors <- err
				return
			}
			if img == nil {
				errors <- err
				return
			}
			if img.Bounds() != image.Rect(0, 0, width, height) {
				errors <- err
				return
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		t.Errorf("Unexpected error in concurrent execution: %v", err)
	}
}
