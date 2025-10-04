package handler

import (
	"fmt"
	"log"
	"text/template"
	"time"

	"github.com/ButterHost69/odoo-hackathon/db"
	"github.com/ButterHost69/odoo-hackathon/errs"
	"github.com/gin-gonic/gin"
)

type ManagerView struct {
	ExpenseID      int
	EmployeeEmail  string
	Description    string
	ExpenseDate    time.Time
	Category       string
	Amount         int
	Remarks        string
	Client_Status  string
	Manager_Status string
}

func RenderManagerPage(ctx *gin.Context, email string) {
	expense_ids, err := db.GetExpenseIDsByManagerEmail(email)
	if err != nil {
		log.Println("[handler.RenderManagerPage] Error While Fetching expense_id:", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
		return
	}

	var expense_info_list []ManagerView
	for _, expense_id := range expense_ids {
		expense_info, err := db.GetExpenseUsingExpenseID(expense_id)
		if err != nil {
			log.Println("[handler.RenderManagerPage] Error While Fetching expense info using expense id:", err)
			fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
			continue
		}

		app_status, err := db.GetApprovalStatusByExpenseID(expense_id)
		if err != nil {
			log.Println("[handler.RenderManagerPage] Error While Fetching expense info using expense id:", err)
			fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
			continue
		}

		expense_info_list = append(expense_info_list, ManagerView{
			ExpenseID: expense_id,
			EmployeeEmail: email,
			Description: expense_info.Description,
			ExpenseDate: expense_info.ExpenseDate,
			Category: expense_info.Category,
			Amount: expense_info.Amount,
			Remarks: "N",
			Client_Status: expense_info.Status,
			Manager_Status: app_status.Status,
		})
	}

	ctx.Header("Content-Type", "text/html")
	tmpl := template.Must(template.ParseFiles("./templates/manager-page.html"))
	err = tmpl.Execute(ctx.Writer, expense_info_list)
	if err != nil {
		log.Println("[handler.RenderUserPage] Error while parsing manager-page.html: ", err)
		fmt.Fprint(ctx.Writer, errs.INTERNAL_SERVER_ERROR_MESSAGE)
	}
}
