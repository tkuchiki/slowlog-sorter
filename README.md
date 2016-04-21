# slowlog-sorter

slowlog-sorter is MySQL slow query log sorter.  

## Installation

Download from https://github.com/tkuchiki/slowlog-sorter/releases

## Usage

Read from stdin or an input file(`-f`).  

```
$ ./slowlog-sorter --help
usage: slowlog-sorter [<flags>]

MySQL slow query log sorter (read from file or stdin).

Flags:
      --help                   Show context-sensitive help (also try --help-long and --help-man).
  -f, --file=FILE              slow query log
  -r, --reverse                reverse the result of comparisons
  -p, --query-pattern=PATTERN  query matching PATTERN
  -s, --sort="query_time"      sort by (query_time, lock_time, rows_sent, rows_examined, time)
  -n, --num=0                  number of lines (0 = all)
      --version                Show application version.
```

## Example

```
$ cat /path/to/mysql-slowquery.log | ./slowlog-sorter -s time -n 10 -p commit -r
Query_time:0.253832     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 09:00:35 +0900 JST      sql:commit;
Query_time:0.164165     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 09:00:35 +0900 JST      sql:commit;
Query_time:0.155609     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 09:00:35 +0900 JST      sql:commit;
Query_time:0.209027     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 09:00:35 +0900 JST      sql:commit;
Query_time:0.143576     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 09:00:35 +0900 JST      sql:commit;
Query_time:0.261648     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 09:00:35 +0900 JST      sql:commit;
Query_time:0.132485     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 08:45:31 +0900 JST      sql:commit;
Query_time:0.177635     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 08:45:31 +0900 JST      sql:commit;
Query_time:0.210171     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 08:43:06 +0900 JST      sql:commit;
Query_time:0.128485     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 08:43:06 +0900 JST      sql:commit;
```

```
$ ./slowlog-sorter -n 5 -p commit -f /path/to/mysql-slowquery.log
Query_time:0.103738     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 08:20:00 +0900 JST      sql:commit;
Query_time:0.107203     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 08:27:04 +0900 JST      sql:commit;
Query_time:0.128485     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 08:43:06 +0900 JST      sql:commit;
Query_time:0.132485     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 08:45:31 +0900 JST      sql:commit;
Query_time:0.137578     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 08:06:21 +0900 JST      sql:commit;
```
