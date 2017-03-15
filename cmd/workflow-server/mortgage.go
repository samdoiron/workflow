package main

import (
	"log"

	"github.com/samdoiron/workflow"
	"gopkg.in/kataras/iris.v6"
)

// RenderMortgageSubmission renders the mortgage submission page
func RenderMortgageSubmission(ctx *iris.Context) {
	ctx.MustRender("mortgage/index.html", nil)
}

// RenderMortgageLogin renders the login page
func RenderMortgageLogin(ctx *iris.Context) {
	ctx.MustRender("mortgage/login.html", nil)
}

type application struct {
	Name              string `form:"name"`
	Phone             string `form:"phone"`
	Address           string `form:"address"`
	EmployerName      string `form:"employer_name"`
	LifeInsuranceName string `form:"life_insurance_name"`
}

// CreateMortgageApplication handle's a POST request to create an application
func CreateMortgageApplication(ctx *iris.Context) {
	var app application
	if err := ctx.ReadForm(&app); err != nil {
		renderInvalidApplication(ctx, err)
		return
	}

	id, err := workflow.SubmitMortgageApplication(workflow.MortgageApplication{
		Name:              app.Name,
		Phone:             app.Phone,
		Address:           app.Address,
		EmployerName:      app.EmployerName,
		LifeInsuranceName: app.LifeInsuranceName,
	})

	if err != nil {
		renderInvalidApplication(ctx, err)
		return
	}

	renderValidApplication(ctx, id)
}

func LogoutMortgage(ctx *iris.Context) {
	ctx.Session().Delete("mortgage_name")
	ctx.Redirect("/mortgage/login")
}

// AuthenticateMortgageLogin authenticates the submitted login form and redirects.
func AuthenticateMortgageLogin(ctx *iris.Context) {
	name := ctx.FormValue("name")
	ctx.Session().Set("mortgage_name", name)
	ctx.Redirect("/mortgage/application")
}

// RenderApplication shows the submitted application if authenticated
func RenderApplication(ctx *iris.Context) {
	name := ctx.Session().GetString("mortgage_name")
	if name == "" {
		ctx.Redirect("/mortgage/login")
		return
	}

	if app, ok := workflow.GetMortgageApplication(name); ok {
		ctx.MustRender("mortgage/application.html", app)
	} else {
		ctx.MustRender("mortgage/no_application.html", nil)
	}
}

type validApplicationPage struct {
	ID int
}

func renderValidApplication(ctx *iris.Context, id int) {
	ctx.MustRender("mortgage/success.html", validApplicationPage{
		ID: id,
	})
}

func renderInvalidApplication(ctx *iris.Context, err error) {
	log.Println("[DEBUG] invalid application:", err)
	ctx.MustRender("mortgage/invalid.html", nil)
}
