package roll

const (
	InsertRollAnswer    = "INSERT INTO roll (id, user_id, question_id, answer) VALUES ($1, $2, $3, $4);"
	CheckUserFilledRoll = "SELECT roll_id FROM roll WHERE roll_id = $1 AND user_id = $2 LIMIT 1;"
	SelectHistStat      = `
	SELECT
		answer, COUNT(user_id) AS frequency
	FROM
		roll
	WHERE
		id = $1 AND question_id = $2
	GROUP BY
		answer
	ORDER BY
		answer;`
)
