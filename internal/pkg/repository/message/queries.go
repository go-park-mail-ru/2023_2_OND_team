package message

const (
	SelectMessageByID     = "SELECT user_from, user_to, content FROM message WHERE id = $1 AND deleted_at IS NULL;"
	SelectMessageFromChat = `SELECT id, user_from, user_to, content
							 FROM message 
							 WHERE deleted_at IS NULL AND (id < $1 OR $1 = 0) AND
							 	(user_from = $2 AND user_to = $3 OR user_from = $3 AND user_to = $2)
							 ORDER BY id DESC
							 LIMIT $4;`

	InsertMessage = "INSERT INTO message (user_from, user_to, content) VALUES ($1, $2, $3) RETURNING id;"

	UpdateMessageContent         = "UPDATE message SET content = $1 WHERE id = $2;"
	UpdateMessageStatusToDeleted = "UPDATE message SET deleted_at = now() WHERE id = $1 AND deleted_at IS NULL;"
)
