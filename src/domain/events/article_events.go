package events

type ArticleCreatedEvent struct {
	articleId string
}

type ArticleUpdatedEvent struct {
	articleId string
}
type ArticleDeletedEvent struct {
	articleId string
}
