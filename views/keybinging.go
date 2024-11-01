package views

import (
	"github.com/jroimartin/gocui"
)

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func InitKeybindings(g *gocui.Gui, av *AppView) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	mainKeybindings := []struct {
		key     interface{}
		handler func(*gocui.Gui, *gocui.View) error
	}{
		{'a', func(g *gocui.Gui, v *gocui.View) error { return av.ShowPopup(g) }},
		{'h', func(g *gocui.Gui, v *gocui.View) error { av.UpdateToPrevDay(g); return nil }},
		{'l', func(g *gocui.Gui, v *gocui.View) error { av.UpdateToNextDay(g); return nil }},
		{'j', func(g *gocui.Gui, v *gocui.View) error { av.UpdateToNextTime(g); return nil }},
		{'k', func(g *gocui.Gui, v *gocui.View) error { av.UpdateToPrevTime(g); return nil }},
		{'H', func(g *gocui.Gui, v *gocui.View) error { av.UpdateToPrevWeek(); return nil }},
		{'L', func(g *gocui.Gui, v *gocui.View) error { av.UpdateToNextWeek(); return nil }},
		{gocui.KeyCtrlD, func(g *gocui.Gui, v *gocui.View) error { av.HideSideView(g); return nil }},
		{gocui.KeyCtrlS, func(g *gocui.Gui, v *gocui.View) error { av.ShowSideView(); return nil }},
	}

	for _, viewName := range weekdayNames {
		for _, kb := range mainKeybindings {
			if err := g.SetKeybinding(viewName, kb.key, gocui.ModNone, kb.handler); err != nil {
				return err
			}
		}
	}

	return nil
}
