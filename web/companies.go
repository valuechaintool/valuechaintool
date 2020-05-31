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
		return
	}
	vcs, err := models.ListVerticals(nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	d := struct {
		PageTitle    string
		Company      models.Company
		CompanyTypes []models.CompanyType
		Verticals    []models.Vertical
	}{
		PageTitle:    "Add Company",
		CompanyTypes: cts,
		Verticals:    vcs,
	}
	c.HTML(http.StatusOK, "companies-form.html", d)
}

// CompaniesCreatePost parses the form from /companies/new page
func CompaniesCreatePost(c *gin.Context) {
	var verticals []models.Vertical
	for _, vc := range c.PostFormArray("verticals[]") {
		id, err := uuid.Parse(vc)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}
		verticals = append(verticals, models.Vertical{ID: id})
	}
	typeID, err := uuid.Parse(c.PostForm("type"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	company := models.Company{
		Name:      c.PostForm("name"),
		Country:   c.PostForm("country"),
		Verticals: verticals,
		TypeID:    typeID,
	}
	if err := models.NewCompany(&company); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, "/companies")
}

// CompaniesList renders the /companies page
func CompaniesList(c *gin.Context) {
	companies, err := models.ListCompanies(nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
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
		return
	}
	company, err := models.GetCompany(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err := company.EagerLoad(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	cts, err := models.ListCompanyTypes(nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	vcs, err := models.ListVerticals(nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
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
			return relationships[cti][i].LeftTier > relationships[cti][j].LeftTier
		})

	}
	cps, err := models.ListCompanies(nil)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
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
	d := struct {
		PageTitle     string
		Company       models.Company
		CompanyTypes  []models.CompanyType
		Relationships map[string][]models.Relationship
		Tiers         []models.Tier
		Verticals     []models.Vertical
		Companies     map[string][]models.Company
	}{
		PageTitle:     fmt.Sprintf("Company %s information", (*company).Name),
		Company:       *company,
		CompanyTypes:  cts,
		Relationships: relationships,
		Tiers:         models.ListTiers(),
		Verticals:     vcs,
		Companies:     companies,
	}
	c.HTML(http.StatusOK, "companies-single.html", d)
}

// CompaniesUpdatePost parses the form from /companies/[ID]/edit page
func CompaniesUpdatePost(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	company, err := models.GetCompany(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	updates := make(map[string]interface{})

	// Name
	if name, ok := c.GetPostForm("name"); ok {
		updates["name"] = name
	}

	// Verticals
	if vcs, ok := c.GetPostFormArray("verticals[]"); ok {
		var verticals []models.Vertical
		for _, vc := range vcs {
			id, err := uuid.Parse(vc)
			if err != nil {
				_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
				return
			}
			verticals = append(verticals, models.Vertical{ID: id})
		}
		updates["verticals"] = verticals
	}

	// Type
	if t, ok := c.GetPostForm("type"); ok {
		typeID, err := uuid.Parse(t)
		if err != nil {
			_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
			return
		}
		updates["type_id"] = typeID
	}

	// Country
	if country, ok := c.GetPostForm("country"); ok {
		updates["country"] = country
	}

	if err := company.Update(updates); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, fmt.Sprintf("/companies/%s?edit_mode=true", id.String()))
}

// CompaniesDelete responds to /companies/[ID]/delete url
func CompaniesDelete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	company, err := models.GetCompany(id)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err := company.Delete(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, "/companies")
}
