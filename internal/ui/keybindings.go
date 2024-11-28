package ui

import (
	"github.com/HubertBel/lazyorg/pkg/views"
	"github.com/jroimartin/gocui"
)

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func InitKeybindings(g *gocui.Gui, av *views.AppView) error {
	g.InputEsc = true

	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '?', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return av.ShowKeybinds(g)
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return av.HandleEscape(g, v)
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return av.HandleEnter(g, v)
	}); err != nil {
		return err
	}

	if err := initMainKeybindings(g, av); err != nil {
		return err
	}
	if err := initNotepadKeybindings(g, av); err != nil {
		return err
	}

	return nil
}

func initMainKeybindings(g *gocui.Gui, av *views.AppView) error {
	mainKeybindings := []struct {
		key     interface{}
		handler func(*gocui.Gui, *gocui.View) error
	}{
		{'a', func(g *gocui.Gui, v *gocui.View) error { return av.ShowPopup(g) }},
		{'h', func(g *gocui.Gui, v *gocui.View) error { av.UpdateToPrevDay(g); return nil }},
		{'l', func(g *gocui.Gui, v *gocui.View) error { av.UpdateToNextDay(g); return nil }},
		{'j', func(g *gocui.Gui, v *gocui.View) error { av.UpdateToNextTime(g); return nil }},
		{'k', func(g *gocui.Gui, v *gocui.View) error { av.UpdateToPrevTime(g); return nil }},
		{gocui.KeyArrowLeft, func(g *gocui.Gui, v *gocui.View) error { av.UpdateToPrevDay(g); return nil }},
		{gocui.KeyArrowRight, func(g *gocui.Gui, v *gocui.View) error { av.UpdateToNextDay(g); return nil }},
		{gocui.KeyArrowDown, func(g *gocui.Gui, v *gocui.View) error { av.UpdateToNextTime(g); return nil }},
		{gocui.KeyArrowUp, func(g *gocui.Gui, v *gocui.View) error { av.UpdateToPrevTime(g); return nil }},
        {'t', func(g *gocui.Gui, v *gocui.View) error { av.JumpToToday(); return nil }},
		{'H', func(g *gocui.Gui, v *gocui.View) error { av.UpdateToPrevWeek(); return nil }},
		{'L', func(g *gocui.Gui, v *gocui.View) error { av.UpdateToNextWeek(); return nil }},
		{'d', func(g *gocui.Gui, v *gocui.View) error { av.DeleteEvent(g); return nil }},
		{'D', func(g *gocui.Gui, v *gocui.View) error { av.DeleteEvents(g); return nil }},
		{gocui.KeyCtrlN, func(g *gocui.Gui, v *gocui.View) error { return av.ChangeToNotepadView(g) }},
		{gocui.KeyCtrlS, func(g *gocui.Gui, v *gocui.View) error { return av.ShowOrHideSideView(g) }},
	}

	for _, viewName := range views.WeekdayNames {
		for _, kb := range mainKeybindings {
			if err := g.SetKeybinding(viewName, kb.key, gocui.ModNone, kb.handler); err != nil {
				return err
			}
		}
	}

	return nil
}

func initNotepadKeybindings(g *gocui.Gui, av *views.AppView) error {
	notepadKeybindings := []struct {
		key     interface{}
		handler func(*gocui.Gui, *gocui.View) error
	}{
		{gocui.KeyCtrlR, func(g *gocui.Gui, v *gocui.View) error { return av.ClearNotepadContent(g) }},
		{gocui.KeyCtrlN, func(g *gocui.Gui, v *gocui.View) error { return av.ReturnToMainView(g) }},
	}
	for _, kb := range notepadKeybindings {
		if err := g.SetKeybinding("notepad", kb.key, gocui.ModNone, kb.handler); err != nil {
			return err
		}
	}

	return nil
}
