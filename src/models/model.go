package models

const (
	ColorOrange = "\033[33m"
	ColorReset  = "\033[0m"
	Production  = "production"
	Staging     = "staging"
	Development = "development"
	Local       = "local"
)

var Status = struct {
	Active      string
	Inactive    string
	Deactivated string
}{
	Active:      "active",
	Inactive:    "inactive",
	Deactivated: "deactivated",
}
var AllModels = []interface{}{&User{}}
