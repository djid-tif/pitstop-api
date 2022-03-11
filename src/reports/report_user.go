package reports

import (
	"net/http"
	"pitstop-api/src/users"
	"pitstop-api/src/utils"
)

func ListUserReports(w http.ResponseWriter, r *http.Request) {
	user, exist := users.LoadUserFromRequest(r)
	if !exist {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}

	if !user.GetRank().CanViewAndReplyToUserReports() {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	reports, err := getAllReports(userReport)
	if err != nil {
		utils.Prettier(w, "échec du listage des reports: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "user reports", reportList{Reports: reports}, http.StatusOK)
}

func ReportUser(w http.ResponseWriter, r *http.Request) {
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

	if id == user.GetId() {
		utils.Prettier(w, "vous ne pouvez pas vous report vous-même !", nil, http.StatusBadRequest)
		return
	}

	target, found := users.LoadUser(id)
	if !found {
		utils.Prettier(w, "Utilisateur non trouvé !", nil, http.StatusBadRequest)
		return
	}

	if !user.GetRank().CanReportUser() {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	reportId, err := createReport(userReport, target.GetId(), user.GetId())
	if err != nil {
		utils.Prettier(w, "échec de la création du report: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	user.Stats().ReportUser(user.GetId())

	utils.Prettier(w, "report créé !", reportCreated{ReportId: reportId}, http.StatusOK)
}
