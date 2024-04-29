package tracerfiber

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"

	"github.com/dalmarcogd/test-defi/pkg/tracer"
	"github.com/dalmarcogd/test-defi/pkg/tracer/obfurl"
)

func NewMiddleware(trc tracer.Tracer, ignorePaths ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		path := c.Path()
		for _, ignorePath := range ignorePaths {
			if strings.Contains(path, ignorePath) {
				return c.Next()
			}
		}

		ctx := c.UserContext()

		obfuscatedURL := obfurl.ObfuscateURL(path)
		method := c.Route().Method
		name := fmt.Sprintf("%v %v", method, obfuscatedURL)

		ctx = trc.Extract(ctx, propagation.HeaderCarrier(c.GetReqHeaders()))
		ctx, span := trc.SpanName(
			ctx,
			name,
		)

		span.SetAttributes(
			attribute.String("http.method", method),
			attribute.String("http.url", string(c.Request().RequestURI())),
			attribute.String("span.type", "web"),
			attribute.String("resource.name", name),
		)

		defer span.End()

		c.SetUserContext(ctx)
		err := c.Next()
		statusCode := c.Response().StatusCode()
		if err != nil {
			span.RecordError(err)
		}
		span.SetAttributes(
			attribute.Int("http.status_code", statusCode),
			attribute.Int("http.response.size", len(c.Response().Body())),
		)
		if statusCode >= http.StatusBadRequest {
			span.RecordError(fmt.Errorf("response code: %d", statusCode))
		}

		return err
	}
}
