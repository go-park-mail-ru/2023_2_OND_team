package board

const (
	InsertBoardQuery                      = "INSERT INTO board (author, title, description, public) VALUES ($1, $2, $3, $4) RETURNING id;"
	SelectBoardAuthorByBoardIdQuery       = "SELECT author FROM board WHERE id = $1 AND deleted_at IS NULL;"
	SelectBoardContributorsByBoardIdQuery = "SELECT user_id FROM contributor WHERE board_id = $1;"
	UpdateBoardByIdQuery                  = "UPDATE board SET title = $1, description = $2, public = $3 WHERE id = $4 AND deleted_at IS NULL;"
	GetContributorBoardsIDs               = "SELECT board_id FROM contributor WHERE user_id = $1;"
	DeleteBoardByIdQuery                  = "UPDATE board SET deleted_at = $1 WHERE id = $2;"
)
