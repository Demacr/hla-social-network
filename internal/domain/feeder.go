package domain

type FeederUsecase interface {
	GetFeedUpdateChannel() chan<- *Post
	GetFeedRebuildChannel() chan<- int
	GetFeedIds(int) ([]int, error)
	RebuildFeeds() error
}
