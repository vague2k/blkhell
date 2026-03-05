package database

import (
	"github.com/vague2k/blkhell/server/data"
)

func (b Band) FullCountry() string {
	return data.Countries[b.Country]
}
