package users

import (
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"pitstop-api/src/database"
	"pitstop-api/src/schemas"
	"pitstop-api/src/users/auth"
	"pitstop-api/src/users/ranks"
	"pitstop-api/src/utils"
	"strconv"
)

func HandleUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId, found := vars["id"]
	if !found {
		utils.Prettier(w, "aucun id fourni !", nil, http.StatusBadRequest)
		return
	}
	intId, err := strconv.Atoi(strId)
	if err != nil {
		utils.Prettier(w, "id invalide !", nil, http.StatusBadRequest)
		return
	}
	id := uint(intId)

	user, found := LoadUser(id)
	if !found {
		utils.Prettier(w, "Utilisateur non trouvé !", nil, http.StatusBadRequest)
		return
	}

	tokenId, err := auth.ExtractIdFromRequest(r)
	if err != nil {
		// exterior viewer
		utils.Prettier(w, "User profile", user.ToPublic(), http.StatusOK)
		return
	}
	if id == tokenId {
		// own profile
		utils.Prettier(w, "Profile information", user.ToPrivate(), http.StatusOK)
		return
	}
	sender, found := LoadUser(tokenId)
	if found {
		if sender.GetRank().CanViewUserProfile() {
			// Admin
			utils.Prettier(w, "User profile", user.ToPrivate(), http.StatusOK)
			return
		}
	}
	// exterior viewer
	utils.Prettier(w, "User profile", user.ToPublic(), http.StatusOK)
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {

	// get id with token
	user, found := LoadUserFromRequest(r)
	if !found {
		utils.Prettier(w, "token invalide !", nil, http.StatusUnauthorized)
		return
	}

	utils.Prettier(w, "Profile information", user.ToPrivate(), http.StatusOK)
}

func SearchUserByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	dbUsers := database.SearchUsersByUsername(username)

	if len(dbUsers) == 0 {
		utils.Prettier(w, "aucun utilisateur trouvé", []uint{}, http.StatusUnauthorized)
		return
	}

	var users []schemas.DbUserPublic

	for i := 0; i < len(dbUsers); i++ {
		users = append(users, dbUsers[i].ToPublic())
	}

	utils.Prettier(w, "Profiles correspondants", users, http.StatusOK)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	// get id with token
	id, errorToken := auth.ExtractIdFromRequest(r)
	if errorToken != nil {
		utils.Prettier(w, errorToken.Error(), nil, http.StatusUnauthorized)
		return
	}

	// load user from database
	user, _ := LoadUser(id)

	err := r.ParseForm()
	if err != nil {
		return
	}

	// check username
	username := r.FormValue("username")
	if invalidUserNameError := utils.IsUserNameValid(username); invalidUserNameError != nil && username != "" {
		utils.Prettier(w, invalidUserNameError.Error(), nil, http.StatusUnauthorized)
		return
	}

	// check email validity
	email := r.FormValue("email")
	if invalidEmailError := utils.IsEmailValid(email); invalidEmailError != nil && email != "" {
		utils.Prettier(w, invalidEmailError.Error(), nil, http.StatusUnauthorized)
		return
	}

	alreadyExistMsg, alreadyExist, err := database.CheckUserExistence(username, email)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	if alreadyExist {
		utils.Prettier(w, alreadyExistMsg, nil, http.StatusBadRequest)
		return
	}

	if username != "" {
		err = user.SetUsername(username)
		if err != nil {
			utils.Prettier(w, "échec de la mise à jour du nom d'utilisateur: "+err.Error(), nil, http.StatusInternalServerError)
			return
		}
	}
	if email != "" {
		err = user.SetEmail(email)
		if err != nil {
			utils.Prettier(w, "échec de la mise à jour de l'email: "+err.Error(), nil, http.StatusInternalServerError)
			return
		}
	}

	utils.Prettier(w, "User update success !", user.ToPrivate(), http.StatusOK)

}

func UpdateOwnPassword(w http.ResponseWriter, r *http.Request) {

	// get id with token
	id, errorToken := auth.ExtractIdFromRequest(r)
	if errorToken != nil {
		utils.Prettier(w, errorToken.Error(), nil, http.StatusUnauthorized)
		return
	}

	// load user from database
	user, _ := LoadUser(id)

	err := r.ParseForm()
	if err != nil {
		return
	}

	// check password validity
	password := r.FormValue("password")
	if invalidPasswordError := utils.IsPasswordValid(password); invalidPasswordError != nil {
		utils.Prettier(w, invalidPasswordError.Error(), nil, http.StatusBadRequest)
		return
	}

	// check hash password
	compareErr := bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(password))
	if compareErr == nil {
		utils.Prettier(w, "it's the same password", nil, http.StatusBadRequest)
		return
	}

	hashPassword, err := auth.HashPassword(password)
	if err != nil {
		utils.Prettier(w, "échec de la mise à jour du mot de passe !", nil, http.StatusInternalServerError)
		return
	}

	// update password
	if setPasswordError := user.SetPassword(hashPassword); setPasswordError != nil {
		utils.Prettier(w, setPasswordError.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.Prettier(w, "mot de passe mis à jour avec succès  !", user.ToPrivate(), http.StatusOK)
}

func GetAvatar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId, found := vars["id"]
	if !found {
		utils.Prettier(w, "aucun id fourni !", nil, http.StatusBadRequest)
		return
	}
	intId, err := strconv.Atoi(strId)
	if err != nil {
		utils.Prettier(w, "id invalide !", nil, http.StatusBadRequest)
		return
	}
	id := uint(intId)
	_, found = LoadUser(id)
	if !found {
		utils.Prettier(w, "Utilisateur non trouvé", nil, http.StatusBadRequest)
		return
	}

	printAvatar(w, strId)
}

