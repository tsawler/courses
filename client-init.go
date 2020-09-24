package clienthandlers

import (
	"github.com/tsawler/goblender/client/clienthandlers/clientdb"
	template_data "github.com/tsawler/goblender/client/clienthandlers/templatedata"
	"github.com/tsawler/goblender/pkg/config"
	"github.com/tsawler/goblender/pkg/driver"
	"github.com/tsawler/goblender/pkg/handlers"
	"log"
)

var app config.AppConfig
var infoLog *log.Logger
var errorLog *log.Logger
var parentDB *driver.DB
var repo *handlers.DBRepo
var dbModel *clientdb.DBModel

const version = "1.3.2"

// ClientInit  initializes client specific code
func ClientInit(conf config.AppConfig, parentDriver *driver.DB, rep *handlers.DBRepo) {
	// conf is the application config, from goBlender
	app = conf
	repo = rep

	// loggers
	infoLog = app.InfoLog
	errorLog = app.ErrorLog

	// In case we need it, we get the database connection from goBlender and save it,
	parentDB = parentDriver

	// enable db for client access
	dbModel = &clientdb.DBModel{DB: parentDB.SQL}
	template_data.NewTemplateData(parentDB.SQL, version)

	// We can access handlers from goBlender, but need to initialize them first.
	if app.Database == "postgresql" {
		handlers.NewPostgresqlHandlers(parentDB, app.ServerName, app.InProduction)
	} else {
		handlers.NewMysqlHandlers(parentDB, app.ServerName, app.InProduction)
	}

	infoLog.Printf("******************************************")
	infoLog.Printf("** %sCourses%s v%s", "\033[31m", "\033[0m", version)

	// Create client middleware
	NewClientMiddleware(app)
}
