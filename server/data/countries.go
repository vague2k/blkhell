package data

import "slices"

func SortCountryMap() []string {
	codes := make([]string, 0, len(Countries))
	for code := range Countries {
		codes = append(codes, code)
	}
	slices.Sort(codes)
	return codes
}

var Countries = map[string]string{
	"US": "United States",
	"CA": "Canada",
	"MX": "Mexico",
	"BR": "Brazil",
	"AR": "Argentina",
	"GB": "United Kingdom",
	"DE": "Germany",
	"FR": "France",
	"ES": "Spain",
	"IT": "Italy",
	"NL": "Netherlands",
	"SE": "Sweden",
	"NO": "Norway",
	"FI": "Finland",
	"JP": "Japan",
	"KR": "South Korea",
	"CN": "China",
	"AU": "Australia",
	"NZ": "New Zealand",
}
