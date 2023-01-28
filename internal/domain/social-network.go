package domain

import (
	"github.com/golang-jwt/jwt"
)

type Profile struct {
	ID        int    `json:"id"                  faker:"-"`
	Name      string `json:"name,omitempty"      faker:"first_name"`
	Surname   string `json:"surname,omitempty"   faker:"last_name"`
	Age       int    `json:"age,omitempty"       faker:"boundary_start=18, boundary_end=60"`
	Sex       string `json:"sex,omitempty"       faker:"oneof: m, f"`
	Interests string `json:"interests,omitempty" faker:"paragraph"`
	City      string `json:"city,omitempty"      faker:"word"`
	Email     string `json:"email,omitempty"     faker:"email"`
	Password  string `json:"password,omitempty"  faker:"password"`
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
	GetProfilesBySearchPrefixes(string, string) ([]Profile, error)
	CreateFriendRequest(int, int) error
	FriendshipRequestAccept(int, int) error
	FriendshipRequestDecline(int, int) error
	GetFriendshipRequests(int) ([]FriendRequest, error)
	GetRelatedProfile(int, int) (*RelatedProfile, error)
}
