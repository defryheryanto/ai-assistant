package repository

const (
	queryFindByJID = `
		SELECT
			id,
			created_at,
			updated_at,
			group_jid,
			is_active,
			registered_by
		FROM whatsapp_groups
		WHERE group_jid = $1;
	`
	queryInsert = `
		INSERT INTO whatsapp_groups
			(created_at, updated_at, group_jid, is_active, registered_by)
		VALUES (CURRENT_TIMESTAMP, CURRENT_TIMESTAMP, $1, $2, $3);
	`
)
