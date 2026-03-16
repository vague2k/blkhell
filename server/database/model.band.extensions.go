package database

import (
	"github.com/vague2k/blkhell/common"
)

func (b Band) FullCountry() string {
	return common.Countries[b.Country]
}
