package subscription

var (
	CreateSubscriptionUser = "INSERT INTO subscription_user (who, whom) values ($1, $2);"
	DeleteSubscriptionUser = "DELETE FROM subscription_user WHERE who = $1 AND whom = $2;"
	GetUserSubscriptions   = `	
		SELECT 
			p.id, p.username, p.avatar, s.who IS NOT NULL AS is_subscribed
		FROM
			subscription_user f
		LEFT JOIN
			profile p ON f.whom = p.id
		LEFT JOIN
			subscription_user s ON f.whom = s.whom AND s.who = $1
		WHERE
			f.who = $2 AND p.deleted_at IS NULL AND f.whom < $3
		ORDER BY
			f.whom DESC
		LIMIT
			$4;`
	GetUserSubscribers = `	
		SELECT
			p.id, p.username, p.avatar, s.who IS NOT NULL AS is_subscribed
		FROM
			subscription_user f
		LEFT JOIN
			profile p ON f.who = p.id
		LEFT JOIN
			subscription_user s ON f.who = s.whom AND s.who = $1
		WHERE
			f.whom = $2 AND p.deleted_at IS NULL AND f.who < $3
		ORDER BY
			f.who DESC
		LIMIT
			$4;`
)
