package controller

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Demacr/otus-hl-socialnetwork/internal/domain"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

type ResponseError struct {
	Message string `json:"message"`
}

type SocialNetworkHandler struct {
	SNUsecase     domain.SocialNetworkUsecase
	FeederUsecase domain.FeederUsecase
	JWTSecret     string
}

func NewSocialNetworkHandler(e *echo.Echo, snuc domain.SocialNetworkUsecase, feeduc domain.FeederUsecase, JWTSecret string) {
	handler := &SocialNetworkHandler{
		SNUsecase:     snuc,
		FeederUsecase: feeduc,
		JWTSecret:     JWTSecret,
	}

	e.POST("/api/registrate", handler.Registrate)
	e.POST("/api/authorize", handler.Authorize)

	r := e.Group("/api/account", middleware.JWT([]byte(handler.JWTSecret)))
	d := e.Group("/api/dialog", middleware.JWT([]byte(handler.JWTSecret)))
	r.GET("/myinfo", handler.GetMyInfo)
	r.GET("/getpeople", handler.GetPeople)
	r.POST("/friend_request", handler.FriendRequest)
	r.POST("/friendship_request_accept", handler.FriendshipRequestAccept)
	r.POST("/friendship_request_decline", handler.FriendshipRequestDecline)
	r.GET("/my_friend_requests", handler.GetFriendshipRequests)
	r.GET("/profile/:id", handler.GetRelatedProfile)
	r.GET("/search", handler.GetProfilesBySearchPrefixes)

	r.POST("/post", handler.CreatePost)
	r.PUT("/post", handler.UpdatePost)
	r.DELETE("/post", handler.DeletePost)
	r.GET("/post/:id", handler.GetPost)
	r.GET("/post/feed", handler.GetFeed)

	d.POST("/:id/send", handler.SendMessage)
	d.GET("/:id/list", handler.GetDialog)
	d.GET("/list", handler.GetDialogList)

	admin := e.Group("/api/admin", middleware.BasicAuthWithConfig(middleware.BasicAuthConfig{
		Validator: func(s1, s2 string, ctx echo.Context) (bool, error) {
			return s1 == "admin" && s2 == "adminpassword", nil
		},
	}))
	admin.GET("/rebuild_feeds", handler.RebuildFeeds)
}

func (h *SocialNetworkHandler) Registrate(c echo.Context) error {
	profile := &domain.Profile{}
	if err := c.Bind(profile); err != nil {
		return c.JSON(http.StatusBadRequest, ResponseError{Message: domain.ErrBadRequest.Error()})
	}

	if err := h.SNUsecase.Registrate(profile); err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	return c.NoContent(http.StatusCreated)
}

func (h *SocialNetworkHandler) Authorize(c echo.Context) error {
	credentials := &domain.Credentials{}
	if err := c.Bind(credentials); err != nil {
		log.Println(err)
		return c.String(http.StatusBadRequest, "Bad json.")
	}

	profile, err := h.SNUsecase.Authorize(credentials)
	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	claims["id"] = profile.ID
	claims["email"] = credentials.Email

	t, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		return errors.Wrap(err, "error during signing token")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}

func (h *SocialNetworkHandler) GetMyInfo(c echo.Context) error {
	// TODO: rework this flow to use id from JWT.
	user := c.Get("user").(*jwt.Token)
	email := user.Claims.(jwt.MapClaims)["email"].(string)

	profile, err := h.SNUsecase.GetProfileInfo(email)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Get profile error.")
	}

	return c.JSON(http.StatusOK, profile)
}

func (h *SocialNetworkHandler) GetPeople(c echo.Context) error {
	user_id := getUserId(c)

	result, err := h.SNUsecase.GetRandomProfiles(user_id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Internal error")
	}

	return c.JSON(http.StatusOK, result)
}

func (h *SocialNetworkHandler) FriendRequest(c echo.Context) error {
	friend_request := &domain.FriendRequest{}
	if err := c.Bind(friend_request); err != nil {
		log.Println(err)
		return c.String(http.StatusBadRequest, "Bad json.")
	}

	id := getUserId(c)

	err := h.SNUsecase.CreateFriendRequest(int(id), friend_request.FriendID)
	if err != nil {
		log.Println(err)
		return c.String(getStatusCode(err), "Cannot create frienship request")
	}

	return c.String(http.StatusOK, "")
}

func (h *SocialNetworkHandler) FriendshipRequestAccept(c echo.Context) error {
	friend_request := &domain.FriendRequest{}
	if err := c.Bind(friend_request); err != nil {
		log.Println(err)
		return c.String(http.StatusBadRequest, "Bad json.")
	}

	id := getUserId(c)

	err := h.SNUsecase.FriendshipRequestAccept(int(id), friend_request.FriendID)
	if err != nil {
		log.Println(err)
		return c.String(getStatusCode(err), "Cannot accept friendship request")
	}

	return c.String(http.StatusOK, "")
}

func (h *SocialNetworkHandler) FriendshipRequestDecline(c echo.Context) error {
	friend_request := &domain.FriendRequest{}
	if err := c.Bind(friend_request); err != nil {
		log.Println(err)
		return c.String(http.StatusBadRequest, "Bad json.")
	}

	id := getUserId(c)

	err := h.SNUsecase.FriendshipRequestDecline(int(id), friend_request.FriendID)
	if err != nil {
		log.Println(err)
		return c.String(getStatusCode(err), "Cannot accept friendship request")
	}

	return c.String(http.StatusOK, "")
}

