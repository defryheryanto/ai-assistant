package repository

const (
	queryFindUserByWhatsAppJID = `
		SELECT
			id,
			name,
			whatsapp_jid,
			role
		FROM users
		WHERE whatsapp_jid = $1;
	`
)
