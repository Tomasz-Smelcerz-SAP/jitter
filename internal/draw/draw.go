package draw

import (
	"fmt"

	"github.com/Tomasz-Smelcerz-SAP/jitter/internal/histogram"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

const (
	graphWidth            float64 = 1000
	verticalMarginTop     float64 = 15
	verticalMarginBottom  float64 = 50
	graphHeight           float64 = 400
	horizontalMarginLeft  float64 = 100
	horizontalMarginRight float64 = 100
	lineThickness         float64 = 1.5
)

func Draw(hist *histogram.Histogram, startLabel, endLabel string, outputFileName string) {
	dc := gg.NewContext(int(graphWidth+horizontalMarginLeft+horizontalMarginRight), int(graphHeight+verticalMarginTop+verticalMarginBottom))

	// Set the background
	dc.SetRGB(0.1, 0.1, 0.1)
	dc.Clear()
	dc.Fill()

	// Draw the histogram
	maxHeight := float64(hist.MaxHeight())

	dc.SetLineWidth(lineThickness)
	dc.SetRGB(float64(0)/255.0, float64(200.0)/255.0, float64(0)/255.0)
	for x := 0; x < hist.BucketCount(); x++ {
		y := (float64(hist.Data()[x]) / maxHeight) * (graphHeight)
		x1 := horizontalMarginLeft + float64(x)
		y1 := graphHeight + verticalMarginTop
		x2 := x1
		y2 := y1 - y
		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	// Draw the marks and labels
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		panic("font!")
	}
	face := truetype.NewFace(font, &truetype.Options{
		Size: 20,
	})
	dc.SetFontFace(face)

	dc.SetRGB(1, 0, 0)
	//// Draw the top horizontal line and label
	text := fmt.Sprintf("%d", hist.MaxHeight())
	topLabelWidth, _ := dc.MeasureString(text)

	dc.DrawStringAnchored(text, 10, 25, 0.0, 0.0)
	dc.DrawLine(10+(topLabelWidth+10), 15, horizontalMarginLeft+graphWidth, 15)

	//// Draw the bottom horizontal line and label
	bottomLabelWidth, _ := dc.MeasureString("0")
	dc.DrawStringAnchored("0", 10, graphHeight+verticalMarginTop+5, 0.0, 0.0)
	dc.DrawLine(10+(bottomLabelWidth+10), graphHeight+verticalMarginTop, horizontalMarginLeft+graphWidth, graphHeight+verticalMarginTop)

	//// Draw the left vertical line and label
	leftLabelWidth, _ := dc.MeasureString(startLabel)
	dc.DrawStringAnchored(startLabel, horizontalMarginLeft-leftLabelWidth/2, graphHeight+verticalMarginTop+30, 0.0, 0.0)
	dc.DrawLine(horizontalMarginLeft, graphHeight+verticalMarginTop, horizontalMarginLeft, graphHeight+verticalMarginTop+10)

	//// Draw the right vertical line and label
	rightLabelWidth, _ := dc.MeasureString(startLabel + "+" + endLabel)
	dc.DrawStringAnchored(startLabel+"+"+endLabel, horizontalMarginLeft+graphWidth-rightLabelWidth/2, graphHeight+verticalMarginTop+30, 0.0, 0.0)
	dc.DrawLine(horizontalMarginLeft+graphWidth, graphHeight+verticalMarginTop, horizontalMarginLeft+graphWidth, graphHeight+verticalMarginTop+10)

	dc.Stroke()
	dc.SavePNG(outputFileName)
}
