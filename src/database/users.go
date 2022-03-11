package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"pitstop-api/src/schemas"
	"pitstop-api/src/utils"
	"time"
)

// CreateUser creates an entry in the database with the given user information and returns any error
func CreateUser(username, email, password string) error { // return at least user id ?
	regDate := time.Now().Unix()
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO users(username, email, password, registration_date, last_password_refresh, settings, stats) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, email, password, regDate, regDate, defaultUserSettings, defaultUserStats)
	if err != nil {
		return err
	}
	return tx.Commit()
}

// CheckUserExistence checks if a user with the given username or password is already registered or not,
// and returns the taken field, the result of the check and an error if their is one
func CheckUserExistence(username, email string) (string, bool, error) {
	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var shouldBeZero int
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", username).Scan(&shouldBeZero)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", false, err
		}
	}
	if shouldBeZero == 1 { // username is already registered
		return "username", true, nil
	}

	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&shouldBeZero)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", false, err
		}
	}
	if shouldBeZero == 1 { // email is already registered
		return "email", true, nil
	}

	return "", false, nil
}

// GetUserById return the user corresponding to the given id in the database and if it was found or not
func GetUserById(id uint) (schemas.DbUser, bool) {
	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT id, username, email, password, registration_date, rank, last_password_refresh, friends, settings, stats FROM users WHERE id = ? LIMIT 1")
	if err != nil {
		utils.PrintError(err)
		return schemas.DbUser{}, false
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		utils.PrintError(err)
		return schemas.DbUser{}, false
	}
	defer rows.Close()

	if !rows.Next() { // No user found
		return schemas.DbUser{}, false
	}

	dbUser := new(schemas.DbUser)
	var jsonFriends, jsonSettings, jsonStats []byte
	err = rows.Scan(&dbUser.Id, &dbUser.Username, &dbUser.Email, &dbUser.Password, &dbUser.RegistrationDate, &dbUser.Rank, &dbUser.LastPasswordRefresh, &jsonFriends, &jsonSettings, &jsonStats)
	if err != nil {
		utils.PrintError(err)
		return schemas.DbUser{}, false
	}
	err = json.Unmarshal(jsonFriends, &dbUser.Friends)
	if err != nil {
		utils.PrintError(err)
		return schemas.DbUser{}, false
	}
	err = json.Unmarshal(jsonSettings, &dbUser.Settings)
	if err != nil {
		utils.PrintError(err)
		return schemas.DbUser{}, false
	}
	err = json.Unmarshal(jsonStats, &dbUser.Stats)
	if err != nil {
		utils.PrintError(err)
		return schemas.DbUser{}, false
	}
	return *dbUser, true
}
func SearchUsersByUsername(username string) []schemas.DbUser {
	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	stmt, err := db.Prepare("SELECT id, username FROM users WHERE username LIKE ? LIMIT 10")
	if err != nil {
		utils.PrintError(err)
		return []schemas.DbUser{}
	}
	defer stmt.Close()
	rows, err := stmt.Query(username + "%")
	if err != nil {
		utils.PrintError(err)
		return []schemas.DbUser{}
	}
	defer rows.Close()
	var dbUsers []schemas.DbUser
	for rows.Next() {
		dbUser := new(schemas.DbUser)
		err = rows.Scan(&dbUser.Id, &dbUser.Username)
		if err != nil {
			utils.PrintError(err)
			return []schemas.DbUser{}
		}
		dbUsers = append(dbUsers, *dbUser)
	}
	return dbUsers
}

// GetUsersBy returns all users matching to the given parameters in the database
func GetUsersBy(selector string, value interface{}) []schemas.DbUser {
	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	stmt, err := db.Prepare(fmt.Sprintf("SELECT id, username, email, password, registration_date, rank, last_password_refresh, friends, settings, stats FROM users WHERE %s = ?", selector))
	if err != nil {
		utils.PrintError(err)
		return []schemas.DbUser{}
	}
	defer stmt.Close()
	rows, err := stmt.Query(value)
	if err != nil {
		utils.PrintError(err)
		return []schemas.DbUser{}
	}
	defer rows.Close()
	var dbUsers []schemas.DbUser
	for rows.Next() {
		dbUser := new(schemas.DbUser)
		var jsonFriends, jsonSettings, jsonStats []byte
		err = rows.Scan(&dbUser.Id, &dbUser.Username, &dbUser.Email, &dbUser.Password, &dbUser.RegistrationDate, &dbUser.Rank, &dbUser.LastPasswordRefresh, &jsonFriends, &jsonSettings, &jsonStats)
		if err != nil {
			utils.PrintError(err)
			return []schemas.DbUser{}
		}
		err = json.Unmarshal(jsonFriends, &dbUser.Friends)
		if err != nil {
			utils.PrintError(err)
			return []schemas.DbUser{}
		}
		err = json.Unmarshal(jsonSettings, &dbUser.Settings)
		if err != nil {
			utils.PrintError(err)
			return []schemas.DbUser{}
		}
		err = json.Unmarshal(jsonStats, &dbUser.Stats)
		if err != nil {
			utils.PrintError(err)
			return []schemas.DbUser{}
		}
		dbUsers = append(dbUsers, *dbUser)
	}
	return dbUsers
}

