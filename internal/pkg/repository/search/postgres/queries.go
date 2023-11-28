package search

// id, username, avatar, is_subcribed, subsCount

const (
	SelectUsersByUsername = `
	SELECT
		p1.id, p1.username, p1.avatar, COUNT(s1.who) AS subscribers, s2.who IS NOT NULL AS is_subscribed
	FROM
		profile p1
	LEFT JOIN
		subscription_user s1 ON p1.id = s1.whom
	LEFT JOIN
		profile p2 ON s1.who = p2.id
	LEFT JOIN
		subscription_user s2 ON p1.id = s2.whom AND s2.who = $1 -- curr user
	WHERE
		p1.deleted_at IS NULL AND p2.deleted_at IS NULL AND p1.username ILIKE $2 -- AND p1.id < $3 --lastID and template
	GROUP BY
		p1.id, p1.username, p1.avatar, s2.who IS NOT NULL
	ORDER BY
		p1.id DESC
	LIMIT
		$3 --count
	OFFSET
		$4;`
	SelectBoardsByTitle = `
	SELECT
		board.id, 
		board.title, 
		board.created_at, 
		board.public,
		COUNT(DISTINCT pin.id) FILTER (WHERE pin.deleted_at IS NULL) AS pins_number,
		COALESCE((ARRAY_AGG(DISTINCT pin.picture) FILTER (WHERE pin.deleted_at IS NULL AND pin.picture is not null))[:3], ARRAY[]::TEXT[]) AS pins
	FROM
		board
	LEFT JOIN
		membership ON board.id = membership.board_id
	LEFT JOIN
		pin ON membership.pin_id = pin.id
	WHERE 
		board.title ILIKE $1 AND (board.public OR board.author = $2 OR $2 IN (SELECT user_id FROM contributor WHERE board_id = board.id))
	GROUP BY
		board.id, board.title, board.created_at
	ORDER BY
		board.id DESC
	LIMIT
		$3
	OFFSET
		$4;`
)
