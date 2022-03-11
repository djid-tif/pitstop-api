package auth

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
	"pitstop-api/src/config"
	"pitstop-api/src/database"
	"pitstop-api/src/mail"
	"pitstop-api/src/schemas"
	//"pitstop-api/src/users"
	"pitstop-api/src/utils"
	"time"
)

var blackListAccessToken []string
var blackListRefreshToken []string

type Data struct {
	RefreshToken string `json:"refresh_token"`
}

func Register(w http.ResponseWriter, r *http.Request) {

	_ = r.ParseForm()

	username := r.FormValue("username")
	if invalidUserNameError := utils.IsUserNameValid(username); invalidUserNameError != nil {
		utils.Prettier(w, invalidUserNameError.Error(), nil, http.StatusBadRequest)
		return
	}

	// Check email validity
	email := r.FormValue("email")
	if invalidEmailError := utils.IsEmailValid(email); invalidEmailError != nil {
		utils.Prettier(w, invalidEmailError.Error(), nil, http.StatusBadRequest)
		return
	}

	takenField, alreadyExist, err := database.CheckAwaitUserExistence(username, email)
	if err != nil {
		utils.Prettier(w, takenField+err.Error(), nil, http.StatusInternalServerError)
		return
	}
	if alreadyExist {
		utils.Prettier(w, "Cet "+takenField+" est déjà enregistré", nil, http.StatusBadRequest) // BadRequest ?
		return
	}

	takenField, alreadyExist, err = database.CheckUserExistence(username, email)
	if err != nil {
		utils.Prettier(w, takenField+err.Error(), nil, http.StatusInternalServerError)
		return
	}
	if alreadyExist {
		utils.Prettier(w, "Cet "+takenField+" est déjà enregistré", nil, http.StatusBadRequest) // BadRequest ?
		return
	}

	// Check password validity
	password := r.FormValue("password")
	if invalidPasswordError := utils.IsPasswordValid(password); invalidPasswordError != nil {
		utils.Prettier(w, invalidPasswordError.Error(), nil, http.StatusBadRequest)
		return
	}

	// Check hash password
	hashedPassword, err := HashPassword(password)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	err = database.CreateAwaitUser(username, email, hashedPassword)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	users := database.GetAwaitUsersByUsername(username)
	if len(users) == 0 {
		utils.Prettier(w, "utilisateur non créé", nil, http.StatusInternalServerError)
		return
	}

	user := users[0]

	// Create the token
	token := jwt.New()
	_ = token.Set("username", user.Username)
	_ = token.Set(jwt.ExpirationKey, time.Now().Add(expirationActivateToken))

	// Sign the token and generate a payload
	signed, err := jwt.Sign(token, jwa.RS256, privKey)
	if err != nil {
		utils.PrintError(err)
		utils.Prettier(w, fmt.Sprintf("échec de la génération de la payload chiffrée: %s\n", err), nil, http.StatusInternalServerError)
		return
	}
	tokenStr := string(signed)

	// Check if Email exist
	err = mail.SendTemplate(username, email, "Vérification de l'email", "confirm_email", schemas.ConfirmEmail{
		Username:    username,
		ConfirmLink: template.URL(config.OriginServerFront + "/activation/" + tokenStr),
		PitStopLink: config.OriginServerFront,
		SendingDate: time.Now().Format("02/01/2006 15:04:05"),
	})
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	utils.Prettier(w, "utilisateur créé !", nil, http.StatusOK)
}

func Activate(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	tokenAwait := vars["token"]

	username, err := ExtractUsernameFromToken(tokenAwait)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusUnauthorized)
		return
	}

	awaitUsers := database.GetAwaitUsersByUsername(username)
	if len(awaitUsers) == 0 {
		utils.Prettier(w, "Utilisateur non trouvé", nil, http.StatusBadRequest)
		return
	}

	awaitUser := awaitUsers[0]

	// create new user
	err = database.CreateUser(awaitUser.Username, awaitUser.Email, awaitUser.Password)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	// delete await user
	err = database.DeleteAwaitUserBy(awaitUser.Username)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "User activated !", nil, http.StatusOK)
}

