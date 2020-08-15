package chart

import (
	"io"

	"github.com/golang/freetype/truetype"
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
	"github.com/wcharczuk/go-chart/drawing"
	"log"
	"math"
)

// PNG returns a new png/raster renderer.
func GetCairoRenderer(cr *cairo.Context, da *gtk.DrawingArea) RendererProvider {
	CairoRenderer := func(width, height int) (Renderer, error) {
		crRenderer := cairoRenderer{}
		crRenderer.cr = cr
		crRenderer.da = da
		return &crRenderer, nil
	}
	cr.SelectFontFace("Sans", cairo.FONT_SLANT_NORMAL,
		cairo.FONT_WEIGHT_BOLD)
	return CairoRenderer
}

type cairoRenderer struct {
	cr           *cairo.Context
	da           *gtk.DrawingArea
	fontColor    drawing.Color
	fillColor    drawing.Color
	strokeColor  drawing.Color
	textRotation float64
}

// ResetStyle should reset any style related settings on the renderer.
func (crRdr *cairoRenderer) ResetStyle() {
	crRdr.ClearTextRotation()
}

// GetDPI gets the DPI for the renderer.
func (crRdr *cairoRenderer) GetDPI() float64 {
	return DefaultDPI
}

// SetDPI sets the DPI for the renderer.
func (crRdr *cairoRenderer) SetDPI(dpi float64) {
}

// SetClassName sets the current class name.
func (crRdr *cairoRenderer) SetClassName(className string) {
}

// SetStrokeColor sets the current stroke color.
func (crRdr *cairoRenderer) SetStrokeColor(c drawing.Color) {
	log.Printf("cairo: setStrokeColor %v", c)
	crRdr.strokeColor = c
}

// SetFillColor sets the current fill color.
func (crRdr *cairoRenderer) SetFillColor(c drawing.Color) {
	log.Printf("cairo: setFillColor %v", c)
	crRdr.fillColor = c
	r, g, b, a := float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0, float64(c.A)/255.0
	log.Printf("cairo: setFillColor (%v, %v, %v, %v)", r, g, b, a)
	crRdr.cr.SetSourceRGBA(r, g, b, a)
}

// SetStrokeWidth sets the stroke width.
func (crRdr *cairoRenderer) SetStrokeWidth(width float64) {
	log.Printf("cairo: setStrokeWidth %v", width)
	crRdr.cr.SetLineWidth(width)
}

// SetStrokeDashArray sets the stroke dash array.
func (crRdr *cairoRenderer) SetStrokeDashArray(dashArray []float64) {
	log.Printf("cairo: setStrokeDashArray %v", dashArray)
}

// MoveTo moves the cursor to a given point.
func (crRdr *cairoRenderer) MoveTo(x, y int) {
	log.Printf("cairo: MoveTo %v %v", x, y)
	crRdr.cr.MoveTo(float64(x), float64(y))
}

// LineTo both starts a shape and draws a line to a given point
// from the previous point.
func (crRdr *cairoRenderer) LineTo(x, y int) {
	log.Printf("cairo: LineTo %v %v", x, y)
	crRdr.cr.LineTo(float64(x), float64(y))
}

// QuadCurveTo draws a quad curve.
// cx and cy represent the bezier "control points".
//https://lists.cairographics.org/archives/cairo/2010-April/019691.html
func (crRdr *cairoRenderer) QuadCurveTo(cx, cy, x, y int) {
	log.Printf("cairo: QuadCurveTo %v %v %v %v", cx, cy, x, y)
	x0, y0 := crRdr.cr.GetCurrentPoint()
	crRdr.cr.CurveTo(
		2.0/3.0*float64(cx)+1.0/3.0*float64(x0),
		2.0/3.0*float64(cy)+1.0/3.0*float64(y0),
		2.0/3.0*float64(cx)+1.0/3.0*float64(x),
		2.0/3.0*float64(cx)+1.0/3.0*float64(y),
		float64(cy), float64(y))
}

// ArcTo draws an arc with a given center (cx,cy)
// a given set of radii (rx,ry), a startAngle and delta (in radians).
func (crRdr *cairoRenderer) ArcTo(cx, cy int, rx, ry, startAngle, delta float64) {
	log.Printf("cairo: ArcTo %v %v %v %v %v %v", cx, cy, rx, ry, startAngle, delta)
	endAngle := startAngle + delta
	clockWise := true
	if delta < 0 {
		clockWise = false
	}
	// normalize
	if clockWise {
		for endAngle < startAngle {
			endAngle += math.Pi * 2.0
		}
	} else {
		for startAngle < endAngle {
			startAngle += math.Pi * 2.0
		}
	}

	crRdr.cr.Save()
	crRdr.cr.Scale(rx/rx, ry/rx)
	crRdr.cr.MoveTo(float64(cx), float64(cy))

	if clockWise {
		crRdr.cr.Arc(float64(cx), float64(cy), rx, startAngle, endAngle)
	} else {
		crRdr.cr.ArcNegative(float64(cx), float64(cy), rx, startAngle, endAngle)
	}
	crRdr.cr.Scale(1.0, 1.0)
	crRdr.cr.Restore()
}

