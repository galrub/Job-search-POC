package ui

import (
	"github.com/galrub/go/jobSearch/internal/database"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type JobDto struct {
	ID            uuid.UUID `json:"id"`
	Company       string    `json:"company"`
	PositionDesc  string    `json:"positionDesc"`
	Remote        bool      `json:"remote"`
	ContractType  string    `json:"contractType"`
	Contacted     bool      `json:"contacted"`
	GeneralStatus string    `json:"generalStatus"`
	CreatedAt     string    `json:"createdAt"`
	UpdatedAt     string    `json:"updatedAt"`
	Comments      string    `json:"comments"`
}

func pgTextToString(t pgtype.Text) string {
	if !t.Valid {
		return ""
	}
	return t.String
}

func JobDtoFactory(j database.Job) JobDto {
	r := JobDto{
		ID:            j.ID,
		Company:       j.Company,
		PositionDesc:  j.PositionDesc,
		Remote:        j.Remote.Valid && j.Remote.Bool,
		ContractType:  j.ContractType,
		Contacted:     j.Contacted.Valid && j.Contacted.Bool,
		GeneralStatus: j.GeneralStatus,
		CreatedAt:     j.CreatedAt.Format("2006-1-2"),
		UpdatedAt:     j.UpdatedAt.Format("2006-1-2"),
		Comments:      pgTextToString(j.Comments),
	}
	return r
}

func MapJobsToDto(j []database.Job) []JobDto {
	r := make([]JobDto, len(j))
	for i, e := range j {
		r[i] = JobDtoFactory(e)
	}
	return r
}
