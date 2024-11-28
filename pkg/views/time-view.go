package views

import (
	"fmt"

	"github.com/HubertBel/lazyorg/internal/utils"
	"github.com/jroimartin/gocui"
)

type TimeView struct {
	*BaseView
	Body   string
	Cursor int
}

func NewTimeView() *TimeView {
	tv := &TimeView{
		BaseView: NewBaseView("time"),
		Cursor:   0,
	}

	return tv
}

func (tv *TimeView) Update(g *gocui.Gui) error {
	v, err := g.SetView(
		tv.Name,
		tv.X,
		tv.Y,
		tv.X+tv.W,
		tv.Y+tv.H,
	)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.FgColor = gocui.ColorGreen
	}

	tv.updateBody(v)

	return nil
}

func (tv *TimeView) updateBody(v *gocui.View) {
	initialTime := 12 - tv.H/4
	tv.Body = ""

	for i := range tv.H {
		var time string

		if i%2 == 0 {
			hour := utils.FormatHour(initialTime, 0)
			time = fmt.Sprintf(" %s - \n", hour)
		} else {
			hour := utils.FormatHour(initialTime, 30)
			time = fmt.Sprintf(" %s \n", hour)
			initialTime++
		}

		if i == tv.Cursor {
            runes := []rune(time)
            runes[0] = '>'
            time = string(runes)
		}

		tv.Body += time
	}

	v.Clear()
	fmt.Fprintln(v, tv.Body)
}

func (tv *TimeView) SetCursor(y int) {
	tv.Cursor = y
}
