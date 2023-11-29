package model

import "source/core/technology/mysql"

type RlPlaylistContent struct{}

type RlsPlaylistContentRecord struct {
	mysql.TableRlsPlaylistContent
}

func (RlsPlaylistContentRecord) TableName() string {
	return mysql.Tables.RlPlaylistContent
}

func (RlPlaylistContent) Create(rl RlsPlaylistContentRecord) (err error) {
	err = mysql.Client.Create(&rl).Error
	return err
}

func (RlPlaylistContent) DeleteAllByPlaylistId(playlistId int64) {
	mysql.Client.Where(RlsPlaylistContentRecord{mysql.TableRlsPlaylistContent{PlaylistId: playlistId}}).Delete(RlsPlaylistContentRecord{})
	return
}

func (RlPlaylistContent) GetByPlaylistId(playlistId int64) (records []RlsPlaylistContentRecord) {
	mysql.Client.Where(RlsPlaylistContentRecord{mysql.TableRlsPlaylistContent{PlaylistId: playlistId}}).Find(&records)
	return
}

func (RlPlaylistContent) DeleteAllByContentId(contentId int64) {
	mysql.Client.Where(RlsPlaylistContentRecord{mysql.TableRlsPlaylistContent{ContentId: contentId}}).Delete(RlsPlaylistContentRecord{})
	return
}
