package handler

import (
	"fmt"
	"log"
	"text/template"

	"github.com/ButterHost69/odoo-hackathon/db"
	"github.com/ButterHost69/odoo-hackathon/errs"
	"github.com/ButterHost69/odoo-hackathon/utils"
	"github.com/gin-gonic/gin"
)

func RenderInitPage(ctx *gin.Context) {
	session_token := utils.GetSessionTokenFromCookie(ctx.Request)

	email, err := db.GetEmailUsingSessionToken(session_token)
	if err != nil {
		if err == errs.ErrSessionTokenDoesNotExist {
			RenderAuthPage(ctx, "")
			fmt.Print("Haa")
		} else {
			log.Println("[handler.RenderInitPage] Error while Getting Email from Session Token:", err)
			fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
		}
		return
	}
	fmt.Print("Yo")
	user, err := db.GetUserDetailsUsingEmail(email)
	if err != nil {
		log.Println("[handler.RenderInitPage] Error while Getting User Details Using Email:", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
		return
	}
	fmt.Print(user)

	if user.Role != "admin" {
		RenderUserPage(ctx)
		return
	}
	users, err := db.GetAllUsersDetailsUsingCompanyID(user.CompanyID)
	if err != nil {
		log.Println("[handler.RenderInitPage] Error while Getting All Users Details Using Company ID:", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
		return
	}
	RenderAdminPage(ctx, users)
}

func RenderAuthPage(ctx *gin.Context, error_msg string) {
	ctx.Header("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("./templates/auth.html"))
	err := tmpl.Execute(ctx.Writer, error_msg)
	if err != nil {
		log.Println("[handler.RenderAuthPage] Error while parsing auth.html:", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
	}
}

func RenderAdminPage(ctx *gin.Context, users []db.User) {
	ctx.Header("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("./templates/admin-page.html"))
	err := tmpl.Execute(ctx.Writer, users)
	if err != nil {
		log.Println("[handler.RenderAdminPage] Error while parsing admin-page.html: ", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
	}
}

func RenderEmployeePage(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("./templates/user-page.html"))
	err := tmpl.Execute(ctx.Writer, nil)
	if err != nil {
		log.Println("[handler.RenderUserPage] Error while parsing user-page.html: ", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
	}
}

func RenderManagerPage(ctx *gin.Context) {
	ctx.Header("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("./templates/manager-page.html"))
	err := tmpl.Execute(ctx.Writer, nil)
	if err != nil {
		log.Println("[handler.RenderUserPage] Error while parsing manager-page.html: ", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
	}
}
