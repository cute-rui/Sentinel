package oauth

import (
	"Sentinel/dao"
	"Sentinel/dao/models"
	"crypto/rsa"
	"strconv"
)

type Service struct {
	keys map[string]*rsa.PublicKey
}

type UserStore interface {
	GetUserByID(string) (*models.User, error)
	GetUserByUsername(string) (*models.User, error)
	ExampleClientID() string
}

type userStore struct {
}

// ExampleClientID is only used in the example server
func (u userStore) ExampleClientID() string {
	return "service"
}

func (u userStore) GetUserByID(id string) (*models.User, error) {
	numID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}

	return dao.FindUserByID(numID)
}

func (u userStore) GetUserByUsername(username string) (*models.User, error) {
	return dao.FindUserByUsername(username)
}
