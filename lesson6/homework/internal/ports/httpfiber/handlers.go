package httpfiber

import (
	"errors"
	"net/http"

	"github.com/gofiber/fiber/v2"

	"homework6/internal/app"
)

func GetErrorHTTPStatus(err error) int {
	if errors.Is(err, app.ErrNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, app.ErrNotAuthor) {
		return http.StatusForbidden
	}
	if errors.Is(err, app.ErrInvalid) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}

// Метод для создания объявления (ad)
func createAd(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody createAdRequest
		err := c.BodyParser(&reqBody)
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		ad, err := a.CreateAd(c.Context(), reqBody.Title, reqBody.Text, reqBody.UserID)

		if err != nil {
			c.Status(GetErrorHTTPStatus(err))
			return c.JSON(AdErrorResponse(err))
		}
		return c.JSON(AdSuccessResponse(ad))
	}
}

// Метод для изменения статуса объявления (опубликовано - Published = true или снято с публикации Published = false)
func changeAdStatus(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody changeAdStatusRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		err = a.ChangeAdStatus(c.Context(), int64(adID), reqBody.UserID, reqBody.Published)

		if err != nil {
			c.Status(GetErrorHTTPStatus(err))
			return c.JSON(GetErrorHTTPStatus(err))
		}

		ad, err := a.GetAdByID(c.Context(), int64(adID))

		if err != nil {
			c.Status(GetErrorHTTPStatus(err))
			return c.JSON(AdErrorResponse(err))
		}

		return c.JSON(AdSuccessResponse(ad))
	}
}

// Метод для обновления текста(Text) или заголовка(Title) объявления
func updateAd(a app.App) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var reqBody updateAdRequest
		if err := c.BodyParser(&reqBody); err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		adID, err := c.ParamsInt("ad_id")
		if err != nil {
			c.Status(http.StatusBadRequest)
			return c.JSON(AdErrorResponse(err))
		}

		err = a.UpdateAd(c.Context(), int64(adID), reqBody.UserID, reqBody.Title, reqBody.Text)

		if err != nil {
			c.Status(GetErrorHTTPStatus(err))
			return c.JSON(GetErrorHTTPStatus(err))
		}

		ad, err := a.GetAdByID(c.Context(), int64(adID))

		if err != nil {
			c.Status(GetErrorHTTPStatus(err))
			return c.JSON(AdErrorResponse(err))
		}

		return c.JSON(AdSuccessResponse(ad))
	}
}
