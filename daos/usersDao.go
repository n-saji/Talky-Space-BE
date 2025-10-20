package daos

import (
	"talky-space-be/models"
)

func (d *Daos) CreateUser(user *models.User) error {
	if err := d.dbConn.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (d *Daos) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	if err := d.dbConn.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Daos) GetUserByID(id string) (*models.User, error) {
	var user models.User
	if err := d.dbConn.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Daos) GetByPhoneNumber(phoneNumber string) (*models.User, error) {
	var user models.User
	if err := d.dbConn.Where("phone_number = ?", phoneNumber).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *Daos) UpdateUser(user *models.User) error {
	if err := d.dbConn.Save(user).Error; err != nil {
		return err
	}
	return nil
}

func (d *Daos) DeleteUser(user *models.User) error {
	if err := d.dbConn.Delete(user).Error; err != nil {
		return err
	}
	return nil
}