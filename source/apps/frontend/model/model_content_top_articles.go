package model

import (
	"source/core/technology/mysql"
)

type ContentTopArticles struct{}

type ContentTopArticlesRecord struct {
	mysql.TableContentTopArticles
}

func (ContentTopArticlesRecord) TableName() string {
	return mysql.Tables.ContentTopArticles
}

func (t *ContentTopArticles) GetContent(domain string, typ string, feedUrl string) (record ContentTopArticlesRecord) {
	query := ContentTopArticlesRecord{mysql.TableContentTopArticles{
		Domain: domain,
		Type:   typ,
	}}
	if typ == "feed" {
		query.FeedUrl = feedUrl
	}
	mysql.Client.Where(query).Find(&record)
	return
}
