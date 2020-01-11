package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type user struct {
	email    string
	password string
}

func checkDump(user string, domain string) string {
	postParam := url.Values{
		"luser":      {user},
		"domain":     {domain},
		"luseropr":   {"0"},
		"domainopr":  {"0"},
		"submitform": {"em"},
	}
	resp, err := http.PostForm("https://pwndb2am4tzkvold.onion.ws/", postParam)
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

func parseDump(rawData string) (users []user) {
	leaks := strings.Split(rawData, "Array")[2:]
	for _, leak := range leaks {
		username := strings.Split(strings.Split(leak, "[luser] => ")[1], "\n")[0]
		domain := strings.Split(strings.Split(leak, "[domain] => ")[1], "\n")[0]
		email := username + "@" + domain
		password := strings.Split(strings.Split(leak, "[password] => ")[1], "\n")[0]
		users = append(users, user{email, password})
	}
	return
}

func main() {
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
	var user, domain string
	flag.StringVar(&user, "user", "", "Username")
	flag.StringVar(&domain, "domain", "", "Domain to check")
	flag.Parse()

	if domain == "" && user == "" {
		flag.Usage()
		fmt.Println()
		log.Fatalln("Please enter a domain or user to check")
	}

	body := checkDump(user, domain)
	rawData := strings.Split(strings.Split(body, "<pre>\n")[1], "</pre>")[0]
	dump := parseDump(rawData)
	if len(dump) < 1 {
		log.Fatalln("No data found")
	}
	fmt.Println("[Users]")
	for _, user := range dump {
		fmt.Println(user.email + ":" + user.password)
	}
}
