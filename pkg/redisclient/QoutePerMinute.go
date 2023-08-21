package pkg

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/gofiber/fiber/v2"
)

const (
	BaseTime = "2006-1-2 15:4:5"
)

func QouteChecker(req *Request, ctx *fiber.Ctx) (bool, error) {
	var Nq int = 1

	if (RedisCacheQoute.Exists(Ctx, "QoutePerMinute|"+req.Uuid)) == false {
		if err := RedisCacheQoute.Set(&cache.Item{
			Ctx:   Ctx,
			Key:   "QoutePerMinute|" + req.Uuid,
			Value: time.Now().Format(BaseTime) + "|" + strconv.Itoa(Nq),
			TTL:   time.Duration(Cfg.Client.ClientBlockMinute) * time.Minute,
		}); err != nil {
			return false, ctx.Status(fiber.StatusOK).JSON(err)
		}
	} else {
		var Value string
		var LastTime string
		var date time.Time

		RedisCacheQoute.Get(Ctx, "QoutePerMinute|"+req.Uuid, &Value)
		Values := strings.Split(Value, "|")
		LastTime = Values[0]
		Nq, _ = strconv.Atoi(Values[1])
		Nq++
		date, _ = time.Parse(BaseTime, LastTime)
		NowTime, _ := time.Parse(BaseTime, time.Now().Format(BaseTime))

		if Nq > Cfg.Client.NumQoutePerMinute && ((NowTime.Unix() - date.Unix()) < int64(Cfg.Client.QoutePerMinute*60)) {
			if date.Unix()+int64(Cfg.Client.QoutePerMinute*60) > NowTime.Unix() {
				return false, ctx.Status(fiber.StatusOK).JSON("You are blocked to post quote")
			}
		}

		if (NowTime.Unix() - date.Unix()) >= int64(Cfg.Client.QoutePerMinute*60) {
			LastTime = time.Now().Format(BaseTime)
			Nq = 1
		}

		if Nq <= Cfg.Client.NumQoutePerMinute || (date.Unix()+int64(Cfg.Client.QoutePerMinute*60) <= NowTime.Unix()) {
			if err := RedisCacheQoute.Set(&cache.Item{
				Ctx:   Ctx,
				Key:   "QoutePerMinute|" + req.Uuid,
				Value: LastTime + "|" + strconv.Itoa(Nq),
				TTL:   time.Duration(Cfg.Client.ClientBlockMinute) * time.Minute,
			}); err != nil {
				return false, ctx.Status(fiber.StatusOK).JSON(err)
			}
		}
	}
	return true, ctx.Status(fiber.StatusOK).JSON("You were allowed to post quote")
}