func GetOwnAvatar(w http.ResponseWriter, r *http.Request) {
	// get id with token
	id, err := auth.ExtractIdFromRequest(r)
	if err != nil {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	printAvatar(w, idToString(id))
}

func UpdateOwnAvatar(w http.ResponseWriter, r *http.Request) {
	err, status := updateAvatar(r)
	if err != nil {
		utils.Prettier(w, fmt.Sprintf("échec de la mise à jour de l'avatar: %v", err), nil, status)
		return
	}
	utils.Prettier(w, "avatar mis à jour avec succès !", nil, http.StatusOK)
}

func DeleteAvatar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId, found := vars["id"]
	if !found {
		utils.Prettier(w, "aucun id fourni !", nil, http.StatusBadRequest)
		return
	}
	intId, err := strconv.Atoi(strId)
	if err != nil {
		utils.Prettier(w, "id invalide !", nil, http.StatusBadRequest)
		return
	}
	id := uint(intId)
	_, found = LoadUser(id)
	if !found {
		utils.Prettier(w, "Utilisateur non trouvé", nil, http.StatusBadRequest)
		return
	}

	askingUser, ok := LoadUserFromRequest(r)
	if !ok {
		utils.Prettier(w, "token invalide !", nil, http.StatusBadRequest)
		return
	}

	if askingUser.GetId() == id || askingUser.GetRank().CanDeleteUsersAvatar() {
		err = deleteAvatar(id)
		if err != nil {
			utils.Prettier(w, "échec de la suppression de l'avatar", nil, http.StatusInternalServerError)
			return
		}
		utils.Prettier(w, "avatar supprimé avec succès !", nil, http.StatusOK)
		return
	}

	utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
}

func DeleteOwnAvatar(w http.ResponseWriter, r *http.Request) {
	// get id with token
	id, err := auth.ExtractIdFromRequest(r)
	if err != nil {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	err = deleteAvatar(id)
	if err != nil {
		utils.Prettier(w, "échec de la suppression de l'avatar: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}
	utils.Prettier(w, "avatar supprimé avec succès !", nil, http.StatusOK)
}

func HasPermission(w http.ResponseWriter, r *http.Request) { // move to new function
	user, exist := LoadUserFromRequest(r)
	if !exist {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	permission, found := vars["permission"]
	if !found {
		utils.Prettier(w, "aucune permission fournie !", nil, http.StatusBadRequest)
		return
	}

	utils.Prettier(w, "user permission", struct {
		HasPermission bool `json:"has_permission"`
	}{HasPermission: ranks.HasPermission(permission, user.GetRank())}, http.StatusOK)
}

func UpdateSettings(w http.ResponseWriter, r *http.Request) {
	user, exist := LoadUserFromRequest(r)
	if !exist {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	toUpdate := vars["setting"]

	switch toUpdate {
	case "darkTheme":
		user.settings.ToggleDarkTheme(user.id)
	case "notifications":
		user.settings.ToggleNotifications(user.id)
	case "blockFriendRequests":
		user.settings.ToggleBlockFriendRequests(user.id)
	default:
		utils.Prettier(w, "paramètre invalide", nil, http.StatusBadRequest)
		return
	}

	utils.Prettier(w, "paramètres mis à jour avec succès !", nil, http.StatusOK)

}

func UpdatePassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	vars := mux.Vars(r)
	token := vars["token"]

	id, err := auth.ExtractIdFromToken(token)

	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	user, exist := LoadUser(id)
	if !exist {
		utils.Prettier(w, "Utilisateur non trouvé !", nil, http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")
	if invalidPasswordError := utils.IsPasswordValid(password); invalidPasswordError != nil {
		utils.Prettier(w, invalidPasswordError.Error(), nil, http.StatusBadRequest)
		return
	}

	hashPassword, err := auth.HashPassword(password)
	if err != nil {
		utils.Prettier(w, "échec de la mise à jour du mot de passe", nil, http.StatusInternalServerError)
		return
	}

	err = user.SetPassword(hashPassword)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.Prettier(w, "password updated", nil, http.StatusOK)
}

func ValidPassword(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	id, err := auth.ExtractIdFromRequest(r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	user, exist := LoadUser(id)
	if !exist {
		utils.Prettier(w, "mot de passe invalide !", nil, http.StatusBadRequest)
		return
	}

	password := r.FormValue("password")
	compareErr := bcrypt.CompareHashAndPassword([]byte(user.GetPassword()), []byte(password))
	if compareErr != nil {
		utils.Prettier(w, "mot de passe invalide !", nil, http.StatusUnauthorized)
		return
	}

	utils.Prettier(w, "valid password !", nil, http.StatusOK)

}
