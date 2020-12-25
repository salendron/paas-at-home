package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func newBool(b bool) *bool {
	return &b
}

func newInt(i int) *int {
	return &i
}

type API struct {
	Storage StorageInterface
}

func (api *API) SetStorage(storage StorageInterface) {
	api.Storage = storage
}

func (api *API) RegisterService(ctx echo.Context) error {
	s := new(Service)
	if err := ctx.Bind(s); err != nil {
		return RaiseError(ctx, fmt.Sprintf("Failed to parse request body. %v", err.Error()), http.StatusBadRequest, ErrorCodeInvalidRequestBody)
	}

	s.IsHealthy = newBool(true)
	s.Address = &(ctx.Request().RemoteAddr)
	s.Latency = newInt(0)

	if err := api.Storage.SaveService(s); err != nil {
		log.Errorf("Failed to register Service. %v", err.Error())
		return err
	}

	return ctx.JSON(http.StatusCreated, s)
}
