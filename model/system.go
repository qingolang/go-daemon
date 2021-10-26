package model

// SystemInfo
type SystemInfo struct {
	RootPath     string `json:"root_path"`
	FlightNumber string `json:"flight_number"`
	TailNo       string `json:"tail_no"`
	Sn           string `json:"sn"`
	IsOffline    int    `json:"is_offline"`
	PaStatus     int    `json:"pa_status"`
	DomainName   string `json:"domain_name"`
}
