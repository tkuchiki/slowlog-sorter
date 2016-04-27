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
      --help                     Show context-sensitive help (also try --help-long and --help-man).
  -f, --file=FILE                slow query log
  -r, --reverse                  reverse the result of comparisons
      --pretty                   pretty print
  -p, --query-pattern=PATTERN    query matching PATTERN
  -s, --sort="query_time"        sort by (query_time, lock_time, rows_sent, rows_examined, time)
  -n, --num=0                    number of lines (0 = all)
      --query-time-begin=TIME    query_time begin
      --query-time-end=TIME      query_time end
      --lock-time-begin=TIME     lock_time begin
      --lock-time-end=TIME       lock_time end
      --rows-sent-begin=NUM      rows_sent begin
      --rows-sent-end=NUM        rows_sent end
      --rows-examined-begin=NUM  rows_examined begin
      --rows-examined-end=NUM    rows_examined end
      --time-begin=TIME          time begin
      --time-end=TIME            time end
      --location=LOCATION        location (default: current location)
      --or                       option conditions (default: and)
      --version                  Show application version.
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

```
$ ./slowlog-sorter -n 5 -p commit -f /path/to/mysql-slowquery.log --pretty
Query_time:0.103738     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 08:20:00 +0900 JST

commit;

Query_time:0.107203     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 08:27:04 +0900 JST

commit;

Query_time:0.128485     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 08:43:06 +0900 JST

commit;

Query_time:0.132485     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 08:45:31 +0900 JST

commit;

Query_time:0.137578     Lock_time:0.000000      Rows_sent:0     Rows_examined:0 time:2016-04-21 08:06:21 +0900 JST

commit;

```

### `--query-time-begin`, `--query-time-end`

- `./slowlog-sorter -f /path/to/mysql-slowquery.log --query-time-begin 0.1`
    - Query_time >= 0.1
- `./slowlog-sorter -f /path/to/mysql-slowquery.log --query-time-end 1.5`
    - Query_time < 1.5
- `./slowlog-sorter -f /path/to/mysql-slowquery.log --query-time-begin 0.1 --query-time-end 1.5`
    - 0.1 <= Query_time < 1.5

### `--lock-time-begin`, `--lock-time-end`

- `./slowlog-sorter -f /path/to/mysql-slowquery.log --lock-time-begin 0.00001`
    - Lock_time >= 0.00001
- `./slowlog-sorter -f /path/to/mysql-slowquery.log --lock-time-end 0.0005`
    - Lock_time < 0.0005
- `./slowlog-sorter -f /path/to/mysql-slowquery.log --lock-time-begin 0.00001 --lock-time-end 0.0005`
    - 0.00001 <= Lock_time < 0.0005

### `--rows-sent-begin`, `--rows-sent-end`

- `./slowlog-sorter -f /path/to/mysql-slowquery.log --rows-sent-begin 1000`
    - Rows_sent >= 1000
- `./slowlog-sorter -f /path/to/mysql-slowquery.log --rows-sent-end 10000`
    - Rows_sent < 10000
- `./slowlog-sorter -f /path/to/mysql-slowquery.log --rows-sent-begin 1000 --rows-sent-end 10000`
    - 1000 <= Rows_sent < 10000

### `--rows-examined-begin`, `--rows-examined-end`

- `./slowlog-sorter -f /path/to/mysql-slowquery.log --rows-examined-begin 10000`
    - Rows_examined >= 10000
- `./slowlog-sorter -f /path/to/mysql-slowquery.log --rows-examined-end 100000`
    - Rows_examined < 100000
- `./slowlog-sorter -f /path/to/mysql-slowquery.log --rows-examined-begin 10000 --rows-examined-end 100000`
    - 10000 <= Rows_examined < 100000

### `--time-begin`, `--time-end`, `--location`

Format is `YYYY-MM-DDThh:mm:ss`, `hh:mm:ss`.

#### local timezone offset: +09:00

- `./slowlog-sorter -f /path/to/mysql-slowquery.log --time-begin "2016-04-21T08:20:00"`
    - Time >= 2016-04-21T08:20:00+09:00
- `./slowlog-sorter -f /path/to/mysql-slowquery.log --time-end "2016-04-21T09:00:00"`
    - Time < 2016-04-21T09:00:00+09:00
- `./slowlog-sorter -f /path/to/mysql-slowquery.log --time-begin "2016-04-21T08:20:00" --time-end "2016-04-21T09:00:00"`
    -  2016-04-21T08:20:00+09:00 <= Time < 2016-04-21T09:00:00+09:00

#### local timezone offset: +09:00, current date: 2016-04-27

- `./slowlog-sorter -f /path/to/mysql-slowquery.log --time-begin "18:30:00"`
    - Time >= 2016-04-27T18:30:00+09:00
- `./slowlog-sorter -f /path/to/mysql-slowquery.log --time-end "19:00:00"`
    - Time < 2016-04-27T19:00:00+09:00
- `./slowlog-sorter -f /path/to/mysql-slowquery.log --time-begin "18:30:00" --time-end "19:00:00"`
    -  2016-04-27T18:30:00+09:00 <= Time < 2016-04-27T19:00:00+09:00

#### local timezone offset: +09:00, `--location UTC`

- `./slowlog-sorter -f /path/to/mysql-slowquery.log --time-begin "2016-04-21T08:20:00" --location UTC`
    - Time >= 2016-04-21T08:20:00+00:00
- `./slowlog-sorter -f /path/to/mysql-slowquery.log --time-end "2016-04-21T09:00:00" --location UTC`
    - Time < 2016-04-21T09:00:00+00:00
- `./slowlog-sorter -f /path/to/mysql-slowquery.log --time-begin "2016-04-21T08:20:00" --time-end "2016-04-21T09:00:00" --location`
    -  2016-04-21T08:20:00+00:00 <= Time < 2016-04-21T09:00:00+00:00

#### `--or`

- `./slowlog-sorter -f /path/to/mysql-slowquery.log --query-time-begin 0.1 --query-time-end 1.5 --query-time-begin 0.1 --rows-examined-begin 10000 --rows-examined-end 100000`
    - 0.1 <= Query_time < 1.5 && 10000 <= Rows_examined < 100000
- `./slowlog-sorter -f /path/to/mysql-slowquery.log --query-time-begin 0.1 --query-time-end 1.5 --query-time-begin 0.1 --rows-examined-begin 10000 --rows-examined-end 100000 --or`
    - 0.1 <= Query_time < 1.5 || 10000 <= Rows_examined < 100000
