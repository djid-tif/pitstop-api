package categories

import (
	"github.com/gorilla/mux"
	"net/http"
	"pitstop-api/src/database"
	"pitstop-api/src/users"
	"pitstop-api/src/utils"
	"strconv"
	"strings"
)

func ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := database.GetAllCategories()
	if err != nil {
		utils.Prettier(w, "error: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}
	utils.Prettier(w, "list of all categories", categories, http.StatusOK)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	strId, found := vars["id"]
	if !found {
		utils.Prettier(w, "aucun id fourni !", nil, http.StatusBadRequest)
		return
	}
	id, err := utils.StringToUint(strId)
	if err != nil {
		utils.Prettier(w, "id invalide !", nil, http.StatusBadRequest)
		return
	}
	category, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "catégorie non trouvée !", nil, http.StatusBadRequest)
		return
	}

	utils.Prettier(w, "category "+strId, category.ToSchema(), http.StatusOK)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	user, exist := users.LoadUserFromRequest(r)
	if !exist {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}

	if !user.GetRank().CanCreateCategory() {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return

	}

	_ = r.ParseForm()
	name := r.FormValue("name")
	if !utils.IsTitleValid(name) {
		utils.Prettier(w, "titre invalide !", nil, http.StatusBadRequest)
		return
	}

	alreadyExist := database.CheckCategoryExistence(name)
	if alreadyExist {
		utils.Prettier(w, "une catégorie avec ce nom existe déjà !", nil, http.StatusBadRequest)
		return
	}

	icon := r.FormValue("icon")
	if len(icon) == 0 {
		utils.Prettier(w, "aucun icon fourni !", nil, http.StatusBadRequest)
		return
	}

	categoryId, err := database.CreateCategory(name, icon)
	if err != nil {
		utils.Prettier(w, "échec de la création de la catégorie : "+err.Error(), nil, http.StatusBadRequest)
		return
	}

	utils.Prettier(w, "catégorie créée !", struct {
		CategoryId uint `json:"category_id"`
	}{CategoryId: categoryId}, http.StatusOK)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	user, exist := users.LoadUserFromRequest(r)
	if !exist {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}

	if !user.GetRank().CanModifyCategory() {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	strId, found := vars["id"]
	if !found {
		utils.Prettier(w, "aucun id fourni !", nil, http.StatusBadRequest)
		return
	}
	id, err := utils.StringToUint(strId)
	if err != nil {
		utils.Prettier(w, "id invalide !", nil, http.StatusBadRequest)
		return
	}
	category, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "catégorie non trouvée !", nil, http.StatusBadRequest)
		return
	}

	_ = r.ParseForm()
	names, _ := r.Form["name"]
	icons, _ := r.Form["icon"]

	if len(names) == 1 {
		if !utils.IsTitleValid(names[0]) {
			utils.Prettier(w, "titre invalide !", nil, http.StatusBadRequest)
			return
		}
		if database.CheckCategoryExistence(names[0]) {
			utils.Prettier(w, "une catégorie avec ce nom existe déjà !", nil, http.StatusBadRequest)
			return
		}
		err = database.UpdateCategoryById(category.GetId(), "name", names[0])
		if err != nil {
			utils.Prettier(w, "échec de la mise à jour du titre: "+err.Error(), nil, http.StatusInternalServerError)
			return
		}
	}

	if len(icons) == 1 {
		err = database.UpdateCategoryById(category.GetId(), "icon", icons[0])
		if err != nil {
			utils.Prettier(w, "échec de la mise à jour de l'icon: "+err.Error(), nil, http.StatusInternalServerError)
			return
		}
	}

	utils.Prettier(w, "catégorie mise à jour avec succès !", nil, http.StatusOK)

}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	user, exist := users.LoadUserFromRequest(r)
	if !exist {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}

	if !user.GetRank().CanDeleteCategory() {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	strId, found := vars["id"]
	if !found {
		utils.Prettier(w, "aucun id fourni !", nil, http.StatusBadRequest)
		return
	}
	id, err := utils.StringToUint(strId)
	if err != nil {
		utils.Prettier(w, "id invalide !", nil, http.StatusBadRequest)
		return
	}
	category, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "catégorie non trouvée !", nil, http.StatusBadRequest)
		return
	}

	err = category.delete()
	if err != nil {
		utils.Prettier(w, "échec de la suppression de la catégorie !", nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "catégorie supprimée avec succès !", nil, http.StatusOK)
}

func GetTopicsOfCategory(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractUintFromRequest("id", r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	category, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "catégorie non trouvée !", nil, http.StatusBadRequest)
		return
	}

	utils.Prettier(w, "topics de la catégorie "+strconv.FormatUint(uint64(id), 10), category.GetTopics(), http.StatusOK)
}

func CreateTopicInCategory(w http.ResponseWriter, r *http.Request) {
	author, exist := users.LoadUserFromRequest(r)
	if !exist {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}
	if !author.GetRank().CanCreateTopic() {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	_ = r.ParseForm()
	name := r.FormValue("name")
	if len(name) == 0 {
		utils.Prettier(w, "aucun titre fourni !", nil, http.StatusBadRequest)
		return
	}
	if !utils.IsTitleValid(name) {
		utils.Prettier(w, "titre du topic invalide !", nil, http.StatusBadRequest)
		return
	}
	tags := strings.Split(r.FormValue("tags"), ",")
	if len(tags) == 1 && tags[0] == "" {
		tags = []string{}
	}

	id, err := utils.ExtractUintFromRequest("id", r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	category, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "catégorie non trouvée !", nil, http.StatusBadRequest)
		return
	}

	alreadyExists := database.CheckTopicExistence(category.id, name)
	if alreadyExists {
		utils.Prettier(w, "un topic avec ce nom existe déjà in this category !", nil, http.StatusBadRequest)
		return
	}

	topicId, err := database.CreateTopic(category.id, name, author.GetId(), tags)
	if err != nil {
		utils.Prettier(w, "échec de la création du topic: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	err = category.appendChild(topicId)
	if err != nil {
		utils.Prettier(w, "échec de l'ajout du topic à la catégorie: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	author.Stats().CreateTopic(author.GetId())

	utils.Prettier(w, "topic créé !", struct {
		TopicId uint `json:"topic_id"`
	}{TopicId: topicId}, http.StatusOK)
}
