package main

import (
	"stock-tracker/cmd"
	"stock-tracker/db"
)

// @title           Swagger calendar API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/
// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io
// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
// @host      localhost:1323
// @BasePath  /api/v1
// @securityDefinitions.apikey  BearerAuth
// @in header
// @name Authorization
// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	closeConn := db.InitPostgres()
	defer closeConn()

	cmd.Execute()
}
