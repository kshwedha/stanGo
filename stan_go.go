package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	c "stan_go/config"
	"strconv"
	"time"
)

type Address struct {
	Street  string `json:"street"`
	Suite   string `json:"suite"`
	City    string `json:"city"`
	Zipcode string `json:"zipcode"`
	Geo     struct {
		Lat string `json:"lat"`
		Lng string `json:"lng"`
	} `json:"geo"`
}

type User struct {
	Username string `json:"username"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Website  string `json:"website"`
	Address  Address
	Age      int
}

var Data map[string]User
var domain string

type Operation interface {
	create()
	update()
	print()
}

func (user User) print() {
	userValue := reflect.ValueOf(user)

	fmt.Println("...........................................")
	for i := 0; i < userValue.NumField(); i++ {
		field := userValue.Field(i)
		fmt.Printf("%s: %v\n", userValue.Type().Field(i).Name, field.Interface())
	}
	fmt.Println("___________________________________________")
}

func printData(data map[string]User) {
	for _, udata := range data {
		udata.print()
	}
}

func (u *User) update() {
	u.Age = rand.Intn(100)
	u.Email = u.Name + "-" + strconv.Itoa(u.Age) + "@" + domain + ".com"
}

func fetchUserRecord() string {
	if len(Data) == 0 {
		return ""
	}
	keys := make([]string, 0, len(Data))
	for key := range Data {
		keys = append(keys, key)
	}
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(keys))
	return keys[index]
}

func doUpdateOperation() {
	user := Data[fetchUserRecord()]
	fmt.Printf("[*] Before user update for %s.\n", user.Username)
	user.print()
	user.update()
	Data[user.Username] = user
	fmt.Printf("[*] After updating the user - %s.\n", user.Username)
	Data[user.Username].print()
}

func doDeleteOperation() {
	user := fetchUserRecord()
	fmt.Printf("[*] No. of users before deleting record %d.\n", len(Data))
	delete(Data, user)
	fmt.Printf("[*] No. of users after deleting record %d.\n", len(Data))
}

func CreateUserData() []User {
	url := "https://jsonplaceholder.typicode.com/users"
	client, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	body, err := io.ReadAll(client.Body)
	if err != nil {
		panic(err)
	}
	defer client.Body.Close()
	var user []User
	err = json.Unmarshal(body, &user)
	if err != nil {
		fmt.Println("Error unmarshaling JSON:", err)
		return []User{}
	}
	return user
}

func prepareDB() {
	fmt.Println("[*] Initialising sqllite-3 db.")
	db := c.SqlCursor()
	defer db.Close()
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		name TEXT,
		age INTEGER,
		email TEXT
	)`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("[*] created the db table for User records insertion.")
}

func main() {
	fmt.Println("[*] D you wanna start(only y/n)?")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	if len(input) > 1 || input != "y" {
		panic("!! You entered a wrong input.")
	}
	domain = "stan"
	Data = make(map[string]User)
	prepareDB()
	go func() { fmt.Println("[*] `Data` is global variable(instruction executed by concurrency).") }()
	fmt.Println("[*] Creating the User data(creation operation).")
	for _, data := range CreateUserData() {
		Data[data.Username] = data
	}
	fmt.Printf("[*] datatype of `Data` is %s", reflect.TypeOf(Data))
	fmt.Printf("[*] len of Data is %d\n", len(Data))
	fmt.Println("[*] Printing each user data(Reading operation).")
	go func() { printData(Data) }()
	time.Sleep(10000000)
	fmt.Println("[*] On update operation.")
	doUpdateOperation()
	fmt.Println("[*] On delete operation.")
	doDeleteOperation()
	time.Sleep(10000000)
	fmt.Println("[*] Finished CRUD operation")
}
