package main

import (
	"Collab/projects"
	"Collab/users"
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"

	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"

	pages "Collab/Pages"
	auth "Collab/auth"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "Livelife4love"
	dbname   = "collab"
)

var (
	templates = template.Must(template.ParseFiles("../templates/welcome-template.html",
		"../templates/collab-page.html", "../templates/practice.html", "../templates/sign-up.html",
		"../templates/success.html", "../templates/login.html"))
	validPath = regexp.MustCompile("^/(hot|home|signUp|login|success)/$")
	db        *sql.DB
	err       error
)

func initDB() {

	connStr := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

}

// renderTemplate takes a response object, the name of an html template and empty interface
// meant to be webpage structs of different types and renders them.
func renderTemplate(response http.ResponseWriter, tmpl string, page interface{}) {
	err := templates.ExecuteTemplate(response, tmpl+".html", page)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
}

func validateHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(response http.ResponseWriter, request *http.Request) {
		fmt.Println(request.URL.Path)
		validatedPath := validPath.FindStringSubmatch(request.URL.Path)
		if validatedPath == nil {
			http.NotFound(response, request)
			return
		}
		fn(response, request, validatedPath[1])
	}
}

// Go application entrypoint
func main() {
	//Instantiate a Welcome struct object and pass in some random information.
	//We shall get the name of the user as a query parameter from the URL

	practicePage := pages.Practice{}

	rtr := mux.NewRouter()

	//Our HTML comes with CSS that go needs to provide when we run the app. Here we tell go to create
	// a handle that looks in the static directory, go then uses the "/static/" as a url that our
	//html can refer to when looking for our css and other files.

	//rtr.Handle("/static/{rest}", //final url can be anything
	//http.StripPrefix("/static/",
	//http.FileServer(http.Dir("../static")))) //Go looks in the relative "static" directory first using http.FileServer(), then matches it to a
	//url of our choice as shown in http.Handle("/static/"). This url is what we need when referencing our css files
	//once the server begins. Our html code would therefore be <link rel="stylesheet"  href="/static/stylesheet/...">
	//It is important to note the url in http.Handle can be whatever we like, so long as we are consistent.

	rtr.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("../static/"))))

	rtr.PathPrefix("/templates/").Handler(
		http.StripPrefix("/templates/", http.FileServer(http.Dir("../templates"))))

	//This method takes in the URL path "/" and a function that takes in a response writer, and a http request.

	rtr.HandleFunc("/", redirectToHome)
	rtr.HandleFunc("/home/", validateHandler(homeHandler))
	rtr.HandleFunc("/practice", func(w http.ResponseWriter, r *http.Request) {
		renderTemplate(w, "practice", practicePage)
	})
	rtr.HandleFunc("/signUp/", validateHandler(signUpHandler))
	rtr.HandleFunc("/success/", validateHandler(success))
	rtr.HandleFunc("/login/", validateHandler(loginHandler))
	rtr.HandleFunc("/on/{id}", collabSpaceConstructor)

	initDB()

	//Start the web server, set the port to listen to 8080. Without a path it assumes localhost
	//Print any errors from starting the webserver using fmt
	fmt.Println("Listening")
	fmt.Println(http.ListenAndServe(":8080", rtr))
}

func loginHandler(response http.ResponseWriter, request *http.Request, title string) {
	if request.Method == "GET" {
		page := &pages.LoginPage{Title: title}
		renderTemplate(response, "login", page)
	} else {
		err := request.ParseForm()
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = auth.Authenticate(request.FormValue("password"), request.FormValue("username"), db)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Redirect(response, request, "/signUp/", http.StatusFound)
				return
			} else if err == bcrypt.ErrMismatchedHashAndPassword {
				http.Redirect(response, request, "/login/", http.StatusFound)
				response.WriteHeader(http.StatusUnauthorized)
				return
			}
			response.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.Redirect(response, request, "/success/", http.StatusFound)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request, title string) {
	hotPage := pages.CollabOnPage{}
	title = "welcome"
	p := projects.GetTrendingProjects(db)

	for _, project := range p {
		hotPage.Title = title
		hotPage.Projects = append(hotPage.Projects, project)
		hotPage.CollabSpacesFollowed = users.GetUserData(db, "remobit1")
	}

	renderTemplate(w, "collab-page", hotPage)
}

func redirectToHome(response http.ResponseWriter, request *http.Request) {
	http.Redirect(response, request, "/home/", http.StatusFound)
}

func signUpHandler(response http.ResponseWriter, request *http.Request, title string) {
	if request.Method == "GET" {
		page := &pages.LoginPage{Title: title}
		renderTemplate(response, "sign-up", page)
	} else {
		err := request.ParseForm()
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = auth.CreateNewUser(request.FormValue("password"), request.FormValue("username"), db); err != nil {
			response.WriteHeader(http.StatusInternalServerError)

			return
		}

		http.Redirect(response, request, "/login/", http.StatusFound)

	}
}

func collabSpaceConstructor(response http.ResponseWriter, request *http.Request) {
	collabSpace := pages.CollabOnPage{}
	vars := mux.Vars(request)
	collabScope := vars["id"]

	p := projects.FilterProjects(db, collabScope)
	

	for _, project := range p {

		collabSpace.Title = project.Name
		collabSpace.Projects = append(collabSpace.Projects, project)
		collabSpace.CollabSpacesFollowed = users.GetUserData(db, "remobit1")
	}

	renderTemplate(response, "collab-page", collabSpace)
}

func success(response http.ResponseWriter, request *http.Request, title string) {
	page := &pages.LoginPage{Title: title}
	renderTemplate(response, "success", page)
}
