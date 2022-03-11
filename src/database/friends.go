package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"pitstop-api/src/utils"
	"strconv"
)

// GetFriendsList returns the list of friends of the given type of the user with the given id
func GetFriendsList(id uint, friendType string) []uint {
	if friendType != "friends" && friendType != "friends_sent" && friendType != "friends_received" {
		utils.PrintError(errors.New("type d'ami invalide: " + friendType))
		return []uint{}
	}

	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	stmt, err := db.Prepare(fmt.Sprintf("SELECT %s FROM users WHERE id = ?", friendType))
	if err != nil {
		utils.PrintError(err)
		return []uint{}
	}
	defer stmt.Close()
	rows, err := stmt.Query(id)
	if err != nil {
		utils.PrintError(err)
		return []uint{}
	}
	defer rows.Close()

	if !rows.Next() {
		utils.PrintError(errors.New("pas de colonne " + friendType + " trouvée pour l'utilisateur " + strconv.FormatUint(uint64(id), 10)))
		return []uint{}
	}

	var jsonList string
	err = rows.Scan(&jsonList)
	if err != nil {
		utils.PrintError(err)
		return []uint{}
	}

	var friendsList []uint
	err = json.Unmarshal([]byte(jsonList), &friendsList)
	if err != nil {
		utils.PrintError(err)
		return []uint{}
	}

	return friendsList
}
