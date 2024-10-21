package handler

import (
	"context"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/umutteroll07/ShortenURL/app/helper"
	"github.com/umutteroll07/ShortenURL/app/model"
	"github.com/umutteroll07/ShortenURL/internal/database"
	"gorm.io/gorm"
)

var ctx = context.Background()

type RequestUrl struct {
	Data string `json:"data"`
}

var db = database.Connection()

func ShortenUrl(c *fiber.Ctx) error {
	var url model.Url
	var body RequestUrl
	rdb := helper.RedisClient()
	ttl := helper.GetTTL()
	expires_at := time.Now().Add(ttl).Format(time.RFC3339)

	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	originalURL := body.Data

	if helper.UniqueOriginalUrlControl(originalURL) {

		if err := db.First(&url, "original_url = ?", originalURL).Error; err != nil {
			return c.Status(500).SendString(err.Error())
		}

		url.Short_url = os.Getenv("BASE_URL") + helper.StringWithCharset(6)
		url.Expires_at = expires_at
		if err := db.Save(&url).Error; err != nil {
			return c.Status(500).SendString(err.Error())
		}

		rdb.Set(ctx, url.Short_url, url.Original_url, ttl)
		return c.JSON(url.Short_url)

	} else {
		urlShort := model.Url{
			Original_url: originalURL,
			Short_url:    os.Getenv("BASE_URL") + helper.StringWithCharset(6),
			Usage_count:  0,
			Expires_at:   expires_at,
		}
		if err := db.Create(&urlShort).Error; err != nil {
			return c.Status(500).SendString(err.Error())
		}
		rdb.Set(ctx, urlShort.Short_url, url.Original_url, ttl)
		return c.JSON(urlShort)
	}

}

func RedirectUrl(c *fiber.Ctx) error {
	var url model.Url
	var redirect_shorten_url = os.Getenv("BASE_URL") + c.Params("shortUrl")
	if err := db.First(&url, "short_url = ?", redirect_shorten_url).Error; err != nil {
		return c.Status(fiber.StatusNotFound).SendString("URL not found")
	}

	  isExpired, err := helper.CheckExpiration(redirect_shorten_url)
	  if err != nil {
	      return c.Status(fiber.StatusInternalServerError).SendString("Redis error: " + err.Error())
	  }
	  if !isExpired {
	      return c.Status(fiber.StatusGone).SendString("URL has expired")
	  }

	SetUsageCount(&url, c)
	var redirect_original_url = url.Original_url
	return c.Redirect(redirect_original_url)

}

func SetUsageCount(url *model.Url, c *fiber.Ctx) error {
	url.Usage_count = url.Usage_count + 1

	if err := db.Save(&url).Error; err != nil {
		return c.Status(500).SendString(err.Error())
	}
	return c.SendString("set usage count is success")
}

func DeleteShortenUrl(c *fiber.Ctx) error {

	var delete_url = os.Getenv("BASE_URL") + c.Params("deleteUrl")

	var url model.Url
	if err := database.Connection().Model(&model.Url{}).Where("short_url = ?", delete_url).First(&url).Error; err != nil {
		if gorm.ErrRecordNotFound == err {
			return c.Status(404).JSON(fiber.Map{
				"error": "Record not found",
			})
		}
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to find item",
		})
	}

	if err := database.Connection().Model(&model.Url{}).Where("short_url = ?", delete_url).Delete(&model.Url{}).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete item",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Item deleted successfully",
	})

}

func GetAllShortenUrl(c *fiber.Ctx) error {

	var urlList []model.Url
	result := db.Find(&urlList)

	if result.Error != nil {
		return c.Status(500).SendString("get all shorten_url is failed")
	}

	return c.JSON(urlList)

}
