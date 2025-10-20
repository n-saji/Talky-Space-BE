package daos

import "talky-space-be/models"

func (d *Daos) CreateMessage(req models.Messages) (*models.Messages, error) {
	result := d.dbConn.Create(&req)
	if result.Error != nil {
		return nil, result.Error
	}
	return &req, nil
}

func (d *Daos) GetMessagesByChatroomID(chatroomID string) ([]models.Messages, error) {
	var messages []models.Messages
	result := d.dbConn.Where("chatroom_id = ?", chatroomID).Find(&messages).Order("created_at")
	if result.Error != nil {
		return nil, result.Error
	}
	return messages, nil
}

func (d *Daos) DeleteMessagesByChatroomID(chatroomID string) error {
	result := d.dbConn.Where("chatroom_id = ?", chatroomID).Delete(&models.Messages{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

