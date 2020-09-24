package templatedata

import (
	"database/sql"
	"fmt"
	"github.com/tsawler/goblender/client/clienthandlers/clientdb"
	"github.com/tsawler/goblender/pkg/templates"
)

var dbModel *clientdb.DBModel
var version string

// NewTemplateData sets database connection
func NewTemplateData(p *sql.DB, v string) {
	dbModel = &clientdb.DBModel{DB: p}
	version = v
}

// AddDefaultData adds default data for templates
func AddDefaultData(td *templates.TemplateData) *templates.TemplateData {
	td.ClientVersion = fmt.Sprintf(`<strong><a href="https://github.com/tsawler/courses/" target="_blank">Courses</a></strong> %s`, version)
	return td
}
