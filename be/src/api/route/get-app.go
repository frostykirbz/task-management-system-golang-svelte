package route

import (
	"database/sql"
	"fmt"

	"backend/api/middleware"
	"backend/api/models"

	"github.com/gin-gonic/gin"
)

// route: /get-all-applications
func GetAllApplications(c *gin.Context) {
	var (
		acronym     sql.NullString
		description sql.NullString
		rNumber     sql.NullInt64
		start       sql.NullString
		end         sql.NullString
	)

	rows, err := middleware.SelectAllApplications()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var data []models.Application

	for rows.Next() {

		if err := rows.Scan(&acronym, &description, &rNumber, &start, &end); err != nil {
			panic(err)
		}

		response := models.Application{
			AppAcronym:  acronym.String,
			Description: description.String,
			Rnumber:     int(rNumber.Int64),
			StartDate:   start.String,
			EndDate:     end.String,
		}
		data = append(data, response)

	}
	c.JSON(200, data)
}

// route: /get-application
func GetApplication(c *gin.Context) {
	var (
		application  models.Application
		description  sql.NullString
		rNumber      sql.NullInt64
		permitCreate sql.NullString
		permitOpen   sql.NullString
		permitToDo   sql.NullString
		permitDoing  sql.NullString
		permitDone   sql.NullString
		start        sql.NullString
		end          sql.NullString
		created      sql.NullString
	)

	if err := c.BindQuery(&application); err != nil {
		fmt.Println(err)
		middleware.ErrorHandler(c, 400, "Bad Request")
		return
	}
	application.AppAcronym = c.Query("app_acronym")
	result := middleware.SelectSingleApplication(application.AppAcronym)

	switch err := result.Scan(&description, &rNumber, &permitCreate, &permitOpen, &permitToDo, &permitDoing, &permitDone, &created, &start, &end); err {
	// Create application
	case sql.ErrNoRows:
		middleware.ErrorHandler(c, 400, "Invalid app acronym")
		return
	}

	query := models.Application{
		AppAcronym:   application.AppAcronym,
		Description:  description.String,
		Rnumber:      int(rNumber.Int64),
		StartDate:    start.String,
		EndDate:      end.String,
		PermitCreate: permitCreate.String,
		PermitOpen:   permitOpen.String,
		PermitToDo:   permitToDo.String,
		PermitDoing:  permitDoing.String,
		PermitDone:   permitDone.String,
		CreatedDate:  created.String,
	}

	c.JSON(200, query)
}
