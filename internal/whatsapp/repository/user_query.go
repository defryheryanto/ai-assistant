package repository

const (
	queryFindUserByJID = `
		SELECT
			id,
			name,
			whatsapp_jid,
			role,
			email
		FROM whatsapp_users
		WHERE whatsapp_jid = $1;
	`
	queryInsertUser = `
		INSERT INTO whatsapp_users
			("name", whatsapp_jid, "role", created_at, updated_at, email)
		VALUES($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $4);
	`
)
