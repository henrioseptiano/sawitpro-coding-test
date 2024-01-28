package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// RegisterTheUser
//
//	@Summary		RegisterTheUser
//	@Description	Register New User
//	@ID				RegisterTheUser
//	@Accept			application/json
//	@Produce		json
//	@Param			user	body	generated.RegisterTheUserJSONRequestBody	true	"Register user JSON Body"
//	@Success		200		{string}	string			"ok"
//	@Router			/register [post]
func (s *Server) RegisterTheUser(ctx echo.Context) error {
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	var regUser generated.RegisterTheUserJSONRequestBody
	json.Unmarshal(body, &regUser)

	if regUser.FullName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "FullName Cannot Be empty!")
	}

	if regUser.PhoneNumber == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Phone Number Cannot Be Empty")
	}

	if regUser.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Password Cannot Be empty")
	}

	if !utils.CheckPhoneNumber(regUser.PhoneNumber) {
		return echo.NewHTTPError(http.StatusBadRequest, "Phone Number Format is not Valid")
	}

	checkUser, err := s.Repository.CheckPhoneNumber(context.Background(), regUser.PhoneNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	if checkUser != 0 {
		return echo.NewHTTPError(http.StatusConflict, errors.New("Phone number already existed"))
	}

	userId := uuid.New()
	hashedPassword, err := utils.HashPassword(regUser.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	user := repository.User{
		UserID:                   userId.String(),
		FullName:                 regUser.FullName,
		PhoneNumber:              regUser.PhoneNumber,
		Password:                 hashedPassword,
		SuccessfullLoginAttempts: 0,
		CreatedAt:                time.Now(),
		UpdatedAt:                time.Now(),
	}

	err = s.Repository.RegisterUser(user)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response := map[string]string{
		"message": "Successfully Registered!",
	}
	return ctx.JSON(http.StatusCreated, response)
}

// LoginUser
//
//	@Summary		LoginUser
//	@Description	Login Existing User
//	@ID				LoginUser
//	@Accept			application/json
//	@Produce		json
//	@Param			user	body	generated.LoginUserJSONRequestBody	true	"Login User JSON Body"
//	@Success		200		{string}	string			"ok"
//	@Router			/login [post]
func (s *Server) LoginUser(ctx echo.Context) error {
	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return err
	}

	var loginUser generated.LoginUserJSONRequestBody
	json.Unmarshal(body, &loginUser)

	if loginUser.PhoneNumber == "" || loginUser.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Please input your phone number and password")
	}

	getUser, err := s.Repository.CheckUser(context.Background(), loginUser.PhoneNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := utils.CheckPassword(loginUser.Password, getUser.Password); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid Password")
	}

	currentTime := time.Now()
	getUser.SuccessfullLoginAttempts += 1
	getUser.LastLogin = &currentTime

	if err := s.Repository.UpdateLoginUser(context.Background(), *getUser); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	accessToken, refreshToken, err := utils.GenerateJWTToken(getUser.UserID, os.Getenv("JWT_SECRET"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	response := map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}

	return ctx.JSON(http.StatusOK, response)
}

// GetProfile
//
//	@Summary		GetProfile
//	@Description	Get User Profile
//	@ID				GetProfile
//	@Accept			application/json
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Success		200		{string}	string			"ok"
//	@Router			/user [get]
func (s *Server) GetProfile(ctx echo.Context) error {
	getAuth := ctx.Request().Header.Get("Authorization")
	tokenString := utils.GetTokenFromAuthHeader(getAuth)
	res, err := utils.DecodeJWTToken(tokenString)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}
	getUser, err := s.Repository.GetUserByUserId(context.Background(), *res)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	dataResponse := map[string]map[string]string{
		"data": {
			"full_name":    getUser.FullName,
			"phone_number": getUser.PhoneNumber,
		},
	}
	return ctx.JSON(http.StatusOK, dataResponse)
}

// UpdateProfile
//
//	@Summary		UpdateProfile
//	@Description	Update User Profile
//	@ID				UpdateProfile
//	@Accept			application/json
//	@Security		ApiKeyAuth
//	@Produce		json
//	@Param			user	body	generated.UpdateProfileJSONRequestBody	true	"Update User Profile JSON Body"
//	@Success		200		{string}	string			"ok"
//	@Router			/user [put]
func (s *Server) UpdateProfile(ctx echo.Context) error {
	getAuth := ctx.Request().Header.Get("Authorization")
	tokenString := utils.GetTokenFromAuthHeader(getAuth)
	res, err := utils.DecodeJWTToken(tokenString)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	body, err := io.ReadAll(ctx.Request().Body)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var updateProfile generated.UpdateProfileJSONRequestBody
	json.Unmarshal(body, &updateProfile)

	if updateProfile.FullName == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "FullName Cannot Be empty!")
	}

	if updateProfile.PhoneNumber == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Phone Number Cannot Be Empty")
	}

	if !utils.CheckPhoneNumber(updateProfile.PhoneNumber) {
		return echo.NewHTTPError(http.StatusBadRequest, "Phone Number Format is not Valid")
	}

	getUser, err := s.Repository.GetUserByUserId(context.Background(), *res)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	checkUser, err := s.Repository.CheckPhoneNumber(context.Background(), updateProfile.PhoneNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	if checkUser != 0 && strings.TrimSpace(updateProfile.PhoneNumber) != strings.TrimSpace(getUser.PhoneNumber) {
		return echo.NewHTTPError(http.StatusConflict, errors.New("Phone number already existed"))
	}
	getUser.PhoneNumber = updateProfile.PhoneNumber
	getUser.FullName = updateProfile.FullName
	err = s.Repository.UpdateUserProfile(context.Background(), *getUser)
	if err != nil {
		return echo.NewHTTPError(http.StatusForbidden, err.Error())
	}

	response := map[string]string{
		"message": "User Profile Successfully Updated!",
	}
	return ctx.JSON(http.StatusAccepted, response)
}
