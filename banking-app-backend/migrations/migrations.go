package migrations

import (
	"github.com/go-bank-backend/helpers"
	"github.com/go-bank-backend/interfaces"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type User struct{
	gorm.Model
	Username   string
	Email      string
	Password   string
}

type Account struct{
	gorm.Model
	Type      string
	Name      string
	Balance   uint
	UserID    uint
}

// Create a slice with two users
// Start a loop that will iterate through our dataset
// Insert the user into the database
// Build a data-structure for the user, and pass that data into the db.Create function
func createAccounts() {
    db := helpers.ConnectDB()

    users := &[2]interfaces.User{
        {Username: "Martin", Email: "martin@martin.com"},
        {Username: "Michael", Email: "michael@michael.com"},
    }

    for i := 0; i < len(users); i++ {
        // Correct one way
        generatedPassword := helpers.HashAndSalt([]byte(users[i].Username))
        user := &interfaces.User{Username: users[i].Username, Email: users[i].Email, Password: generatedPassword}
        db.Create(&user)

        account := &interfaces.Account{Type: "Daily Account", Name: string(users[i].Username + "'s" + " account"), Balance: uint(10000 * int(i+1)), UserID: user.ID}
        db.Create(&account)
    }
    defer db.Close()
}


func Migrate() {
    User := &interfaces.User{}
    Account := &interfaces.Account{}
    db := helpers.ConnectDB()
    db.AutoMigrate(&User, &Account)
    defer db.Close()
    
    createAccounts()
}
