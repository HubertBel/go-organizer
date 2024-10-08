package views

import "github.com/jroimartin/gocui"

type View interface {
	Update(g *gocui.Gui) error
	SetProperties(x, y, w, h int)
    GetName() string
    GetProperties() (int, int, int, int)
	AddChild(name string, child View)
	GetChild(name string) (View, bool)
    ClearChildren(g *gocui.Gui) error
	Children() map[string]View
}

type BaseView struct {
	Name     string
	X, Y, W, H int
	children   map[string]View
}

func NewBaseView(name string) *BaseView {
	return &BaseView{
		Name:     name,
		children: make(map[string]View),
	}
}

func (bv *BaseView) SetProperties(x, y, w, h int) {
	bv.X, bv.Y, bv.W, bv.H = x, y, w, h
}

func (bv *BaseView) GetProperties() (int, int, int, int) {
	return bv.X, bv.Y, bv.W, bv.H
}

func (bv *BaseView) GetName() string {
	return bv.Name
}

func (bv *BaseView) ClearChildren(g *gocui.Gui) error {
    for _, v := range bv.children {
        if err := g.DeleteView(v.GetName()); err != nil {
            return err
        }
    }
    clear(bv.children)

    return nil
}

func (bv *BaseView) AddChild(name string, child View) {
	bv.children[name] = child
}

func (bv *BaseView) GetChild(name string) (View, bool) {
	child, ok := bv.children[name]
	return child, ok
}

func (bv *BaseView) Children() map[string]View {
	return bv.children
}

func (bv *BaseView) UpdateChildren(g *gocui.Gui) error {
    
	for _, child := range bv.children {
		if err := child.Update(g); err != nil {
			return err
		}
	}

    return nil
}
