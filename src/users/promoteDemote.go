package users

import (
	"net/http"
	"pitstop-api/src/users/auth"
	"pitstop-api/src/users/ranks"
	"pitstop-api/src/utils"
)

func Promote(w http.ResponseWriter, r *http.Request) {
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

	if !adminUser.rank.CanPromoteDemoteUser() {
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

	if user.rank.GetId() == 2 {
		utils.Prettier(w, "already moderator !", nil, http.StatusBadRequest)
		return
	}

	if user.rank.GetId() == 0 {
		utils.Prettier(w, "is banned !", nil, http.StatusBadRequest)
		return
	}

	if user.rank.GetId() == 3 {
		utils.Prettier(w, "is admin !", nil, http.StatusUnauthorized)
		return
	}

	err = user.SetRank(ranks.GetRankById(2))
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	utils.Prettier(w, "user promoted !", nil, http.StatusOK)
}

func Demote(w http.ResponseWriter, r *http.Request) {
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

	if !adminUser.rank.CanPromoteDemoteUser() {
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

	if user.rank.GetId() == 1 {
		utils.Prettier(w, "is not moderator!", nil, http.StatusBadRequest)
		return
	}

	if user.rank.GetId() == 0 {
		utils.Prettier(w, "is banned !", nil, http.StatusBadRequest)
		return
	}

	if user.rank.GetId() == 3 {
		utils.Prettier(w, "is admin !", nil, http.StatusUnauthorized)
		return
	}

	err = user.SetRank(ranks.GetRankById(1))
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	utils.Prettier(w, "user demoted !", nil, http.StatusOK)
}
