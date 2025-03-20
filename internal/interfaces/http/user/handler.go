// interfaces/http/user/handler.go
package user

import (
	"github.com/GoSimplicity/LinkMe/internal/app/user/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	userSvc service.UserService
}

func NewHandler(userSvc service.UserService) *Handler {
	return &Handler{userSvc: userSvc}
}

func (h *Handler) Register(c *gin.Context) {}

func (h *Handler) Login(c *gin.Context) {}

func (h *Handler) LoginSMS(c *gin.Context) {}

func (h *Handler) SendSMS(c *gin.Context) {}

func (h *Handler) SendEmail(c *gin.Context) {}

func (h *Handler) RefreshToken(c *gin.Context) {}

func (h *Handler) Logout(c *gin.Context) {}

func (h *Handler) GetProfile(c *gin.Context) {}

func (h *Handler) UpdateProfile(c *gin.Context) {}

func (h *Handler) ChangePassword(c *gin.Context) {}

func (h *Handler) DeleteUser(c *gin.Context) {}

func (h *Handler) List(c *gin.Context) {}

func (h *Handler) GetUserById(c *gin.Context) {}

func (h *Handler) UpdateUserById(c *gin.Context) {}

func (h *Handler) DeleteUserById(c *gin.Context) {}
