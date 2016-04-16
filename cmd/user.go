package cmd

import (
	"fmt"
	"github.com/codegangsta/cli"
	"os"
	"rakewire/model"
)

// UserList lists users in the system.
func UserList(c *cli.Context) {

	db, _, err := initConfig(c)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}
	defer closeDatabase(db)

	var users model.Users
	db.Select(func(tx model.Transaction) error {
		users = model.U.Range(tx)
		return nil
	})

	for _, user := range users {
		fmt.Printf("%s: %s\n", user.ID, user.Username)
	}

}
