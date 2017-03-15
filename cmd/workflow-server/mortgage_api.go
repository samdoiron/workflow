package main

import (
	"log"

	"github.com/samdoiron/workflow"
	"gopkg.in/kataras/iris.v6"
)

type APIResponse struct {
	OK      bool
	Message string
}

func APIPutEmployerInfo(ctx *iris.Context) {
	id, err := ctx.ParamInt64("id")
	if err != nil {
		apiBadRequest(ctx, "Invalid ;d")
		return
	}

	var info workflow.MortgageEmployerInfo
	if err := ctx.ReadJSON(&info); err != nil {
		apiBadRequest(ctx, "Malformed Request")
		return
	}

	err = workflow.SetEmployerInfo(id, info)
	if err != nil {
		log.Println("[DEBUG] Failed to set employer info", err)
		apiBadRequest(ctx, "Failed to record data")
		return
	}

	ctx.JSON(iris.StatusOK, APIResponse{
		OK:      true,
		Message: "Successfully recorded",
	})
}

func apiBadRequest(ctx *iris.Context, message string) {
	ctx.JSON(iris.StatusBadRequest, APIResponse{
		OK:      false,
		Message: message,
	})
}
