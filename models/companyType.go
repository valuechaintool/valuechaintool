package models

type CompanyType int

const (
	PartnerCompany CompanyType = 1 + iota
	CustomerCompany
)
