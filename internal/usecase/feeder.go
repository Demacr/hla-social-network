package usecase

import (
	"log"

	"github.com/Demacr/otus-hl-socialnetwork/internal/domain"
	"github.com/Demacr/otus-hl-socialnetwork/internal/storages"
	"github.com/pkg/errors"
)

type feederUsecase struct {
	chUpdate  chan *domain.Post
	chRebuild chan int
	db        storages.SocialNetworkRepository
	cache     storages.CacheRepository
}

func NewFeederUsecase(db storages.SocialNetworkRepository, cache storages.CacheRepository) domain.FeederUsecase {
	chUpdate := make(chan *domain.Post, 1000)
	chRebuild := make(chan int, 1000)

	fuc := &feederUsecase{
		chUpdate:  chUpdate,
		chRebuild: chRebuild,
		cache:     cache,
		db:        db,
	}

	go func() {
	L:
		for {
			var post *domain.Post
			var ok bool
			if post, ok = <-chUpdate; ok {
				log.Println("DEBUG: Usecase.Feeder.ChanHandler: Get new post")
				friendsList, err := db.GetFriends(post.ProfileID)
				if err != nil {
					log.Println("Usecase.Feeder.Database.GetFriends")
					continue L
				}

				log.Println("DEBUG: Usecase.Feeder.ChanHandler: Friendlist", friendsList)

				for _, friend := range friendsList {
					err := cache.AddToFeed(friend, post.ID)
					if err != nil {
						log.Println(errors.Wrap(err, "Usecase.Feeder.ChanHandler"))
					}
					log.Printf("DEBUG: Usecase.Feeder.ChanHandler: Add %d post to %d friend\n", post.ID, friend)
				}
			}
		}
	}()

	go func() {
		for {
			profileID := <-chRebuild
			if err := fuc.rebuildFeed(profileID); err != nil {
				log.Println(err)
			}
		}
	}()

	return fuc
}

func (fuc *feederUsecase) GetFeedUpdateChannel() chan<- *domain.Post {
	return fuc.chUpdate
}

func (fuc *feederUsecase) GetFeedRebuildChannel() chan<- int {
	return fuc.chRebuild
}

func (fuc *feederUsecase) GetFeedIds(profileID int) ([]int, error) {
	result, err := fuc.cache.GetFeed(profileID)
	if err != nil {
		return nil, errors.Wrap(err, "Usecase.Feeder.GetFeedIds")
	}

	return result, nil
}

func (fuc *feederUsecase) RebuildFeeds() error {
	lastId, err := fuc.db.GetLastProfileId()
	if err != nil {
		return errors.Wrap(err, "Usecase.Feeder.RebuildFeeds.Database.GetLastProfileId")
	}

	for i := 1; i <= lastId; i++ {
		if err = fuc.rebuildFeed(i); err != nil {
			return errors.Wrap(err, "Usecase.Feeder.RebuildFeeds.rebuildFeed")
		}
	}

	return nil
}

func (fuc *feederUsecase) rebuildFeed(profileID int) error {
	postIds, err := fuc.db.GetFeedLastN(profileID, 1000)
	if err != nil {
		return errors.Wrap(err, "Usecase.Feeder.RebuildFeed.Database.GetFeedLastN")
	}

	err = fuc.cache.RebuildFeed(profileID, postIds...)
	if err != nil {
		return errors.Wrap(err, "Usecase.Feeder.RebuildFeed.RebuildFeed")
	}

	return nil
}
