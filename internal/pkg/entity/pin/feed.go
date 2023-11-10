package pin

type FeedPin struct {
	Condition
	Pins []Pin `json:"pins"`
}

type Condition struct {
	MinID int `json:"minID"`
	MaxID int `json:"maxID"`
}

type FeedPinConfig struct {
	Condition
	Count      int
	UserID     int
	BoardID    int
	Protection protection
	Liked      bool
	Deleted    bool
}

type protection int8

const (
	_ protection = iota
	FeedProtectionPublic
	FeedProtectionPrivate
	FeedAll
)
