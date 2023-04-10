package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Email    string
	Password string
}
type Task struct {
	ID       int
	UserId   int
	Title    string
	DueDate  string
	Category int
	IsDone   bool
}
type Category struct {
	ID     int
	Title  string
	UserID int
	Color  string
}

var userStorage []User
var authenticatedUser *User
var taskStorage []Task
var categoryStorage []Category

const userStoragePath = "users.txt"

func main() {

	fmt.Println("This is a simple todo project for practice")
	command := flag.String("command", "no-command", "a command flag to run \n"+
		"--command=register-user for register a new user\n"+
		"--command=create-task for create a new task\n"+
		"--command=create-category for create a new category\n"+
		"--command=login for login if you have already an account\n"+
		"--command=task-list for list of your tasks\n"+
		"--command=exit for exit from app command loop")
	flag.Parse()
	loadUsers()
	for {
		runCommand(*command)

		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Please inter your next command:")
		scanner.Scan()
		*command = scanner.Text()
	}

}
func runCommand(command string) {
	/*if command != "register-user" && command != "exit" && authenticatedUser == nil {
		login()
	}*/
	/*if authenticatedUser == nil {
		return
	}*/
	switch command {
	case "login":
		login()
	case "create-category":
		createCategory()
	case "create-task":
		createTask()
	case "register-user":
		registerUser()
	case "task-list":
		listOfTasks()
	case "exit":
		os.Exit(0)
	case "no-command":
		fmt.Println("your command is empty, please run: ./todo -h to help")
		os.Exit(0)
	default:
		fmt.Println("Command is not valid:", command)

	}
}
func createTask() {
	if authenticatedUser != nil {
		authenticatedUser.Print()
	}
	var name, deadlineTime, category string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter your task name:")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("Please enter your task category:")
	scanner.Scan()
	category = scanner.Text()
	categoryID, err := strconv.Atoi(category)
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}
	isFound := false
	for _, c := range categoryStorage {
		if c.ID == categoryID && c.UserID == authenticatedUser.ID {
			isFound = true
			break
		}
	}
	if !isFound {
		fmt.Println("Sorry we cant create a task because you haven't a valid category")
		return
	}
	fmt.Println("Please enter your task deadlineTime:")
	scanner.Scan()
	deadlineTime = scanner.Text()

	task := Task{
		ID:       len(taskStorage) + 1,
		Title:    name,
		Category: categoryID,
		DueDate:  deadlineTime,
		IsDone:   false,
		UserId:   authenticatedUser.ID,
	}
	taskStorage = append(taskStorage, task)
	fmt.Println("task:", name, category, deadlineTime)

}
func registerUser() {
	var email, password, name string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Please enter your name:")
	scanner.Scan()
	name = scanner.Text()

	fmt.Println("Please enter your email:")
	scanner.Scan()
	email = scanner.Text()

	fmt.Println("Please enter your password:")
	scanner.Scan()
	password = scanner.Text()

	user := User{
		ID:       len(userStorage) + 1,
		Name:     name,
		Email:    email,
		Password: password,
	}

	userStorage = append(userStorage, user)
	writeUserToFile(user)
	fmt.Printf("userStorage: %+v\n", userStorage)

}
func createCategory() {
	var title, color string
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Please enter your category title:")
	scanner.Scan()
	title = scanner.Text()

	fmt.Println("Please enter your category color:")
	scanner.Scan()
	color = scanner.Text()

	c := Category{
		ID:     len(categoryStorage) + 1,
		Title:  title,
		Color:  color,
		UserID: authenticatedUser.ID,
	}
	categoryStorage = append(categoryStorage, c)
}
func login() {
	/* get the email and password from user
	if they were correct then user can log in
	*/
	fmt.Println("You must log in first.")
	scn := bufio.NewScanner(os.Stdin)
	fmt.Println("Please inter your email:")
	scn.Scan()
	email := scn.Text()

	fmt.Println("Please inter your password:")
	scn.Scan()
	password := scn.Text()
	userExist := false
	for _, user := range userStorage {
		if user.Email == email && user.Password == password {
			userExist = true
			authenticatedUser = &user
			fmt.Println("you're logged in.")
			break
		}
	}
	if !userExist {
		fmt.Println("Username or password is incorrect.")
	}
}

func (u User) Print() {
	fmt.Println("User", u.ID, u.Email)
}

func listOfTasks() {
	for _, task := range taskStorage {
		if task.UserId == authenticatedUser.ID {
			fmt.Println(task)
		}
	}
}
func loadUsers() {
	file, err := os.OpenFile(userStoragePath, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println("Error(Open file:)", err)
	}

	data := make([]byte, 1024)

	_, oErr := file.Read(data)
	if oErr != nil {
		fmt.Println("Error(Read file:)", oErr)
	}

	dataStr := string(data)
	userSlice := strings.Split(dataStr, "\n")
	for _, u := range userSlice {
		if u[0] != '{' && u[len(u)-1] != '}' {
			continue
		}
		userStruct := User{}
		uError := json.Unmarshal([]byte(u), &userStruct)
		if uError != nil {
			fmt.Println("Unmarshalling error:", uError)
		}
	}

}
func writeUserToFile(user User) {
	var file *os.File
	file, err := os.OpenFile(userStoragePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err)

		return
	}
	var data []byte
	data = []byte(fmt.Sprintf("id: %d, name: %s, email: %s, password: %s\n ",
		user.ID, user.Name, user.Email, user.Password))

	var jErr error

	data, jErr = json.Marshal(user)
	if jErr != nil {
		fmt.Println("Marshaling Error:", jErr)
	}

	_, wErr := file.Write(data)
	if wErr != nil {
		fmt.Println("Error:", wErr)
	}
	cErr := file.Close()
	if cErr != nil {
		fmt.Println("Error:", cErr)
	}

}
