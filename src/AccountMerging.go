package main

import (
	"encoding/json"
	"fmt"
	"github.com/ihebu/dsu"
	"io"
	"os"
)

// Account represents a single account
type Account struct {
	Application json.Number `json:"application"`
	Emails      []string    `json:"emails"`
	Name        string      `json:"name"`
}

// Person represents an individual
type Person struct {
	Applications []string `json:"applications"`
	Emails       []string `json:"emails"`
	Name         string   `json:"name"`
}

func (p Person) mergePeople(a Account) Person {
	p.Name = a.Name
	p.Applications = append(p.Applications, string(a.Application))

	uniqueEmail := make(map[string]bool)

	for _, email := range p.Emails {
		uniqueEmail[email] = true
	}
	for _, email := range a.Emails {
		uniqueEmail[email] = true
	}

	temp := make([]string, 0, len(uniqueEmail))

	for key := range uniqueEmail {
		temp = append(temp, key)
	}
	p.Emails = temp

	return p
}

func mergeAccounts(accounts []Account) []Person {
	d := dsu.New()
	emailGroup := make(map[string]int)

	for i, account := range accounts {
		d.Add(i)
		for _, email := range account.Emails {
			val, exists := emailGroup[email]
			if exists {
				d.Union(i, val)
			}
			emailGroup[email] = i
		}
	}

	personGroup := make(map[int]Person)

	for i, account := range accounts {
		rep := d.Find(i)
		_, ok := personGroup[rep.(int)]
		if !ok {
			personGroup[rep.(int)] = Person{}
		}
		personGroup[rep.(int)] = personGroup[rep.(int)].mergePeople(account)
	}

	var people []Person
	for i := range personGroup {
		people = append(people, personGroup[i])
	}

	return people
}

func readFromFile(filePath string) ([]Account, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	raw, fErr := io.ReadAll(file)
	if fErr != nil {
		return nil, fErr
	}

	var accounts []Account

	mErr := json.Unmarshal(raw, &accounts)
	if mErr != nil {
		return nil, mErr
	}

	return accounts, nil
}

func main() {
	filePath := os.Args[1]

	accounts, err := readFromFile(filePath)
	if err != nil {
		fmt.Println("error reading from file at path")
		panic(err)
	}

	people := mergeAccounts(accounts)

	output, err := json.Marshal(people)
	if err != nil {
		fmt.Println("error marshalling output")
		panic(err)
	}

	fmt.Println(string(output))
}
