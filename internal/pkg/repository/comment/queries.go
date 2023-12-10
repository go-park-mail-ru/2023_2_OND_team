package comment

const (
	InsertNewComment = "INSERT INTO comment (author, pin_id, content) VALUES ($1, $2, $3) RETURNING id;"

	UpdateCommentOnDeleted = "UPDATE comment SET deleted_at = now() WHERE id = $1;"

	SelectCommentByID     = "SELECT author, pin_id, content FROM comment WHERE id = $1 AND deleted_at IS NULL;"
	SelectCommentsByPinID = `SELECT c.id, p.id, p.username, p.avatar, c.content
							 FROM comment AS c INNER JOIN profile AS p
							 ON c.author = p.id
							 WHERE c.pin_id = $1 AND (c.id < $2 OR $2 = 0) AND c.deleted_at IS NULL
							 ORDER BY c.id DESC
							 LIMIT $3;`
)
