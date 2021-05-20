package projects

import (
	"database/sql"
	"fmt"

	// for postgres sql driver
	"github.com/lib/pq"
)

// Project stores information on Collab projects
type Project struct {
	Name              string
	ProjectImgs       []string
	DefaultImageIndex int
	Tags              []string
	Description       string
	Followers         int64
}

func (p *Project) getDefaultImage(defaultImageIndex int) string {
	return p.ProjectImgs[defaultImageIndex]
}

// GetTrendingProjects returns the most recent projects, maximum 10.
func GetTrendingProjects(db *sql.DB) []Project {
	projects := []Project{}

	sqlQuery := `SELECT * FROM projects ORDER BY followers desc LIMIT 20;`
	rows, err := db.Query(sqlQuery)

	defer rows.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	for rows.Next() {
		project := Project{}
		rows.Scan(&project.Name, pq.Array(&project.ProjectImgs),
			&project.DefaultImageIndex, pq.Array(&project.Tags),
			&project.Description, &project.Followers)

		projects = append(projects, project)
	}

	return projects
}

// FilterProjects returs a slice of projects that have the collabSpace string as their first tag.
func FilterProjects(db *sql.DB, collabSpace string) []Project {
	projects := []Project{}
	fmt.Println(collabSpace)
	sqlQuery := `SELECT * FROM projects WHERE tags[1]=$1;`
	rows, err := db.Query(sqlQuery, collabSpace)

	defer rows.Close()

	if err != nil {
		fmt.Println(err.Error())
	}

	for rows.Next() {
		project := Project{}
		rows.Scan(&project.Name, pq.Array(&project.ProjectImgs),
			&project.DefaultImageIndex, pq.Array(&project.Tags),
			&project.Description, &project.Followers)

		projects = append(projects, project)
	}

	return projects
}
