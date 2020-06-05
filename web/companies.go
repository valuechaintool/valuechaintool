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
	if !isAllowed(c, models.WildCardResource, "createCompany") {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("operation not allowed"))
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
	if !isAllowed(c, models.WildCardResource, "createCompany") {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("operation not allowed"))
		return
	}
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
	userID, ok := c.Get("userID")
	if !ok {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	if err := models.NewCompany(&company, userID.(uuid.UUID)); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, "/companies")
}

// CompaniesList renders the /companies page
func CompaniesList(c *gin.Context) {
	companies, err := models.ListCompanies(nil, true)
	if err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	visibleCompanies := []models.Company{}
	for k := range companies {
		if isAllowed(c, companies[k].ID, "readCompany") {
			visibleCompanies = append(visibleCompanies, companies[k])
		}
	}
	d := struct {
		PageTitle          string
		Companies          []models.Company
		CanCreateCompanies bool
	}{
		PageTitle:          "HomePage",
		Companies:          visibleCompanies,
		CanCreateCompanies: isAllowed(c, models.WildCardResource, "createCompany"),
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
	if !isAllowed(c, id, "readCompany") {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("operation not allowed"))
		return
	}
	company, err := models.GetCompany(id, true)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	if err := company.EagerLoad(); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	owners, err := company.Owners()
	if err != nil {
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

	cps, err := models.ListCompanies(nil, false)
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
		PageTitle              string
		Company                models.Company
		Owners                 []models.User
		CompanyTypes           []models.CompanyType
		Relationships          map[string][]models.Relationship
		Tiers                  []models.Tier
		Verticals              []models.Vertical
		Companies              map[string][]models.Company
		CanUpdateCompany       bool
		CanDeleteCompany       bool
		CanCreateRelationships bool
		CanUpdateRelationships bool
		CanDeleteRelationships bool
	}{
		PageTitle:              fmt.Sprintf("Company %s information", (*company).Name),
		Company:                *company,
		Owners:                 owners,
		CompanyTypes:           cts,
		Relationships:          relationships,
		Tiers:                  models.ListTiers(),
		Verticals:              vcs,
		Companies:              companies,
		CanUpdateCompany:       isAllowed(c, id, "updateCompany"),
		CanDeleteCompany:       isAllowed(c, id, "deleteCompany"),
		CanCreateRelationships: isAllowed(c, id, "createRelationship"),
		CanUpdateRelationships: isAllowed(c, id, "updateRelationship"),
		CanDeleteRelationships: isAllowed(c, id, "deleteRelationship"),
	}
	c.HTML(http.StatusOK, "companies-single.html", d)
}

// CompaniesRead renders the /companies/[ID] page
func CompaniesChangelog(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	if !isAllowed(c, id, "readCompany") {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("operation not allowed"))
		return
	}
	company, err := models.GetCompany(id, true)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}
	changes, err := models.ListChangesByCompany(id, true)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	// Change descriptions
	changesDescriptions := []string{}
	for _, change := range changes {
		changesDescriptions = append(changesDescriptions, changeDescription(change))
	}

	d := struct {
		PageTitle           string
		Company             models.Company
		Changes             []models.Change
		ChangesDescriptions []string
	}{
		PageTitle:           fmt.Sprintf("Company %s history", (*company).Name),
		Company:             *company,
		Changes:             changes,
		ChangesDescriptions: changesDescriptions,
	}
	c.HTML(http.StatusOK, "companies-changelog.html", d)
}

func changeDescription(change models.Change) string {
	switch change.Type {
	case models.CompanyCreated:
		return "the company was created"
	case models.CompanyUpdated:
		return fmt.Sprintf("the value of '%s' was changed from '%s' to '%s'", change.Key, change.PreviousValue, change.NewValue)
	case models.CompanyDeleted:
		return "the company was deleted"
	case models.VerticalAdded:
		return fmt.Sprintf("the vertical '%s' was added", change.NewValue)
	case models.VerticalRemoved:
		return fmt.Sprintf("the vertical '%s' was removed", change.PreviousValue)
	case models.RelationshipCreated:
		return fmt.Sprintf("the relationship with '%s' was created", change.Relationship.RightCompany.Name)
	case models.RelationshipUpdated:
		return fmt.Sprintf("the relationship with '%s' had the field '%s' changed from '%s' to '%s'", change.Relationship.RightCompany.Name, change.Key, change.PreviousValue, change.NewValue)
	case models.RelationshipDeleted:
		return fmt.Sprintf("the relationship with '%s' was deleted", change.Relationship.RightCompany.Name)
	default:
		return "unimplemented"
	}
}

// CompaniesUpdatePost parses the form from /companies/[ID]/edit page
func CompaniesUpdatePost(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	if !isAllowed(c, id, "updateCompany") {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("operation not allowed"))
		return
	}
	company, err := models.GetCompany(id, true)
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

	// User
	userID, ok := c.Get("userID")
	if !ok {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	if err := company.Update(updates, userID.(uuid.UUID)); err != nil {
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
	if !isAllowed(c, id, "deleteCompany") {
		_ = c.AbortWithError(http.StatusUnauthorized, fmt.Errorf("operation not allowed"))
		return
	}
	company, err := models.GetCompany(id, true)
	if err != nil {
		_ = c.AbortWithError(http.StatusNotFound, err)
		return
	}

	// User
	userID, ok := c.Get("userID")
	if !ok {
		_ = c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}

	if err := company.Delete(userID.(uuid.UUID)); err != nil {
		_ = c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Redirect(http.StatusFound, "/companies")
}
