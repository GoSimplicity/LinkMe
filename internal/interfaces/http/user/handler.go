package user

import (
	"github.com/GoSimplicity/LinkMe/internal/app/user/service"
	"github.com/GoSimplicity/LinkMe/internal/interfaces/http/user/dto"
	"github.com/GoSimplicity/LinkMe/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

func (h *UserHandler) Register(ctx *gin.Context) {
	var req dto.SignUpReq

	utils.HandleRequest(ctx, req, func() (interface{}, error) {
		return nil, h.userSvc.Create(ctx, req)
	})
}

func (h *UserHandler) Login(ctx *gin.Context) {
	var req dto.LoginReq

	utils.HandleRequest(ctx, req, func() (interface{}, error) {
		return h.userSvc.Login(ctx, req)
	})
}

func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	var req dto.LoginSMSReq

	utils.HandleRequest(ctx, req, func() (interface{}, error) {
		return h.userSvc.LoginSMS(ctx, req)
	})
}

func (h *UserHandler) SendSMS(ctx *gin.Context) {
	var req dto.SendSMSReq

	utils.HandleRequest(ctx, req, func() (interface{}, error) {
		return nil, h.userSvc.SendSMS(ctx, req)
	})
}

func (h *UserHandler) SendEmail(ctx *gin.Context) {
	var req dto.SendEmailReq

	utils.HandleRequest(ctx, req, func() (interface{}, error) {
		return nil, h.userSvc.SendEmail(ctx, req)
	})
}

func (h *UserHandler) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshTokenReq

	utils.HandleRequest(ctx, req, func() (interface{}, error) {
		return h.userSvc.RefreshToken(ctx, req)
	})
}

func (h *UserHandler) Logout(ctx *gin.Context) {
	var req dto.LogoutReq

	utils.HandleRequest(ctx, req, func() (interface{}, error) {
		return nil, h.userSvc.Logout(ctx, req)
	})
}

func (h *UserHandler) GetProfile(ctx *gin.Context) {
	var req dto.GetProfileReq

	utils.HandleRequest(ctx, req, func() (interface{}, error) {
		return h.userSvc.GetProfile(ctx, req)
	})
}

func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	var req dto.UpdateProfileReq

	utils.HandleRequest(ctx, req, func() (interface{}, error) {
		return nil, h.userSvc.UpdateProfile(ctx, req)
	})
}

func (h *UserHandler) ChangePassword(ctx *gin.Context) {
	var req dto.ChangePasswordReq

	utils.HandleRequest(ctx, req, func() (interface{}, error) {
		return nil, h.userSvc.ChangePassword(ctx, req)
	})
}

func (h *UserHandler) DeleteUser(ctx *gin.Context) {
	var req dto.DeleteUserReq

	utils.HandleRequest(ctx, req, func() (interface{}, error) {
		return nil, h.userSvc.DeleteUser(ctx, req)
	})
}

func (h *UserHandler) List(ctx *gin.Context) {
	var req dto.ListUserReq

	utils.HandleRequest(ctx, req, func() (interface{}, error) {
		return h.userSvc.List(ctx, req, req.Page, req.Size)
	})
}

func (h *UserHandler) GetUser(ctx *gin.Context) {
	var req dto.GetUserReq

	utils.HandleRequest(ctx, req, func() (interface{}, error) {
		return h.userSvc.GetUser(ctx, req)
	})
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	var req dto.UpdateUserReq

	utils.HandleRequest(ctx, req, func() (interface{}, error) {
		return nil, h.userSvc.UpdateUser(ctx, req)
	})
}
