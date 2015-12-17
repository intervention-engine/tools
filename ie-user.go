package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/intervention-engine/ie/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func main() {

	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		printUsage()
		return
	}

	command := args[0]

	switch command {
	default:
		fmt.Println("Invalid command.")
		printUsage()
		return
	case "add":
		if len(args) == 3 {
			username := args[1]
			password := args[2]
			addUser(username, password)
		} else {
			fmt.Println("Invalid number of arguments for 'add' command.")
			printUsage()
			return
		}
	case "add_file":
		if len(args) == 2 {
			filepath := args[1]
			loadUsersFromFile(filepath)
		} else {
			fmt.Println("Invalid number of arguments for 'add_file' command.")
			printUsage()
			return
		}
	case "remove":
		if len(args) == 2 {
			username := args[1]
			removeUser(username)
		} else {
			fmt.Println("Invalid number of arguments for 'remove' command.")
			printUsage()
			return
		}
	case "remove_all":
		if len(args) == 1 {
			removeAllUsers()
		} else {
			fmt.Println("Invalid number of arguments for 'remove_all' command.")
			printUsage()
			return
		}
	case "change_pass":
		if len(args) == 3 {
			username := args[1]
			password := args[2]
			updatePassword(username, password)
		} else {
			fmt.Println("Invalid number of arguments for 'change_pass' command.")
			printUsage()
			return
		}
	}
}

func addUser(username, password string) {
	newuser := &models.User{
		Username: username,
		ID:       bson.NewObjectId(),
	}
	newuser.SetPassword(password)

	userCollection := getUserCollection()

	if userExists(username, userCollection) {
		fmt.Printf("User %s already exists.\n", username)
		return
	}

	err := userCollection.Insert(newuser)
	if err != nil {
		panic(err)
	} else {
		fmt.Printf("Successfully added user %s.\n", username)
	}
}

func loadUsersFromFile(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		rawentry := scanner.Text()
		userentry := strings.Split(rawentry, ",")
		username := userentry[0]
		password := userentry[1]
		addUser(username, password)
	}
}

func removeUser(username string) {
	userCollection := getUserCollection()

	if userExists(username, userCollection) {
		err := userCollection.Remove(bson.M{"username": username})
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("Successfully removed user %s.\n", username)
		}
	} else {
		fmt.Printf("User %s not found.\n", username)
	}
}

func removeAllUsers() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("YOU ARE ABOUT TO REMOVE ALL USERS. ENTER 'Y' TO CONTINUE: ")
	resp, _ := reader.ReadString('\n')
	if resp == "Y\n" {
		userCollection := getUserCollection()
		userCollection.DropCollection()
		fmt.Println("Successfully removed all users.")
	} else {
		fmt.Println("User removal cancelled.")
	}

}

func updatePassword(username, password string) {
	userCollection := getUserCollection()

	if userExists(username, userCollection) {
		updated := &models.User{
			Username: username,
		}
		updated.SetPassword(password)

		_, err := userCollection.Upsert(bson.M{"username": username}, updated)
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("Successfully updated password of user %s.\n", username)
		}
	} else {
		fmt.Printf("User %s not found.\n", username)
	}
}

func getUserCollection() *mgo.Collection {
	mongoHost := os.Getenv("MONGO_PORT_27017_TCP_ADDR")
	if mongoHost == "" {
		mongoHost = "localhost"
	}

	session, err := mgo.Dial(mongoHost)
	if err != nil {
		panic(err)
	}

	return session.DB("fhir").C("users")
}

func userExists(username string, collection *mgo.Collection) bool {
	count, err := collection.Find(bson.M{"username": username}).Count()
	if err != nil {
		panic(err)
	}
	return count > 0
}

func printUsage() {
	usageStatement := `Usage: command <arguments> (function)
	------
	add <username> <password> (add single user)
	add_file <filepath> (add users from comma separated file)
	remove <username> (remove single user)
	remove_all (remove all users)
	change_pass <username> <password> (change user's password)`

	fmt.Println(usageStatement)
}
