package main

import (
	"fmt"

	"github.com/gocql/gocql"
)

type User struct {
	FirstName string
	LastName  string
}

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
	var userList []User
	m := map[string]interface{}{}
	query := "SELECT id,age,firstname,lastname,city,email FROM users"
	iterable := Session.Query(query).Iter()
	for iterable.MapScan(m) {
		userList = append(userList, User{
			FirstName: m["firstname"].(string),
			LastName:  m["lastname"].(string),
		})
		m = map[string]interface{}{}
	}
	fmt.Println("User List Data", userList)
}
