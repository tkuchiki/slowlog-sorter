package main

import (
	"fmt"
	"github.com/handlename/go-mysql-slowlog-parser"
	"sort"
	"time"
)

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
	for _, s := range slowlogs {
		outputRow(s, c)
	}
}

func outputRow(s slowlog.Parsed, c Config) {
	t := time.Unix(s.Datetime, 0).In(c.Location).String()

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
