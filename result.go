package witai

import "errors"

var (
	ErrInvalidResult = errors.New("invalid result")
)

type Result struct {
	Text     string     `json:"_text"`
	MsgID    string     `json:"msg_id"`
	Outcomes []*Outcome `json:"outcomes"`
}

func (r *Result) IsValid() bool {
	return r.Text != "" && len(r.Outcomes) > 0
}

type Outcome struct {
	Text       string  `json:"_text"`
	Confidence float32 `json:"confidence"`
	Intent     string  `json:"intent"`
}
