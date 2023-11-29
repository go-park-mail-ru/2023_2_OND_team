package user

var (
	InsertNewUser = "INSERT INTO profile (username, password, email) VALUES ($1, $2, $3);"

	SelectAuthByUsername         = "SELECT id, password, email FROM profile WHERE username = $1;"
	SelectUsernameAndAvatar      = "SELECT username, avatar FROM profile WHERE id = $1;"
	SelectUserDataExceptPassword = "SELECT username, email, avatar, name, surname, about_me FROM profile WHERE id = $1;"

	UpdateAvatarProfile    = "UPDATE profile SET avatar = $1 WHERE id = $2;"
	SelectUserIdByUsername = "SELECT id FROM profile WHERE username = $1;"
	SelectLastUserID       = "SELECT id FROM profile ORDER BY id DESC LIMIT 1;"
	CheckUserExistence     = "SELECT username FROM profile WHERE id = $1 AND deleted_at IS NULL;"
	GetUserInfo            = `
		SELECT
			p1.id, p1.username, p1.avatar, COALESCE(p1.name, '') name, COALESCE(p1.surname, '') surname, COALESCE(p1.about_me, '') about_me, s2.who IS NOT NULL as is_subscribed, COUNT(s1.who) subscribers
		FROM
			profile p1
		LEFT JOIN 
			subscription_user s1 ON p1.id = s1.whom
		LEFT JOIN
			profile p2 ON s1.who = p2.id
		LEFT JOIN
			subscription_user s2 ON s1.whom = s2.whom AND s2.who = $1
		WHERE
			p1.id = $2 AND p1.deleted_at IS NULL AND p2.deleted_at IS NULL
		GROUP BY
			p1.id, p1.username,p1.avatar, p1.name, p1.surname, p1.about_me, s2.who IS NOT NULL;
		`
	GetProfileInfo = `
		SELECT
			p1.id, p1.username, p1.avatar, COUNT(s.who) subscribers
		FROM
			profile p1
		LEFT JOIN 
			subscription_user s ON p1.id = s.whom
		LEFT JOIN
			profile p2 ON s.who = p2.id
		WHERE
			p1.id = $1 AND p1.deleted_at IS NULL AND p2.deleted_at IS NULL
		GROUP BY
			p1.id, p1.username, p1.avatar;
	`
)
