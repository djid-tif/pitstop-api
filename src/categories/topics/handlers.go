package topics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pitstop-api/src/config"
	"pitstop-api/src/database"
	"pitstop-api/src/mail"
	"pitstop-api/src/schemas"
	"pitstop-api/src/users"
	"pitstop-api/src/utils"
	"strconv"
	"strings"
	"time"
)

func GetTopic(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractUintFromRequest("id", r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	topic, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "topic non trouvé !", nil, http.StatusBadRequest)
		return
	}

	utils.Prettier(w, "topic "+strconv.FormatUint(uint64(id), 10), topic.ToSchema(), http.StatusOK)
}

func UpdateTopic(w http.ResponseWriter, r *http.Request) {
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

	topic, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "topic non trouvé !", nil, http.StatusBadRequest)
		return
	}

	if topic.GetAuthorId() == user.GetId() {
		if !user.GetRank().CanModifyOwnTopic() {
			utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
			return
		}
	} else {
		if !user.GetRank().CanModifyUsersTopics() {
			utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
			return
		}
	}

	_ = r.ParseForm()
	if names, exist := r.Form["name"]; exist {
		name := names[0]
		if name == topic.name {
			utils.Prettier(w, "le nouveau nom est le même que l'ancien !", nil, http.StatusBadRequest)
			return
		}
		if !utils.IsTitleValid(name) {
			utils.Prettier(w, "titre du topic invalide !", nil, http.StatusBadRequest)
			return
		}
		if database.CheckTopicExistence(topic.GetParentId(), name) {
			utils.Prettier(w, "un topic avec ce nom existe déjà in this category !", nil, http.StatusBadRequest)
			return
		}
		err = database.UpdateTopicById(topic.id, "name", name)
		if err != nil {
			utils.Prettier(w, "échec de la mise à jour du topic: "+err.Error(), nil, http.StatusInternalServerError)
			return
		}
	}
	if strTags, exist := r.Form["tags"]; exist {
		tags := strings.Split(strTags[0], ",")
		if len(tags) == 1 && tags[0] == "" {
			tags = []string{}
		}
		jsonTags, err := json.Marshal(tags)
		if err != nil {
			utils.Prettier(w, "échec du parsing des tags: "+err.Error(), nil, http.StatusInternalServerError)
			return
		}
		err = database.UpdateTopicById(topic.id, "tags", string(jsonTags))
		if err != nil {
			utils.Prettier(w, "échec de la mise à jour du topic: "+err.Error(), nil, http.StatusInternalServerError)
			return
		}
	}

	user.Stats().ModifyTopic(user.GetId())

	utils.Prettier(w, "topic mis à jour avec succès !", nil, http.StatusOK)
}

