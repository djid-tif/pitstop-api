package reports

import (
	"net/http"
	"pitstop-api/src/categories/topics/posts"
	"pitstop-api/src/users"
	"pitstop-api/src/utils"
)

func ListPostReports(w http.ResponseWriter, r *http.Request) {
	user, exist := users.LoadUserFromRequest(r)
	if !exist {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}

	if !user.GetRank().CanViewAndReplyToPostReports() {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	reports, err := getAllReports(postReport)
	if err != nil {
		utils.Prettier(w, "échec du listage des reports: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "post reports", reportList{Reports: reports}, http.StatusOK)
}

func ReportPost(w http.ResponseWriter, r *http.Request) {
	user, exist := users.LoadUserFromRequest(r)
	if !exist {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}

	id, err := utils.ExtractUintFromRequest("id", r)
	if err != nil {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}

	target, found := posts.LoadFromId(id)
	if !found {
		utils.Prettier(w, "post non trouvé !", nil, http.StatusBadRequest)
		return
	}

	if target.GetAuthorId() == user.GetId() {
		utils.Prettier(w, "vous ne pouvez report votre propre post !", nil, http.StatusBadRequest)
		return
	}

	if !user.GetRank().CanReportPost() {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	reportId, err := createReport(postReport, target.GetId(), user.GetId())
	if err != nil {
		utils.Prettier(w, "échec de la création du report: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	user.Stats().ReportPost(user.GetId())

	utils.Prettier(w, "report créé !", reportCreated{ReportId: reportId}, http.StatusOK)
}
