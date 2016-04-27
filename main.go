package main

import (
	"log"
	"os"
	"regexp"
	"time"

	"github.com/handlename/go-mysql-slowlog-parser"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Config struct {
	Reverse           bool
	Include           string
	Exclude           string
	Line              int
	Pretty            bool
	QueryTimeBegin    string
	QueryTimeEnd      string
	LockTimeBegin     string
	LockTimeEnd       string
	RowsSentBegin     string
	RowsSentEnd       string
	RowsExaminedBegin string
	RowsExaminedEnd   string
	TimeBegin         string
	TimeEnd           string
	RangeOr           bool
	Location          *time.Location
}

func currentLocation() *time.Location {
	zone, offset := time.Now().In(time.Local).Zone()

	return time.FixedZone(zone, offset)
}

var (
	file     = kingpin.Flag("file", "slow query log").Short('f').String()
	reverse  = kingpin.Flag("reverse", "reverse the result of comparisons").Short('r').Bool()
	pretty   = kingpin.Flag("pretty", "pretty print").Bool()
	include  = kingpin.Flag("include", "don't exclude query matching PATTERN").Short('i').PlaceHolder("PATTERN").String()
	exclude  = kingpin.Flag("exclude", "exclude query matching PATTERN").Short('e').PlaceHolder("PATTERN").String()
	sortBy   = kingpin.Flag("sort", "sort by (query_time, lock_time, rows_sent, rows_examined, time)").Short('s').Default("query_time").String()
	nol      = kingpin.Flag("num", "number of lines (0 = all)").Short('n').Default("0").Int()
	qtBegin  = kingpin.Flag("query-time-begin", "query_time begin").PlaceHolder("TIME").String()
	qtEnd    = kingpin.Flag("query-time-end", "query_time end").PlaceHolder("TIME").String()
	ltBegin  = kingpin.Flag("lock-time-begin", "lock_time begin").PlaceHolder("TIME").String()
	ltEnd    = kingpin.Flag("lock-time-end", "lock_time end").PlaceHolder("TIME").String()
	rsBegin  = kingpin.Flag("rows-sent-begin", "rows_sent begin").PlaceHolder("NUM").String()
	rsEnd    = kingpin.Flag("rows-sent-end", "rows_sent end").PlaceHolder("NUM").String()
	reBegin  = kingpin.Flag("rows-examined-begin", "rows_examined begin").PlaceHolder("NUM").String()
	reEnd    = kingpin.Flag("rows-examined-end", "rows_examined end").PlaceHolder("NUM").String()
	tBegin   = kingpin.Flag("time-begin", "time begin").PlaceHolder("TIME").String()
	tEnd     = kingpin.Flag("time-end", "time end").PlaceHolder("TIME").String()
	location = kingpin.Flag("location", "location (default: current location)").String()
	rangeOr  = kingpin.Flag("or", "option conditions (default: and)").Bool()
)

func main() {
	kingpin.CommandLine.Help = "MySQL slow query log sorter (read from file or stdin)."
	kingpin.Version("0.2.1")
	kingpin.Parse()

	fileinfo, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}

	var f *os.File
	if fileinfo.Mode()&os.ModeNamedPipe == 0 {
		f, err = os.Open(*file)
		defer f.Close()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		f = os.Stdin
	}

	parser := slowlog.NewParser()
	c := Config{
		Reverse:           *reverse,
		Include:           *include,
		Exclude:           *exclude,
		Line:              *nol,
		Pretty:            *pretty,
		QueryTimeBegin:    *qtBegin,
		QueryTimeEnd:      *qtEnd,
		LockTimeBegin:     *ltBegin,
		LockTimeEnd:       *ltEnd,
		RowsSentBegin:     *rsBegin,
		RowsSentEnd:       *rsEnd,
		RowsExaminedBegin: *reBegin,
		RowsExaminedEnd:   *reEnd,
		TimeBegin:         *tBegin,
		TimeEnd:           *tEnd,
		RangeOr:           *rangeOr,
	}

	if *location == "" {
		c.Location = currentLocation()
	} else {
		var loc *time.Location
		loc, err = time.LoadLocation(*location)
		if err != nil {
			log.Fatal(err)
		}

		c.Location = loc
	}

	slowlogs := SlowLogs{}

	reInc := regexp.MustCompile(c.Include)
	reExc := regexp.MustCompile(c.Exclude)
	i := 1
	for parsed := range parser.Parse(f) {

		if parsed.Sql == "" {
			continue
		}

		if c.Exclude != "" && reExc.Match([]byte(parsed.Sql)) {
			continue
		}

		if c.Include != "" && !reInc.Match([]byte(parsed.Sql)) {
			continue
		}

		if !VerifyRange(parsed, c) {
			continue
		}

		slowlogs = append(slowlogs, parsed)

		if c.Line != 0 && c.Line <= i {
			break
		}

		i++
	}

	switch *sortBy {
	case "query_time":
		SortByQueryTime(slowlogs, c)
	case "lock_time":
		SortByLockTime(slowlogs, c)
	case "rows_sent":
		SortByRowsSent(slowlogs, c)
	case "rows_examined":
		SortByRowsExamined(slowlogs, c)
	case "time":
		SortByTime(slowlogs, c)
	default:
		log.Fatal("Invalid sort key")
	}
}
