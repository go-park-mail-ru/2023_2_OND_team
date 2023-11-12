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
	userID     int
	boardID    int
	Protection protection
	Liked      bool
	Deleted    bool
	hasUser    bool
	hasBoard   bool
}

func (cfg *FeedPinConfig) SetBoard(boardID int) {
	cfg.boardID = boardID
	cfg.hasBoard = true
}

func (cfg *FeedPinConfig) SetUser(userID int) {
	cfg.userID = userID
	cfg.hasUser = true
}

func (cfg *FeedPinConfig) Board() (int, bool) {
	return cfg.boardID, cfg.hasBoard
}

func (cfg *FeedPinConfig) User() (int, bool) {
	return cfg.userID, cfg.hasUser
}

type protection int8

const (
	_ protection = iota
	FeedProtectionPublic
	FeedProtectionPrivate
	FeedAll
)
