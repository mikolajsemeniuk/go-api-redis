package order

import (
	"errors"
	"time"

	"github.com/goccy/go-json"

	"github.com/google/uuid"
)

// Description
type Description string

func (d *Description) UnmarshalJSON(data []byte) error {
	var value string

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	if 3 < len(value) && len(value) < 25 {
		*d = Description(value)
		return nil
	}

	return errors.New("Description has to be between 4 and 25 characters")
}

// Status
type Status uint8

const (
	Accepted Status = iota + 1
	Pending
	Done
	Rejected
)

func (s *Status) UnmarshalJSON(data []byte) error {
	var value uint

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	switch Status(value) {
	case Accepted, Pending, Done, Rejected:
		*s = Status(value)
		return nil
	}

	return errors.New("Status has to be in range of 1 to 4")
}

// Timestamp
type Timestamp time.Time

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	*t = Timestamp(time.Now())
	return nil
}

// Deprecated, remove later
type Order struct {
	Key         uuid.UUID   `json:"key"         example:"7e1efed0-e86a-4286-8c4a-a609751f9a0b"`
	Description Description `json:"description" example:"Delivered"`
	Status      Status      `json:"status"      example:"1"`
	Updated     Timestamp   `json:"updated"     example:"2022-04-17T12:25:43.511Z"`
}
