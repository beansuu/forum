package models

// define the structure of a user
type User struct {
}

*CreateUser(user User) error


// retrieves a user based on their email
func GetUserByEmail(email string)(*User,error)