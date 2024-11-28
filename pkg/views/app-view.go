package views

import (
	"time"

	"github.com/HubertBel/lazyorg/internal/calendar"
	"github.com/HubertBel/lazyorg/internal/database"
	"github.com/HubertBel/lazyorg/internal/utils"
	"github.com/jroimartin/gocui"
	"github.com/nsf/termbox-go"
)

var (
	MainViewWidthRatio = 0.8
	SideViewWidthRatio = 0.2
)

type AppView struct {
	*BaseView

	Database *database.Database
	Calendar *calendar.Calendar
}

func NewAppView(g *gocui.Gui, db *database.Database) *AppView {
	now := time.Now()
	t := time.Date(now.Year(), now.Month(), now.Day(), 12, 0, 0, 0, now.Location())

	c := calendar.NewCalendar(calendar.NewDay(t))

	av := &AppView{
		BaseView: NewBaseView("app"),
		Database: db,
		Calendar: c,
	}

	av.AddChild("title", NewTitleView(c))
	av.AddChild("popup", NewEvenPopup(g, c, db))
	av.AddChild("main", NewMainView(c))
	av.AddChild("side", NewSideView(c, db))
	av.AddChild("keybinds", NewKeybindsView())

	return av
}

func (av *AppView) Layout(g *gocui.Gui) error {
	return av.Update(g)
}

func (av *AppView) Update(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	av.SetProperties(0, 1, maxX-1, maxY-1)

	v, err := g.SetView(
		av.Name,
		av.X,
		av.Y,
		av.X+av.W,
		av.Y+av.H,
	)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
	}

	if err = av.updateEventsFromDatabase(); err != nil {
		return err
	}

	av.updateChildViewProperties()

	if err = av.UpdateChildren(g); err != nil {
		return err
	}

	if err = av.updateCurrentView(g); err != nil {
		return err
	}

	return nil
}

func (av *AppView) updateEventsFromDatabase() error {
	for _, v := range av.Calendar.CurrentWeek.Days {
		clear(v.Events)

		var err error
		events, err := av.Database.GetEventsByDate(v.Date)
		if err != nil {
			return err
		}

		v.Events = events
		v.SortEventsByTime()
	}

	return nil
}

func (av *AppView) ShowOrHideSideView(g *gocui.Gui) error {
	if sideView, ok := av.GetChild("side"); ok {
		if err := sideView.ClearChildren(g); err != nil {
			return err
		}
		SideViewWidthRatio = 0.0
		MainViewWidthRatio = 1.0

		av.children.Delete("side")
		return g.DeleteView("side")
	}

	SideViewWidthRatio = 0.2
	MainViewWidthRatio = 0.8

	av.AddChild("side", NewSideView(av.Calendar, av.Database))

	return nil
}

func (av *AppView) JumpToToday() {
	av.Calendar.JumpToToday()
}

func (av *AppView) UpdateToNextWeek() {
	av.Calendar.UpdateToNextWeek()
}

func (av *AppView) UpdateToPrevWeek() {
	av.Calendar.UpdateToPrevWeek()
}

func (av *AppView) UpdateToNextDay(g *gocui.Gui) {
	av.Calendar.UpdateToNextDay()
	av.updateCurrentView(g)
}

func (av *AppView) UpdateToPrevDay(g *gocui.Gui) {
	av.Calendar.UpdateToPrevDay()
	av.updateCurrentView(g)
}

func (av *AppView) UpdateToNextTime(g *gocui.Gui) {
	_, height := g.CurrentView().Size()
	if _, y := g.CurrentView().Cursor(); y < height-1 {
		av.Calendar.UpdateToNextTime()
	}
}

func (av *AppView) UpdateToPrevTime(g *gocui.Gui) {
	if _, y := g.CurrentView().Cursor(); y > 0 {
		av.Calendar.UpdateToPrevTime()
	}
}

func (av *AppView) ChangeToNotepadView(g *gocui.Gui) error {
	_, err := g.SetCurrentView("notepad")
	if err != nil {
		return err
	}
	if view, ok := av.FindChildView("notepad"); ok {
		if notepadView, ok := view.(*NotepadView); ok {
            notepadView.IsActive = true
		}
	}

	return nil
}

func (av *AppView) ClearNotepadContent(g *gocui.Gui) error {
	if view, ok := av.FindChildView("notepad"); ok {
		if notepadView, ok := view.(*NotepadView); ok {
			return notepadView.ClearContent(g)
		}
	}

	return nil
}

func (av *AppView) SaveNotepadContent(g *gocui.Gui) error {
	if view, ok := av.FindChildView("notepad"); ok {
		if notepadView, ok := view.(*NotepadView); ok {
			return notepadView.SaveContent(g)
		}
	}

	return nil
}

func (av *AppView) ReturnToMainView(g *gocui.Gui) error {
	if err := av.SaveNotepadContent(g); err != nil {
		return err
	}
	if view, ok := av.FindChildView("notepad"); ok {
		if notepadView, ok := view.(*NotepadView); ok {
            notepadView.IsActive = false
		}
	}

	viewName := WeekdayNames[av.Calendar.CurrentDay.Date.Weekday()]
	g.SetCurrentView(viewName)
	return av.updateCurrentView(g)
}

