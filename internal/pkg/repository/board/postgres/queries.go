package board

const (
	InsertBoardQuery                      = "INSERT INTO board (author, title, description, public) VALUES ($1, $2, $3, $4) RETURNING id;"
	SelectBoardAuthorByBoardIdQuery       = "SELECT author FROM board WHERE id = $1;"
	SelectBoardContributorsByBoardIdQuery = "SELECT user_id FROM contributor WHERE board_id = $1;"
)
