package main

import (
	"fmt"

	"github.com/gosuri/uilive"
)

type Tui struct {
	processImageWriter *uilive.Writer
	writerImageWriter  *uilive.Writer
}

func InitTui() *Tui {

	procTui := uilive.New()
	writeTui := uilive.New()

	procTui.Start() 
	writeTui.Start()

	return &Tui{
		processImageWriter: procTui,
		writerImageWriter: writeTui,
	}
}

func (tui *Tui) UpdateProcessText(text string) { 
	fmt.Fprintln(tui.processImageWriter, text)
}

func (tui *Tui) UpdateWriterText(text string) { 
	fmt.Fprintln(tui.writerImageWriter, text)
} 

func (tui *Tui) CloseTui() { 
	tui.processImageWriter.Stop()
	tui.writerImageWriter.Stop()
}
