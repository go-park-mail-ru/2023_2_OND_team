package message

const (
	SelectMessageByID = "SELECT user_from, user_to, content FROM message WHERE id = $1 AND deleted_at IS NULL;"
	SelectUserChats   = `SELECT max(message.id) AS mmid, profile.id, username, avatar 
						 FROM message INNER JOIN profile ON (user_to = $1 AND user_from = profile.id) OR (user_to = profile.id AND user_from = $1)
						 WHERE (message.id < $2 OR $2 = 0)
						 GROUP BY profile.id
						 ORDER BY mmid DESC
						 LIMIT $3;`
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
