package templatedata

import (
	"database/sql"
	"github.com/tsawler/goblender/client/clienthandlers/clientdb"
	"github.com/tsawler/goblender/pkg/templates"
)

var dbModel *clientdb.DBModel

// NewTemplateData sets database connection
func NewTemplateData(p *sql.DB) {
	dbModel = &clientdb.DBModel{DB: p}
}

// AddDefaultData adds default data for templates
func AddDefaultData(td *templates.TemplateData) *templates.TemplateData {
	td.RowSets["courses_version"] = "1.0.0"
	return td
}
