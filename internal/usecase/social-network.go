package usecase

import (
	"github.com/Demacr/otus-hl-socialnetwork/internal/domain"
	"github.com/Demacr/otus-hl-socialnetwork/internal/storages"
	"github.com/pkg/errors"
)

type socialNetworkUsecase struct {
	repo        storages.SocialNetworkRepository
	cache       storages.CacheRepository
	feedChan    chan<- *domain.Post
	rebuildChan chan<- int
}

func NewSocialNetworkUsecase(snRepo storages.SocialNetworkRepository, cacheRepo storages.CacheRepository, feeder domain.FeederUsecase) domain.SocialNetworkUsecase {
	return &socialNetworkUsecase{
		repo:     snRepo,
		cache:    cacheRepo,
		feedChan: feeder.GetFeedUpdateChannel(),
	}
}

func (sn *socialNetworkUsecase) Registrate(profile *domain.Profile) error {
	if profile.Email == "" && profile.Password == "" {
		return errors.New("email or password missed")
	}
	err := sn.repo.WriteProfile(profile)
	if err != nil {
		err = errors.Wrap(err, "creating profile")
		return err
	}

	return nil
}

func (sn *socialNetworkUsecase) Authorize(credentials *domain.Credentials) (*domain.Profile, error) {
	result, err := sn.repo.CheckCredentials(credentials)
	if err != nil {
		err = errors.Wrap(err, "authorization failed")
		return nil, err
	}
	if !result {
		return nil, domain.ErrWrongCredentials
	}

	return sn.GetProfileInfo(credentials.Email)
}

func (sn *socialNetworkUsecase) GetProfileInfo(email string) (*domain.Profile, error) {
	result, err := sn.repo.GetProfileByEmail(email)
	if err != nil {
		err = errors.Wrap(err, "get profile error")
		return nil, err
	}
	return result, nil
}

func (sn *socialNetworkUsecase) GetRandomProfiles(id int) ([]domain.Profile, error) {
	result, err := sn.repo.GetRandomProfiles(id)
	if err != nil {
		err = errors.Wrap(err, "getting random profiles")
		return nil, err
	}
	return result, nil
}

func (sn *socialNetworkUsecase) GetProfilesBySearchPrefixes(first_name string, last_name string) ([]domain.Profile, error) {
	result, err := sn.repo.GetProfilesBySearchPrefixes(first_name, last_name)
	if err != nil {
		err = errors.Wrap(err, "searching profiles")
		return nil, err
	}
	return result, nil
}

func (sn *socialNetworkUsecase) CreateFriendRequest(from, to int) error {
	created, err := sn.repo.CreateFriendRequest(from, to)
	if err != nil {
		err = errors.Wrap(err, "creating friendship request")
		return err
	}
	if !created {
		return domain.ErrFriendshipRequestExists
	}

	return nil
}

func (sn *socialNetworkUsecase) FriendshipRequestAccept(id1, id2 int) error {
	accepted, err := sn.repo.AcceptFriendship(id1, id2)
	if err != nil {
		return errors.Wrap(err, "accepting friendship request")
	}
	if !accepted {
		return domain.ErrFriendshipRequestNotExists
	}

	sn.rebuildChan <- id1

	return nil
}

func (sn *socialNetworkUsecase) FriendshipRequestDecline(id1, id2 int) error {
	declined, err := sn.repo.DeclineFriendship(id1, id2)
	if err != nil {
		return errors.Wrap(err, "declining friendship request")
	}
	if !declined {
		return domain.ErrFriendshipRequestNotExists
	}

	return nil
}

func (sn *socialNetworkUsecase) GetFriendshipRequests(id int) ([]domain.FriendRequest, error) {
	fr, err := sn.repo.GetFriendRequests(id)
	if err != nil {
		err = errors.Wrap(err, "cannot get friendship requests")
		return nil, err
	}

	return fr, nil
}

func (sn *socialNetworkUsecase) GetRelatedProfile(id, related_id int) (*domain.RelatedProfile, error) {
	result, err := sn.repo.GetRelatedProfileById(id, related_id)
	if err != nil {
		err = errors.Wrap(err, "getting related profile")
		return nil, err
	}

	return result, nil
}

func (sn *socialNetworkUsecase) CreatePost(profileId int, post *domain.Post) error {
	post.ProfileId = profileId

	//TODO: add validation for empty posts
	id, err := sn.repo.CreatePost(profileId, post)
	if err != nil {
		err = errors.Wrap(err, "creating post")
		return err
	}

	post.Id = id

	err = sn.cache.SetPost(post)
	if err != nil {
		return errors.Wrap(err, "Usecase.CreatePost.Cache.SetPost")
	}

	sn.feedChan <- post

	return nil
}
func (sn *socialNetworkUsecase) UpdatePost(profileId int, post *domain.Post) error {
	post.ProfileId = profileId

	err := sn.repo.UpdatePost(profileId, post)
	if err != nil {
		err = errors.Wrap(err, "updating post")
		return err
	}

	err = sn.cache.SetPost(post)
	if err != nil {
		return errors.Wrap(err, "Usecase.UpdatePost.Cache.SetPost")
	}

	return nil
}
func (sn *socialNetworkUsecase) DeletePost(profileId int, post *domain.Post) error {
	err := sn.repo.DeletePost(profileId, post)
	if err != nil {
		err = errors.Wrap(err, "deleting post")
		return err
	}

	err = sn.cache.DeletePost(post.Id)
	if err != nil {
		return errors.Wrap(err, "Usecase.DeletePost.Cache.DeletePost")
	}

	return nil
}
func (sn *socialNetworkUsecase) GetPost(postId int) (*domain.Post, error) {
	post, err := sn.cache.GetPost(postId)
	if err != nil {
		return nil, errors.Wrap(err, "Usecase.GetPost.Cache.GetPost")
	}

	if *post != (domain.Post{}) {
		return post, nil
	}

	post, err = sn.repo.GetPost(postId)
	if err != nil {
		return nil, errors.Wrap(err, "getting post")
	}

	err = sn.cache.SetPost(post)
	if err != nil {
		return nil, errors.Wrap(err, "Usecase.GetPost.Cache.SetPost")
	}

	return post, nil
}
