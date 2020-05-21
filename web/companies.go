package web

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sort"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/valuechaintool/valuechaintool/models"
)

// CompaniesCreate renders the /companies/new page
func CompaniesCreate(w http.ResponseWriter, r *http.Request) {
	cts, err := models.ListCompanyTypes(nil)
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	scs, err := models.ListSectors(nil)
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	d := struct {
		PageTitle    string
		Company      models.Company
		CompanyTypes []models.CompanyType
		Sectors      []models.Sector
	}{
		PageTitle:    "Add Company",
		CompanyTypes: cts,
		Sectors:      scs,
	}
	funcMap := template.FuncMap{
		"uts": func(u uuid.UUID) string {
			return u.String()
		},
	}
	t, err := template.Must(template.ParseFiles("web/layout.html")).Funcs(funcMap).ParseFiles("web/companies-form.html")
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
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
	sectorID, err := uuid.Parse(r.FormValue("sector"))
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	typeID, err := uuid.Parse(r.FormValue("type"))
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	company := models.Company{
		Name:     r.FormValue("name"),
		Country:  r.FormValue("country"),
		SectorID: sectorID,
		TypeID:   typeID,
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
	cts, err := models.ListCompanyTypes(nil)
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	relationships := make(map[string][]models.Relationship)
	for _, ct := range cts {
		relationships[ct.ID.String()] = []models.Relationship{}
	}
	for _, r := range company.Relationships {
		relationships[r.RightCompany.TypeID.String()] = append(relationships[r.RightCompany.TypeID.String()], r)
	}
	for cti := range relationships {
		sort.SliceStable(relationships[cti], func(i, j int) bool {
			return relationships[cti][i].Tier > relationships[cti][j].Tier
		})

	}
	cps, err := models.ListCompanies(nil)
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	companies := make(map[string][]models.Company)
	for _, ct := range cts {
		companies[ct.ID.String()] = []models.Company{}
	}
	for _, cp := range cps {
		companies[cp.TypeID.String()] = append(companies[cp.TypeID.String()], cp)
	}
	for cti := range companies {
		sort.SliceStable(companies[cti], func(i, j int) bool {
			return companies[cti][i].Name < companies[cti][j].Name
		})

	}
	funcMap := template.FuncMap{
		"uts": func(u uuid.UUID) string {
			return u.String()
		},
	}
	t, err := template.Must(template.ParseFiles("web/layout.html")).Funcs(funcMap).ParseFiles("web/companies-single.html")
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	d := struct {
		PageTitle     string
		Company       models.Company
		CompanyTypes  []models.CompanyType
		Relationships map[string][]models.Relationship
		Tiers         []models.Tier
		Companies     map[string][]models.Company
	}{
		PageTitle:     fmt.Sprintf("Company %s information", (*company).Name),
		Company:       *company,
		CompanyTypes:  cts,
		Relationships: relationships,
		Tiers:         models.ListTiers(),
		Companies:     companies,
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
	cts, err := models.ListCompanyTypes(nil)
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	scs, err := models.ListSectors(nil)
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	d := struct {
		PageTitle    string
		Company      models.Company
		CompanyTypes []models.CompanyType
		Sectors      []models.Sector
	}{
		PageTitle:    "Add Company",
		Company:      *company,
		CompanyTypes: cts,
		Sectors:      scs,
	}
	funcMap := template.FuncMap{
		"uts": func(u uuid.UUID) string {
			return u.String()
		},
	}
	t, err := template.Must(template.ParseFiles("web/layout.html")).Funcs(funcMap).ParseFiles("web/companies-form.html")
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
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
	sectorID, err := uuid.Parse(r.FormValue("sector"))
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
		return
	}
	typeID, err := uuid.Parse(r.FormValue("type"))
	if err != nil {
		log.Println(err)
		if _, err := w.Write([]byte(err.Error())); err != nil {
			log.Println(err)
			return
		}
	}
	// Fill the data
	company.Name = r.FormValue("name")
	company.SectorID = sectorID
	company.TypeID = typeID
	company.Country = r.FormValue("country")

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
