package server

import (
	httpPkg "gallery-service/pkg/http"
	"gallery-service/pkg/zap"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/juju/errors"
	"net"
	"net/http"
	"strings"
	"time"
)

func (s *server) start(errorChan chan<- error) {
	s.configure()
	s.routes()
	defer close(errorChan)

	if err := s.fiber.Listen(
		net.JoinHostPort(
			s.cfg.App.API.Rest.Host,
			s.cfg.App.API.Rest.Port,
		),
	); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errorChan <- err
	}
}

func (s *server) configure() {
	if s.cfg.App.API.Rest.Setting.Debug {
		s.fiber.Use(s.handleErrors(s.log, s.cfg.App.API.Rest.Setting.Debug))
	}

	// Add middlewares
	s.fiber.Use(recover.New())
	s.fiber.Use(s.mw.RequestLoggerMiddleware())
	s.fiber.Use(s.mw.Recovery()) // Recover middlewares
	s.fiber.Use(logger.New(logger.Config{
		Format:     "${pid} [${ip}]:${port} ${status} - ${method} ${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "Asia/Ho_Chi_Minh",
	}))
	s.fiber.Use(compress.New(compress.Config{
		Level: gzipLevel,
		Next: func(c *fiber.Ctx) bool {
			return strings.Contains(c.Path(), "swagger")
		},
	}))
	s.fiber.Use(helmet.New())
	s.fiber.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))
	s.fiber.Use(limiter.New(limiter.Config{
		Next: func(c *fiber.Ctx) bool {
			return strings.Contains(c.Path(), "swagger")
		},
		LimitReached: func(c *fiber.Ctx) error {
			return httpPkg.NewTooManyRequestError(c, "too many requests", s.cfg.App.API.Rest.Setting.DebugErrorsResponse)
		},
		Max:        10,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(ctx *fiber.Ctx) string {
			return ctx.IP()
		},
	}))
}

func (s *server) handleErrors(log zap.Logger, debug bool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Call the next handler
		err := c.Next()

		if err != nil {
			var (
				code       = http.StatusInternalServerError
				msg        string
				errorStack string
			)

			var he *fiber.Error
			if errors.As(err, &he) {
				code = he.Code
				msg = he.Message
			} else {
				msg = err.Error()
				switch true {
				case errors.Is(err, errors.BadRequest):
					code = http.StatusBadRequest
				case errors.Is(err, errors.Forbidden):
					code = http.StatusForbidden
				case errors.Is(err, errors.Unauthorized):
					code = http.StatusUnauthorized
				case errors.Is(err, errors.NotFound):
					code = http.StatusNotFound
				case errors.Is(err, errors.AlreadyExists):
					code = http.StatusConflict
				}

				if debug {
					errorStack = errors.ErrorStack(err)
				}

				if code == fiber.StatusInternalServerError {
					log.Error("(handleErrors) An error occurred", err)
				}

				response := httpPkg.NewRestError(code, msg, errorStack, s.cfg.App.API.Rest.Setting.Debug)

				return c.Status(code).JSON(response)
			}
		}

		return nil
	}
}
