# gocql-sample

sample code to use gocql

## How to run

```bash
$ export CASSANDRA_HOST="127.0.0.1"
$ export CASSANDRA_PORT=9042
$ export CASSANDRA_USER=""
$ export CASSANDRA_PASSWORD=""
```

```bash
$ go build
$ ./gocql-sample -cql "select * from keyspace.table"
2019/12/18 14:24:34 cql: select cluster_name, release_version from system.local
  row: map[cluster_name:test release_version:3.11.2]
2019/12/18 14:24:34 cql: select * from keyspace.table
  row: map[...]
```
