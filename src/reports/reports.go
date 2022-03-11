package reports

import (
	"net/http"
	"pitstop-api/src/users"
	"pitstop-api/src/utils"
	"strconv"
)

func GetReport(w http.ResponseWriter, r *http.Request) {
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

	report, found := loadReportFromId(id)
	if !found {
		utils.Prettier(w, "report not found !", nil, http.StatusBadRequest)
		return
	}

	if user.GetId() != report.authorId && !report.hasAccess(user) {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	utils.Prettier(w, "report "+strconv.FormatUint(uint64(id), 10), report.ToSchema(), http.StatusOK)
}

func ReplyToReport(w http.ResponseWriter, r *http.Request) {
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

	report, found := loadReportFromId(id)
	if !found {
		utils.Prettier(w, "report not found !", nil, http.StatusBadRequest)
		return
	}

	if report.authorId != user.GetId() && !report.hasAccess(user) {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	_ = r.ParseForm()
	content := r.FormValue("content")
	if len(content) == 0 {
		utils.Prettier(w, "aucun contenu fourni !", nil, http.StatusBadRequest)
		return
	}

	messageId, err := report.createReportMessage(user.GetId(), content)
	if err != nil {
		utils.Prettier(w, "échec de la création du message de report: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "message de report créé", reportMessageCreated{ReportMessageId: messageId}, http.StatusOK)
}

func GetReportMessage(w http.ResponseWriter, r *http.Request) {
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

	message, found := loadReportMessageFromId(id)
	if !found {
		utils.Prettier(w, "message not found !", nil, http.StatusBadRequest)
		return
	}

	if user.GetId() != message.authorId && !message.hasAccess(user) {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	utils.Prettier(w, "message "+strconv.FormatUint(uint64(id), 10), message.ToSchema(), http.StatusOK)
}

func CloseReport(w http.ResponseWriter, r *http.Request) {
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

	report, found := loadReportFromId(id)
	if !found {
		utils.Prettier(w, "report not found !", nil, http.StatusBadRequest)
		return
	}

	if !report.hasAccess(user) {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	err = report.close()
	if err != nil {
		utils.Prettier(w, "échec de la fermeture du report: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "report fermé !", nil, http.StatusOK)
}
