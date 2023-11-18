package pin

type Tag struct {
	ID    int    `json:"-"`
	Title string `json:"title"`
}
