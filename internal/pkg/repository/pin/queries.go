package pin

var (
	SelectAfterIdWithLimit = "SELECT id, picture FROM pin WHERE id > $1 ORDER BY id LIMIT $2;"

	InsertLikePinFromUser = "INSERT INTO like_pin (pin_id, user_id) VALUES ($1, $2);"

	UpdatePinSetStatusDelete = "UPDATE pin SET deleted_at = now() WHERE id = $1 AND author = $2;"

	DeleteLikePinFromUser = "DELETE FROM like_pin WHERE pin_id = $1 AND user_id = $2;"
)