func (av *AppView) DeleteEvent(g *gocui.Gui) {
	_, y := g.CurrentView().Cursor()

	if view, ok := av.FindChildView(WeekdayNames[av.Calendar.CurrentDay.Date.Weekday()]); ok {
		if dayView, ok := view.(*DayView); ok {
			if view, ok := dayView.IsOnEvent(y); ok {
				if eventView, ok := view.(*EventView); ok {
					av.Database.DeleteEvent(eventView.Event.Id)
				}
			}
		}
	}
}

func (av *AppView) DeleteEvents(g *gocui.Gui) {
	_, y := g.CurrentView().Cursor()

	if view, ok := av.FindChildView(WeekdayNames[av.Calendar.CurrentDay.Date.Weekday()]); ok {
		if dayView, ok := view.(*DayView); ok {
			if view, ok := dayView.IsOnEvent(y); ok {
				if eventView, ok := view.(*EventView); ok {
					av.Database.DeleteEventsByName(eventView.Event.Name)
				}
			}
		}
	}
}

func (av *AppView) ShowPopup(g *gocui.Gui) error {
	if view, ok := av.GetChild("popup"); ok {
		if popupView, ok := view.(*EventPopupView); ok {
			view.SetProperties(
				av.X+(av.W-PopupWidth)/2,
				av.Y+(av.H-PopupHeight)/2,
				PopupWidth,
				PopupHeight,
			)
			return popupView.Show(g)
		}
	}
	return nil
}

func (av *AppView) HandleEnter(g *gocui.Gui, v *gocui.View) error {
	if view, ok := av.GetChild("popup"); ok {
		if popupView, ok := view.(*EventPopupView); ok {
			if popupView.IsVisible {
				return popupView.AddEvent(g, v)
			}
		}
	}
    return nil
}

func (av *AppView) HandleEscape(g *gocui.Gui, v *gocui.View) error {
	if view, ok := av.GetChild("popup"); ok {
		if popupView, ok := view.(*EventPopupView); ok {
			if popupView.IsVisible {
				return popupView.Close(g, v)
			}
		}
	}
	if view, ok := av.FindChildView("notepad"); ok {
		if notepadView, ok := view.(*NotepadView); ok {
            if notepadView.IsActive {
                return av.ReturnToMainView(g)
            }
		}
	}
	return nil
}

func (av *AppView) ShowKeybinds(g *gocui.Gui) error {
	if view, ok := av.GetChild("popup"); ok {
		if popupView, ok := view.(*EventPopupView); ok {
			if popupView.IsVisible {
				return nil
			}
		}
	}

	if view, ok := av.GetChild("keybinds"); ok {
		if keybindsView, ok := view.(*KeybindsView); ok {
			if keybindsView.IsVisible {
				keybindsView.IsVisible = false
				return g.DeleteView(keybindsView.Name)
			}

			keybindsView.IsVisible = true
			keybindsView.SetProperties(
				av.X+(av.W-KeybindsWidth)/2,
				av.Y+(av.H-KeybindsHeight)/2,
				KeybindsWidth,
				KeybindsHeight,
			)

			return keybindsView.Update(g)
		}
	}

	return nil
}

func (av *AppView) updateChildViewProperties() {
	mainViewWidth := int(float64(av.W-1) * MainViewWidthRatio)
	sideViewWidth := int(float64(av.W) * SideViewWidthRatio)

	if titleView, ok := av.GetChild("title"); ok {
		titleView.SetProperties(
			av.X+sideViewWidth+1,
			av.Y,
			mainViewWidth,
			TitleViewHeight,
		)
	}

	if mainView, ok := av.GetChild("main"); ok {
		y := av.Y + TitleViewHeight + 1
		mainView.SetProperties(
			av.X+sideViewWidth+1,
			y,
			mainViewWidth,
			av.H-y,
		)
	}

	if sideView, ok := av.GetChild("side"); ok {
		sideView.SetProperties(
			av.X,
			av.Y,
			sideViewWidth,
			av.H,
		)
	}
}

func (av *AppView) updateCurrentView(g *gocui.Gui) error {
	if view, ok := av.GetChild("popup"); ok {
		if popupView, ok := view.(*EventPopupView); ok {
			if popupView.IsVisible {
				return nil
			}
		}
	}

	if g.CurrentView() != nil && g.CurrentView().Name() == "notepad" {
		return nil
	}

	viewName := WeekdayNames[av.Calendar.CurrentDay.Date.Weekday()]
	var y int
	var hoveredView View
	g.Cursor = true

	if view, ok := av.FindChildView("time"); ok {
		if timeView, ok := view.(*TimeView); ok {
			y = utils.TimeToPosition(av.Calendar.CurrentDay.Date, timeView.Body)
		}
	}
	g.Cursor = true

	if view, ok := av.FindChildView(viewName); ok {
		if dayView, ok := view.(*DayView); ok {
			if eventView, ok := dayView.IsOnEvent(y); ok {
				hoveredView = eventView
			} else {
				hoveredView = dayView
			}
		}
	}

	if view, ok := av.FindChildView("hover"); ok {
		if hoverView, ok := view.(*HoverView); ok {
			hoverView.CurrentView = hoveredView
			hoverView.Update(g)
		}
	}

	g.SetCurrentView(viewName)
	g.CurrentView().BgColor = gocui.Attribute(termbox.ColorBlack)
	g.CurrentView().SetCursor(1, y)

	return nil
}
