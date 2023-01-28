package storages

import "github.com/Demacr/otus-hl-socialnetwork/internal/domain"

type SocialNetworkRepository interface {
	WriteProfile(profile *domain.Profile) error
	GetProfileByEmail(email string) (*domain.Profile, error)
	GetRelatedProfileById(id, related_id int) (*domain.RelatedProfile, error)
	CheckCredentials(credentials *domain.Credentials) (bool, error)
	CreateFriendRequest(id, friend_id int) (bool, error)
	GetRandomProfiles(exclude_id int) ([]domain.Profile, error)
	GetProfilesBySearchPrefixes(first_name string, last_name string) ([]domain.Profile, error)
	GetFriendRequests(id int) ([]domain.FriendRequest, error)
	AcceptFriendship(id, friend_id int) (bool, error)
	DeclineFriendship(id, friend_id int) (bool, error)
}
