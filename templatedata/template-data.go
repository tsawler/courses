package templatedata

import (
	"database/sql"
	"github.com/tsawler/goblender/client/clienthandlers/clientdb"
	"github.com/tsawler/goblender/pkg/templates"
)

var dbModel *clientdb.DBModel

func NewTemplateData(p *sql.DB) {
	dbModel = &clientdb.DBModel{DB: p}
}

// AddDefaultData adds default data for templates
func AddDefaultData(td *templates.TemplateData) *templates.TemplateData {
	return td
}