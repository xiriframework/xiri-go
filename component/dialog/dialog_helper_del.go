package dialog

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xiriframework/xiri-go/component/core"
	"github.com/xiriframework/xiri-go/component/url"
	"github.com/xiriframework/xiri-go/response"
)

// HandleDelRequest handles single-item delete operations using the standard dialog pattern.
//
//   - GET: shows the delete confirmation dialog built with NewDialogDelete
//   - POST: executes deleteFunc and returns the success response (500 on error)
//
// For custom button texts or a warning dialog use HandleDelRequestWithDialog.
func HandleDelRequest(
	c echo.Context,
	confirmMessage string,
	urlPath string,
	deleteFunc func() error,
	successResponse response.SuccessResponse,
	ctx *core.UiContext,
) error {
	dlg := NewDialogDelete(confirmMessage, url.NewUrl(urlPath), nil, nil, nil, nil)
	return HandleDelRequestWithDialog(c, dlg, deleteFunc, successResponse, ctx)
}

// HandleDelRequestWithDialog works like HandleDelRequest, but shows the given dialog on GET.
//
// Typically dlg is built with NewDialogDelete or NewDialogWarning and custom okText/closeText.
// The dialog is built by the caller on every request, POST included, so keep it cheap.
//
//   - GET: shows dlg
//   - POST: executes deleteFunc and returns the success response (500 on error)
func HandleDelRequestWithDialog(
	c echo.Context,
	dlg Dialog,
	deleteFunc func() error,
	successResponse response.SuccessResponse,
	ctx *core.UiContext,
) error {
	if c.Request().Method == http.MethodGet {
		return c.JSON(http.StatusOK, dlg.Print(ctx))
	}

	if err := deleteFunc(); err != nil {
		slog.Error("delete failed", "error", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "delete failed",
		})
	}

	return c.JSON(http.StatusOK, successResponse)
}
