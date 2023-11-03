package board

const (
	InsertBoardQuery                      = "INSERT INTO board (author, title, description, public) VALUES ($1, $2, $3, $4) RETURNING id;"
	SelectBoardAuthorByBoardIdQuery       = "SELECT author FROM board WHERE id = $1;"
	SelectBoardContributorsByBoardIdQuery = "SELECT user_id FROM contributor WHERE board_id = $1;"
	UpdateBoardByIdQuery                  = "UPDATE board SET title = $1, description = $2, public = $3 WHERE id = $4;"
	GetContributorBoardsIDs               = "SELECT board_id FROM contributor WHERE user_id = $1;"
)
