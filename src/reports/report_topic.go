package reports

import (
	"net/http"
	"pitstop-api/src/categories/topics"
	"pitstop-api/src/users"
	"pitstop-api/src/utils"
)

func ListTopicReports(w http.ResponseWriter, r *http.Request) {
	user, exist := users.LoadUserFromRequest(r)
	if !exist {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}

	if !user.GetRank().CanViewAndReplyToTopicReports() {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	reports, err := getAllReports(topicReport)
	if err != nil {
		utils.Prettier(w, "échec du listage des reports: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "topic reports", reportList{Reports: reports}, http.StatusOK)
}

func ReportTopic(w http.ResponseWriter, r *http.Request) {
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

	target, found := topics.LoadFromId(id)
	if !found {
		utils.Prettier(w, "topic non trouvé !", nil, http.StatusBadRequest)
		return
	}

	if target.GetAuthorId() == user.GetId() {
		utils.Prettier(w, "vous ne pouvez report votre propre topic !", nil, http.StatusBadRequest)
		return
	}

	if !user.GetRank().CanReportTopic() {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	reportId, err := createReport(topicReport, target.GetId(), user.GetId())
	if err != nil {
		utils.Prettier(w, "échec de la création du report: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "report créé !", reportCreated{ReportId: reportId}, http.StatusOK)
}
