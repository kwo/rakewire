package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/howeyc/gopass"
	"os"
	"rakewire/model"
)

// UserAdd adds a user
func UserAdd(c *cli.Context) {

	var username string
	var rolestr string
	if c.NArg() == 2 {
		username = c.Args().First()
		rolestr = c.Args()[1]
	} else {
		cli.ShowCommandHelp(c, c.Command.Name)
		os.Exit(1)
	}

	db, err := initCmd(c)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	defer closeDatabase(db)

	var user *model.User
	if err := db.Select(func(tx model.Transaction) error {
		user = model.U.GetByUsername(tx, username)
		return nil
	}); err != nil {
		fmt.Printf("Error retrieving user: %s\n", err.Error())
		os.Exit(1)
	}

	if user != nil {
		fmt.Printf("Username already exists. Cannot add new user %s.\n", username)
		os.Exit(1)
	}

	var password string
	fmt.Printf("password: ")
	if pass, err := gopass.GetPasswd(); err == nil {
		password = string(pass)
	} else {
		fmt.Printf("Cannot read password: %s\n", err.Error())
		os.Exit(1)
	}

	var password2 string
	fmt.Printf("confirm password: ")
	if pass, err := gopass.GetPasswd(); err == nil {
		password2 = string(pass)
	} else {
		fmt.Printf("Cannot read password: %s\n", err.Error())
		os.Exit(1)
	}

	if password != password2 {
		fmt.Println("Passwords do not match.")
		os.Exit(1)
	}

	if err := db.Update(func(tx model.Transaction) error {
		user := model.U.New(username, password)
		user.SetRoles(rolestr)
		return model.U.Save(tx, user)
	}); err != nil {
		fmt.Printf("Error adding user: %s\n", err.Error())
		os.Exit(1)
	}

	fmt.Printf("User added: %s\n", username)

}
