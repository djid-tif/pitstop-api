package friends

import (
	"github.com/gorilla/mux"
	"net/http"
	"pitstop-api/src/database"
	"pitstop-api/src/users"
	_ "pitstop-api/src/users/auth"
	"pitstop-api/src/utils"
)

func AddFriend(w http.ResponseWriter, r *http.Request) {
	user, exist := users.LoadUserFromRequest(r)
	if !exist {
		utils.Prettier(w, "token invalide !", nil, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	idStr, found := vars["id"]
	if !found {
		utils.Prettier(w, "aucun id fourni !", nil, http.StatusBadRequest)
		return
	}

	id, err := utils.StringToUint(idStr)
	if err != nil {
		utils.Prettier(w, "id invalide !", nil, http.StatusBadRequest)
		return
	}

	if id == user.GetId() {
		utils.Prettier(w, "vous ne pouvez envoyer une demande d'ami à vous-même !", nil, http.StatusBadRequest)
		return
	}

	fList := user.GetFriendsList()
	if utils.ContainsUint(id, fList) {
		utils.Prettier(w, "vous avez déjà cet ami !", nil, http.StatusBadRequest)
		return
	}

	sentList := database.GetFriendsList(user.GetId(), typeFriendSent)
	if utils.ContainsUint(id, sentList) {
		utils.Prettier(w, "vous avez déjà envoyé une demande d'ami à cet utilisateur !", nil, http.StatusBadRequest)
		return
	}

	friend, found := users.LoadUser(id)
	if !found {
		utils.Prettier(w, "cet utilisateur n'existe pas !", nil, http.StatusBadRequest)
		return
	}

	if friend.Settings().AreFriendRequestsBlocked() {
		utils.Prettier(w, "cet utilisateur n'accepte pas les demandes d'ami !", nil, http.StatusOK)
		return
	}

	receivedList := database.GetFriendsList(user.GetId(), typeFriendReceived)
	if utils.ContainsUint(id, receivedList) {

		err = confirmFriendship(user, friend)
		if err != nil {
			utils.Prettier(w, "échec de la confirmation de l'amitié: "+err.Error(), nil, http.StatusInternalServerError)
			return
		}
		utils.Prettier(w, "vous avez accepté cette demande d'ami !", nil, http.StatusOK)
		return
	}

	err = sendFriendRequest(user, friend)
	if err != nil {
		utils.Prettier(w, "error: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "demande d'ami envoyée avec succès !", nil, http.StatusOK)
}

func RemoveFriend(w http.ResponseWriter, r *http.Request) {
	user, exist := users.LoadUserFromRequest(r)
	if !exist {
		utils.Prettier(w, "token invalide !", nil, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	idStr, found := vars["id"]
	if !found {
		utils.Prettier(w, "aucun id fourni !", nil, http.StatusBadRequest)
		return
	}

	id, err := utils.StringToUint(idStr)
	if err != nil {
		utils.Prettier(w, "id invalide !", nil, http.StatusBadRequest)
		return
	}

	if utils.ContainsUint(id, database.GetFriendsList(user.GetId(), typeFriendSent)) {
		if err = removeFriend(user.GetId(), id, typeFriendSent); err != nil {
			utils.Prettier(w, "échec de l'annulation de la demande d'ami !", nil, http.StatusInternalServerError)
			return
		}
		if err = removeFriend(id, user.GetId(), typeFriendReceived); err != nil {
			utils.Prettier(w, "échec de l'annulation de la demande d'ami !", nil, http.StatusInternalServerError)
			return
		}
		utils.Prettier(w, "demande annulée avec succès !", nil, http.StatusOK)
		return
	}

	if utils.ContainsUint(id, database.GetFriendsList(user.GetId(), typeFriendReceived)) {
		if err = removeFriend(user.GetId(), id, typeFriendReceived); err != nil {
			utils.Prettier(w, "échec du refus de la demande d'ami !", nil, http.StatusInternalServerError)
			return
		}
		if err = removeFriend(id, user.GetId(), typeFriendSent); err != nil {
			utils.Prettier(w, "échec du refus de la demande d'ami !", nil, http.StatusInternalServerError)
			return
		}
		utils.Prettier(w, "requête refusée avec succès !", nil, http.StatusOK)
		return
	}

	if !utils.ContainsUint(id, user.GetFriendsList()) {
		utils.Prettier(w, "vous n'avez pas cet ami !", nil, http.StatusBadRequest)
		return
	}

	if err = removeFriend(user.GetId(), id, typeFriendConfirmed); err != nil {
		utils.Prettier(w, "échec du retrait de cet ami !", nil, http.StatusInternalServerError)
		return
	}
	if err = removeFriend(id, user.GetId(), typeFriendConfirmed); err != nil {
		utils.Prettier(w, "échec du retrait de cet ami !", nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "ami supprimé avec succès !", nil, http.StatusOK)
}

func ListFriends(w http.ResponseWriter, r *http.Request) {
	user, exist := users.LoadUserFromRequest(r)
	if !exist {
		utils.Prettier(w, "token invalide !", nil, http.StatusUnauthorized)
		return
	}

	utils.Prettier(w, "friends lists", friendsLists{
		ConfirmedFriends: user.GetFriendsList(),
		SentFriends:      database.GetFriendsList(user.GetId(), typeFriendSent),
		ReceivedFriends:  database.GetFriendsList(user.GetId(), typeFriendReceived),
	}, http.StatusOK)
}
