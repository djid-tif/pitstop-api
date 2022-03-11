package friends

import (
	"encoding/json"
	"errors"
	"pitstop-api/src/database"
	"pitstop-api/src/users"
	"pitstop-api/src/utils"
)

const (
	typeFriendConfirmed = "friends"
	typeFriendSent      = "friends_sent"
	typeFriendReceived  = "friends_received"
)

type friendsLists struct {
	ConfirmedFriends []uint `json:"confirmed_friends"`
	SentFriends      []uint `json:"sent_friends"`
	ReceivedFriends  []uint `json:"received_friends"`
}

func addFriend(self, other uint, friendType string) error {
	if !isFriendTypeValid(friendType) {
		return errors.New("type d'ami invalide: " + friendType)
	}

	return setFriendsList(self, append(database.GetFriendsList(self, friendType), other), friendType)
}

func removeFriend(self, other uint, friendType string) error {
	if !isFriendTypeValid(friendType) {
		return errors.New("type d'ami invalide: " + friendType)
	}

	sList := database.GetFriendsList(self, friendType)
	newList, removed := utils.RemoveUintFromList(other, sList)
	if removed {
		return setFriendsList(self, newList, friendType)
	}

	return nil
}

func setFriendsList(id uint, list []uint, friendType string) error {
	if !isFriendTypeValid(friendType) {
		return errors.New("type d'ami invalide: " + friendType)
	}

	jsonList, err := json.Marshal(list)
	if err != nil {
		return err
	}

	return database.UpdateUserById(id, friendType, jsonList)
}

func sendFriendRequest(self, other users.User) error {
	err := addFriend(self.GetId(), other.GetId(), typeFriendSent)
	if err != nil {
		return err
	}
	return addFriend(other.GetId(), self.GetId(), typeFriendReceived)
}

func confirmFriendship(self, other users.User) error {
	if err := removeFriend(self.GetId(), other.GetId(), typeFriendSent); err != nil {
		return err
	}
	if err := removeFriend(self.GetId(), other.GetId(), typeFriendReceived); err != nil {
		return err
	}
	if err := removeFriend(other.GetId(), self.GetId(), typeFriendSent); err != nil {
		return err
	}
	if err := removeFriend(other.GetId(), self.GetId(), typeFriendReceived); err != nil {
		return err
	}

	if err := addFriend(self.GetId(), other.GetId(), typeFriendConfirmed); err != nil {
		return err
	}
	if err := addFriend(other.GetId(), self.GetId(), typeFriendConfirmed); err != nil {
		return err
	}

	return nil
}

func isFriendTypeValid(friendType string) bool {
	return friendType == typeFriendConfirmed || friendType == typeFriendSent || friendType == typeFriendReceived
}
