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

func (epv *EventPopupView) Show(g *gocui.Gui) error {
	if epv.IsVisible {
		return nil
	}

	form := component.NewForm(g, "New Event", epv.X, epv.Y, epv.W, epv.H)

	form.AddInputField("Name", LabelWidth, FieldWidth)
	form.AddInputField("Time", LabelWidth, FieldWidth).SetText(epv.Calendar.CurrentDay.Date.Format(TimeFormat))
	form.AddInputField("Location", LabelWidth, FieldWidth)
	form.AddInputField("Duration", LabelWidth, FieldWidth).SetText("1.0")
	form.AddInputField("Frequency", LabelWidth, FieldWidth).SetText("7")
	form.AddInputField("Occurence", LabelWidth, FieldWidth).SetText("1")
	form.AddInputField("Description", LabelWidth, FieldWidth)

	form.AddButton("Add", epv.AddEvent)
	form.AddButton("Cancel", epv.Close)

	form.SetCurrentItem(0)

	epv.Form = form
	epv.IsVisible = true

	if err := epv.initKeybindings(); err != nil {
		return err
	}

	form.Draw()

	return nil
}

func (epv *EventPopupView) AddEvent(g *gocui.Gui, v *gocui.View) error {

	name := epv.Form.GetFieldText("Name")
	time, _ := time.Parse(TimeFormat, epv.Form.GetFieldText("Time"))
	location := epv.Form.GetFieldText("Location")

	duration, _ := strconv.ParseFloat(epv.Form.GetFieldText("Duration"), 64)
	frequency, _ := strconv.Atoi(epv.Form.GetFieldText("Frequency"))
	occurence, _ := strconv.Atoi(epv.Form.GetFieldText("Occurence"))

	description := epv.Form.GetFieldText("Description")

	event := calendar.NewEvent(name, description, location, time, duration, frequency, occurence)

	if _, err := epv.Database.AddRecurringEvents(event); err != nil {
		return err
	}

	return epv.Close(g, v)
}

func (epv *EventPopupView) Close(g *gocui.Gui, v *gocui.View) error {
	epv.IsVisible = false
	return epv.Form.Close(g, v)
}

func (epv *EventPopupView) initKeybindings() error {
	if err := epv.Form.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, epv.Close); err != nil {
		return err
	}

	if err := epv.Form.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, epv.AddEvent); err != nil {
		return err
	}

	return nil
}
