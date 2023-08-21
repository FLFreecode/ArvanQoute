package api

import (
	pkg "github.com/arvan/qoute/pkg/redisclient"
	"github.com/gofiber/fiber/v2"
)

func middleware(c *fiber.Ctx) error {
	return c.Next()
}

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"ok": true})
}

func info(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"name":      "Arvan-Qoute-Service",
		"desc":      "contact business logic service",
		"tech-used": "go 1.21",
		"version":   "1.0.1",
	})
}

func add(c *fiber.Ctx) error {
	var req pkg.Request

	req.Uuid = c.Params("uuid")
	req.UserName = c.Params("username")

	if err := c.BodyParser(&req); err != nil || len(req.Uuid) <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON or params is empty",
		})
	}

	Result, err := pkg.RepetitiousCheck(&req, c)
	if Result != true {
		zlogger.Info().Msgf("Repetitious Qoute  %+v", string(req.Uuid))
		return err
	} else {
		zlogger.Info().Msgf("Success Qoute  %+v", string(req.Uuid))
	}

	Result, err = pkg.QouteChecker(&req, c)
	if Result != true {
		zlogger.Info().Msgf("You are blocked to post quote %+v", string(req.Uuid))
		return err
	} else {
		zlogger.Info().Msgf("You were allowed to post quote   %+v", string(req.Uuid))
	}

	Result, err = pkg.VolumeChecker(&req, c)
	if Result != true {
		zlogger.Info().Msgf("You are blocked to post quote [Date Policy] %+v", string(req.Uuid))
		return err
	} else {
		zlogger.Info().Msgf("You were allowed to post quote [Date Policy]  %+v", string(req.Uuid))
	}

	err = pkg.StoreQoute(&req, c)
	if err != nil {
		zlogger.Info().Msgf("Qoute Stored Unsuccessfully %+v", string(req.Uuid))
		return err
	} else {
		zlogger.Info().Msgf("Qoute Stored Successfully   %+v", string(req.Uuid))
		return c.Status(fiber.StatusOK).JSON("Qoute Stored Successfully")
	}
}
