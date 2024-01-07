package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/ScarletTanager/basilisk/basilisk/model"
	"github.com/labstack/echo/v4"
)

const (
	ParamModelID = "id"

	ContextKeyModel = "model"
)

// RetrieveModelMiddlware returns a middleware function configured for:
//   - Extracting model ID from URI
//   - Checking if model exists, and:
//   - Adding model to context if it exists, or
//   - Returning a 404 if it does not
func RetrieveModelMiddleware(rm *model.RunningModels) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var (
				id  int
				err error
			)

			if id, err = strconv.Atoi(c.Param(ParamModelID)); err != nil {
				return c.JSON(http.StatusBadRequest, &model.ModelsError{Message: fmt.Sprintf("%s is not a valid model id", c.Param(ParamModelID))})
			}

			if id < 0 || id > len(rm.Classifiers)-1 {
				return c.JSON(http.StatusNotFound, &model.ModelsError{Message: "Model not found"})
			}

			c.Set(ContextKeyModel, rm.Classifiers[id])
			return next(c)
		}
	}
}

type AllowedHeaders []string

func (ah AllowedHeaders) String() string {
	var b strings.Builder
	for hi := 0; hi < (len(ah) - 1); hi++ {
		b.WriteString(ah[hi] + " | ")
	}
	b.WriteString(ah[len(ah)-1])

	return b.String()
}

func CheckContentTypeMiddleware(allowed AllowedHeaders) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// An empty slice means allow anything
			if len(allowed) == 0 {
				return next(c)
			}

			headers := c.Request().Header
			if contentTypeHeader := headers.Get(echo.HeaderContentType); contentTypeHeader == "" {
				return c.JSON(http.StatusBadRequest, &model.ModelsError{Message: "Missing required header: Content-type"})
			} else {
				for _, allowedHeader := range allowed {
					if contentTypeHeader == allowedHeader {
						return next(c)
					}
				}

				return c.JSON(http.StatusUnsupportedMediaType, &model.ModelsError{Message: fmt.Sprintf("Unsupported media type: must be one of %s, found %s", allowed.String(), contentTypeHeader)})
			}
		}
	}
}
