package main

import (
	"bytes"
	"os"
	"sync"
	"time"

	"golang.org/x/term"
)

const (
	noCmd uint8 = iota
	forwardCmd
	backwardCmd
	toFirstCmd
	toLastCmd
	exitCmd
	escCmd
)

func readInput(cmdCh chan uint8, stopCh chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	// switch stdin into 'raw' mode  to read each byte at a time
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	defer term.Restore(int(os.Stdin.Fd()), oldState)

	buf := make([]byte, 1)
	cmdBuf := make([]byte, 0, 128)
	for {
		select {
		case <-stopCh:
			return
		default:
			if _, err := os.Stdin.Read(buf); err == nil {
				cmdBuf = append(cmdBuf, buf[0])
				cmd := feed(cmdBuf)
				if cmd != noCmd {
					switch cmd {
					case escCmd:
						cmdBuf = cmdBuf[:0]
					case exitCmd:
						cmdCh <- exitCmd
						return
					default:
						cmdCh <- cmd
						cmdBuf = cmdBuf[:0]
					}
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

type CmdBind struct {
	Key  uint8
	Code []byte
}

var commands = []*CmdBind{
	{Key: forwardCmd, Code: []byte{0x6b}},       // k
	{Key: backwardCmd, Code: []byte{0x68}},      // h
	{Key: toLastCmd, Code: []byte{0x47}},        // G
	{Key: toFirstCmd, Code: []byte{0x67, 0x67}}, // gg
	{Key: exitCmd, Code: []byte{0x71}},          // q
	{Key: exitCmd, Code: []byte{0x3}},           // CtrC
}

func feed(b []byte) uint8 {
	if b[len(b)-1] == 0x1b {
		return escCmd
	}
	for _, k := range commands {
		if bytes.Equal(k.Code, b) {
			return k.Key
		}
	}
	return noCmd
}
