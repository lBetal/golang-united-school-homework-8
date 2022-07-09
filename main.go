package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
)

type Arguments map[string]string

type User struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func Perform(args Arguments, writer io.Writer) error {
	if _, ok := args["operation"]; !ok || args["operation"] == "" {
		return errors.New("-operation flag has to be specified")
	}
	fileName, ok := args["fileName"]
	if !ok || fileName == "" {
		return fmt.Errorf("-fileName flag has to be specified")
	}

	switch args["operation"] {
	case "list":
		users, err := readFile(fileName)
		if err != nil {
			return err
		}

		usersBytes, err := json.Marshal(users)
		if err != nil {
			return err
		}
		_, err = writer.Write(usersBytes)
		if err != nil {
			return err
		}

		return nil
	case "add":
		user := User{}
		item, ok := args["item"]
		if !ok || item == "" {
			return fmt.Errorf("-item flag has to be specified")
		}
		err := json.Unmarshal([]byte(item), &user)
		if err != nil {
			return err
		}

		users, err := readFile(fileName)
		if err != nil {
			return err
		}

		for _, val := range users {
			if val.Id == user.Id {
				errorText := fmt.Sprintf("Item with id %s already exists", user.Id)
				_, err = writer.Write([]byte(errorText))
				return err
			}
		}

		users = append(users, user)
		return writeFile(users, fileName)
	case "findById":
		id, ok := args["id"]
		if !ok || id == "" {
			return fmt.Errorf("-id flag has to be specified")
		}
		users, err := readFile(fileName)
		if err != nil {
			return err
		}

		for _, user := range users {
			if user.Id == id {
				userBytes, err := json.Marshal(user)
				if err != nil {
					return err
				}

				_, err = writer.Write(userBytes)
				return err
			}
		}

		return nil
	case "remove":
		id, ok := args["id"]
		if !ok || id == "" {
			return fmt.Errorf("-id flag has to be specified")
		}
		users, err := readFile(fileName)
		if err != nil {
			return err
		}

		for i, user := range users {
			if user.Id == id {
				users = append(users[:i], users[i+1:]...)

				return writeFile(users, fileName)
			}
		}

		errorText := fmt.Sprintf("Item with id %s not found", id)

		_, err = writer.Write([]byte(errorText))
		return err
	default:
		return fmt.Errorf("Operation %s not allowed!", args["operation"])
	}
}

func readFile(fileName string) ([]User, error) {
	users := []User{}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, err
	}
	file.Close()
	return users, nil
}

func writeFile(users []User, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	usersBody, err := json.Marshal(users)
	if err != nil {
		return err
	}
	file.Truncate(0)
	file.Seek(0, 0)
	_, err = file.Write(usersBody)
	file.Close()
	return err
}

func parseArgs() Arguments {
	args := Arguments{}

	return args
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
