package models

type Message struct {
	Model
	UserID     int    `json:"from_user_id" redis:"from_user_id"`
	TargetID   int    `json:"to_user_id" redis:"to_user_id"`
	Content    string `json:"content" redis:"content"`
	ActionType string `json:"action_type"`
	CreatedAt  int    `json:"create_time" redis:"create_time"`
}

func ChatMessageCreat(m *Message) error {
	return DB.Model(&Message{}).Create(m).Error
}

//func (m Message) GetLastMessage(db *gorm.DB) (*Message, error) {
//	var lastMessage *Message
//	var err error
//	if err = db.Where("target_id = ?", m.TargetID).Last(&lastMessage).Error; err != nil {
//		return nil, err
//	}
//	return lastMessage, nil
//}

func GetMessageByID(user_id, target_id int) (*[]Message, error) {
	var messages []Message
	if err := DB.Model(&Message{}).Where("user_id = ? AND target_id = ?", user_id, target_id).Find(&messages).Error; err != nil {
		return nil, err
	}
	return &messages, nil
}
