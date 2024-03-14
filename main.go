package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/k0l1br1/converter/bins"
)

const (
	exitOk    = 0
	exitError = 1
)

func errorPrint(err error) {
	os.Stderr.WriteString(err.Error() + "\n")
}

func errorWrap(msg string, err error) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func run() int {
	if len(os.Args) < 2 {
		errorPrint(fmt.Errorf("path to bins file is required"))
		return exitError
	}

	var err error
	var off int
	if len(os.Args) > 2 {
		off, err = strconv.Atoi(os.Args[2])
		if err != nil {
			errorPrint(errorWrap("parse offset", err))
			return exitError
		}
	}

	bs, err := bins.ReadAll(os.Args[1])
	if err != nil {
		errorPrint(errorWrap("read bins file", err))
		return exitError
	}
	if len(bs) == 0 {
		errorPrint(fmt.Errorf("bins file is empty"))
		return exitError
	}

	printer := NewPrinter(os.Stdout)
	defer printer.BreakLine()
	printer.Print(bs, off)

	cmdCh := make(chan uint8, 128)
	stopCh := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)
	go readInput(cmdCh, stopCh, &wg)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGWINCH)
	for {
		select {
		case c := <-cmdCh:
			switch c {
			case forwardCmd:
				off++
				if off+printLimit > len(bs) {
					off--
				} else {
                    printer.Reprint(bs, off)
				}
			case backwardCmd:
				off--
				if off < 0 {
					off++
				} else {
                    printer.Reprint(bs, off)
				}
            case toFirstCmd:
                off = 0
                printer.Reprint(bs, off)
            case toLastCmd:
                off = len(bs) - printLimit
                printer.Reprint(bs, off)
			case exitCmd:
				return exitOk
			}
		case s := <-sigChan:
			switch s {
			case syscall.SIGWINCH:
				 // TODO win size?
			default:
				close(stopCh)
				wg.Wait()
				return exitOk
			}
		}
	}
}

func main() {
	os.Exit(run())
}
