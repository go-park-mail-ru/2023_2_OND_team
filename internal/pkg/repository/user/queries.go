package user

var (
	InsertNewUser = "INSERT INTO profile (username, password, email) VALUES ($1, $2, $3);"

	SelectAuthByUsername    = "SELECT id, password, email FROM profile WHERE username = $1;"
	SelectUsernameAndAvatar = "SELECT username, avatar FROM profile WHERE id = $1;"

	UpdateAvatarProfile = "UPDATE profile SET avatar = $1 WHERE id = $2;"
)
