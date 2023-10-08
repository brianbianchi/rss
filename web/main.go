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
	Code  string
	Email string
	Subs  map[data.Feed]bool
	Error string
}

func main() {
	db := data.InitDb()
	defer db.Close()

	http.HandleFunc("/success/", func(w http.ResponseWriter, r *http.Request) {
		code := strings.TrimPrefix(r.URL.Path, "/success/")
		if code == "" {
			fmt.Fprintln(w, http.StatusNotFound)
		}

		path := util.GetRootPath()
		t, err := template.ParseFiles(fmt.Sprint(path, "web/success.html"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		t.Execute(w, code)
	})

	http.HandleFunc("/unsubscribe/", func(w http.ResponseWriter, r *http.Request) {
		code := strings.TrimPrefix(r.URL.Path, "/unsubscribe/")
		data.DeleteSubs(db, code)
		data.DeleteUser(db, code)
		fmt.Fprint(w, "unsubbed")
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		code := strings.TrimPrefix(r.URL.Path, "/")
		if r.Method == "GET" {
			serveForm(w, db, code, "")
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

			if code == "" {
				existingCode, _ := data.GetUserByEmail(db, email)
				if existingCode != "" {
					serveForm(w, db, "", "Your email is already registered.")
					return
				}
				newUserCode, err := data.CreateUser(db, email)
				if err != nil {
					serveForm(w, db, "", "Failed to create a new user.")
					return
				}
				err = data.CreateSubs(db, links, newUserCode)
				if err != nil {
					serveForm(w, db, "", "Failed to create new subs.")
					return
				}
				http.Redirect(w, r, fmt.Sprint("/success/", newUserCode), http.StatusSeeOther)
			}
			if code != "" {
				email, _ := data.GetUserByCode(db, code)
				if email == "" {
					serveForm(w, db, code, "Your link isn't correct.")
					return
				}
				err = data.DeleteSubs(db, code)
				if err != nil {
					serveForm(w, db, code, "Failed to delete subs.")
					return
				}
				err = data.CreateSubs(db, links, code)
				if err != nil {
					serveForm(w, db, code, "Failed to create subs.")
					return
				}
				http.Redirect(w, r, fmt.Sprint("/success/", code), http.StatusSeeOther)
			}
		}
	})

	port := os.Getenv("PORT")
	fmt.Println("Serving on port", port)
	http.ListenAndServe(fmt.Sprint(":", port), nil)
}

func serveForm(w http.ResponseWriter, db *sql.DB, code string, errorDisplay string) {
	path := util.GetRootPath()
	template, err := template.ParseFiles(fmt.Sprint(path, "web/form.html"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pageData, err := getPageData(db, code, errorDisplay)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = template.Execute(w, pageData)
	if err != nil {
		fmt.Println(err)
	}

}

func getPageData(db *sql.DB, code string, errorDisplay string) (PageData, error) {
	pageData := PageData{}
	feeds, err := data.GetFeeds(db)
	if err != nil {
		return pageData, err
	}
	subs := make(map[data.Feed]bool)
	for _, f := range feeds {
		subs[f] = false
	}
	if code != "" {
		pageData.Code = code
		email, err := data.GetUserByCode(db, code)
		if err != nil {
			return pageData, err
		}
		pageData.Email = email
		userSub, err := data.GetSubs(db, code)
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
	pageData.Error = errorDisplay
	return pageData, nil
}
