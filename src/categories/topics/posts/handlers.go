package posts

import (
	"encoding/json"
	"net/http"
	"pitstop-api/src/database"
	"pitstop-api/src/users"
	"pitstop-api/src/utils"
	"strconv"
)

func GetPost(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractUintFromRequest("id", r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	post, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "post non trouvé !", nil, http.StatusBadRequest)
		return
	}

	utils.Prettier(w, "post "+strconv.FormatUint(uint64(id), 10), post.ToSchema(), http.StatusOK)
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	user, found := users.LoadUserFromRequest(r)
	if !found {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}

	id, err := utils.ExtractUintFromRequest("id", r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	post, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "post non trouvé !", nil, http.StatusBadRequest)
		return
	}

	if post.GetAuthorId() == user.GetId() {
		if !user.GetRank().CanModifyOwnPost() {
			utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
			return
		}
	} else {
		if !user.GetRank().CanModifyUsersPosts() {
			utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
			return
		}
	}

	_ = r.ParseForm()
	content := r.FormValue("content")
	if len(content) == 0 {
		utils.Prettier(w, "aucun contenu fourni !", nil, http.StatusBadRequest)
		return
	}

	err = database.UpdatePostById(post.id, "content", content)
	if err != nil {
		utils.Prettier(w, "échec de la mise à jour du post: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	user.Stats().ModifyPost(user.GetId())

	utils.Prettier(w, "post mis à jour avec succès !", nil, http.StatusOK)
}

func LikePost(w http.ResponseWriter, r *http.Request) {
	user, found := users.LoadUserFromRequest(r)
	if !found {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}

	if !user.GetRank().CanReact() {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	id, err := utils.ExtractUintFromRequest("id", r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	post, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "post non trouvé !", nil, http.StatusBadRequest)
		return
	}

	if user.GetId() == post.GetId() {
		utils.Prettier(w, "vous ne pouvez pas aimer votre propre post !", nil, http.StatusBadRequest)
		return
	}

	dislikeRemoved := false
	finalLikes, likeRemoved := utils.RemoveUintFromList(user.GetId(), post.GetLikes())
	if !likeRemoved { // already liked
		finalLikes = append(finalLikes, user.GetId())
		var finalDislikes []uint
		finalDislikes, dislikeRemoved = utils.RemoveUintFromList(user.GetId(), post.GetDislikes())
		if dislikeRemoved {
			jsonDislikes, err := json.Marshal(finalDislikes)
			if err != nil {
				utils.Prettier(w, "échec de la réaction au post: "+err.Error(), nil, http.StatusInternalServerError)
				return
			}

			err = database.UpdatePostById(post.GetId(), "dislikes", string(jsonDislikes))
			if err != nil {
				utils.Prettier(w, "échec de la réaction au post: "+err.Error(), nil, http.StatusInternalServerError)
				return
			}
		}
	}

	jsonLikes, err := json.Marshal(finalLikes)
	if err != nil {
		utils.Prettier(w, "échec de la réaction au post: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	err = database.UpdatePostById(post.GetId(), "likes", string(jsonLikes))
	if err != nil {
		utils.Prettier(w, "échec de la réaction au post: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	if likeRemoved {
		utils.Prettier(w, "post plus aimé avec succès !", nil, http.StatusOK)
	} else {
		if !dislikeRemoved && !utils.ContainsUint(user.GetId(), post.dislikes) {
			user.Stats().AddReaction(user.GetId())
		}
		utils.Prettier(w, "post aimé avec succès !", nil, http.StatusOK)
	}
}

func DislikePost(w http.ResponseWriter, r *http.Request) {
	user, found := users.LoadUserFromRequest(r)
	if !found {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}

	if !user.GetRank().CanReact() {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	id, err := utils.ExtractUintFromRequest("id", r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	post, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "post non trouvé !", nil, http.StatusBadRequest)
		return
	}

	if user.GetId() == post.GetId() {
		utils.Prettier(w, "vous ne pouvez pas ne pas aimer votre propre post !", nil, http.StatusBadRequest)
		return
	}

	var likeRemoved bool
	finalDislikes, dislikeRemoved := utils.RemoveUintFromList(user.GetId(), post.GetDislikes())
	if !dislikeRemoved { // already disliked
		finalDislikes = append(finalDislikes, user.GetId())
		var finalLikes []uint
		finalLikes, likeRemoved = utils.RemoveUintFromList(user.GetId(), post.GetLikes())
		if likeRemoved {
			jsonLikes, err := json.Marshal(finalLikes)
			if err != nil {
				utils.Prettier(w, "échec de la réaction au post: "+err.Error(), nil, http.StatusInternalServerError)
				return
			}

			err = database.UpdatePostById(post.GetId(), "likes", string(jsonLikes))
			if err != nil {
				utils.Prettier(w, "échec de la réaction au post: "+err.Error(), nil, http.StatusInternalServerError)
				return
			}
		}
	}

	jsonDislikes, err := json.Marshal(finalDislikes)
	if err != nil {
		utils.Prettier(w, "échec de la réaction au post: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	err = database.UpdatePostById(post.GetId(), "dislikes", string(jsonDislikes))
	if err != nil {
		utils.Prettier(w, "échec de la réaction au post: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	if dislikeRemoved {
		utils.Prettier(w, "post plus pas aimé avec succès !", nil, http.StatusOK)
	} else {
		if !likeRemoved && !utils.ContainsUint(user.GetId(), post.dislikes) {
			user.Stats().AddReaction(user.GetId())
		}
		utils.Prettier(w, "post pas aimé avec succès !", nil, http.StatusOK)
	}
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	user, found := users.LoadUserFromRequest(r)
	if !found {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}

	id, err := utils.ExtractUintFromRequest("id", r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	post, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "post non trouvé !", nil, http.StatusBadRequest)
		return
	}

	if post.GetAuthorId() == user.GetId() {
		if !user.GetRank().CanDeleteOwnPost() {
			utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
			return
		}
	} else {
		if !user.GetRank().CanDeleteUsersPosts() {
			utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
			return
		}
	}

	err = post.Delete(true)
	if err != nil {
		utils.Prettier(w, "échec de la suppression du post: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "post supprimé avec succès !", nil, http.StatusOK)
}

func GetMostLikedPosts(w http.ResponseWriter, r *http.Request) {
	var postLimit uint = 6

	posts, err := database.GetMostLikedPosts(postLimit)
	if err != nil {
		utils.Prettier(w, "échec de la récupération des posts les plus aimés: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "posts les plus aimés", posts, http.StatusOK)
}
