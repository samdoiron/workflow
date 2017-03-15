package main

import (
	"log"
	"strconv"

	"github.com/samdoiron/workflow"
	"gopkg.in/kataras/iris.v6"
)

const sessionKey = "employee_id"

// RenderEmployerLogin renders the login for the employer portal
func RenderEmployerLogin(ctx *iris.Context) {
	ctx.MustRender("employer/login.html", nil)
}

// AuthenticateEmployee tries to log in an employee from a form submission
func AuthenticateEmployee(ctx *iris.Context) {
	name := ctx.FormValue("account_id")
	password := ctx.FormValue("account_password")
	if workflow.AuthenticateEmployee(name, password) {
		ctx.Session().Set(sessionKey, name)
		ctx.Redirect("/employer/home")
	} else {
		ctx.Redirect("/employer")
	}
}

type HomePage struct {
	Employee workflow.Employee
	Error    string
	Info     string
}

// RenderEmployerHome renders the main page of the portal for an authenticated
// user.
func RenderEmployerHome(ctx *iris.Context) {
	id := ctx.Session().GetString(sessionKey)
	if id == "" {
		ctx.Redirect("/employer")
		return
	}

	if employee, ok := workflow.GetEmployee(id); ok {
		ctx.MustRender("employer/home.html", HomePage{
			Employee: employee,
			Error:    ctx.Session().GetFlashString("error"),
			Info:     ctx.Session().GetFlashString("info"),
		})
	} else {
		ctx.Redirect("/employer")
	}
}

// SubmitEmployeeInfo submits employee information for a mortgage application
func SubmitEmployeeInfo(ctx *iris.Context) {
	employee, ok := ensureLogin(ctx)
	if !ok {
		ctx.Redirect("/employer/home")
	}

	appID, err := strconv.Atoi(ctx.FormValue("application_id"))
	if err != nil {
		ctx.Session().SetFlash("error", "Invalid application id")
		ctx.Redirect("/employer/home")
		return
	}

	response, err := workflow.PutEmployerInfo(appID, workflow.MortgageEmployerInfo{
		Name:           employee.Name,
		YearlySalary:   employee.YearlySalary,
		YearsOfService: employee.YearsOfService,
		Position:       employee.Position,
	})
	if err != nil {
		log.Println("[DEBUG] Failed to put employer info:", err)
		ctx.Redirect("/employer/home")
	}

	if response.OK {
		log.Println("[DEBUG] writing info flash")
		ctx.Session().SetFlash("info", "Successfully submitted")
		ctx.Redirect("/employer/home")
	} else {
		ctx.Session().SetFlash("error", "Failed to submit application")
		ctx.Redirect("/employer/home")
	}
}

func ensureLogin(ctx *iris.Context) (workflow.Employee, bool) {
	id := ctx.Session().GetString(sessionKey)
	if id == "" {
		return workflow.Employee{}, false
	}

	return workflow.GetEmployee(id)
}
