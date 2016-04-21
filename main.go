package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"time"

	"github.com/handlename/go-mysql-slowlog-parser"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Config struct {
	Reverse bool
	Pattern string
	Line    int
	Pretty  bool
}

type SlowLogs []slowlog.Parsed

type ByQueryTime struct{ SlowLogs }
type ByLockTime struct{ SlowLogs }
type ByRowsSent struct{ SlowLogs }
type ByRowsExamined struct{ SlowLogs }
type ByTime struct{ SlowLogs }

func (s SlowLogs) Len() int      { return len(s) }
func (s SlowLogs) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s ByQueryTime) Less(i, j int) bool { return s.SlowLogs[i].QueryTime < s.SlowLogs[j].QueryTime }
func (s ByLockTime) Less(i, j int) bool  { return s.SlowLogs[i].LockTime < s.SlowLogs[j].LockTime }
func (s ByRowsSent) Less(i, j int) bool  { return s.SlowLogs[i].RowsSent < s.SlowLogs[j].RowsSent }
func (s ByRowsExamined) Less(i, j int) bool {
	return s.SlowLogs[i].RowsExamined < s.SlowLogs[j].RowsExamined
}
func (s ByTime) Less(i, j int) bool { return s.SlowLogs[i].Datetime < s.SlowLogs[j].Datetime }

func SortByQueryTime(s SlowLogs, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByQueryTime{s}))
	} else {
		sort.Sort(ByQueryTime{s})
	}
	Output(s, c)
}

func SortByLockTime(s SlowLogs, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByLockTime{s}))
	} else {
		sort.Sort(ByLockTime{s})
	}
	Output(s, c)
}

func SortByRowsSent(s SlowLogs, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByRowsSent{s}))
	} else {
		sort.Sort(ByRowsSent{s})
	}
	Output(s, c)
}

func SortByRowsExamined(s SlowLogs, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByRowsExamined{s}))
	} else {
		sort.Sort(ByRowsExamined{s})
	}
	Output(s, c)
}

func SortByTime(s SlowLogs, c Config) {
	if c.Reverse {
		sort.Sort(sort.Reverse(ByTime{s}))
	} else {
		sort.Sort(ByTime{s})
	}
	Output(s, c)
}

func Output(slowlogs SlowLogs, c Config) {
	i := 0
	for _, s := range slowlogs {
		if c.Pattern == "" {
			outputRow(s, c)
			i++
		} else {
			re := regexp.MustCompile(c.Pattern)
			if ok := re.Match([]byte(s.Sql)); ok {
				outputRow(s, c)
				i++
			}
		}

		if c.Line != 0 && c.Line <= i {
			break
		}
	}
}

func outputRow(s slowlog.Parsed, c Config) {
	t := time.Unix(s.Datetime, 0)

	if c.Pretty {
		fmt.Printf("Query_time:%f\tLock_time:%f\tRows_sent:%d\tRows_examined:%d\ttime:%s\n\n",
			s.QueryTime, s.LockTime, s.RowsSent, s.RowsExamined, t,
		)
		fmt.Println(s.Sql)
		fmt.Println("")
	} else {
		fmt.Printf("Query_time:%f\tLock_time:%f\tRows_sent:%d\tRows_examined:%d\ttime:%s\tsql:%s\n",
			s.QueryTime, s.LockTime, s.RowsSent, s.RowsExamined, t, s.Sql,
		)
	}
}

var (
	file    = kingpin.Flag("file", "slow query log").Short('f').String()
	reverse = kingpin.Flag("reverse", "reverse the result of comparisons").Short('r').Bool()
	pretty  = kingpin.Flag("pretty", "pretty print").Bool()
	pattern = kingpin.Flag("query-pattern", "query matching PATTERN").Short('p').PlaceHolder("PATTERN").String()
	sortBy  = kingpin.Flag("sort", "sort by (query_time, lock_time, rows_sent, rows_examined, time)").Short('s').Default("query_time").String()
	nol     = kingpin.Flag("num", "number of lines (0 = all)").Short('n').Default("0").Int()
)

func main() {
	kingpin.CommandLine.Help = "MySQL slow query log sorter (read from file or stdin)."
	kingpin.Version("0.1.1")
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
		Reverse: *reverse,
		Pattern: *pattern,
		Line:    *nol,
		Pretty:  *pretty,
	}

	slowlogs := SlowLogs{}
	for parsed := range parser.Parse(f) {
		if parsed.Sql == "" {
			continue
		}
		slowlogs = append(slowlogs, parsed)
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
