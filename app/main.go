package main

import (
	"sample-service/app/esim"
	"sample-service/app/extension"
	"sample-service/app/infra"
	"sample-service/app/packages"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/swagger"
	"github.com/robfig/cron/v3"
)

// @title LightSIM Vendor API
// @version 1.0
// @description This is a LightSIM Vendor API
// @BasePath /
func main() {

	infra.LoadEnv()
	infra.DbConnect()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(recover.New())

	// Routes
	app.Get("vendor-packages", packages.VendorPackagesHandler)
	app.Get("load-packages", packages.LoadPackages)

	// Common vendor interface routes
	app.Post("esim", esim.CreateEsim)
	app.Get("esim", esim.GetEsim)

	app.Post("extension", extension.AddExtensionTopup)
	app.Get("extension", extension.GetExtension)

	app.Get("package", packages.GetPackage)
	app.Patch("package/:package_id", packages.UpdatePackage)
	app.Get("package/:iccid", packages.GetPackagesTopup)

	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Cron Jobs
	c := cron.New()
	c.AddFunc("@every 24h", func() { packages.LoadPackagesService("") })
	c.Start()

	app.Listen(":3000")
}
