package view

import (
	"image/color"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type paragraphBlock struct {
	Timestamp string
	Text      string
}

func parseContentBlocks(content string) []paragraphBlock {
	if strings.TrimSpace(content) == "" {
		return nil
	}

	re := regexp.MustCompile(`\[(\d{2}/\d{2} \d{2}:\d{2})\]`)
	indices := re.FindAllStringIndex(content, -1)
	matches := re.FindAllStringSubmatch(content, -1)

	if len(indices) == 0 {
		return []paragraphBlock{{Timestamp: "", Text: strings.TrimSpace(content)}}
	}

	var blocks []paragraphBlock
	for i, idx := range indices {
		ts := matches[i][1]
		textStart := idx[1]
		if textStart < len(content) && content[textStart] == '\n' {
			textStart++
		}

		var textEnd int
		if i+1 < len(indices) {
			textEnd = indices[i+1][0]
		} else {
			textEnd = len(content)
		}

		text := strings.TrimRight(content[textStart:textEnd], "\n")
		if text != "" {
			blocks = append(blocks, paragraphBlock{Timestamp: ts, Text: text})
		}
	}

	return blocks
}

func buildParagraphWidget(block paragraphBlock) fyne.CanvasObject {
	var dateLabel *canvas.Text
	if block.Timestamp != "" {
		dateLabel = canvas.NewText("["+block.Timestamp+"] | Pedro Paulino", color.NRGBA{R: 130, G: 130, B: 130, A: 255})
		dateLabel.TextSize = 11
		dateLabel.TextStyle.Italic = true
	}

	textLabel := widget.NewLabel(block.Text)
	textLabel.Wrapping = fyne.TextWrapWord

	if dateLabel != nil {
		return container.NewVBox(dateLabel, textLabel)
	}
	return container.NewVBox(textLabel)
}
