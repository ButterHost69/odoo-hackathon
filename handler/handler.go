package handler

import (
	"fmt"
	"log"

	"github.com/ButterHost69/odoo-hackathon/db"
	"github.com/ButterHost69/odoo-hackathon/errs"
	"github.com/ButterHost69/odoo-hackathon/utils"
	"github.com/gin-gonic/gin"
)

func CreateCompany(ctx *gin.Context) {
	// Set the Content-Type header to "text/html"
	ctx.Header("Content-Type", "text/html")
	company_name := ctx.PostForm("company-name")
	email := ctx.PostForm("email")
	password := ctx.PostForm("password")
	country := ctx.PostForm("country")

	currency_symbol, err := utils.GetCurrencyUsingCountryName(country)
	if err != nil {
		log.Println("[handler.CreateCompany] Error while Getting Currency using Country Name:", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
		return
	}

	company_id, err := db.InsertNewCompany(company_name, country, currency_symbol, email, []db.ManagerInfo{})
	if err != nil {
		log.Println("[handler.CreateCompany] Error while Inserting New Company:", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
		return
	}

	err = db.InsertNewUserAccount(email, company_name+" Admin", "admin", "", "", company_id)
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

	RenderAuthPage(ctx, "")
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

	new_token, err := db.GenerateToken()
	if err != nil {
		log.Println("[handler.Login] Error while Generating New Token:", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
		return
	}

	err = db.UpdateSessionTokenInAuthDB(email, new_token)
	if err != nil {
		log.Println("[handler.Login] Error while Updating Session Token:", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
		return
	}

	utils.SetSessionTokenInCookie(ctx.Writer, new_token)

	users, err := db.GetAllUsersDetailsUsingCompanyID(user.CompanyID)
	if err != nil {
		log.Println("[handler.Login] Error while Getting All Users Details Using Company ID:", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
		return
	}

	if user.Role == "admin" {
		RenderAdminPage(ctx, users)
	} else {
		RenderUserPage(ctx)
	}
}
