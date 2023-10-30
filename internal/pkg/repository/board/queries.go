package board

const (
	CreateBoardQuery = "INSERT INTO board (author, title, description, public) VALUES ($1 $2 $3 $4) RETURNING id;"
)