// UpdateUserById updates in the database the user with the given id the given value and returns any error
func UpdateUserById(id uint, field string, value interface{}) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare(fmt.Sprintf("UPDATE users SET %s = ? WHERE id = ?", field))
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(value, id)
	return err
}

// DeleteUserById deletes the user with the given id from the database and returns any error
func DeleteUserById(id uint) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	return err
}

//AWAIt

func CreateAwaitUser(username, email, password string) error { // return at least user id ?
	regDate := time.Now().Unix()
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	stmt, err := tx.Prepare("INSERT INTO awaitUsers(username, email, password, registration_date) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, email, password, regDate)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func GetAwaitUsersByUsername(username string) []schemas.DbAwaitUser {
	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	stmt, err := db.Prepare(fmt.Sprintf("SELECT * FROM awaitUsers WHERE username = ?"))
	if err != nil {
		utils.PrintError(err)
		return []schemas.DbAwaitUser{}
	}
	defer stmt.Close()
	rows, err := stmt.Query(username)
	if err != nil {
		utils.PrintError(err)
		return []schemas.DbAwaitUser{}
	}
	defer rows.Close()
	var dbAwaitUsers []schemas.DbAwaitUser
	for rows.Next() {
		dbAwaitUser := new(schemas.DbAwaitUser)
		err = rows.Scan(&dbAwaitUser.Id, &dbAwaitUser.Username, &dbAwaitUser.Email, &dbAwaitUser.Password, &dbAwaitUser.RegistrationDate)
		if err != nil {
			utils.PrintError(err)
			return []schemas.DbAwaitUser{}
		}
		dbAwaitUsers = append(dbAwaitUsers, *dbAwaitUser)
	}
	return dbAwaitUsers
}

func GetAwaitUsersBy(selector string, value interface{}) []schemas.DbAwaitUser {
	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	stmt, err := db.Prepare(fmt.Sprintf("SELECT * FROM awaitUsers WHERE %s = ?", selector))
	if err != nil {
		utils.PrintError(err)
		return []schemas.DbAwaitUser{}
	}
	defer stmt.Close()
	rows, err := stmt.Query(value)
	if err != nil {
		utils.PrintError(err)
		return []schemas.DbAwaitUser{}
	}
	defer rows.Close()
	var dbAwaitUsers []schemas.DbAwaitUser
	for rows.Next() {
		dbAwaitUser := new(schemas.DbAwaitUser)
		err = rows.Scan(&dbAwaitUser.Id, &dbAwaitUser.Username, &dbAwaitUser.Email, &dbAwaitUser.Password, &dbAwaitUser.RegistrationDate)
		if err != nil {
			utils.PrintError(err)
			return []schemas.DbAwaitUser{}
		}
		dbAwaitUsers = append(dbAwaitUsers, *dbAwaitUser)
	}
	return dbAwaitUsers
}

func CheckAwaitUserExistence(username, email string) (string, bool, error) {
	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	var shouldBeZero int
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM awaitUsers WHERE username = ?)", username).Scan(&shouldBeZero)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", false, err
		}
	}
	if shouldBeZero == 1 { // username is already registered
		return "username", true, nil
	}

	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM awaitUsers WHERE email = ?)", email).Scan(&shouldBeZero)
	if err != nil {
		if err != sql.ErrNoRows {
			return "", false, err
		}
	}
	if shouldBeZero == 1 { // email is already registered
		return "email", true, nil
	}

	return "", false, nil
}

func DeleteAwaitUserBy(username string) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()
	stmt, err := db.Prepare("DELETE FROM awaitUsers WHERE username = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(username)
	return err
}
