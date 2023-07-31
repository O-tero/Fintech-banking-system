package users

import (
	"time"

	"github.com/go-bank-backend/helpers"
	"github.com/go-bank-backend/interfaces"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

// Refactor prepareToken
func prepareToken(user *interfaces.User) string {
	tokenContent := jwt.MapClaims{
		"user.id": user.ID,
		"expiry":  time.Now().Add(time.Minute * 60).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tokenContent)
	token, err := jwtToken.SignedString([]byte("TokenPassword"))
	helpers.HandleErr(err)

	return token
}

// Refactor prepareResponse
func prepareResponse(user *interfaces.User, accounts []interfaces.ResponseAccount) map[string]interface{} {
	responseuser := &interfaces.ResponseUser{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Accounts: accounts,
	}
	var token = prepareToken(user)
	var response = map[string]interface{}{"message": "all is fine"}
	response["jwt"] = token
	response["data"] = responseuser

	return response
}

func Login(username string, pass string) map[string]interface{} {
	// Add validation to login
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: pass, Valid: "password"},
		})
	if valid {
		// connect DB
		db := helpers.ConnectDB()
		user := &interfaces.User{}
		if db.Where("username = ? ", username).First(&user).RecordNotFound() {
			return map[string]interface{}{"message": "User not found"}
		}
		// verify password
		passErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass))

		if passErr == bcrypt.ErrMismatchedHashAndPassword && passErr != nil {
			return map[string]interface{}{"message": "Wrong password"}
		}
		// Find accounts for user
		accounts := []interfaces.ResponseAccount{}
		db.Table("accounts").Select("id, name, balance").Where("user_id = ?", user.ID).Scan(&accounts)

		defer db.Close()

		var response = prepareResponse(user, accounts)

		return response
	} else {
		return map[string]interface{}{"message": "not valid values"}
	}
}

// Register function
func Register(username string, email string, pass string) map[string]interface{} {
	// Add validation to the register
	valid := helpers.Validation(
		[]interfaces.Validation{
			{Value: username, Valid: "username"},
			{Value: email, Valid: "email"},
			{Value: pass, Valid: "password"},
		})
	if valid {
		//Create registration logic
		// Connect DB
		db := helpers.ConnectDB()
		// Create user and account
		generatedPassword := helpers.HashAndSalt([]byte(pass))
		user := &interfaces.User{Username: username, Email: email, Password: generatedPassword}
		db.Create(&user)

		account := &interfaces.Account{Type: "Daily Account", Name: string(username + "'s" + " account "), Balance: 0, UserID: user.ID}
		db.Create(&user)

		// Prepare response and token
		defer db.Close()
		accounts := []interfaces.ResponseAccount{}
		respAccount := interfaces.ResponseAccount{ID: account.ID, Name: account.Name, Balance: int(account.Balance)}
		var response = prepareResponse(user, accounts)

		return response

	} else {
		return map[string]interface{}{"message": "not valid values"}
	}
}
