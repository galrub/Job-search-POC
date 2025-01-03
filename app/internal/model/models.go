package model

import "time"

type LoginData struct {
	Email string `json:"email"`
	Passw string `json:"password"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

type JobEditParams struct {
	Company       string    `json:"company"`
	PositionDesc  string    `json:"positionDesc"`
	Remote        bool      `json:"remote"`
	ContractType  string    `json:"contractType"`
	Contacted     bool      `json:"contacted"`
	CreatedAt     time.Time `json:"createdAt"`
	Comments      string    `json:"comments"`
	GeneralStatus string    `json:"generalStatus"`
	ID            string    `json:"id"`
}

type ContractTypeSelect struct {
	Current string
	Options []string
}
