package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gocql/gocql"
)

var (
	CHOST   = os.Getenv("CASSANDRA_HOST")
	CPORT   = os.Getenv("CASSANDRA_PORT")
	CUSER   = os.Getenv("CASSANDRA_USER")
	CPASS   = os.Getenv("CASSANDRA_PASSWORD")
	CCAPATH = os.Getenv("CASSANDRA_CA_PATH")
)

func selectRows(session *gocql.Session, cql string) (err error) {
	log.Printf("cql: %+v\n", cql)
	iter := session.Query(cql).Iter()
	defer iter.Close()

	rows, err := iter.SliceMap()
	if err != nil {
		return fmt.Errorf("failed to query: %w", err)
	}

	for _, row := range rows {
		fmt.Printf("  row: %+v\n", row)
	}

	return
}

func newSslOptions() (opts *gocql.SslOptions) {
	config := &tls.Config{
		ServerName:         CHOST,
		InsecureSkipVerify: false,
	}
	opts = &gocql.SslOptions{
		Config:                 config,
		EnableHostVerification: true,
		CaPath:                 CCAPATH,
	}
	return
}

func getClusterConfig() (cluster *gocql.ClusterConfig) {
	port, _ := strconv.Atoi(CPORT)
	cluster = gocql.NewCluster(CHOST)
	cluster.CQLVersion = "3.4.4"
	cluster.Port = port
	cluster.Consistency = gocql.LocalOne
	cluster.SerialConsistency = gocql.LocalSerial
	cluster.ProtoVersion = 4
	cluster.Timeout = 3 * time.Second
	cluster.ConnectTimeout = 3 * time.Second

	if CPASS != "" {
		cluster.Authenticator = gocql.PasswordAuthenticator{
			Username: CUSER,
			Password: CPASS,
		}
	}

	if CCAPATH != "" {
		cluster.SslOpts = newSslOptions()
	}

	return cluster
}

var cql = flag.String("cql", "", "specify cql statement")

func main() {
	flag.Parse()

	cluster := getClusterConfig()
	session, err := cluster.CreateSession()
	if err != nil {
		log.Printf("failed to create session: %+v\n", err)
		return
	}
	defer session.Close()

	selectRows(session, "select cluster_name, release_version from system.local")

	if *cql != "" {
		if err := selectRows(session, *cql); err != nil {
			log.Printf("%+v\n", err)
		}
	}
}
