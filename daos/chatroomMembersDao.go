package daos

import "talky-space-be/models"

func (d *Daos) CreateChatroomMember(chatroomMember *models.ChatroomMember) error {
	if err := d.dbConn.Create(chatroomMember).Error; err != nil {
		return err
	}
	return nil
}

func (d *Daos) GetChatroomMembersByChatroomID(chatroomID string) ([]models.ChatroomMember, error) {
	var members []models.ChatroomMember
	if err := d.dbConn.Where("chatroom_id = ?", chatroomID).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (d *Daos) DeleteChatroomMember(chatroomID string, userID string) error {
	if err := d.dbConn.Delete(&models.ChatroomMember{}, "chatroom_id = ? AND user_id = ?", chatroomID, userID).Error; err != nil {
		return err
	}
	return nil
}	

func (d *Daos) IsUserInChatroom(chatroomID string, userID string) (bool, error) {
	var count int64
	if err := d.dbConn.Model(&models.ChatroomMember{}).Where("chatroom_id = ? AND user_id = ?", chatroomID, userID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

