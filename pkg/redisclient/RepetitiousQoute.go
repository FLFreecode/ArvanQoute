package pkg

import (
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/gofiber/fiber/v2"
)

func RepetitiousCheck(req *Request, ctx *fiber.Ctx) (bool, error) {
	if (RedisCacheCheck.Exists(Ctx, req.Uuid+req.Qoute)) == false {
		if err := RedisCacheQoute.Set(&cache.Item{
			Ctx:   Ctx,
			Key:   req.Uuid + req.Qoute,
			Value: req.Uuid,
			TTL:   time.Hour,
		}); err != nil {
			return false, ctx.Status(fiber.StatusForbidden).JSON(err)
		}
		return true, ctx.Status(fiber.StatusForbidden).JSON("Qoute is valid")

	}
	return false, ctx.Status(fiber.StatusForbidden).JSON("Qoute is repetitious")
}
