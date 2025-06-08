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
	queryInsert = `
		INSERT INTO users
			("name", whatsapp_jid, "role", created_at, updated_at, email)
		VALUES($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $4);
	`
)
