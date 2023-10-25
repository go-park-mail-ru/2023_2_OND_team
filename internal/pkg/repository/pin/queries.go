package pin

var (
	SelectAfterIdWithLimit = "SELECT id, picture FROM pin WHERE id > $1 ORDER BY id LIMIT $2;"
)
