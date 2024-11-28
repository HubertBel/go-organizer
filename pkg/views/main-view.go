package views

import (
	"github.com/HubertBel/lazyorg/internal/calendar"
	"github.com/HubertBel/lazyorg/internal/utils"
	"github.com/jroimartin/gocui"
)

type MainView struct {
	*BaseView

	Calendar *calendar.Calendar
}

func NewMainView(c *calendar.Calendar) *MainView {
	mv := &MainView{
		BaseView: NewBaseView("main"),
		Calendar: c,
	}

	tv := NewTimeView()
	mv.AddChild("time", tv)
	mv.AddChild("week", NewWeekView(c, tv))

	return mv
}

func (mv *MainView) Update(g *gocui.Gui) error {
	v, err := g.SetView(
		mv.Name,
		mv.X,
		mv.Y,
		mv.X+mv.W,
		mv.Y+mv.H,
	)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.FgColor = gocui.AttrBold
	}

	mv.updateChildViewProperties()

	if err = mv.UpdateChildren(g); err != nil {
		return err
	}

	return nil
}

func (mv *MainView) updateChildViewProperties() {
	if view, ok := mv.GetChild("time"); ok {
		if timeView, ok := view.(*TimeView); ok {
			y := utils.TimeToPosition(mv.Calendar.CurrentDay.Date, timeView.Body)
			timeView.SetCursor(y)
			timeView.SetProperties(
				mv.X+1,
				mv.Y+1,
				TimeViewWidth,
				mv.H-2,
			)
		}
	}

	if weekView, ok := mv.GetChild("week"); ok {
		weekView.SetProperties(
			mv.X+TimeViewWidth+1,
			mv.Y,
			mv.W-TimeViewWidth-1,
			mv.H,
		)
	}

	if titleView, ok := mv.GetChild("title"); ok {
		titleView.SetProperties(
			mv.X,
			mv.Y,
			mv.W,
			TitleViewHeight,
		)
	}
}