func (h *SocialNetworkHandler) GetFriendshipRequests(c echo.Context) error {
	id := getUserId(c)

	friend_requests, err := h.SNUsecase.GetFriendshipRequests(int(id))
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Can't get friendship requests")
	}

	if len(friend_requests) > 0 {
		return c.JSON(http.StatusOK, friend_requests)
	}
	return c.JSON(http.StatusNoContent, nil)
}

func (h *SocialNetworkHandler) GetRelatedProfile(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad profile id")
	}

	related_id := c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"].(float64)

	result, err := h.SNUsecase.GetRelatedProfile(id, int(related_id))
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Can't get profile info")
	}

	return c.JSON(http.StatusOK, result)
}

func (h *SocialNetworkHandler) GetProfilesBySearchPrefixes(c echo.Context) error {
	fn := ""
	ln := ""
	err := echo.QueryParamsBinder(c).
		String("firstName", &fn).
		String("lastName", &ln).
		BindError()
	if err != nil {
		log.Println(errors.Wrap(err, "searching profiles error"))
		return c.String(http.StatusBadRequest, "Bad request")
	}

	result, err := h.SNUsecase.GetProfilesBySearchPrefixes(fn, ln)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Can't get searched profiles")
	}

	if len(result) == 0 {
		return c.JSON(http.StatusNoContent, nil)
	}

	return c.JSON(http.StatusOK, result)
}

func (h *SocialNetworkHandler) CreatePost(c echo.Context) error {
	id := getUserId(c)

	post := &domain.Post{}
	if err := c.Bind(post); err != nil {
		log.Println(err)
		return c.String(http.StatusBadRequest, "Bad json.")
	}

	err := h.SNUsecase.CreatePost(int(id), post)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "")
	}

	return c.NoContent(http.StatusCreated)
}

func (h *SocialNetworkHandler) UpdatePost(c echo.Context) error {
	id := getUserId(c)

	post := &domain.Post{}
	if err := c.Bind(post); err != nil {
		log.Println(err)
		return c.String(http.StatusBadRequest, "Bad json.")
	}

	err := h.SNUsecase.UpdatePost(int(id), post)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "")
	}

	return c.NoContent(http.StatusOK)
}

func (h *SocialNetworkHandler) DeletePost(c echo.Context) error {
	id := getUserId(c)

	post := &domain.Post{}
	if err := c.Bind(post); err != nil {
		log.Println(err)
		return c.String(http.StatusBadRequest, "Bad json.")
	}

	err := h.SNUsecase.DeletePost(int(id), post)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "")
	}

	return c.NoContent(http.StatusOK)
}

func (h *SocialNetworkHandler) GetPost(c echo.Context) error {
	post_id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Bad post id")
	}

	post, err := h.SNUsecase.GetPost(post_id)
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "")
	}

	return c.JSON(http.StatusOK, post)
}

func (h *SocialNetworkHandler) GetFeed(c echo.Context) error {
	id := getUserId(c)

	feed, err := h.FeederUsecase.GetFeedIds(int(id))
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, feed)
}

func (h *SocialNetworkHandler) RebuildFeeds(c echo.Context) error {
	err := h.FeederUsecase.RebuildFeeds()
	if err != nil {
		log.Println(err)
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusOK)
}

func (h *SocialNetworkHandler) SendMessage(c echo.Context) error {
	id := getUserId(c)

	friendId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(errors.Wrap(err, "Controller.SendMessage.Atoi"))
		return c.NoContent(http.StatusBadRequest)
	}

	message := &domain.Message{}
	if err := c.Bind(message); err != nil {
		log.Println(errors.Wrap(err, "Controller.SendMessage.Bind"))
		return c.String(http.StatusBadRequest, "Bad json.")
	}
	message.From = int(id)
	message.To = friendId
	message.Timestamp = time.Now()

	if err = h.SNUsecase.SendMessage(message); err != nil {
		log.Println(errors.Wrap(err, "Controller.SendMessage.Usecase.SendMessage"))
		return c.NoContent(http.StatusInternalServerError)
	}

	return nil
}

func (h *SocialNetworkHandler) GetDialog(c echo.Context) error {
	id := getUserId(c)

	friendId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		log.Println(errors.Wrap(err, "Controller.GetDialog.Atoi"))
		return c.NoContent(http.StatusBadRequest)
	}

	messages, err := h.SNUsecase.GetDialog(int(id), friendId)
	if err != nil {
		log.Println(errors.Wrap(err, "Controller.GetDialog.Usecase.GetDialog"))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, messages)
}

func (h *SocialNetworkHandler) GetDialogList(c echo.Context) error {
	id := getUserId(c)

	previewList, err := h.SNUsecase.GetDialogList(id)
	if err != nil {
		log.Println(errors.Wrap(err, "Controller.GetDialogList.Usecase.GetDialogList"))
		return c.NoContent(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, previewList)
}

func getUserId(c echo.Context) int {
	return int(c.Get("user").(*jwt.Token).Claims.(jwt.MapClaims)["id"].(float64))
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	log.Println(err)
	switch err {
	case domain.ErrWrongCredentials:
		return http.StatusInternalServerError
	case domain.ErrFriendshipRequestExists:
		return http.StatusNotAcceptable
	case domain.ErrFriendshipRequestNotExists:
		return http.StatusNotAcceptable
	default:
		return http.StatusInternalServerError
	}
}
