package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// User is the struct that contains leaked credentials
type User struct {
	email    string
	password string
}

// CheckDump checks pwndb for leaked credentials
func CheckDump(user string, domain string) string {
	postParam := url.Values{
		"luser":      {user},
		"domain":     {domain},
		"luseropr":   {"0"},
		"domainopr":  {"0"},
		"submitform": {"em"},
	}
	client := http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.PostForm("https://pwndb2am4tzkvold.onion.ws/", postParam)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	return string(body)
}

// ParseDump parses the data returned from CheckDump
func ParseDump(rawData string) (users []User) {
	leaks := strings.Split(rawData, "Array")[2:]
	for _, leak := range leaks {
		username := strings.Split(strings.Split(leak, "[luser] => ")[1], "\n")[0]
		domain := strings.Split(strings.Split(leak, "[domain] => ")[1], "\n")[0]
		email := username + "@" + domain
		password := strings.Split(strings.Split(leak, "[password] => ")[1], "\n")[0]
		users = append(users, User{email, password})
	}
	return
}

func init() {
	log.Println()
	fmt.Println(`                                         /$$ /$$      
                                        | $$| $$      
  /$$$$$$  /$$  /$$  /$$ /$$$$$$$   /$$$$$$$| $$$$$$$ 
 /$$__  $$| $$ | $$ | $$| $$__  $$ /$$__  $$| $$__  $$
| $$  \ $$| $$ | $$ | $$| $$  \ $$| $$  | $$| $$  \ $$
| $$  | $$| $$ | $$ | $$| $$  | $$| $$  | $$| $$  | $$
| $$$$$$$/|  $$$$$/$$$$/| $$  | $$|  $$$$$$$| $$$$$$$/
| $$____/  \_____/\___/ |__/  |__/ \_______/|_______/ 
| $$                                                  
| $$                                                  
|__/                                                  ` + "\n")
}

func main() {
	var user, domain, totalString string
	flag.StringVar(&user, "user", "", "Username")
	flag.StringVar(&domain, "domain", "", "Domain to check")
	flag.Parse()

	if domain == "" && user == "" {
		flag.Usage()
		fmt.Println()
		log.Fatalln("Please enter a domain or user to check")
	}

	body := CheckDump(user, domain)
	// Errors if data not returned correctly
	if !strings.Contains(body, "<pre>") {
		log.Fatalln("Error contacting pwndb")
	}
	rawData := strings.Split(strings.Split(body, "<pre>\n")[1], "</pre>")[0]
	dump := ParseDump(rawData)
	if len(dump) < 1 {
		log.Fatalln("No data found")
	} else if len(dump) == 1 {
		totalString = "[1 User]\n"
	} else {
		totalString = fmt.Sprintf("[%d Users]\n", len(dump))
	}
	fmt.Printf(totalString)
	for _, user := range dump {
		fmt.Println(user.email + ":" + user.password)
	}
}
