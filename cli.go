package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/olekukonko/tablewriter"

	"io"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Run invokes the CLI with the given arguments.
func (cli *CLI) Run(args []string) int {
	var (
		f        string
		graduate bool
		version  bool
	)

	// Define option flag parse
	flags := flag.NewFlagSet(Name, flag.ContinueOnError)
	flags.SetOutput(cli.errStream)

	flags.BoolVar(&graduate, "graduate", false, "Set graduate school mode flag.")
	flags.StringVar(&f, "f", "text", "Setting output format.")
	flags.BoolVar(&version, "version", false, "Print version information and quit.")

	// Parse commandline flag
	if err := flags.Parse(args[1:]); err != nil {
		return ExitCodeError
	}

	// Show version
	if version {
		fmt.Fprintf(cli.errStream, "%s version %s\n", Name, Version)
		return ExitCodeOK
	}

	// Select ungraduate/graduate school
	var url string
	if graduate {
		url = GraduateURL
	} else {
		url = UndergraduateURL
	}

	var l LecturesSlice
	l.Lectures = GetLectures(url)

	if f == "json" {
		b, err := json.Marshal(l)
		if err != nil {
			fmt.Fprintf(cli.errStream, "%v\n", err)
			return ExitCodeError
		}

		fmt.Fprintf(cli.outStream, "%s\n", b)
	} else if f == "text" {
		table := tablewriter.NewWriter(cli.outStream)
		table.SetHeader([]string{"クラス", "日時", "時限", "科目", "担当教員", "備考"})

		for _, lecture := range l.Lectures {
			table.Append(lecture.ToArray())
		}

		table.Render()
	}

	return ExitCodeOK
}
