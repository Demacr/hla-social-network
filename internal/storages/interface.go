package storages

import "github.com/Demacr/otus-hl-socialnetwork/internal/domain"

type CacheRepository interface {
	GetPost(post_id int) (*domain.Post, error)
	SetPost(*domain.Post) error
	DeletePost(post_id int) error
	AddToFeed(profile_id int, post_id ...int) error
	RebuildFeed(profileId int, post_id ...int) error
	GetFeed(profileId int) ([]int, error)
}

type SocialNetworkRepository interface {
	WriteProfile(profile *domain.Profile) error
	GetProfileByEmail(email string) (*domain.Profile, error)
	GetRelatedProfileById(id, related_id int) (*domain.RelatedProfile, error)
	GetLastProfileId() (int, error)
	// Friendship section.
	CreateFriendRequest(id, friend_id int) (bool, error)
	GetRandomProfiles(exclude_id int) ([]domain.Profile, error)
	GetProfilesBySearchPrefixes(first_name string, last_name string) ([]domain.Profile, error)
	GetFriendRequests(id int) ([]domain.FriendRequest, error)
	AcceptFriendship(id, friend_id int) (bool, error)
	DeclineFriendship(id, friend_id int) (bool, error)
	GetFriends(id int) ([]int, error)
	// Post section.
	CreatePost(profile_id int, post *domain.Post) (int, error)
	UpdatePost(profile_id int, post *domain.Post) error
	DeletePost(profile_id int, post *domain.Post) error
	GetPost(post_id int) (*domain.Post, error)
	GetFeedLastN(int, int) ([]int, error)
	// Dialogs section.
	CreateMessage(*domain.Message) error
	GetDialog(id1 int, id2 int) ([]*domain.Message, error)
	GetDialogList(id int) ([]*domain.DialogPreview, error)
}

type FeedManager interface {
	GetFeedUpdateChannel() chan<- *domain.Post
}