// Close finalizes a shape as drawn by LineTo.
func (crRdr *cairoRenderer) Close() {
	log.Printf("cairo: Close()")
	crRdr.cr.ClosePath()
}

// Stroke strokes the path.
func (crRdr *cairoRenderer) Stroke() {
	log.Printf("cairo: Stroke()")
	c := crRdr.strokeColor
	r, g, b, a := float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0, float64(c.A)/255.0
	crRdr.cr.SetSourceRGBA(r, g, b, a)
	crRdr.cr.Stroke()
}

// Fill fills the path, but does not stroke.
func (crRdr *cairoRenderer) Fill() {
	log.Printf("cairo: Fill()")
	c := crRdr.fillColor
	r, g, b, a := float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0, float64(c.A)/255.0
	crRdr.cr.SetSourceRGBA(r, g, b, a)

	crRdr.cr.Fill()
}

// FillStroke fills and strokes a path.
func (crRdr *cairoRenderer) FillStroke() {
	log.Printf("cairo: FillStroke()")

	c := crRdr.fillColor
	r, g, b, a := float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0, float64(c.A)/255.0
	crRdr.cr.SetSourceRGBA(r, g, b, a)
	crRdr.cr.FillPreserve()

	c = crRdr.strokeColor
	r, g, b, a = float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0, float64(c.A)/255.0
	crRdr.cr.SetSourceRGBA(r, g, b, a)
	crRdr.cr.Stroke()
}

// Circle draws a circle at the given coords with a given radius.
func (crRdr *cairoRenderer) Circle(radius float64, x, y int) {
	log.Printf("cairo: Circle %v %v %v", radius, x, y)
	crRdr.cr.MoveTo(float64(x), float64(y))
	crRdr.cr.Arc(float64(x), float64(y), radius, 0, 2.0*math.Pi)
}

// SetFont sets a font for a text field.
func (crRdr *cairoRenderer) SetFont(ft *truetype.Font) {
	//log.Printf("cairo: SetFont %+v", ft)
}

// SetFontColor sets a font's color
func (crRdr *cairoRenderer) SetFontColor(c drawing.Color) {
	log.Printf("cairo: SetFontColor %+v", c)
	crRdr.fontColor = c
}

// SetFontSize sets the font size for a text field.
func (crRdr *cairoRenderer) SetFontSize(size float64) {
	log.Printf("cairo: SetFontSize %+v", size)
	crRdr.cr.SetFontSize(size)
}

// Text draws a text blob.
func (crRdr *cairoRenderer) Text(body string, x, y int) {
	log.Printf("cairo: Text %v %v %v", body, x, y)
	crRdr.cr.Save()
	c := crRdr.fontColor
	r, g, b, a := float64(c.R)/255.0, float64(c.G)/255.0, float64(c.B)/255.0, float64(c.A)/255.0
	crRdr.cr.SetSourceRGBA(r, g, b, a)
	crRdr.cr.Rotate(crRdr.textRotation)
	crRdr.cr.MoveTo(float64(x), float64(y))
	crRdr.cr.ShowText(body)
	crRdr.cr.Restore()
}

// MeasureText measures text.
func (crRdr *cairoRenderer) MeasureText(body string) Box {
	log.Printf("cairo: MeasureText %v", body)
	extent := crRdr.cr.TextExtents(body)
	w, h := extent.Width, extent.Height

	textBox := Box{
		Top:    0,
		Left:   0,
		Right:  int(w),
		Bottom: int(h),
	}

	if crRdr.textRotation == 0 {
		return textBox
	}

	return textBox.Corners().Rotate(RadiansToDegrees(crRdr.textRotation)).Box()
}

// SetTextRotatation sets a rotation for drawing elements.
func (crRdr *cairoRenderer) SetTextRotation(radians float64) {
	log.Printf("cairo: SetTextRotation ", radians)
	crRdr.textRotation = radians
}

// ClearTextRotation clears rotation.
func (crRdr *cairoRenderer) ClearTextRotation() {
	log.Printf("cairo: ClearTextRotation ")
	crRdr.textRotation = 0.0
}

// Save writes the image to the given writer.
func (crRdr *cairoRenderer) Save(w io.Writer) error {
	log.Printf("cairo: Save ")
	return nil
}
