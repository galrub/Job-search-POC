package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/galrub/go/jobSearch/config"
	"github.com/galrub/go/jobSearch/internal/database"
	"github.com/galrub/go/jobSearch/internal/logger"
	"github.com/galrub/go/jobSearch/internal/model"
	"github.com/galrub/go/jobSearch/internal/services"
	"github.com/galrub/go/jobSearch/internal/utils"
	"github.com/galrub/go/jobSearch/ui"
	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func GetMainPage(w http.ResponseWriter, r *http.Request) {
	err := utils.RenderFragment("index.html", nil, &w)
	if err != nil {
		logger.LOG.Err(err).Msg("Error parsing the index Page")
		http.Error(w, "something went wrong", http.StatusInternalServerError)
	}
}

func GetJobListFragment(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("email")
	// for testinng
	if config.DevMode() && email == "" {
		email = "galrub@gmail.com"
	}
	jobs, err := services.GetJobs(email)
	if err != nil {
		logger.LOG.Err(err).Msg("Error getting the jobs list")
		http.Error(w, "something went wrong", http.StatusInternalServerError) //TODO: render an error frgment
		return
	}
	dtos := ui.MapJobsToDto(jobs)
	err = utils.RenderFragment("jobsList.html", dtos, &w)
	if err != nil {
		logger.LOG.Err(err).Msg("Error parsing the jobs fragment")
		http.Error(w, "something went wrong", http.StatusInternalServerError) //TODO: render an error frgment
	}
}

func GetLoginFragment(w http.ResponseWriter, r *http.Request) {
	email, res := utils.LoginVerification(w, r)
	if res == false {
		utils.RenderFragment("loginForm.html", email, &w)
		return
	}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["idStrentity"] = email
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		logger.LOG.Fatal().Err(err).Msg("cannot encode generate jwt token")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cookie := http.Cookie{
		Name:     "authentication",
		Value:    t,
		Expires:  time.Now().Add(2 * time.Hour),
		HttpOnly: true,
	}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusOK)
	r.Header.Set("email", email)
	GetJobListFragment(w, r)
}

func DeleteJobWithJobsListFragment(w http.ResponseWriter, r *http.Request) {
	idStr, err := uuid.FromString(r.PathValue("id"))
	if err != nil {
		logger.LOG.Err(err).Msg("failed to parse uuidStr")
	} else {
		err = services.DeleteJob(idStr)
		if err != nil {
			logger.LOG.Err(err).Msg("failed to delete Job")
		}
	}
	GetJobListFragment(w, r)
	return
}

func CreateOrUpdateJobWithJobsListFragment(w http.ResponseWriter, r *http.Request) {
	email := r.Header.Get("email")
	// for testinng
	if config.DevMode() && email == "" {
		email = "galrub@gmail.com"
	}
	// load values
	idStr := r.PathValue("id")
	company := r.FormValue("company")
	position := r.FormValue("position")
	remote := r.FormValue("remote") == "on"
	contractType := r.FormValue("type")
	contacted := r.FormValue("contacted") == "on"
	status := r.FormValue("status")
	createDate, err := time.Parse("2006-1-2", r.FormValue("initDate"))
	if err != nil {
		logger.LOG.Err(err).Msg("eror parsing created date, will use current date")
		createDate = time.Now()
	}
	comments := r.FormValue("comments")
	var id uuid.UUID
	if idStr == "" || idStr == "nan" {
		logger.LOG.Debug().Msg("no id passed, creating new job entry")
		job := database.InsertJobParams{
			Company:       company,
			PositionDesc:  position,
			Remote:        pgtype.Bool{Bool: remote, Valid: true},
			ContractType:  contractType,
			GeneralStatus: status,
			Contacted:     pgtype.Bool{Bool: false, Valid: true},
			CreatedAt:     createDate,
			Comments:      pgtype.Text{String: comments, Valid: true},
		}
		err = services.SaveJob(&job, email)
		if err != nil {
			logger.LOG.Err(err).Msg("faild to insert new job entry")
		}
	} else if id, err = uuid.FromString(idStr); err == nil {
		job := database.UpdateJobParams{
			Company:       company,
			PositionDesc:  position,
			Remote:        pgtype.Bool{Bool: remote, Valid: true},
			ContractType:  contractType,
			GeneralStatus: status,
			CreatedAt:     createDate,
			Contacted:     pgtype.Bool{Bool: contacted, Valid: true},
			ID:            id,
			Comments:      pgtype.Text{String: comments, Valid: true},
		}
		err = services.UpdateJob(job)
	}
	if err != nil {
		logger.LOG.Err(err).Msgf("error updating job id %s", idStr)
	}
	GetJobListFragment(w, r)
}

func PrepareJobEditForNew(w http.ResponseWriter, r *http.Request) {
	logger.LOG.Debug().Msg("preapring a new job edit")
	job := model.JobEditParams{
		Company:       "",
		PositionDesc:  "",
		Remote:        false,
		ContractType:  "",
		GeneralStatus: "",
		CreatedAt:     time.Now(),
		Contacted:     false,
		ID:            "nan",
		Comments:      "",
	}

	data := struct {
		JobData      model.JobEditParams
		IsNew        bool
		ContractType model.ContractTypeSelect
	}{
		JobData: job,
		IsNew:   true,
		ContractType: model.ContractTypeSelect{
			Options: []string{"B2B", "Fulltime", "Hybrid"},
			Current: "B2B",
		},
	}
	err := utils.RenderFragments(data, &w, "jobEdit.html", "ContractTypeCombo.html")
	if err != nil {
		logger.LOG.Err(err).Msg("error creating jobEdit template")
	}
}

func PrepareJobEditForExiting(w http.ResponseWriter, r *http.Request) {
	logger.LOG.Debug().Msg("preapring a exiting job edit")
	idStr := r.PathValue("id")
	id, err := uuid.FromString(idStr)
	if err != nil {
		logger.LOG.Err(err).Msg("Error parsing uuid")
		http.Error(w, "something went wrong", http.StatusInternalServerError) //TODO: render an error frgment
		return
	}
	j, err := services.GetJob(id)
	if err != nil {
		logger.LOG.Err(err).Msg("Error getting the job data")
		http.Error(w, "something went wrong", http.StatusInternalServerError) //TODO: render an error frgment
		return
	}
	job := model.JobEditParams{
		Company:       j.Company,
		PositionDesc:  j.PositionDesc,
		Remote:        j.Remote.Bool,
		ContractType:  j.ContractType,
		GeneralStatus: j.GeneralStatus,
		CreatedAt:     j.CreatedAt,
		Contacted:     j.Contacted.Bool,
		ID:            idStr,
		Comments:      j.Comments.String,
	}

	data := struct {
		JobData      model.JobEditParams
		IsNew        bool
		ContractType model.ContractTypeSelect
	}{
		JobData: job,
		IsNew:   true,
		ContractType: model.ContractTypeSelect{
			Options: []string{"B2B", "Fulltime", "Hybrid"},
			Current: "B2B",
		},
	}
	err = utils.RenderFragments(data, &w, "jobEdit.html", "ContractTypeCombo.html")
	if err != nil {
		logger.LOG.Err(err).Msg("error creating jobEdit template")
	}

}
