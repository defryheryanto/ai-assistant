package whatsapp

type CreateGroupParams struct {
	GroupJID     string `json:"group_jid"`
	RegisteredBy string `json:"registered_by"`
}

type CreateUserParams struct {
	Name  string   `json:"name"`
	Phone string   `json:"phone"`
	Role  UserRole `json:"role"`
	Email string   `json:"email"`
}
