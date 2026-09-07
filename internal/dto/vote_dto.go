package dto

type CreateVoteRequest struct {
	Type string `json:"type"`
}

type VoteResponse struct {
	ConfirmationsCount int `json:"confirmationsCount"`
	InvalidationsCount int `json:"invalidationsCount"`
}
