package user

import (
	"github.com/GoSimplicity/LinkMe/internal/app/user/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

func (h *UserHandler) Register(c *gin.Context) {}

func (h *UserHandler) Login(c *gin.Context) {}

func (h *UserHandler) LoginSMS(c *gin.Context) {}

func (h *UserHandler) SendSMS(c *gin.Context) {}

func (h *UserHandler) SendEmail(c *gin.Context) {}

func (h *UserHandler) RefreshToken(c *gin.Context) {}

func (h *UserHandler) Logout(c *gin.Context) {}

func (h *UserHandler) GetProfile(c *gin.Context) {}

func (h *UserHandler) UpdateProfile(c *gin.Context) {}

func (h *UserHandler) ChangePassword(c *gin.Context) {}

func (h *UserHandler) DeleteUser(c *gin.Context) {}

func (h *UserHandler) List(c *gin.Context) {}

func (h *UserHandler) GetUserById(c *gin.Context) {}

func (h *UserHandler) UpdateUserById(c *gin.Context) {}

func (h *UserHandler) DeleteUserById(c *gin.Context) {}
