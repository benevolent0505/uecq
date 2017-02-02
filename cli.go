package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/text/encoding/japanese"
	"io"
	"regexp"
	"strconv"
	"time"
)

// Exit codes are int values that represent an exit code for a particular error.
const (
	ExitCodeOK    int = 0
	ExitCodeError int = 1 + iota
)

// URLs are.
const (
	UndergraduateURL = "http://kyoumu.office.uec.ac.jp/kyuukou/kyuukou.html"
	GraduateURL      = "http://kyoumu.office.uec.ac.jp/kyuukou/kyuukou2.html"
)

// CLI is the command line object
type CLI struct {
	// outStream and errStream are the stdout and stderr
	// to write message from the CLI.
	outStream, errStream io.Writer
}

// Lesson stands for lesson of UEC.
type Lesson struct {
	Class   string    `json:"class"`
	Date    time.Time `json:"date"`
	Period  int       `json:"period"`
	Subject string    `json:"subject"`
	Teacher string    `json:"teacher"`
	Remark  string    `json:"remark"`
}

// LessonsSlice is struct for json.
type LessonsSlice struct {
	Lessons []Lesson `json:"lessons"`
}

// ToArray converts Lesson struct to an array.
func (l Lesson) ToArray() []string {
	lesson := []string{
		l.Class,
		l.Date.String(),
		strconv.Itoa(l.Period),
		l.Subject,
		l.Teacher,
		l.Remark,
	}

	return lesson
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

	flags.BoolVar(&graduate, "graduate", false, "For graduate school mode")
	flags.StringVar(&f, "f", "json", "Setting output format.")
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

	var l LessonsSlice
	l.Lessons = getLessons(url)

	if f == "json" {
		b, err := json.Marshal(l)
		if err != nil {
			fmt.Fprintf(cli.errStream, "%v\n", err)
			return ExitCodeError
		}

		fmt.Println(string(b))
	} else if f == "text" {
		table := tablewriter.NewWriter(cli.outStream)
		table.SetHeader([]string{"クラス", "日時", "時限", "科目", "担当教員", "備考"})

		for _, lesson := range l.Lessons {
			table.Append(lesson.ToArray())
		}

		table.Render()
	}

	return ExitCodeOK
}

func getLessons(url string) []Lesson {
	doc, _ := goquery.NewDocument(url)
	lessons := []Lesson{}
	doc.Find("table > tbody > tr").Next().Each(func(_ int, row *goquery.Selection) {
		elems := make([]string, 6)
		row.Find("td").Each(func(i int, s *goquery.Selection) {
			raw, _ := s.Html()
			if raw != string(0xA0) { // nbsp避け
				text, _ := japanese.ShiftJIS.NewDecoder().String(raw)
				elems[i] = text
			}
		})

		period, _ := strconv.Atoi(elems[2])

		monthAndDay := regexp.MustCompile("(\\d+)月(\\d+)日").FindStringSubmatch(elems[1])
		month, _ := strconv.Atoi(monthAndDay[1])
		day, _ := strconv.Atoi(monthAndDay[2])
		loc, _ := time.LoadLocation("Asia/Tokyo")
		date := time.Date(time.Now().Year(), time.Month(month), day, 0, 0, 0, 0, loc)

		lessons = append(lessons, Lesson{
			Class:   elems[0],
			Date:    date,
			Period:  period,
			Subject: elems[3],
			Teacher: elems[4],
			Remark:  elems[5],
		})
	})

	return lessons
}
