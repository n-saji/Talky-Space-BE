package daos

import "talky-space-be/models"

func (d *Daos) CreateChatroom(chatroom *models.Chatroom) error {
	if err := d.dbConn.Create(chatroom).Error; err != nil {
		return err
	}
	return nil
}

func (d *Daos) GetChatroomByID(id string) (*models.Chatroom, error) {
	var chatroom models.Chatroom
	if err := d.dbConn.Where("id = ?", id).First(&chatroom).Error; err != nil {
		return nil, err
	}
	return &chatroom, nil
}

func (d *Daos) UpdateChatroom(chatroom *models.Chatroom) error {
	if err := d.dbConn.Save(chatroom).Error; err != nil {
		return err
	}
	return nil
}

func (d *Daos) DeleteChatroom(id string) error {
	if err := d.dbConn.Delete(&models.Chatroom{}, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (d *Daos) CheckChatroomExistForSenderReceiver(senderID string, receiverID string) (*models.Chatroom, error) {
	var chatroom *models.Chatroom
	if err := d.dbConn.Raw(`select c.* from chatrooms c
		JOIN chatroom_members cm1 ON c.id = cm1.chatroom_id
		JOIN chatroom_members cm2 ON c.id = cm2.chatroom_id
		WHERE c.is_group = false AND cm1.user_id = ? AND cm2.user_id = ? LIMIT 1`, senderID, receiverID).Scan(&chatroom).Error; err != nil {
		return nil, err
	}
	if chatroom.ID.String() == "" {
		return nil, nil
	}

	return chatroom, nil
}
