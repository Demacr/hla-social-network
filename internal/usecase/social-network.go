package usecase

import (
	"github.com/Demacr/otus-hl-socialnetwork/internal/domain"
	"github.com/Demacr/otus-hl-socialnetwork/internal/storages"
	"github.com/pkg/errors"
)

type socialNetworkUsecase struct {
	repo storages.SocialNetworkRepository
}

func NewSocialNetworkUsecase(snRepo storages.SocialNetworkRepository) domain.SocialNetworkUsecase {
	return &socialNetworkUsecase{
		repo: snRepo,
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
