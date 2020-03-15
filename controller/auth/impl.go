package auth

import (
	"net/http"
	"time"

	"github.com/acm-uiuc/core/context"
	"github.com/acm-uiuc/core/service"
)

type AuthController struct {
	svc *service.Service
}

func New(svc *service.Service) *AuthController {
	return &AuthController{
		svc: svc,
	}
}

func (controller *AuthController) GetOAuthRedirect(ctx *context.Context) error {
	provider := ctx.Param("provider")

	uri, err := controller.svc.Auth.GetOAuthRedirect(provider)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Invalid Provider")
	}

	return ctx.Redirect(http.StatusFound, uri)
}

func (controller *AuthController) GetOAuthRedirectLanding(ctx *context.Context) error {
	params := struct {
		Provider string
	}{
		Provider: "google",
	}

	return ctx.Render(http.StatusOK, "redirect", params)
}

func (controller *AuthController) GetToken(ctx *context.Context) error {
	provider := ctx.Param("provider")

	req := struct {
		Code string `json:"code"`
	}{}

	err := ctx.Bind(&req)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Failed Bind")
	}

	token, err := controller.svc.Auth.Authorize(provider, req.Code)
	if err != nil {
		return ctx.String(http.StatusBadRequest, "Failed Token Generation")
	}

	ctx.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    token.Token,
		Expires:  time.Unix(token.Expiration, 0),
		Path:     "/",
		HttpOnly: true,
	})

	return ctx.JSON(http.StatusOK, token)
}
