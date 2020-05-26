package web

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/valuechaintool/valuechaintool/models"
)

// CompaniesCreate renders the /companies/new page
func CompaniesCreate(c *gin.Context) {
	cts, err := models.ListCompanyTypes(nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	scs, err := models.ListSectors(nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
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
	c.HTML(http.StatusOK, "companies-form.html", d)
}

// CompaniesCreatePost parses the form from /companies/new page
func CompaniesCreatePost(c *gin.Context) {
	sectorID, err := uuid.Parse(c.PostForm("sector"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
	}
	typeID, err := uuid.Parse(c.PostForm("type"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
	}
	company := models.Company{
		Name:     c.PostForm("name"),
		Country:  c.PostForm("country"),
		SectorID: sectorID,
		TypeID:   typeID,
	}
	if err := models.NewCompany(&company); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Redirect(http.StatusFound, "/companies")
}

// CompaniesList renders the /companies page
func CompaniesList(c *gin.Context) {
	companies, err := models.ListCompanies(nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	d := struct {
		PageTitle string
		Companies []models.Company
	}{
		PageTitle: "HomePage",
		Companies: companies,
	}
	c.HTML(http.StatusOK, "companies-list.html", d)
}

// CompaniesRead renders the /companies/[ID] page
func CompaniesRead(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
	}
	company, err := models.GetCompany(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
	}
	if err := company.EagerLoad(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	cts, err := models.ListCompanyTypes(nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
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
		_ = c.AbortWithError(http.StatusInternalServerError, err)
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
	c.HTML(http.StatusOK, "companies-single.html", d)
}

// CompaniesUpdate renders the /companies/[ID]/edit page
func CompaniesUpdate(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
	}
	company, err := models.GetCompany(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
	}
	cts, err := models.ListCompanyTypes(nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	scs, err := models.ListSectors(nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
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
	c.HTML(http.StatusOK, "companies-form.html", d)
}

// CompaniesUpdatePost parses the form from /companies/[ID]/edit page
func CompaniesUpdatePost(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
	}
	company, err := models.GetCompany(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
	}
	sectorID, err := uuid.Parse(c.PostForm("sector"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
	}
	typeID, err := uuid.Parse(c.PostForm("type"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
	}
	// Fill the data
	company.Name = c.PostForm("name")
	company.SectorID = sectorID
	company.TypeID = typeID
	company.Country = c.PostForm("country")

	if err := company.Save(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/companies/%s", id.String()))
}

// CompaniesDelete responds to /companies/[ID]/delete url
func CompaniesDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
	}
	company, err := models.GetCompany(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
	}
	if err := company.Delete(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
	}
	c.Redirect(http.StatusFound, "/companies")
}
