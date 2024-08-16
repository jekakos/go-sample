package packages

import (
	"github.com/gofiber/fiber/v2"
)

func LoadPackages(c *fiber.Ctx) error {

	reset := c.Query("reset")
	result, err := LoadPackagesService(reset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.SendString(result)
}
