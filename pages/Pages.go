package pages

import (
	"Collab/projects"
)

// Welcome struct represents the structure of the welcome page
type Welcome struct {
	Name string
	Time string
}

// CollabOnPage is the landing page of Collab
type CollabOnPage struct {
	Title                string
	Projects             []projects.Project
	CollabSpacesFollowed []string
}

// LoginPage allows users to login
type LoginPage struct {
	Title string
}

// Practice is for practice
type Practice struct {
}
