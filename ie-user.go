package main

import (
	"flag"
	"fmt"
	"github.com/intervention-engine/ie/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"os"
)

func main() {
	usernameFlag := flag.String("username", "", "username to add")
	passwordFlag := flag.String("password", "", "password to add")
	flag.Parse()

	username := *usernameFlag
	password := *passwordFlag

	newuser := &models.User{
		Username: username,
		ID:       bson.NewObjectId(),
	}
	newuser.SetPassword(password)

	mongoHost := os.Getenv("MONGO_PORT_27017_TCP_ADDR")
	if mongoHost == "" {
		mongoHost = "localhost"
	}

	session, err := mgo.Dial(mongoHost)
	if err != nil {
		panic(err)
	}

	count, err := session.DB("fhir").C("users").Find(bson.M{"username": username}).Count()
	if count > 0 {
		fmt.Printf("User %s already exists.\n", username)
		return
	}

	err = session.DB("fhir").C("users").Insert(newuser)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Successfully added user %s.\n", username)
	}
}
