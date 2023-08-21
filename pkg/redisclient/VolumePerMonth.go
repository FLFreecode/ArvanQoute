package pkg

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/cache/v9"
	"github.com/gofiber/fiber/v2"
)

const (
	BaseDate = "2006-1-2"
)

func VolumeChecker(req *Request, ctx *fiber.Ctx) (bool, error) {
	var Sv int = len(req.Qoute)
	if (RedisCacheVolume.Exists(Ctx, "SentVolumeQoute|"+req.Uuid)) == false {
		if err := RedisCacheVolume.Set(&cache.Item{
			Ctx:   Ctx,
			Key:   "SentVolumeQoute|" + req.Uuid,
			Value: time.Now().Format(BaseDate) + "|" + strconv.Itoa(Sv),
			TTL:   time.Duration(Cfg.Client.AmountOfVolumeBlocking * 24 * int(time.Hour)), //30 * 24 * time.Hour,
		}); err != nil {
			return false, ctx.Status(fiber.StatusOK).JSON(err)
		}
	} else {

		var Value string
		var LastDate string
		var date time.Time

		RedisCacheVolume.Get(Ctx, "SentVolumeQoute|"+req.Uuid, &Value)
		Values := strings.Split(Value, "|")
		LastDate = Values[0]
		LSv, _ := strconv.Atoi(Values[1])
		Sv = Sv + LSv
		date, _ = time.Parse(BaseDate, LastDate)
		NowTime, _ := time.Parse(BaseDate, time.Now().Format(BaseDate))

		if Sv > Cfg.Client.VolumeQoute*1024 && ((NowTime.Unix() - date.Unix()) < int64(Cfg.Client.AmountOfDailyVolume*24*360)) {
			if date.Unix()+int64(Cfg.Client.AmountOfDailyVolume*24*360) > NowTime.Unix() {
				return false, ctx.Status(fiber.StatusOK).JSON("You are blocked to post quote [Date Policy]")
			}
		}

		if (NowTime.Unix() - date.Unix()) >= int64(Cfg.Client.AmountOfDailyVolume*24*360) {
			LastDate = time.Now().Format(BaseDate)
			Sv = 0
		}

		if Sv <= Cfg.Client.VolumeQoute*1024 || (date.Unix()+int64(Cfg.Client.AmountOfDailyVolume*24*360) <= NowTime.Unix()) {
			if err := RedisCacheVolume.Set(&cache.Item{
				Ctx:   Ctx,
				Key:   "SentVolumeQoute|" + req.Uuid,
				Value: LastDate + "|" + strconv.Itoa(Sv),
				TTL:   time.Duration(Cfg.Client.AmountOfVolumeBlocking * 24 * int(time.Hour)), //30 * 24 * time.Hour,
			}); err != nil {
				return false, ctx.Status(fiber.StatusOK).JSON(err)
			}
		}

	}

	return true, ctx.Status(fiber.StatusOK).JSON("You were allowed to post quote [Date Policy]")
}
