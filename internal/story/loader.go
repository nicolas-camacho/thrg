package story

type StoryData struct {
	HolderName          string    `json:"holderName"`
	Title               string    `json:"title"`
	Description         string    `json:"description"`
	MisfortuneThreshold float64   `json:"misfortuneThreshold"`
	Acts                []ActData `json:"acts"`
}

type ActData struct {
	Order   int          `json:"order"`
	Text    string       `json:"text"`
	Options []OptionData `json:"options"`
}

type OptionData struct {
	Text         string            `json:"text"`
	NextActOrder *int              `json:"nextActOrder"`
	Consequences []ConsequenceData `json:"consequences"`
}

type ConsequenceData struct {
	Type  string  `json:"type"`
	Value float64 `json:"value"`
}
