package object

import "time"

type Status struct {
	// The text content
	Content string

	// The internal id
	ID int64

	// The account that posted this status
	PostedBy *Account

	// The time this status was posted
	CreateAt time.Time `json:"create_at,omitempty" db:"create_at"`

	AttachedMedias []*AttachedMedia
}

type AttachedMedia struct {
	// content of the attached media
	Content Media
	// description of the media in a status
	Description string
}

func (s *Status) SetCreateAt() {
	s.CreateAt = time.Now()
}
