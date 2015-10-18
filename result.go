package witai

import (
	"encoding/json"
	"errors"
)

var (
	// ErrInvalidResult is returned when the returned wit.ai result
	// does not contain a valid response.
	ErrInvalidResult = errors.New("invalid result")
)

// Result contains the wit.ai result with a list of outcomes.
type Result struct {
	Text     string     `json:"_text"`
	MsgID    string     `json:"msg_id"`
	Outcomes []*Outcome `json:"outcomes"`
}

// NewResult transforms an input string into a Result object.
// It returns an error if there's any issue parsing the input.
func NewResult(in string) (*Result, error) {
	var r Result
	if err := json.Unmarshal([]byte(in), &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// IsValid returns true if the Result struct contains at least one
// proper outcome with a valid result.
func (r *Result) IsValid() bool {
	return r.Text != "" && len(r.Outcomes) > 0
}

// A Outcome points to each each wit.ai outcome contained within
// the returned Result for a given text or voice query.
type Outcome struct {
	Text       string  `json:"_text"`
	Confidence float32 `json:"confidence"`
	Intent     string  `json:"intent"`
}
