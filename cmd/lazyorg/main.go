package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/HubertBel/lazyorg/internal/cli"
	"github.com/HubertBel/lazyorg/internal/database"
	"github.com/HubertBel/lazyorg/internal/ui"
	"github.com/HubertBel/lazyorg/pkg/views"
	"github.com/jroimartin/gocui"
)

func main() {
    cli.InitFlags()

    homeDir, err := os.UserHomeDir()
    if err != nil {
        log.Fatal(err)
    }
    dbDirPath := filepath.Join(homeDir, ".local", "share", "lazyorg")
    dbFilePath := filepath.Join(dbDirPath, "data.db")

    err = os.MkdirAll(dbDirPath, 0755)
    if err != nil {
        log.Fatal(err)
    }

	database := &database.Database{}
	err = database.InitDatabase(dbFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer database.CloseDatabase()

	g, err := gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	av := views.NewAppView(g, database)
	g.SetManager(av)

	if err := ui.InitKeybindings(g, av); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
