package views

import (
	"strconv"
	"time"

	"github.com/HubertBel/lazyorg/internal/calendar"
	"github.com/HubertBel/lazyorg/internal/database"
	component "github.com/j-04/gocui-component"
	"github.com/jroimartin/gocui"
)

type EventPopupView struct {
	*BaseView
	Form     *component.Form
	Calendar *calendar.Calendar
	Database *database.Database

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

func (epv *EventPopupView) NewEventForm(g *gocui.Gui, title, name, time, location, duration, frequency, occurence, description string) *component.Form {
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

func (epv *EventPopupView) EditEventForm(g *gocui.Gui, title, name, time, location, duration, description string) *component.Form {
	form := component.NewForm(g, title, epv.X, epv.Y, epv.W, epv.H)

	form.AddInputField("Name", LabelWidth, FieldWidth).SetText(name)
	form.AddInputField("Time", LabelWidth, FieldWidth).SetText(time)
	form.AddInputField("Location", LabelWidth, FieldWidth).SetText(location)
	form.AddInputField("Duration", LabelWidth, FieldWidth).SetText(duration)
	form.AddInputField("Description", LabelWidth, FieldWidth).SetText(description)

	return form
}

func (epv *EventPopupView) ShowNewEventPopup(g *gocui.Gui) error {
	if epv.IsVisible {
		return nil
	}

	epv.Form = epv.NewEventForm(g, "New Event", "", epv.Calendar.CurrentDay.Date.Format(TimeFormat), "", "1.0", "7", "1", "")
	epv.Form.AddButton("Add", epv.AddEvent)
	epv.Form.AddButton("Cancel", epv.Close)

	if err := epv.setPopupKeybind(g, gocui.KeyEsc, gocui.ModNone, epv.Close); err != nil {
		return err
	}
	if err := epv.setPopupKeybind(g, gocui.KeyEnter, gocui.ModNone, epv.AddEvent); err != nil {
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

	event := eventView.Event

	epv.Form = epv.EditEventForm(g,
		"Edit Event",
		event.Name,
		event.Time.Format(TimeFormat),
		event.Location,
		strconv.FormatFloat(event.DurationHour, 'f', -1, 64),
		event.Description,
	)

	editHandler := func(g *gocui.Gui, v *gocui.View) error {
        return epv.EditEvent(g, v, event)
	}

	epv.Form.AddButton("Edit", editHandler)
	epv.Form.AddButton("Cancel", epv.Close)

	if err := epv.setPopupKeybind(g, gocui.KeyEsc, gocui.ModNone, epv.Close); err != nil {
		return err
	}
	if err := epv.setPopupKeybind(g, gocui.KeyEnter, gocui.ModNone, editHandler); err != nil {
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

func (epv *EventPopupView) AddEvent(g *gocui.Gui, v *gocui.View) error {
	if !epv.IsVisible {
		return nil
	}

	event := epv.CreateEventFromInputs()
	events := event.GetReccuringEvents()

	for _, v := range events {
		if _, err := epv.Database.AddEvent(v); err != nil {
			return err
		}
	}

	return epv.Close(g, v)
}

func (epv *EventPopupView) EditEvent(g *gocui.Gui, v *gocui.View, event *calendar.Event) error {
	if !epv.IsVisible {
		return nil
	}

	newEvent := epv.CreateEventFromInputs()
	newEvent.Id = event.Id

	if err := epv.Database.UpdateEventById(event.Id, newEvent); err != nil {
		return err
	}

	return epv.Close(g, v)
}

func (epv *EventPopupView) Close(g *gocui.Gui, v *gocui.View) error {
	epv.IsVisible = false
	return epv.Form.Close(g, v)
}

func (epv *EventPopupView) setPopupKeybind(g *gocui.Gui, key interface{}, mod gocui.Modifier, handler func(g *gocui.Gui, v *gocui.View) error) error {
	for _, item := range epv.Form.GetItems() {
		if err := g.SetKeybinding(item.GetLabel(), key, mod, handler); err != nil {
			return err
		}
	}

	return nil
}
