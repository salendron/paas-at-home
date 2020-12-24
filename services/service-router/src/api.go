package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type APIImpl struct {
}

func (api *APIImpl) RegisterService(ctx echo.Context) error {
	fmt.Println(ctx.Request().RemoteAddr)
	s := new(Service)
	if err := ctx.Bind(s); err != nil {
		return err
	}
	return ctx.JSON(http.StatusCreated, s)
}
