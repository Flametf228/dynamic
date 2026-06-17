package dto

type StatsResponse struct {
	Products   int `json:"products"`
	Clients    int `json:"clients"`
	Brands     int `json:"brands"`
	Categories int `json:"categories"`
}
