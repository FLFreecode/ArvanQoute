package pkg

import "github.com/gofiber/fiber/v2"

func StoreQoute(req *Request, ctx *fiber.Ctx) error {
	_, err := RedisClient.SAdd(Ctx, req.Uuid, req.Qoute).Result()
	if err != nil {
		return ctx.Status(fiber.StatusOK).JSON("Qoute Stored Unsuccessfully")
	}
	return ctx.Status(fiber.StatusOK).JSON("Qoute Stored Successfully")
}
