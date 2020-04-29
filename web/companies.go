package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/valuechaintool/valuechaintool/models"
)

// CompaniesCreate renders the /companies/new page
func CompaniesCreate(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("web/layout.html", "web/companies-form.html"))
	d := struct {
		PageTitle string
		Company   models.Company
	}{
		PageTitle: "Add Company",
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

// CompaniesCreatePost parses the form from /companies/new page
func CompaniesCreatePost(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	company := models.Company{
		Name:             r.FormValue("name"),
		Industry:         r.FormValue("industry"),
		Country:          r.FormValue("country"),
		SalesManager:     r.FormValue("salesManager"),
		TechnicalManager: r.FormValue("technicalManager"),
		Notes:            r.FormValue("notes"),
	}
	switch t := r.FormValue("type"); t {
	case "partner":
		company.Type = models.PartnerCompany
	case "customer":
		company.Type = models.CustomerCompany
	}
	if err := models.NewCompany(&company); err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	http.Redirect(w, r, "/companies", 302)
}

// CompaniesList renders the /companies page
func CompaniesList(w http.ResponseWriter, r *http.Request) {
	companies, err := models.ListCompanies(nil)
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	t := template.Must(template.ParseFiles("web/layout.html", "web/companies-list.html"))
	d := struct {
		PageTitle string
		Companies []models.Company
	}{
		PageTitle: "HomePage",
		Companies: companies,
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

// CompaniesRead renders the /companies/[ID] page
func CompaniesRead(w http.ResponseWriter, r *http.Request) {
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
	company, err := models.GetCompany(id)
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	if err := company.EagerLoad(); err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	t := template.Must(template.ParseFiles("web/layout.html", "web/companies-single.html"))
	d := struct {
		PageTitle string
		Company   models.Company
	}{
		PageTitle: fmt.Sprintf("Company %s information", (*company).Name),
		Company:   *company,
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

// CompaniesUpdate renders the /companies/[ID]/edit page
func CompaniesUpdate(w http.ResponseWriter, r *http.Request) {
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
	company, err := models.GetCompany(id)
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	t := template.Must(template.ParseFiles("web/layout.html", "web/companies-form.html"))
	d := struct {
		PageTitle string
		Company   models.Company
	}{
		PageTitle: "Add Company",
		Company:   *company,
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

// CompaniesUpdatePost parses the form from /companies/[ID]/edit page
func CompaniesUpdatePost(w http.ResponseWriter, r *http.Request) {
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
	company, err := models.GetCompany(id)
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}

	// Fill the data
	company.Name = r.FormValue("name")
	company.Industry = r.FormValue("industry")
	switch t := r.FormValue("type"); t {
	case "partner":
		company.Type = models.PartnerCompany
	case "customer":
		company.Type = models.CustomerCompany
	}
	company.Country = r.FormValue("country")
	company.SalesManager = r.FormValue("salesManager")
	company.TechnicalManager = r.FormValue("technicalManager")
	company.Notes = r.FormValue("notes")

	if err := company.Save(); err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	http.Redirect(w, r, "/companies", 302)
}

// CompaniesDelete responds to /companies/[ID]/delete url
func CompaniesDelete(w http.ResponseWriter, r *http.Request) {
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
	company, err := models.GetCompany(id)
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	if err := company.Delete(); err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	http.Redirect(w, r, "/companies", 302)
}
