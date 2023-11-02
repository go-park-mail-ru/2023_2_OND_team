package pin

var (
	SelectWithExcludeLimit  = "SELECT id, picture FROM pin WHERE public AND (id < $1 OR id > $2) ORDER BY id DESC LIMIT $3;"
	SelectPinByID           = "SELECT author, title, description, picture, public, deleted_at FROM pin WHERE id = $1;"
	SelectCountLikePin      = "SELECT COUNT(*) FROM like_pin WHERE pin_id = $1;"
	SelectPinByIDWithAuthor = `SELECT author, title, description, picture, public, pin.deleted_at, username, avatar 
	                           FROM pin INNER JOIN profile ON author = profile.id WHERE pin.id = $1;`
	SelectTagsByPinID = `SELECT tag.title FROM pin INNER JOIN pin_tag ON pin.id = pin_tag.pin_id
							   INNER JOIN tag ON pin_tag.tag_id = tag.id WHERE pin.id = $1;`
	SelectCheckAvailability = `SELECT EXISTS (SELECT FROM pin INNER JOIN membership 
											  ON pin.id = membership.pin_id INNER JOIN board 
			                                  ON membership.board_id = board.id INNER JOIN contributor 
			                                  ON board.id = contributor.board_id 
			                                  WHERE pin.id = $1 AND (board.author = $2 OR contributor.user_id = $2));`

	InsertLikePinFromUser       = "INSERT INTO like_pin (pin_id, user_id) VALUES ($1, $2);"
	InsertLikePinFromUserAtomic = `INSERT INTO like_pin (pin_id, user_id)
							 	   SELECT $1, $2 WHERE (
							 			SELECT public OR author = $2 OR
											EXISTS (SELECT FROM pin 
													INNER JOIN membership 
													ON pin.id = membership.pin_id 
												    INNER JOIN board 
													ON membership.board_id = board.id 
											        INNER JOIN contributor 
													ON board.id = contributor.board_id 
													WHERE contributor.user_id = $2 AND pin.id = $1)
								   		FROM pin WHERE id = $1
							       );`

	UpdatePinSetStatusDelete = "UPDATE pin SET deleted_at = now() WHERE id = $1 AND author = $2 AND deleted_at IS NULL;"

	DeleteLikePinFromUser = "DELETE FROM like_pin WHERE pin_id = $1 AND user_id = $2;"
	DeleteAllTagsFromPin  = "DELETE FROM pin_tag WHERE pin_id = $1;"
)
