package object

type Media struct {
	// internal id
	ID int64

	// The URL where the content of this media is saved.
	URL string

	// Type of this media. For example image, video, sound.
	MediaType string `db:"type"`
}
