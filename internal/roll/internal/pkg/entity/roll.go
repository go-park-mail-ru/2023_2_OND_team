package roll

type RollAnswer struct {
	UserID     int
	RollID     int
	QuestionID int
	Answer     string
}

type HistStatObj struct {
	Answer    string
	Frequency int
}
