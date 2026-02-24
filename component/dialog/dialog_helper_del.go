package dialog

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// HandleDelRequest handles single-item delete operations using the standard dialog pattern.
//
// GET request: Shows delete confirmation dialog using NewDialogDelete
// POST request: Executes delete operation and returns type-safe success response
func HandleDelRequest(
	c echo.Context,
	confirmMessage string,
	urlPath string,
	deleteFunc func() error,
	successResponse response.SuccessResponse,
	translator core.TranslateFunc,
) error {
	if c.Request().Method == http.MethodGet {
		dlg := NewDialogDelete(
			confirmMessage,
			url.NewUrl(urlPath),
			nil,
			nil,
			nil,
			nil,
			translator,
		)

		return c.JSON(http.StatusOK, dlg.Print(translator))
	}

	if err := deleteFunc(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, successResponse)
}
