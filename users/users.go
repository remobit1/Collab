package users

import (
	"database/sql"

	"github.com/lib/pq"
)

// GetUserData gets user data from database
func GetUserData(db *sql.DB, user string) []string {
	var followedCollabSpaces []string
	sqlQuery := `SELECT followed_collab_spaces FROM users WHERE username = $1;`

	row := db.QueryRow(sqlQuery, user)

	row.Scan(pq.Array(&followedCollabSpaces))
	return followedCollabSpaces
}
