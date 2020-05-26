package models

import (
	"errors"

	"github.com/google/uuid"
)

type Role struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Capabilities []string  `json:"capabilities"`
}

func ListRoles(filters map[string]interface{}) ([]Role, error) {
	return roles, nil
}

func GetRole(id uuid.UUID) (*Role, error) {
	for _, role := range roles {
		if role.ID == id {
			return &role, nil
		}
	}
	return nil, errors.New("role not found")
}

var roles = []Role{
	{
		ID:   uuid.Must(uuid.Parse("28a98657-2c5b-435e-bfe5-18081066af8d")),
		Name: "Viewer",
		Capabilities: []string{
			"readCompany",
		},
	},
	{
		ID:   uuid.Must(uuid.Parse("2e30ca74-216d-4b3f-b6dd-7e52d1714dde")),
		Name: "Creator",
		Capabilities: []string{
			"createCompany",
		},
	},
	{
		ID:   uuid.Must(uuid.Parse("64f75ce2-63d4-4d9f-b4e7-2f6aec5b9906")),
		Name: "Editor",
		Capabilities: []string{
			"createCompany",
			"readCompany",
			"updateCompany",
			"createRelationship",
			"updateRelationship",
			"deleteReleationship",
		},
	},
	{
		ID:   uuid.Must(uuid.Parse("5a2dbf8e-8ba8-4ca5-ac2d-cc11f1f0fb2d")),
		Name: "Owner",
		Capabilities: []string{
			"createCompany",
			"readCompany",
			"updateCompany",
			"deleteCompany",
			"createRelationship",
			"updateRelationship",
			"deleteReleationship",
		},
	},
	{
		ID:   uuid.Must(uuid.Parse("23b0f96f-bd32-4421-96fa-e3bad618740c")),
		Name: "Administrator",
		Capabilities: []string{
			"createCompany",
			"readCompany",
			"updateCompany",
			"deleteCompany",
			"createRelationship",
			"updateRelationship",
			"deleteReleationship",
			"readUser",
			"updateUser",
			"deleteUser",
			"managePermission",
		},
	},
}
