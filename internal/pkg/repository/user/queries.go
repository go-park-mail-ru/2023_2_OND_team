package user

var (
	InsertNewUser = "INSERT INTO profile (username, password, email) VALUES ($1, $2, $3);"

	SelectAuthByUsername         = "SELECT id, password, email FROM profile WHERE username = $1;"
	SelectUsernameAndAvatar      = "SELECT username, avatar FROM profile WHERE id = $1;"
	SelectUserDataExceptPassword = "SELECT username, email, avatar, name, surname FROM profile WHERE id = $1;"

	UpdateAvatarProfile    = "UPDATE profile SET avatar = $1 WHERE id = $2;"
	SelectUserIdByUsername = "SELECT id FROM profile WHERE username = $1;"
	SelectLastUserID       = "SELECT id FROM profile ORDER BY id DESC LIMIT 1;"
)
