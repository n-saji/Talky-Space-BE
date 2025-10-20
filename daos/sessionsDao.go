package daos

import "talky-space-be/models"

func (d *Daos) SaveSession(session *models.Sessions) error {
	if err := d.dbConn.Create(session).Error; err != nil {
		return err
	}
	return nil
}
