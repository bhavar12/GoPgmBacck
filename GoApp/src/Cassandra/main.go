package main

import (
	"fmt"
	"github.com/gocql/gocql"
)

//Session for cassandra connection
var Session *gocql.Session

func main() {
	var err error
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "streamdemoapi"
	Session, err = cluster.CreateSession()
	if err != nil {
		panic(err)
	}
	str := "Hello"

	fmt.Println("cassandra " + str + "init done")
}
