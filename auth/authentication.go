package authentication

import (
	"database/sql"

	// for postgres database
	_ "github.com/lib/pq"

	"golang.org/x/crypto/bcrypt"
)

//Credentials are security credentials
type Credentials struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

// Authenticate processes logins
func Authenticate(password string, username string, db *sql.DB) (err error) {
	creds := &Credentials{Password: password, Username: username}
	sqlQuery := `select password from users where username=$1;`

	result := db.QueryRow(sqlQuery, creds.Username)

	storedCreds := &Credentials{}

	err = result.Scan(&storedCreds.Password)
	if err != nil {
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
		return err
	}

	return nil
}

// CreateNewUser creates new user in database
func CreateNewUser(password string, username string, db *sql.DB) (err error) {

	creds := &Credentials{Password: password, Username: username}
	sqlQuery := `insert into users (username, password) values ($1, $2);`

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)

	if _, err := db.Query(sqlQuery, creds.Username, string(hashedPassword)); err != nil {
		return err
	}

	return nil
}
