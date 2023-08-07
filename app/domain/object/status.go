package object

import "time"

type Status struct {
	// The text content
	Content string

	// The internal id
	ID int64

	// The account that posted this status
	PostedBy *Account `json:"account"`

	// The time this status was posted
	CreateAt time.Time `json:"create_at,omitempty" db:"create_at"`

	AttachedMedias []*AttachedMedia `json:"media_attachments"`
}

// メディアを別に投稿して、それを参照するという構造から考えて、このようなデータの持ちからが自然だと思ったのですが
// 求められているJSONの形状と違うという問題を解消し終わるまえに時間切れになりました
// DBのスキーマも同様にMediaとAttachmentが1:多の関係だったので従属性を考えて分けました
type AttachedMedia struct {
	// content of the attached media
	Content Media
	// description of the media in a status
	Description string
}

func (s *Status) SetCreateAt() {
	s.CreateAt = time.Now()
}
