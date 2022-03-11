package database

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"pitstop-api/src/schemas"
	"pitstop-api/src/utils"
	"time"
)

func GetTopicById(id uint) (schemas.DbTopic, bool) {
	db, err := connect()
	if err != nil {
		return schemas.DbTopic{}, false
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, parent_id, name, slug, author_id, creation_date, posts, pinned, locked, tags FROM topics WHERE id = ? LIMIT 1")
	if err != nil {
		utils.PrintError(err)
		return schemas.DbTopic{}, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		utils.PrintError(err)
		return schemas.DbTopic{}, false
	}
	defer rows.Close()

	if !rows.Next() { // No topic found
		return schemas.DbTopic{}, false
	}

	topic := new(schemas.DbTopic)
	var jsonPosts, jsonTags string
	err = rows.Scan(&topic.Id, &topic.ParentId, &topic.Name, &topic.Slug, &topic.AuthorId, &topic.CreationDate, &jsonPosts, &topic.PinnedBy, &topic.LockDate, &jsonTags)
	if err != nil {
		utils.PrintError(err)
		return schemas.DbTopic{}, false
	}

	topic.Pinned = topic.PinnedBy > 0
	topic.Locked = topic.LockDate > 0

	err = json.Unmarshal([]byte(jsonPosts), &topic.Posts)
	if err != nil {
		return schemas.DbTopic{}, false
	}
	err = json.Unmarshal([]byte(jsonTags), &topic.Tags)
	if err != nil {
		return schemas.DbTopic{}, false
	}

	return *topic, true
}

func CheckTopicExistence(parentId uint, name string) bool {
	db, err := connect()
	if err != nil {
		utils.PrintError(errors.New("échec de la vérification de l'éxistence du topic: " + err.Error()))
		return true
	}
	defer db.Close()

	var shouldBeZero int
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM topics WHERE parent_id = ? AND slug = ?)", parentId, utils.Slugify(name)).Scan(&shouldBeZero)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.PrintError(errors.New("échec de la vérification de l'éxistence du topic: " + err.Error()))
			return true
		}
	}
	return shouldBeZero == 1
}

func CreateTopic(parentId uint, name string, authorId uint, tags []string) (uint, error) {
	slug := utils.Slugify(name)
	creationDate := time.Now().Unix()
	jsonTags, err := json.Marshal(tags)
	if err != nil {
		return 0, err
	}

	db, err := connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare("INSERT INTO topics(parent_id, name, slug, author_id, creation_date, tags) VALUES(?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(parentId, name, slug, authorId, creationDate, jsonTags)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(id), tx.Commit()
}

func UpdateTopicById(id uint, field string, value interface{}) error {
	db, err := connect()
	if err != nil {
		utils.PrintError(errors.New("échec de la mise à jour du topic: " + err.Error()))
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(fmt.Sprintf("UPDATE topics SET %s = ? WHERE id = ?", field))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(value, id)
	return err
}

func DeleteTopic(id uint, deleteFromParent bool) error {
	db, err := connect()
	if err != nil {
		utils.PrintError(err)
		utils.PrintError(errors.New("échec de la suppression du topic: " + err.Error()))
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM topics WHERE id = ?", id)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if deleteFromParent {
		db, err := connect()
		if err != nil {
			utils.PrintError(err)
			utils.PrintError(errors.New("échec de la suppression du topic: " + err.Error()))
			return err
		}
		defer db.Close()

		rows, err := db.Query("SELECT id, topics FROM categories")
		if err != nil {
			utils.PrintError(err)
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var categoryId uint
			var jsonTopics string
			err = rows.Scan(&categoryId, &jsonTopics)
			if err != nil {
				return err
			}
			var topics []uint
			err = json.Unmarshal([]byte(jsonTopics), &topics)
			if err != nil {
				return err
			}
			for i, topic := range topics {
				if topic == id {
					topics = append(topics[:i], topics[i+1:]...)
					newJsonTopics, err := json.Marshal(topics)
					go func(newValue string) {
						time.Sleep(100 * time.Millisecond)
						err := UpdateCategoryById(categoryId, "topics", newValue)
						if err != nil {
							utils.PrintError(err)
						}
					}(string(newJsonTopics))
					return err
				}
			}
		}
		return errors.New("topic supprimé mais pas retiré de sa catégorie car elle n'a pas été trouvée")
	}

	return nil
}
