package subscription

var (
	CreateSubscriptionUser = "INSERT INTO subscription_user (who, whom) values ($1, $2);"
	DeleteSubscriptionUser = "DELETE FROM subscription_user WHERE who = $1 AND whom = $2;"
)
