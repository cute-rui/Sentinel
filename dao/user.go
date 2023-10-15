package dao

import (
	"Sentinel/dao/models"
	"errors"
	"gorm.io/gorm"
	"log"
)

func FindUser(email string, username string) (bool, *models.User, error) {
	var user models.User

	if err := Instance.Database.Where(&models.User{Email: email, Username: username}).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, nil
		}

		return false, nil, err
	}

	return true, &user, nil
}

func FindUserByID(id int) (*models.User, error) {
	var user models.User

	if err := Instance.Database.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func FindUserByUsername(username string) (*models.User, error) {
	var user models.User

	if err := Instance.Database.Where(&models.User{Username: username}).First(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func CreateVerify(email, token string) error {
	verify := models.Verify{
		Email:  email,
		Verify: token,
	}

	return Instance.Database.Create(verify).Error
}

func MatchVerify(email, token string) bool {
	verify := models.Verify{Email: email}

	if err := Instance.Database.Where(&models.Verify{Verify: token}).First(&verify).Error; err != nil {
		return false
	}

	if verify.Email != email {
		return false
	}

	defer DeleteVerify(verify.ID)

	return true
}

func DeleteVerify(id int) {
	if err := Instance.Database.Delete(&models.Verify{}, id).Error; err != nil {
		log.Println(err)
	}
}

func CreateUser(user *models.User) error {
	return Instance.Database.Create(user).Error
}

func IsAdmin(uid int) bool {
	var user models.User
	if err := Instance.Database.First(&user, uid).Error; err != nil {
		return false
	}

	return user.IsAdmin
}

func SetApproved(uid int) error {
	user := models.User{
		ID: uid,
	}

	return Instance.Database.Model(&user).Update(`permission`, ``).Error
}
