package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"pitstop-api/src/schemas"
	"pitstop-api/src/utils"
	"time"
)

func GetUserPosts(userId uint) []uint {
	db, err := connect()
	if err != nil {
		return []uint{}
	}
	defer db.Close()

	rows, err := db.Query("SELECT id FROM posts WHERE author_id = ?", userId)
	if err != nil {
		utils.PrintError(err)
		return []uint{}
	}

	posts := []uint{}
	for rows.Next() {
		var postId uint
		err = rows.Scan(&postId)
		if err != nil {
			continue
		}
		posts = append(posts, postId)
	}

	return posts
}

func GetPostById(id uint) (schemas.DbPost, bool) {
	db, err := connect()
	if err != nil {
		return schemas.DbPost{}, false
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, parent_id, content, author_id, creation_date, likes, dislikes FROM posts WHERE id = ? LIMIT 1")
	if err != nil {
		utils.PrintError(err)
		return schemas.DbPost{}, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		utils.PrintError(err)
		return schemas.DbPost{}, false
	}
	defer rows.Close()

	if !rows.Next() { // No post found
		return schemas.DbPost{}, false
	}

	post := new(schemas.DbPost)
	var jsonLikes, jsonDislikes string
	err = rows.Scan(&post.Id, &post.ParentId, &post.Content, &post.AuthorId, &post.CreationDate, &jsonLikes, &jsonDislikes)
	if err != nil {
		utils.PrintError(err)
		return schemas.DbPost{}, false
	}

	err = json.Unmarshal([]byte(jsonLikes), &post.Likes)
	if err != nil {
		return schemas.DbPost{}, false
	}
	err = json.Unmarshal([]byte(jsonDislikes), &post.Dislikes)
	if err != nil {
		return schemas.DbPost{}, false
	}

	return *post, true
}

func CreatePost(parentId uint, content string, authorId uint) (uint, error) {
	creationDate := time.Now().Unix()

	db, err := connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare("INSERT INTO posts(parent_id, content, author_id, creation_date) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(parentId, content, authorId, creationDate)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(id), tx.Commit()
}

func UpdatePostById(id uint, field string, value interface{}) error {
	db, err := connect()
	if err != nil {
		utils.PrintError(errors.New("échec de la mise à jour du post: " + err.Error()))
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(fmt.Sprintf("UPDATE posts SET %s = ? WHERE id = ?", field))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(value, id)
	return err
}

func DeletePost(id uint, deleteFromParent bool) error {
	db, err := connect()
	if err != nil {
		utils.PrintError(err)
		utils.PrintError(errors.New("échec de la suppression du post: " + err.Error()))
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM posts WHERE id = ?", id)
	if err != nil {
		utils.PrintError(err)
		return err
	}

	if deleteFromParent {
		db, err := connect()
		if err != nil {
			utils.PrintError(err)
			utils.PrintError(errors.New("échec de la suppression du post: " + err.Error()))
			return err
		}
		defer db.Close()

		rows, err := db.Query("SELECT id, posts FROM topics")
		if err != nil {
			utils.PrintError(err)
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var topicId uint
			var jsonPosts string
			err = rows.Scan(&topicId, &jsonPosts)
			if err != nil {
				return err
			}
			var posts []uint
			err = json.Unmarshal([]byte(jsonPosts), &posts)
			if err != nil {
				return err
			}
			for i, topic := range posts {
				if topic == id {
					posts = append(posts[:i], posts[i+1:]...)
					newJsonPosts, err := json.Marshal(posts)
					go func(newValue string) {
						time.Sleep(100 * time.Millisecond)
						err := UpdateTopicById(topicId, "posts", newValue)
						if err != nil {
							utils.PrintError(err)
						}
					}(string(newJsonPosts))
					return err
				}
			}
		}
		return errors.New("post supprimé mais pas retiré de son topic car il n'a pas été trouvé")
	}

	return nil
}

func GetMostLikedPosts(limit uint) ([]uint, error) {
	db, err := connect()
	if err != nil {
		return []uint{}, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id FROM posts ORDER BY likes LIMIT ?", limit)
	if err != nil {
		return []uint{}, err
	}

	posts := []uint{}
	for rows.Next() {
		var postId uint
		err = rows.Scan(&postId)
		if err != nil {
			return []uint{}, err
		}
		posts = append(posts, postId)
	}

	return posts, nil
}
