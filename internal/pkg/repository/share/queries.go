package share

const (
	SelectLink = `SELECT board_id, role_id, user_id
				  FROM link LEFT JOIN access_link 
					ON link_id = id
				  WHERE id = $1;`

	InsertNewLink = "INSERT INTO link (board_id, role_id) VALUES ($1, $2) RETURNING id;"
)