func Login(w http.ResponseWriter, r *http.Request) {

	var dbUsers []schemas.DbUser
	var dbAwaitUsers []schemas.DbAwaitUser

	_ = r.ParseForm()

	// check login (username or email) validity
	login := r.FormValue("login")
	if invalidEmailError := utils.IsEmailValid(login); invalidEmailError == nil {

		dbUsers = database.GetUsersBy("email", login)
		dbAwaitUsers = database.GetAwaitUsersBy("email", login)

	} else if invalidUsernameError := utils.IsUserNameValid(login); invalidUsernameError == nil {

		dbUsers = database.GetUsersBy("username", login)
		dbAwaitUsers = database.GetAwaitUsersBy("username", login)

	} else {
		utils.Prettier(w, "login ou mot de passe incorrect !", nil, http.StatusUnauthorized)
		return
	}

	// get password
	password := r.FormValue("password")

	// verify if account was activate
	if len(dbAwaitUsers) != 0 {
		utils.Prettier(w, "veuillez valider votre email", nil, http.StatusUnauthorized)
		return
	}

	// check if user exist
	if len(dbUsers) == 0 {
		utils.Prettier(w, "login ou mot de passe incorrect !", nil, http.StatusUnauthorized)
		return
	}

	// get user
	dbUser := dbUsers[0]

	// check hash password
	compareErr := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(password))
	if compareErr != nil {
		utils.Prettier(w, "login ou mot de passe incorrect !", nil, http.StatusUnauthorized)
		return
	}

	// Create the RefreshToken
	refreshToken := jwt.New()
	_ = refreshToken.Set("id", dbUser.Id)
	_ = refreshToken.Set(jwt.ExpirationKey, time.Now().Add(expirationRefreshToken))

	// Sign the token and generate a payload
	signed, err := jwt.Sign(refreshToken, jwa.RS256, privKey)
	if err != nil {
		utils.PrintError(err)
		utils.Prettier(w, fmt.Sprintf("échec de la génération de la payload chiffrée: %s\n", err), nil, http.StatusInternalServerError)
		return
	}
	refreshTokenStr := string(signed)

	// Create the accessToken
	accessToken := jwt.New()
	_ = accessToken.Set("id", dbUser.Id)
	_ = accessToken.Set("username", dbUser.Username)
	_ = accessToken.Set("email", dbUser.Email)
	_ = accessToken.Set(jwt.ExpirationKey, time.Now().Add(expirationAccessToken))

	// Sign the token and generate a payload
	signed, err = jwt.Sign(accessToken, jwa.RS256, privKey)
	if err != nil {
		utils.PrintError(err)
		utils.Prettier(w, fmt.Sprintf("échec de la génération de la payload chiffrée: %s\n", err), nil, http.StatusInternalServerError)
		return
	}
	accessTokenStr := string(signed)

	utils.Prettier(w, "Token generated !", struct {
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
	}{refreshTokenStr, accessTokenStr}, http.StatusOK)
}

func DecryptToken(w http.ResponseWriter, r *http.Request) {
	id, _ := ExtractIdFromRequest(r)
	utils.Prettier(w, "Good !", id, http.StatusOK)
}

