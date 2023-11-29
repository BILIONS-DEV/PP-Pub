package mysql

type TableRlsPlaylistContent struct {
	Id         int64 `gorm:"column:id" json:"id"`
	PlaylistId int64 `gorm:"column:playlist_id" json:"playlist_id"`
	ContentId  int64 `gorm:"column:content_id" json:"content_id"`
}
