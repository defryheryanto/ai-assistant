package whatsapp

type CreateGroupParams struct {
	GroupJID     string `json:"group_jid"`
	RegisteredBy string `json:"registered_by"`
}
