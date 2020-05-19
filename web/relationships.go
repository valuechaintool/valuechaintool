package web

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/valuechaintool/valuechaintool/models"
)

// RelationshipsCreate renders the /relationships/new page
func RelationshipsCreate(w http.ResponseWriter, r *http.Request) {
	relationship := models.Relationship{}
	if uid, ok := r.URL.Query()["left_id"]; ok {
		id, err := uuid.Parse(uid[0])
		if err != nil {
			log.Println(err)
			if _, err := w.Write([]byte(err.Error())); err != nil {
				log.Println(err)
				return
			}
			return
		}
		relationship.LeftID = id
	}
	companies, err := models.ListCompanies(nil)
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	t := template.Must(template.ParseFiles("web/layout.html", "web/relationships-form.html"))
	d := struct {
		PageTitle    string
		Companies    []models.Company
		Relationship models.Relationship
	}{
		PageTitle:    "New relationship",
		Companies:    companies,
		Relationship: relationship,
	}
	if err := t.Execute(w, d); err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
}

// RelationshipsCreatePost parses the form from /relationships/new page
func RelationshipsCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	tier, err := strconv.Atoi(r.FormValue("tier"))
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	relationship := models.Relationship{
		LeftID:  uuid.MustParse(r.FormValue("left")),
		RightID: uuid.MustParse(r.FormValue("right")),
		Tier:    tier,
	}
	if err := models.NewRelationship(&relationship); err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	http.Redirect(w, r, "/companies/"+r.FormValue("left"), 302)
}

// RelationshipsDelete responds to /relationships/[ID]/delete url
func RelationshipsDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	relationship, err := models.GetRelationship(id)
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	leftId := relationship.LeftID
	if err := relationship.Delete(); err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	http.Redirect(w, r, "/companies/"+leftId.String(), 302)
}