func Forgot(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		return
	}

	// get email from request
	email := r.FormValue("email")
	if invalidEmailMsg := utils.IsEmailValid(email); invalidEmailMsg != nil {
		utils.Prettier(w, "email invalide !", nil, http.StatusBadRequest)
		return
	}

	// get user from database
	users := database.GetUsersBy("email", email)

	// check if user exist
	if len(users) == 0 {
		utils.Prettier(w, "Si le compte a été trouvé, un email a été envoyé", nil, http.StatusOK) // (actually no)
		return
	}

	user := users[0]

	// Create the token
	token := jwt.New()
	_ = token.Set("id", user.Id)
	_ = token.Set(jwt.ExpirationKey, time.Now().Add(expirationResetToken))

	// Sign the token and generate a payload
	signed, err := jwt.Sign(token, jwa.RS256, privKey)
	if err != nil {
		utils.PrintError(err)
		utils.Prettier(w, fmt.Sprintf("échec de la génération de la payload chiffrée: %s\n", err), nil, http.StatusInternalServerError)
		return
	}
	tokenStr := string(signed)

	//send mail to the user
	err = mail.SendTemplate(user.Username, email, "Modification du mot de passe", "reset_password", schemas.ResetPassword{
		Username:    user.Username,
		ResetLink:   template.URL(config.OriginServerFront + "/mot-de-passe-oublie/" + tokenStr),
		PitStopLink: config.OriginServerFront,
		SendingDate: time.Now().Format("02/01/2006 15:04:05"),
	})
	if err != nil {
		return
	}
	//

	utils.Prettier(w, "Si le compte a été trouvé, un email a été envoyé", nil, http.StatusOK)
}

func Create(w http.ResponseWriter, r *http.Request) {
	creds := map[string]string{}
	_ = r.ParseForm()
	creds["username"], creds["email"], creds["password"] = r.FormValue("username"), r.FormValue("email"), r.FormValue("password")
	for name, cred := range creds {
		if len(cred) == 0 {
			utils.Prettier(w, "no "+name+" provided !", nil, http.StatusBadRequest)
			return
		}
	}
	password, err := HashPassword(creds["password"])
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	err = database.CreateUser(creds["username"], creds["email"], password)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusInternalServerError)
		return
	}
	utils.Prettier(w, "utilisateur créé !", nil, http.StatusOK)
}

func Profile(w http.ResponseWriter, r *http.Request) {
	id, _ := ExtractIdFromRequest(r)
	utils.Prettier(w, "Good !", id, http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		utils.Prettier(w, "token invalide", nil, http.StatusBadRequest)
		return
	}

	accessToken, err := ExtractToken(r)
	if err != nil {
		utils.Prettier(w, "token invalide (bearer manquant)", nil, http.StatusBadRequest)
	}

	refreshToken := r.FormValue("refresh_token")

	blackListAccessToken = append(blackListAccessToken, accessToken)
	blackListRefreshToken = append(blackListRefreshToken, refreshToken)

	utils.Prettier(w, "logout success", nil, http.StatusOK)
}

func Refresh(w http.ResponseWriter, r *http.Request) {

	var data Data

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	refreshToken := data.RefreshToken

	token, err := jwt.ParseString(refreshToken, jwt.WithKeySet(keySet), jwt.UseDefaultKey(true))
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}
	id, ok := token.Get("id")
	if !ok {
		utils.Prettier(w, "token invalide", nil, http.StatusBadRequest)
		return
	}

	if !time.Now().UTC().Before(token.Expiration()) {
		utils.Prettier(w, "token expiré", nil, http.StatusBadRequest)
		return
	}

	floatId := id.(float64)
	if floatId < 1 {
		utils.Prettier(w, "token invalide", nil, http.StatusBadRequest)
		return
	}

	user, exist := database.GetUserById(uint(floatId))

	if !exist {
		utils.Prettier(w, "token invalide", nil, http.StatusBadRequest)
		return
	}

	// Create the accessToken
	accessToken := jwt.New()
	_ = accessToken.Set("id", user.Id)
	_ = accessToken.Set("username", user.Username)
	_ = accessToken.Set("email", user.Email)
	_ = accessToken.Set(jwt.ExpirationKey, time.Now().Add(expirationAccessToken))

	// Sign the token and generate a payload
	signed, err := jwt.Sign(accessToken, jwa.RS256, privKey)
	if err != nil {
		utils.PrintError(err)
		utils.Prettier(w, fmt.Sprintf("échec de la génération de la payload chiffrée: %s\n", err), nil, http.StatusInternalServerError)
		return
	}

	accessTokenStr := string(signed)
	utils.Prettier(w, "Token generated !", struct {
		AccessToken string `json:"access_token"`
	}{accessTokenStr}, http.StatusOK)

}
