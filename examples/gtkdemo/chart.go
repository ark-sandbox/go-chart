package main

import (
	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
	"log"
	_ "math"
	"time"

	gchart "github.com/ark-sandbox/go-chart"
)

type GtkGoChart struct {
	width, height int
	GtkCanvas     *gtk.DrawingArea
	gochart       gchart.GoChart
}

func GtkGoChartNew(width, height int, chart gchart.GoChart) *GtkGoChart {
	gtkchart := GtkGoChart{}
	gtkchart.width = width
	gtkchart.height = height
	gtkchart.GtkCanvas, _ = gtk.DrawingAreaNew()
	gtkchart.gochart = chart

	//Call Draw method by connecting to draw event.
	gtkchart.GtkCanvas.Connect("draw", gtkchart.Draw)
	gtkchart.GtkCanvas.SetSizeRequest(width, height)
	return &gtkchart
}

func (chart *GtkGoChart) Draw(da *gtk.DrawingArea, cr *cairo.Context) {
	log.Println("Draw is called")

	cairoChart := gchart.GetCairoRenderer(cr, da)
	start := time.Now()
	chart.gochart.Render(cairoChart, nil)
	log.Println("***** time it took to render:", time.Since(start))
}

var logger gchart.Logger

func BarChart() *gtk.DrawingArea {
	chart := gchart.BarChart{
		Title: "Test Bar Chart",
		Background: gchart.Style{
			Padding: gchart.Box{
				Top: 40,
			},
		},
		Height:   512,
		BarWidth: 60,
		Bars: []gchart.Value{
			{Value: 5.25, Label: "Blue"},
			{Value: 4.88, Label: "Green"},
			{Value: 4.74, Label: "Gray"},
			{Value: 3.22, Label: "Orange"},
			{Value: 3, Label: "Test"},
			{Value: 2.27, Label: "??"},
			{Value: 1, Label: "!!"},
		},
	}
	gtkchart := GtkGoChartNew(1024, 500, chart)
	return gtkchart.GtkCanvas
}

func Annotate() *gtk.DrawingArea {
	chart := gchart.Chart{
		Series: []gchart.Series{
			gchart.ContinuousSeries{
				XValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
				YValues: []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			},
			gchart.AnnotationSeries{
				Annotations: []gchart.Value2{
					{XValue: 1.0, YValue: 1.0, Label: "One"},
					{XValue: 2.0, YValue: 2.0, Label: "Two"},
					{XValue: 3.0, YValue: 3.0, Label: "Three"},
					{XValue: 4.0, YValue: 4.0, Label: "Four"},
					{XValue: 5.0, YValue: 5.0, Label: "Five"},
				},
			},
		},
	}
	gtkchart := GtkGoChartNew(1024, 500, chart)
	return gtkchart.GtkCanvas
}

func Donut() *gtk.DrawingArea {
	chart := gchart.DonutChart{
		Width:  512,
		Height: 512,
		Values: []gchart.Value{
			{Value: 5, Label: "Blue"},
			{Value: 5, Label: "Green"},
			{Value: 4, Label: "Gray"},
			{Value: 4, Label: "Orange"},
			{Value: 3, Label: "Deep Blue"},
			{Value: 3, Label: "test"},
		},
	}
	gtkchart := GtkGoChartNew(512, 512, chart)
	return gtkchart.GtkCanvas
}
func PieChart() *gtk.DrawingArea {
	chart := gchart.PieChart{
		Width:  512,
		Height: 512,
		Values: []gchart.Value{
			{Value: 5, Label: "Blue"},
			{Value: 5, Label: "Green"},
			{Value: 4, Label: "Gray"},
			{Value: 4, Label: "Orange"},
			{Value: 3, Label: "Deep Blue"},
			{Value: 3, Label: "??"},
			{Value: 1, Label: "!!"},
		},
	}
	gtkchart := GtkGoChartNew(512, 512, chart)
	return gtkchart.GtkCanvas
}
func main() {
	gtk.Init(nil)

	logger = gchart.NewLogger()
	// gui boilerplate
	win, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)

	scrollWin, _ := gtk.ScrolledWindowNew(nil, nil)
	scrollWin.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_AUTOMATIC)

	grid, err := gtk.GridNew()
	if err != nil {
		log.Fatal("Unable to create grid:", err)
	}
	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	//piechart := PieChart()
	grid.Add(PieChart())
	grid.Add(PieChart())

	topLabel, err := gtk.LabelNew("Grid End")
	if err != nil {
		log.Fatal("Unable to create label:", err)
	}
	_ = topLabel
	//grid.Add(topLabel)
	//win.Add(grid)
	win.Add(scrollWin)

	flowbox, err := gtk.FlowBoxNew()
	flowbox.SetVAlign(gtk.ALIGN_START)
	flowbox.SetSelectionMode(gtk.SELECTION_NONE)
	flowbox.Add(BarChart())
	flowbox.Add(Annotate())
	flowbox.Add(Donut())
	flowbox.Add(PieChart())

	scrollWin.Add(flowbox)

	win.SetTitle("Arrow keys")
	win.Connect("destroy", gtk.MainQuit)
	win.ShowAll()
	win.SetSizeRequest(300, 600)

	gtk.Main()
}
