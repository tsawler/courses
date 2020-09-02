package clienthandlers

import (
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

// ClientInit gives client code access to goBlender configuration
func ClientInit(conf config.AppConfig, parentDriver *driver.DB, rep *handlers.DBRepo) {
	// conf is the application config, from goBlender
	app = conf
	repo = rep

	// If we have additional databases (external to this application) we set the connection here.
	// The connection is specified in goBlender preferences.
	//conn := app.AlternateConnection

	// loggers
	infoLog = app.InfoLog
	errorLog = app.ErrorLog

	// In case we need it, we get the database connection from goBlender and save it,
	parentDB = parentDriver

	// We can access handlers from goBlender, but need to initialize them first.
	if app.Database == "postgresql" {
		handlers.NewPostgresqlHandlers(parentDB, app.ServerName, app.InProduction)
	} else {
		handlers.NewMysqlHandlers(parentDB, app.ServerName, app.InProduction)
	}

	// Set a different template for home page, if needed.
	//repo.SetHomePageTemplate("client-sample.page.tmpl")

	// Set a different template for inside pages, if needed.
	//repo.SetDefaultPageTemplate("client-sample.page.tmpl")

	// Create client middleware
	NewClientMiddleware(app)
}
