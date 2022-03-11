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

func GetAllCategories() ([]uint, error) {
	db, err := connect()
	if err != nil {
		return []uint{}, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT id FROM categories")
	if err != nil {
		utils.PrintError(err)
		return []uint{}, err
	}
	defer rows.Close()

	categories := []uint{}
	for rows.Next() {
		var id uint
		err = rows.Scan(&id)
		if err != nil {
			utils.PrintError(err)
			return []uint{}, err
		}
		categories = append(categories, id)
	}

	return categories, nil
}

func GetCategory(id uint) (schemas.DbCategory, bool) {
	db, err := connect()
	if err != nil {
		return schemas.DbCategory{}, false
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, name, slug, creation_date, topics, icon FROM categories WHERE id = ? LIMIT 1")
	if err != nil {
		utils.PrintError(err)
		return schemas.DbCategory{}, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		utils.PrintError(err)
		return schemas.DbCategory{}, false
	}
	defer rows.Close()

	if !rows.Next() { // No category found
		return schemas.DbCategory{}, false
	}

	category := new(schemas.DbCategory)
	var jsonTopics string
	err = rows.Scan(&category.Id, &category.Name, &category.Slug, &category.CreationDate, &jsonTopics, &category.Icon)
	if err != nil {
		return schemas.DbCategory{}, false
	}

	err = json.Unmarshal([]byte(jsonTopics), &category.Topics)
	if err != nil {
		return schemas.DbCategory{}, false
	}

	return *category, true
}

func CreateCategory(name, icon string) (uint, error) {
	creationDate := time.Now().Unix()
	slug := utils.Slugify(name)

	db, err := connect()
	if err != nil {
		return 0, err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare("INSERT INTO categories(name, slug, creation_date, icon) VALUES(?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	res, err := stmt.Exec(name, slug, creationDate, icon)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return uint(id), tx.Commit()
}

func CheckCategoryExistence(name string) bool {
	db, err := connect()
	if err != nil {
		utils.PrintError(errors.New("échec de la vérification de l'éxistence de la catégorie: " + err.Error()))
		return true
	}
	defer db.Close()

	var shouldBeZero int
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE slug = ?)", utils.Slugify(name)).Scan(&shouldBeZero)
	if err != nil {
		if err != sql.ErrNoRows {
			utils.PrintError(errors.New("échec de la vérification de l'éxistence de la catégorie: " + err.Error()))
			return true
		}
	}
	return shouldBeZero == 1
}

func UpdateCategoryById(id uint, field string, value interface{}) error {
	db, err := connect()
	if err != nil {
		utils.PrintError(errors.New("échec de la mise à jour de la catégorie: " + err.Error()))
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare(fmt.Sprintf("UPDATE categories SET %s = ? WHERE id = ?", field))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(value, id)
	return err
}

func DeleteCategory(id uint) error {
	db, err := connect()
	if err != nil {
		utils.PrintError(errors.New("échec de la suppression de la catégorie: " + err.Error()))
		return err
	}
	defer db.Close()

	_, err = db.Exec("DELETE FROM categories WHERE id = ?", id)
	return err
}

func GetLastTopic(id uint) uint {
	db, err := connect()
	if err != nil {
		utils.PrintError(errors.New("échec de la récupération du dernier topic: " + err.Error()))
		return 0
	}
	defer db.Close()

	rows, err := db.Query("SELECT id FROM topics WHERE parent_id = ? ORDER BY creation_date DESC LIMIT 1", id)
	if err != nil {
		utils.PrintError(errors.New("échec de la récupération du dernier topic: " + err.Error()))
		return 0
	}

	if !rows.Next() {
		return 0
	}

	var lastTopic uint
	err = rows.Scan(&lastTopic)
	if err != nil {
		utils.PrintError(errors.New("échec de la récupération du dernier topic: " + err.Error()))
		return 0
	}

	return lastTopic
}
