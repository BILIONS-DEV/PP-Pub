package mysql

import "time"

type TableQizPosts struct {
	ID                  int64     `gorm:"column:id"`
	PostAuthor          int64     `gorm:"column:post_author"`
	PostDate            time.Time `gorm:"column:post_date"`
	PostDateGMT         time.Time `gorm:"column:post_date_gmt"`
	PostContent         string    `gorm:"column:post_content"`
	PostTitle           string    `gorm:"column:post_title"`
	PostExcerpt         string    `gorm:"column:post_excerpt"`
	PostStatus          string    `gorm:"column:post_status"`
	CommentStatus       string    `gorm:"column:comment_status"`
	PingStatus          string    `gorm:"column:ping_status"`
	PostPassword        string    `gorm:"column:post_password"`
	PostName            string    `gorm:"column:post_name"`
	ToPing              string    `gorm:"column:to_ping"`
	Pinged              string    `gorm:"column:pinged"`
	PostModified        time.Time `gorm:"column:post_modified"`
	PostModifiedGmt     time.Time `gorm:"column:post_modified_gmt"`
	PostContentFiltered string    `gorm:"column:post_content_filtered"`
	PostParent          int64     `gorm:"column:post_parent"`
	Guid                string    `gorm:"column:guid"`
	MenuOrder           int64     `gorm:"column:menu_order"`
	PostType            string    `gorm:"column:post_type"`
	PostMimeType        string    `gorm:"column:post_mime_type"`
	CommentCount        int64     `gorm:"column:comment_count"`
	CacheSync           string    `gorm:"column:cache_sync"`
}

func (TableQizPosts) TableName() string {
	return Tables.QizPosts
}
