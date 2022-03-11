package users

import (
	"net/http"
	"pitstop-api/src/database"
	"pitstop-api/src/users/auth"
	"pitstop-api/src/utils"
)

func BanUser(w http.ResponseWriter, r *http.Request) {

	adminId, err := auth.ExtractIdFromRequest(r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	adminUser, exist := LoadUser(adminId)
	if !exist {
		utils.Prettier(w, "vous n'avez pas la permission", nil, http.StatusUnauthorized)
		return
	}

	if !adminUser.rank.CanBanUser() {
		utils.Prettier(w, "vous n'avez pas la permission", nil, http.StatusUnauthorized)
		return
	}

	id, err := utils.ExtractUintFromRequest("id", r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	user, exist := LoadUser(id)
	if !exist {
		utils.Prettier(w, "user not exist", nil, http.StatusUnauthorized)
		return
	}

	if user.rank.GetId() == 0 {
		utils.Prettier(w, "already ban !", nil, http.StatusBadRequest)
		return
	}

	if id == adminId {
		utils.Prettier(w, "Chehhh !", nil, http.StatusUnauthorized)
		return
	}

	err = database.UpdateUserById(id, "rank", 0)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	utils.Prettier(w, "user banned !", nil, http.StatusOK)

}

func UnBanUser(w http.ResponseWriter, r *http.Request) {

	adminId, err := auth.ExtractIdFromRequest(r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	adminUser, exist := LoadUser(adminId)
	if !exist {
		utils.Prettier(w, "vous n'avez pas la permission", nil, http.StatusUnauthorized)
		return
	}

	if !adminUser.rank.CanBanUser() {
		utils.Prettier(w, "vous n'avez pas la permission", nil, http.StatusUnauthorized)
		return
	}

	id, err := utils.ExtractUintFromRequest("id", r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	user, exist := LoadUser(id)
	if !exist {
		utils.Prettier(w, "user not exist", nil, http.StatusUnauthorized)
		return
	}

	if user.rank.GetId() != 0 {
		utils.Prettier(w, "is not ban !", nil, http.StatusBadRequest)
		return
	}

	err = database.UpdateUserById(id, "rank", 1)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	utils.Prettier(w, "user unbanned !", nil, http.StatusOK)

}
