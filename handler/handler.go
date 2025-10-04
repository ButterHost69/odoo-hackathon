package handler

import (
	"fmt"
	"log"

	"github.com/ButterHost69/odoo-hackathon/db"
	"github.com/ButterHost69/odoo-hackathon/errs"
	"github.com/gin-gonic/gin"
)

func CreateCompany(ctx *gin.Context) {
	// Set the Content-Type header to "text/html"
	ctx.Header("Content-Type", "text/html")
	company_name := ctx.PostForm("company-name")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	country := ctx.PostForm("country")

	// TODO: Validation

	err := db.InsertNewCompany(company_name, country, "$", email, []db.ManagerInfo{})
	if err != nil {
		log.Println("[handler.CreateCompany] Error while Inserting New Company:", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
		return
	}

	err = db.InsertNewUserAccount(email, company_name+" Admin", "Admin", "", "", 3)
	if err != nil {
		log.Println("[handler.CreateCompany] Error while Inserting New User:", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
		return
	}

	err = db.InsertNewRecordInAuthDB(email, password)
	if err != nil {
		log.Println("[handler.CreateCompany] Error while Inserting New Record in Auth Table:", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
		return
	}

	RenderAdminPage(ctx)
}

func Login(ctx *gin.Context) {
	// Set the Content-Type header to "text/html"
	ctx.Header("Content-Type", "text/html")

	email := ctx.PostForm("email")
	password := ctx.PostForm("password")

	correct_pass := db.GetPasswordByEmailFromAuth(email)
	if correct_pass == "" || password != correct_pass {
		RenderAuthPage(ctx, errs.ErrInvalidCredentials.Error())
		return
	}

	user, err := db.GetUserDetailsUsingEmail(email)
	if err != nil {
		log.Println("[handler.Login] Error while Getting User Details Using Email:", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
		return
	}

	if user.Role == "admin" {
		RenderAdminPage(ctx)
	} else {
		RenderUserPage(ctx)
	}
}
