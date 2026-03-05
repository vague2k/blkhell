package database

import (
	"fmt"
)

func (r Release) CatalogNo() string {
	return fmt.Sprintf("BH%s", r.Number)
}
