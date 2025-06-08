package user

type User struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	WhatsAppJID string `json:"whatsapp_jid"`
	Role        Role   `json:"role"`
}
