package main

import (
	"fmt"

	"github.com/samdoiron/workflow"
	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
	"gopkg.in/kataras/iris.v6/adaptors/sessions"
	"gopkg.in/kataras/iris.v6/adaptors/view"
)

func main() {
	app := iris.New()

	session := sessions.New(sessions.Config{})

	app.Adapt(
		httprouter.New(),
		view.HTML("./templates", ".html"),
		iris.DevLogger(),
		session,
	)

	app.StaticWeb("/static", "./static")
	employer := app.Party("/employer").Layout("employer/layout.html")
	{
		employer.Get("/", RenderEmployerLogin)
		employer.Post("/login", AuthenticateEmployee)
		employer.Get("/home", RenderEmployerHome)
		employer.Post("/submit-to-mortgage", SubmitEmployeeInfo)
	}

	config := workflow.LoadConfig()
	app.Listen(fmt.Sprintf(":%d", config.Server.Port))
}
