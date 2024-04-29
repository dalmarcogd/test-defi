package recovermdw

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func NewMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) (err error) {
		panicked := true

		defer func() {
			if r := recover(); r != nil || panicked {
				if c.Response().StatusCode() == 0 {
					err := c.SendStatus(http.StatusInternalServerError)
					if err != nil {
						zap.L().Error("error marshalling the response", zap.Error(err))
					}
				}
			}
		}()

		err = c.Next()
		if err != nil {
			return err
		}

		panicked = false
		return nil
	}
}
