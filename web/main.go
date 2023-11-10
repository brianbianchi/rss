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

		if r.Method == "GET" {
			serveForm(w, db, pageData)
			return
		}

		if r.Method == "POST" {
			err := r.ParseForm()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			email := r.FormValue("email")
			links := r.Form["links"]

			if urlCode == "" {
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
			if urlCode != "" {
				email, _ := data.GetUserByCode(db, urlCode)
				if email == "" {
					pageData.Error = "Your link isn't correct."
					serveForm(w, db, pageData)
					return
				}
				err = data.DeleteSubs(db, urlCode)
				if err != nil {
					pageData.Error = "Failed to delete subs."
					serveForm(w, db, pageData)
					return
				}
				err = data.CreateSubs(db, links, urlCode)
				if err != nil {
					pageData.Error = "Failed to create subs."
					serveForm(w, db, pageData)
					return
				}
				pageData.SuccessUpdated = true
				serveForm(w, db, pageData)
			}
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
		fmt.Println(err)
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
