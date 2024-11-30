package views

import (
	"strconv"
	"time"

	"github.com/HubertBel/lazyorg/internal/calendar"
	"github.com/HubertBel/lazyorg/internal/database"
	"github.com/j-04/gocui-component"
	"github.com/jroimartin/gocui"
)

type EventPopupView struct {
	*BaseView
	Form     *component.Form
	Calendar *calendar.Calendar
	Database *database.Database

    EventToEdit *calendar.Event

	IsVisible bool
}

func NewEvenPopup(g *gocui.Gui, c *calendar.Calendar, db *database.Database) *EventPopupView {

	epv := &EventPopupView{
		BaseView:  NewBaseView("popup"),
		Form:      nil,
		Calendar:  c,
		Database:  db,
		IsVisible: false,
	}

	return epv
}

func (epv *EventPopupView) Update(g *gocui.Gui) error {
	return nil
}

func (epv *EventPopupView) NewForm(g *gocui.Gui, title, name, time, location, duration, frequency, occurence, description string) *component.Form {
	form := component.NewForm(g, title, epv.X, epv.Y, epv.W, epv.H)

	form.AddInputField("Name", LabelWidth, FieldWidth).SetText(name)
	form.AddInputField("Time", LabelWidth, FieldWidth).SetText(time)
	form.AddInputField("Location", LabelWidth, FieldWidth).SetText(location)
	form.AddInputField("Duration", LabelWidth, FieldWidth).SetText(duration)
	form.AddInputField("Frequency", LabelWidth, FieldWidth).SetText(frequency)
	form.AddInputField("Occurence", LabelWidth, FieldWidth).SetText(occurence)
	form.AddInputField("Description", LabelWidth, FieldWidth).SetText(description)

	return form
}

func (epv *EventPopupView) ShowNewEventPopup(g *gocui.Gui) error {
	if epv.IsVisible {
		return nil
	}

	epv.Form = epv.NewForm(g, "New Event", "", epv.Calendar.CurrentDay.Date.Format(TimeFormat), "", "1.0", "7", "1", "")
	epv.Form.AddButton("Add", epv.AddEvent)
	epv.Form.AddButton("Cancel", epv.Close)
	if err := epv.initKeybindings(g, true); err != nil {
		return err
	}

	epv.Form.SetCurrentItem(0)
	epv.IsVisible = true
	epv.Form.Draw()

	return nil
}

func (epv *EventPopupView) ShowEditEventPopup(g *gocui.Gui, eventView *EventView) error {
	if epv.IsVisible {
		return nil
	}

	epv.EventToEdit = eventView.Event

	epv.Form = epv.NewForm(g,
		"Edit Event",
		epv.EventToEdit.Name,
		epv.EventToEdit.Time.Format(TimeFormat),
		epv.EventToEdit.Location, strconv.FormatFloat(epv.EventToEdit.DurationHour, 'f', -1, 64),
		strconv.Itoa(epv.EventToEdit.FrequencyDay),
		strconv.Itoa(epv.EventToEdit.Occurence),
		epv.EventToEdit.Description,
	)
	epv.Form.AddButton("Edit", epv.EditEvent)
	if epv.EventToEdit.Occurence > 1 {
		epv.Form.AddButton("Edit All", epv.AddEvent)
	}
	epv.Form.AddButton("Cancel", epv.Close)
	if err := epv.initKeybindings(g, false); err != nil {
		return err
	}

	epv.Form.SetCurrentItem(0)
	epv.IsVisible = true
	epv.Form.Draw()

	return nil
}

func (epv *EventPopupView) CreateEventFromInputs() *calendar.Event {
	name := epv.Form.GetFieldText("Name")
	time, _ := time.Parse(TimeFormat, epv.Form.GetFieldText("Time"))
	location := epv.Form.GetFieldText("Location")

	duration, _ := strconv.ParseFloat(epv.Form.GetFieldText("Duration"), 64)
	frequency, _ := strconv.Atoi(epv.Form.GetFieldText("Frequency"))
	occurence, _ := strconv.Atoi(epv.Form.GetFieldText("Occurence"))

	description := epv.Form.GetFieldText("Description")

	return calendar.NewEvent(name, description, location, time, duration, frequency, occurence)
}

func (epv *EventPopupView) EditEvent(g *gocui.Gui, v *gocui.View) error {
	if !epv.IsVisible {
		return nil
	}

	newEvent := epv.CreateEventFromInputs()

	if err := epv.Database.EditEventById(epv.EventToEdit.Id, newEvent); err != nil {
		return err
	}

	return epv.Close(g, v)
}

func (epv *EventPopupView) AddEvent(g *gocui.Gui, v *gocui.View) error {
	if !epv.IsVisible {
		return nil
	}

	event := epv.CreateEventFromInputs()

	if _, err := epv.Database.AddRecurringEvents(event); err != nil {
		return err
	}

	return epv.Close(g, v)
}

func (epv *EventPopupView) Close(g *gocui.Gui, v *gocui.View) error {
	epv.IsVisible = false
	return epv.Form.Close(g, v)
}

func (epv *EventPopupView) initKeybindings(g *gocui.Gui, isNewEvent bool) error {
	for _, item := range epv.Form.GetItems() {
        if err := g.SetKeybinding(item.GetLabel(), gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
            return epv.Close(g, v)
        }); err != nil {
            return err
        }

        if err := g.SetKeybinding(item.GetLabel(), gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
            if isNewEvent {
                return epv.AddEvent(g, v)
            } else {
                return epv.EditEvent(g, v)
            }
        }); err != nil {
            return err
        }
	}

	return nil
}
