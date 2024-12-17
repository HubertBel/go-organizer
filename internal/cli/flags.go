package cli

import (
	"flag"
	"fmt"
	"os"
)

var version = "v1.1.1" // TODO

var (
	Debug bool
    ShowVersion bool
)

func InitFlags() {
    flag.BoolVar(&Debug, "debug", false, "enable debug mode")
    flag.BoolVar(&Debug, "d", false, "")

    flag.BoolVar(&ShowVersion, "version", false, "show application version")
    flag.BoolVar(&ShowVersion, "v", false, "")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "lazyorg\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n")
		fmt.Fprintf(os.Stderr, "  %s [options]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Options:\n")
        flag.PrintDefaults()
	}

    flag.Parse()

    if Debug {
        fmt.Println("Debug mode enabled")
    }

    if ShowVersion {
        fmt.Println(version)
        os.Exit(0)
    }
}
