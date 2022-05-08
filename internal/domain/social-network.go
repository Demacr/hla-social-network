package domain

import (
	"github.com/dgrijalva/jwt-go"
)

type Profile struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Age       int    `json:"age"`
	Sex       string `json:"sex"`
	Interests string `json:"interests"`
	City      string `json:"city"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type RelatedProfile struct {
	Profile
	IsFriend      bool `json:"is_friend"`
	IsRequestSent bool `json:"is_request_sent"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	UserId uint
	jwt.StandardClaims
}

type FriendRequest struct {
	FriendID int `json:"friendID"`
}

type SocialNetworkUsecase interface {
	Registrate(*Profile) error
	Authorize(*Credentials) (*Profile, error)
	GetProfileInfo(string) (*Profile, error)
	GetRandomProfiles(int) ([]Profile, error)
	CreateFriendRequest(int, int) error
	FriendshipRequestAccept(int, int) error
	FriendshipRequestDecline(int, int) error
	GetFriendshipRequests(int) ([]FriendRequest, error)
	GetRelatedProfile(int, int) (*RelatedProfile, error)
}
