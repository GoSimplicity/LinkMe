package api

//
//import (
//	. "LinkMe/internal/constants"
//	"LinkMe/internal/domain"
//	"LinkMe/internal/service"
//	svcmocks "LinkMe/internal/service/mocks"
//	"LinkMe/pkg/ginp"
//	"bytes"
//	"encoding/json"
//	"github.com/gin-gonic/gin"
//	"github.com/stretchr/testify/assert"
//	"go.uber.org/mock/gomock"
//	"net/http"
//	"net/http/httptest"
//	"testing"
//)
//
//func TestUserHandler_SignUp(t *testing.T) {
//	testCases := []struct {
//		name       string
//		mock       func(ctrl *gomock.Controller) service.UserService
//		reqBuilder func(t *testing.T) *http.Request
//		wantCode   int
//		wantBody   ginp.Result
//	}{
//		{
//			name:     "注册成功",
//			wantCode: 200,
//			wantBody: ginp.Result{
//				Code: http.StatusOK,
//				Msg:  UserSignUpSuccess,
//			},
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				userSvc := svcmocks.NewMockUserService(ctrl)
//				userSvc.EXPECT().
//					SignUp(gomock.Any(), domain.User{
//						Email:    "1234@qq.com",
//						Password: "Abc@1234",
//					}).Return(nil)
//				return userSvc
//			},
//			reqBuilder: func(t *testing.T) *http.Request {
//				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
//					"email": "1234@qq.com",
//					"password": "Abc@1234",
//					"confirmPassword": "Abc@1234"
//				}`)))
//				req.Header.Set("Content-Type", "application/json")
//				assert.NoError(t, err)
//				return req
//			},
//		},
//		{
//			name:     "无效的请求负载",
//			wantCode: http.StatusBadRequest,
//			wantBody: ginp.Result{
//				Msg: "无效的请求负载",
//			},
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				userSvc := svcmocks.NewMockUserService(ctrl)
//				return userSvc
//			},
//			reqBuilder: func(t *testing.T) *http.Request {
//				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
//				"email": "1234@qq.com",
//				"password": "Abc@1234",
//			}`)))
//				req.Header.Set("Content-Type", "application/json")
//				assert.NoError(t, err)
//				return req
//			},
//		},
//		{
//			name:     "非法邮箱格式",
//			wantCode: http.StatusOK,
//			wantBody: ginp.Result{
//				Code: UserInvalidInput,
//				Msg:  UserEmailFormatError,
//			},
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				userSvc := svcmocks.NewMockUserService(ctrl)
//				return userSvc
//			},
//			reqBuilder: func(t *testing.T) *http.Request {
//				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
//				"email": "123",
//				"password": "Abc@1234",
//				"confirmPassword": "Abc@1234"
//			}`)))
//				req.Header.Set("Content-Type", "application/json")
//				assert.NoError(t, err)
//				return req
//			},
//		},
//		{
//			name:     "两次输入密码不对",
//			wantCode: http.StatusOK,
//			wantBody: ginp.Result{
//				Code: UserInvalidInput,
//				Msg:  UserPasswordMismatchError,
//			},
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				userSvc := svcmocks.NewMockUserService(ctrl)
//				return userSvc
//			},
//			reqBuilder: func(t *testing.T) *http.Request {
//				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
//				"email": "1234@qq.com",
//				"password": "Abc@1234",
//				"confirmPassword": "Abc@12345"
//			}`)))
//				req.Header.Set("Content-Type", "application/json")
//				assert.NoError(t, err)
//				return req
//			},
//		},
//		{
//			name:     "密码必须包含字母、数字、特殊字符，并且不少于八位",
//			wantCode: http.StatusOK,
//			wantBody: ginp.Result{
//				Code: UserInvalidInput,
//				Msg:  UserPasswordFormatError,
//			},
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				userSvc := svcmocks.NewMockUserService(ctrl)
//				return userSvc
//			},
//			reqBuilder: func(t *testing.T) *http.Request {
//				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
//				"email": "1234@qq.com",
//				"password": "123",
//				"confirmPassword": "123"
//			}`)))
//				req.Header.Set("Content-Type", "application/json")
//				assert.NoError(t, err)
//				return req
//			},
//		},
//		{
//			name:     "邮箱冲突",
//			wantCode: 200,
//			wantBody: ginp.Result{
//				Code: UserDuplicateEmail,
//				Msg:  UserEmailConflictError,
//			},
//			mock: func(ctrl *gomock.Controller) service.UserService {
//				userSvc := svcmocks.NewMockUserService(ctrl)
//				userSvc.EXPECT().SignUp(gomock.Any(), domain.User{
//					Email:    "1234@qq.com",
//					Password: "Abc@1234",
//				}).Return(service.ErrDuplicateEmail)
//				return userSvc
//			},
//			reqBuilder: func(t *testing.T) *http.Request {
//				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
//				"email": "1234@qq.com",
//				"password": "Abc@1234",
//				"confirmPassword": "Abc@1234"
//			}`)))
//				req.Header.Set("Content-Type", "application/json")
//				assert.NoError(t, err)
//				return req
//			},
//		},
//	}
//	// 遍历所有测试场景并执行测试
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			ctrl := gomock.NewController(t)
//			defer ctrl.Finish()
//			userSvc := tc.mock(ctrl)
//			hdl := NewUserHandler(userSvc, nil, nil)
//			server := gin.Default()
//			hdl.RegisterRoutes(server)
//			req := tc.reqBuilder(t)
//			// 创建httptest.Recorder以捕获服务器响应
//			recorder := httptest.NewRecorder()
//			// 发送请求到服务器并接收响应
//			server.ServeHTTP(recorder, req)
//			// 对响应进行断言，验证状态码和响应体是否符合预期
//			assert.Equal(t, tc.wantCode, recorder.Code)
//			wantBody, err := json.Marshal(tc.wantBody)
//			assert.NoError(t, err)
//			assert.Equal(t, string(wantBody), recorder.Body.String())
//		})
//	}
//}
