package views

import (
	"fmt"

	"github.com/HubertBel/go-organizer/cmd/types"
	"github.com/jroimartin/gocui"
)


type HoverView struct {
	*BaseView

	Calendar *types.Calendar
    CurrentView View
}

func NewHoverView(c *types.Calendar) *HoverView {
	hv := &HoverView{
        BaseView: NewBaseView("hover"),
        Calendar: c,
	}

	return hv
}

func (hv *HoverView) Update(g *gocui.Gui) error {
	v, err := g.SetView(
		hv.Name,
		hv.X,
		hv.Y,
		hv.X+hv.W,
		hv.Y+hv.H,
	)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
        v.Wrap = true
	}

    hv.updateTitle(v)
    hv.updateBody(v)

	return nil
}

func (hv *HoverView) updateTitle(v *gocui.View) {
    if view, ok := hv.CurrentView.(*DayView); ok {
        v.Title = view.Day.FormatTimeAndHour()
    } else if view, ok := hv.CurrentView.(*EventView); ok {
        v.Title = view.Event.FormatTimeAndName()
    }
}

func (hv *HoverView) updateBody(v *gocui.View) {
    v.Clear()
    if view, ok := hv.CurrentView.(*DayView); ok {
        v.FgColor = gocui.AttrBold | gocui.ColorYellow
        fmt.Fprintln(v)
        fmt.Fprintln(v, "\nEvents :")
        fmt.Fprintln(v, "---------")
        fmt.Fprintln(v, view.Day.FormatBody())
    } else if view, ok := hv.CurrentView.(*EventView); ok {
        fmt.Fprintln(v, view.Y)
        fmt.Fprintln(v, view.Y+view.H)
    }
}