func PinTopic(w http.ResponseWriter, r *http.Request) {
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

	if !user.GetRank().CanModifyCategory() { // can modify category = can also pin a topic
		utils.Prettier(w, "vous n'avez pas la permmlllission !", nil, http.StatusUnauthorized)
		return
	}

	topic, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "topic non trouvé !", nil, http.StatusBadRequest)
		return
	}

	if topic.IsPinned() {
		utils.Prettier(w, "ce topic est déjà épinglé !", nil, http.StatusBadRequest)
		return
	}

	err = database.UpdateTopicById(topic.id, "pinned", user.GetId())
	if err != nil {
		utils.Prettier(w, "échec de l'épinglage du topic: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "topic épinglé avec succès !", nil, http.StatusOK)
}

func UnpinTopic(w http.ResponseWriter, r *http.Request) {
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

	if !user.GetRank().CanModifyCategory() { // can modify category = can also pin a topic
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	topic, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "topic non trouvé !", nil, http.StatusBadRequest)
		return
	}

	if !topic.IsPinned() {
		utils.Prettier(w, "ce topic n'est pas épinglé !", nil, http.StatusBadRequest)
		return
	}

	err = database.UpdateTopicById(topic.id, "pinned", 0)
	if err != nil {
		utils.Prettier(w, "échec du désépinglage du topic: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "topic désépinglé avec succès !", nil, http.StatusOK)
}

func LockTopic(w http.ResponseWriter, r *http.Request) {
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

	topic, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "topic non trouvé !", nil, http.StatusBadRequest)
		return
	}

	if topic.IsLocked() {
		utils.Prettier(w, "ce topic est déjà verrouillé !", nil, http.StatusBadRequest)
		return
	}

	if topic.GetAuthorId() == user.GetId() {
		if !user.GetRank().CanLockOwnTopic() {
			utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
			return
		}
	} else {
		if !user.GetRank().CanLockUsersTopic() {
			utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
			return
		}
	}

	err = database.UpdateTopicById(topic.id, "locked", time.Now().Unix())
	if err != nil {
		utils.Prettier(w, "échec du verrouillage du topic: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	user.Stats().LockTopic(user.GetId())

	utils.Prettier(w, "topic verrouillé avec succès !", nil, http.StatusOK)
}

func UnlockTopic(w http.ResponseWriter, r *http.Request) {
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

	topic, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "topic non trouvé !", nil, http.StatusBadRequest)
		return
	}

	if !topic.IsLocked() {
		utils.Prettier(w, "ce topic n'est pas verrouillé !", nil, http.StatusBadRequest)
		return
	}

	if topic.GetAuthorId() == user.GetId() {
		if !user.GetRank().CanLockOwnTopic() {
			utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
			return
		}
	} else {
		if !user.GetRank().CanLockUsersTopic() {
			utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
			return
		}
	}

	err = database.UpdateTopicById(topic.id, "locked", 0)
	if err != nil {
		utils.Prettier(w, "échec du déverrouillage du topic: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "topic déverrouillé avec succès !", nil, http.StatusOK)
}

func DeleteTopic(w http.ResponseWriter, r *http.Request) {
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

	topic, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "topic non trouvé !", nil, http.StatusBadRequest)
		return
	}

	if topic.GetAuthorId() == user.GetId() {
		if !user.GetRank().CanDeleteOwnTopic() {
			utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
			return
		}
	} else {
		if !user.GetRank().CanDeleteUsersTopics() {
			utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
			return
		}
	}

	err = topic.Delete(true)
	if err != nil {
		utils.Prettier(w, "échec de la suppression du topic: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	utils.Prettier(w, "topic supprimé avec succès !", nil, http.StatusOK)
}

func GetPostsOfTopic(w http.ResponseWriter, r *http.Request) {
	id, err := utils.ExtractUintFromRequest("id", r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	topic, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "topic non trouvé !", nil, http.StatusBadRequest)
		return
	}

	utils.Prettier(w, "posts of topic "+strconv.FormatUint(uint64(id), 10), topic.GetPosts(), http.StatusOK)
}

func CreatePostInTopic(w http.ResponseWriter, r *http.Request) {
	author, exist := users.LoadUserFromRequest(r)
	if !exist {
		utils.Prettier(w, "id invalide !", nil, http.StatusUnauthorized)
		return
	}
	if !author.GetRank().CanCreatePost() {
		utils.Prettier(w, "vous n'avez pas la permission !", nil, http.StatusUnauthorized)
		return
	}

	_ = r.ParseForm()
	content := r.FormValue("content")
	if len(content) == 0 {
		utils.Prettier(w, "aucun contenu fourni !", nil, http.StatusBadRequest)
		return
	}

	id, err := utils.ExtractUintFromRequest("id", r)
	if err != nil {
		utils.Prettier(w, err.Error(), nil, http.StatusBadRequest)
		return
	}

	topic, found := LoadFromId(id)
	if !found {
		utils.Prettier(w, "topic non trouvé !", nil, http.StatusBadRequest)
		return
	}

	if topic.IsLocked() {
		utils.Prettier(w, "ce topic est verrouillé !", nil, http.StatusBadRequest)
		return
	}

	postId, err := database.CreatePost(topic.id, content, author.GetId())
	if err != nil {
		utils.Prettier(w, "échec de la création du post: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	err = topic.appendChild(postId)
	if err != nil {
		utils.Prettier(w, "échec de l'ajout du post au topic: "+err.Error(), nil, http.StatusInternalServerError)
		return
	}

	author.Stats().CreateTopic(author.GetId())

	go func() {
		if author.GetId() != topic.GetAuthorId() {
			topicAuthor, found := users.LoadUser(topic.GetId())
			if found && topicAuthor.Settings().AreNotificationsActive() {
				subject := fmt.Sprintf("%s a répondu à votre topic !", author.GetUsername())
				err = mail.SendTemplate(topicAuthor.GetUsername(), topicAuthor.GetEmail(), subject, "notification", schemas.ReplyToTopic{
					Email: schemas.Email{
						Username:    topicAuthor.GetUsername(),
						PitStopLink: config.OriginServerFront,
						SendingDate: time.Now().Format("02/01/2006 15:04:05"),
					},
					TopicName:   topic.GetName(),
					ReplierName: author.GetUsername(),
					TopicLink:   config.OriginServerFront + "/forum/" + strconv.FormatUint(uint64(topic.parentId), 10) + "/topic/" + strconv.FormatUint(uint64(topic.id), 10),
				})
				if err != nil {
					utils.PrintError(err)
				}
			}
		}
	}()

	utils.Prettier(w, "post créé !", struct {
		PostId uint `json:"post_id"`
	}{PostId: postId}, http.StatusOK)
}
