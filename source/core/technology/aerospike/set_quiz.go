package aerospike

type SetQuiz struct {
	ListQuiz []Quiz `json:"list_quiz"`
}

type Quiz struct {
	ID           int64      `json:"id"`
	Title        string     `json:"title"`
	Illustration string     `json:"illustration"`
	ContentType  int64      `json:"content_type"`
	Category     string     `json:"category"`
	Questions    []Question `json:"questions"`
}

type Question struct {
	ID      int64  `json:"id"`
	Ques    string `json:"ques"`
	Ask     []Ask  `json:"ask"`
	Correct int64  `json:"correct"`
	ImgBG   string `json:"imgBG"`
}

type Ask struct {
	ID  int64  `json:"id"`
	Val string `json:"val"`
	Img string `json:"img"`
}
