package v1

import (
	"auth/internal/contracts/open_api/v1/dto"
	"auth/libs/context"
	"fmt"
	"github.com/gin-gonic/gin"
)

type AuthContract struct {
	auth authService
}

func NewAuthContract(auth authService) *AuthContract {
	return &AuthContract{auth: auth}
}

func (ac *AuthContract) SignIn(ginC *gin.Context) {
	var (
		httpErr  *HttpError
		response = &dto.SignInResponse{}
	)
	defer func() {
		processError(ginC, httpErr)
		if httpErr == nil {
			responseOk(ginC, response)
		}
	}()

	req := dto.SignInRequest{}

	err := ginC.ShouldBind(&req)
	if err != nil {
		httpErr = NewHttpError(400, "request data", nil, err)
		return
	}

	accessToken, refreshToken, user, err := ac.auth.SignIn(context.FromGinC(ginC), req.Username, req.Password)
	if err != nil {
		httpErr = NewHttpError(400, "sign in", nil, err)
		return
	}
	if user == nil {
		httpErr = NewHttpError(400, "sign in", nil, fmt.Errorf("user not found"))
		return
	}

	response = dto.BuildSignInResponse(accessToken, refreshToken, *user)

	return
}

func (ac *AuthContract) SignUp(ginC *gin.Context) {
	var (
		httpErr  *HttpError
		response = &dto.SignUpResponse{}
	)
	defer func() {
		processError(ginC, httpErr)
		if httpErr == nil {
			responseOk(ginC, response)
		}
	}()

	req := dto.SignUpRequest{}

	err := ginC.ShouldBind(&req)
	if err != nil {
		httpErr = NewHttpError(400, "request data", nil, err)
		return
	}

	err = ac.auth.SignUp(context.FromGinC(ginC), req.Username, req.Password)
	if err != nil {
		httpErr = NewHttpError(400, "sign up", nil, err)
		return
	}

	return
}
