package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/brianbianchi/rss/data"
	"github.com/brianbianchi/rss/util"
)

type PageData struct {
	Code           string
	Email          string
	Subs           map[data.Feed]bool
	SuccessCreated bool
	SuccessUpdated bool
	Error          string
}

func main() {
	db := data.InitDb()
	defer db.Close()

	http.HandleFunc("/unsubscribe/", unsubscribeHandler(db))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		urlCode := strings.TrimPrefix(r.URL.Path, "/")
		pageData := PageData{Code: urlCode}

		switch r.Method {
		case http.MethodGet:
			serveForm(w, db, pageData)

		case http.MethodPost:
			handleFormSubmission(w, r, db, pageData)

		default:
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		}
	})

	port := os.Getenv("PORT")
	fmt.Println("Serving on port", port)
	http.ListenAndServe(fmt.Sprint(":", port), nil)
}

func unsubscribeHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := strings.TrimPrefix(r.URL.Path, "/unsubscribe/")
		data.DeleteSubs(db, code)
		data.DeleteUser(db, code)
		fmt.Fprint(w, "We deleted your data. You are now unsubscribed. We regret seeing you leave.")
	}
}

func handleFormSubmission(w http.ResponseWriter, r *http.Request, db *sql.DB, pageData PageData) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := r.FormValue("email")
	links := r.Form["links"]

	if pageData.Code == "" {
		handleNewUser(w, db, email, links, pageData)
	} else {
		handleExistingUser(w, db, links, pageData)
	}
}

func handleNewUser(w http.ResponseWriter, db *sql.DB, email string, links []string, pageData PageData) {
	existingCode, _ := data.GetUserByEmail(db, email)
	if existingCode != "" {
		pageData.Error = "Your email is already registered."
		serveForm(w, db, pageData)
		return
	}

	newUserCode, err := data.CreateUser(db, email)
	if err != nil {
		pageData.Error = "Failed to create a new user."
		serveForm(w, db, pageData)
		return
	}

	err = data.CreateSubs(db, links, newUserCode)
	if err != nil {
		pageData.Error = "Failed to create new subs."
		serveForm(w, db, pageData)
		return
	}

	pageData.Code = newUserCode
	pageData.SuccessCreated = true
	serveForm(w, db, pageData)
}

func handleExistingUser(w http.ResponseWriter, db *sql.DB, links []string, pageData PageData) {
	email, _ := data.GetUserByCode(db, pageData.Code)
	if email == "" {
		pageData.Error = "Your link isn't correct."
		serveForm(w, db, pageData)
		return
	}

	err := data.DeleteSubs(db, pageData.Code)
	if err != nil {
		pageData.Error = "Failed to delete subs."
		serveForm(w, db, pageData)
		return
	}

	err = data.CreateSubs(db, links, pageData.Code)
	if err != nil {
		pageData.Error = "Failed to create subs."
		serveForm(w, db, pageData)
		return
	}

	pageData.SuccessUpdated = true
	serveForm(w, db, pageData)
}

func serveForm(w http.ResponseWriter, db *sql.DB, pageData PageData) {
	path := util.GetRootPath()
	template, err := template.ParseFiles(fmt.Sprint(path, "web/index.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pageData, err = getPageData(db, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = template.Execute(w, pageData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getPageData(db *sql.DB, pageData PageData) (PageData, error) {
	feeds, err := data.GetFeeds(db)
	if err != nil {
		return pageData, err
	}

	subs := make(map[data.Feed]bool)
	for _, f := range feeds {
		subs[f] = false
	}

	if pageData.Code != "" {
		email, err := data.GetUserByCode(db, pageData.Code)
		if err != nil {
			return pageData, err
		}

		pageData.Email = email
		userSub, err := data.GetSubs(db, pageData.Code)
		if err != nil {
			return pageData, err
		}

		for _, link := range userSub {
			for key := range subs {
				if key.Link == link {
					subs[key] = true
				}
			}
		}
	}
	pageData.Subs = subs
	return pageData, nil
}
