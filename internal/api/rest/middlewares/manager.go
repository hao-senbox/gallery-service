package middlewares

import (
	"context"
	"fmt"
	"gallery-service/config"
	"gallery-service/internal/pkg/apicall"
	httpPkg "gallery-service/pkg/http"
	"gallery-service/pkg/zap"
	"github.com/gofiber/fiber/v2"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"runtime/debug"
	"strings"
	"time"
)

type MiddlewareManager interface {
	RequestLoggerMiddleware() fiber.Handler
	Auth(*api.Client) fiber.Handler
	Recovery() fiber.Handler
}

type middlewareManager struct {
	log zap.Logger
	cfg *config.Config
}

func NewMiddlewareManager(log zap.Logger, cfg *config.Config) *middlewareManager {
	return &middlewareManager{log: log, cfg: cfg}
}

func (mw *middlewareManager) RequestLoggerMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		start := time.Now()

		// Call the next handler
		err := ctx.Next()

		req := ctx.Request()
		res := ctx.Response()
		status := res.StatusCode()
		size := int64(len(res.Body()))
		duration := time.Since(start)

		if !mw.checkIgnoredURI(string(req.RequestURI()), mw.cfg.App.API.Rest.Setting.IgnoreLogUrls) {
			mw.log.HttpMiddlewareAccessLogger(ctx.Method(), string(req.RequestURI()), status, size, duration)
		}

		return err
	}
}

func (mw *middlewareManager) checkIgnoredURI(requestURI string, uriList []string) bool {
	for _, s := range uriList {
		if strings.Contains(requestURI, s) {
			return true
		}
	}
	return false
}

// Auth checks for a valid token in the request headers
func (mw *middlewareManager) Auth(client *api.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			mw.log.Warn("Authorization header missing")
			return httpPkg.ErrorCtxResponse(c, errors.New("missing authorization header"), mw.cfg.App.API.Rest.Setting.DebugErrorsResponse)
		}

		// Split the header into parts (e.g., "Bearer token")
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			mw.log.Warn("Invalid authorization header format")
			return httpPkg.ErrorCtxResponse(c, errors.New("invalid authorization header format"), mw.cfg.App.API.Rest.Setting.DebugErrorsResponse)
		}

		call, err := apicall.NewGoMainServiceAPI(client)
		if err != nil {
			return errors.New("failed to create apicall")
		}

		// Get the current context.Context from Fiber
		userCtx := c.UserContext()

		user, err := call.GetUserByToken(userCtx, authHeader)
		if err != nil {
			mw.log.Errorf("failed to get user: {%v}", err)
			return httpPkg.ErrorCtxResponse(c, errors.New("failed to get user"), mw.cfg.App.API.Rest.Setting.DebugErrorsResponse)
		}

		// Add a value to the context
		userCtx = context.WithValue(userCtx, "current_user", user)
		userCtx = context.WithValue(userCtx, "current_token", authHeader)
		// Set the updated context back to Fiber
		c.SetUserContext(userCtx)

		// If the token is valid, proceed to the next handler
		return c.Next()
	}
}

func (mw *middlewareManager) Recovery() fiber.Handler {
	return func(c *fiber.Ctx) error {
		defer func() {
			if err := recover(); err != nil {
				mw.log.Errorf("Recovered from panic error: {%v}", err)

				// Print the stack trace to the console
				stackTrace := string(debug.Stack())
				fmt.Printf("Panic stack trace: \n%s\n", stackTrace)

				_ = httpPkg.ErrorCtxResponse(c, errors.Errorf("%v", err), mw.cfg.App.API.Rest.Setting.DebugErrorsResponse)
			}
		}()
		return c.Next()
	}
}
