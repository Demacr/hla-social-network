package controller

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Demacr/otus-hl-socialnetwork/internal/domain"
	smock "github.com/Demacr/otus-hl-socialnetwork/internal/storages/mocks"
	"github.com/labstack/echo/v4"
)

func TestNewSocialNetworkHandler(t *testing.T) {
	type args struct {
		e         *echo.Echo
		snuc      domain.SocialNetworkUsecase
		feeduc    domain.FeederUsecase
		JWTSecret string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"First test",
			args{
				echo.New(),
				nil,
				nil,
				"JWTSecret",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewSocialNetworkHandler(tt.args.e, tt.args.snuc, tt.args.feeduc, tt.args.JWTSecret)
		})
	}
}

func TestSocialNetworkHandler_Registrate(t *testing.T) {
	e := echo.New()
	mockRepo := smock.NewSocialNetworkRepository(t)

	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			"Registrate test",
			fields{},
			args{
				e.NewContext(
					httptest.NewRequest(
						http.MethodPost, "/",
						strings.NewReader(`{"email": "test@test.com", "password": "testpass"}`)),
					httptest.NewRecorder()),
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.Registrate(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.Registrate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_Authorize(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.Authorize(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.Authorize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_GetMyInfo(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.GetMyInfo(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.GetMyInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_GetPeople(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.GetPeople(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.GetPeople() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_FriendRequest(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.FriendRequest(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.FriendRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_FriendshipRequestAccept(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.FriendshipRequestAccept(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.FriendshipRequestAccept() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_FriendshipRequestDecline(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.FriendshipRequestDecline(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.FriendshipRequestDecline() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_GetFriendshipRequests(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.GetFriendshipRequests(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.GetFriendshipRequests() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_GetRelatedProfile(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.GetRelatedProfile(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.GetRelatedProfile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_GetProfilesBySearchPrefixes(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.GetProfilesBySearchPrefixes(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.GetProfilesBySearchPrefixes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_CreatePost(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.CreatePost(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.CreatePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_UpdatePost(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.UpdatePost(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.UpdatePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_DeletePost(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.DeletePost(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.DeletePost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_GetPost(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.GetPost(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.GetPost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_GetFeed(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.GetFeed(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.GetFeed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_RebuildFeeds(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.RebuildFeeds(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.RebuildFeeds() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_SendMessage(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.SendMessage(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.SendMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_GetDialog(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.GetDialog(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.GetDialog() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSocialNetworkHandler_GetDialogList(t *testing.T) {
	type fields struct {
		SNUsecase     domain.SocialNetworkUsecase
		FeederUsecase domain.FeederUsecase
		JWTSecret     string
	}
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SocialNetworkHandler{
				SNUsecase:     tt.fields.SNUsecase,
				FeederUsecase: tt.fields.FeederUsecase,
				JWTSecret:     tt.fields.JWTSecret,
			}
			if err := h.GetDialogList(tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("SocialNetworkHandler.GetDialogList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getUserId(t *testing.T) {
	type args struct {
		c echo.Context
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getUserId(tt.args.c); got != tt.want {
				t.Errorf("getUserId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getStatusCode(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getStatusCode(tt.args.err); got != tt.want {
				t.Errorf("getStatusCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
