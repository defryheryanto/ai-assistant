package repository

const (
	queryFindUserByWhatsAppJID = `
		SELECT
			id,
			name,
			whatsapp_jid,
			role,
			email
		FROM users
		WHERE whatsapp_jid = $1;
	`
)
