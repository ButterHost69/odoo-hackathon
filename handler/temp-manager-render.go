package handler

import (
	"fmt"
	"log"
	"text/template"

	"github.com/ButterHost69/odoo-hackathon/db"
	"github.com/ButterHost69/odoo-hackathon/errs"
	"github.com/gin-gonic/gin"
)


func RenderManagerPage(ctx *gin.Context, email string) {
	
	ctx.Header("Content-Type", "text/html")

	tmpl := template.Must(template.ParseFiles("./templates/manager-page.html"))
	err := tmpl.Execute(ctx.Writer, nil)
	if err != nil {
		log.Println("[handler.RenderUserPage] Error while parsing manager-page.html: ", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
	}
}
