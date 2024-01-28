package handler

import (
	"bytes"
	"encoding/json"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/SawitProRecruitment/UserService/utils"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestRegisterTheUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := repository.NewMockRepositoryInterface(ctrl)

	server := &Server{
		Repository: mockRepository,
	}

	tests := []struct {
		name                string
		requestBody         interface{}
		mockOutput          int
		expectedCode        int
		expectedError       bool
		isInputValidated    bool
		isPhoneNumberUnique bool
	}{
		{
			name: "Successful Registration",
			requestBody: map[string]interface{}{
				"full_name":    "John Doe",
				"phone_number": "+62234567890",
				"Password":     "password",
			},
			mockOutput:          0,
			expectedCode:        http.StatusCreated,
			expectedError:       false,
			isInputValidated:    true,
			isPhoneNumberUnique: true,
		},
		{
			name: "Empty Fullname",
			requestBody: map[string]interface{}{
				"full_name":    "",
				"phone_number": "+62234567890",
				"Password":     "password",
			},
			mockOutput:          0,
			expectedCode:        http.StatusCreated,
			expectedError:       true,
			isInputValidated:    false,
			isPhoneNumberUnique: false,
		},
		{
			name: "Empty Phone Number",
			requestBody: map[string]interface{}{
				"full_name":    "John Doe",
				"phone_number": "",
				"Password":     "password",
			},
			mockOutput:          0,
			expectedCode:        http.StatusCreated,
			expectedError:       true,
			isInputValidated:    false,
			isPhoneNumberUnique: false,
		},
		{
			name: "Empty Password",
			requestBody: map[string]interface{}{
				"full_name":    "John Doe",
				"phone_number": "+62234567890",
				"Password":     "",
			},
			mockOutput:          0,
			expectedCode:        http.StatusCreated,
			expectedError:       true,
			isInputValidated:    false,
			isPhoneNumberUnique: false,
		},
		{
			name: "Invalid Phone Number",
			requestBody: map[string]interface{}{
				"full_name":    "John Doe",
				"phone_number": "2234567890",
				"Password":     "password",
			},
			mockOutput:          0,
			expectedCode:        http.StatusCreated,
			expectedError:       true,
			isInputValidated:    false,
			isPhoneNumberUnique: false,
		},
		{
			name: "Phone Number Existed",
			requestBody: map[string]interface{}{
				"full_name":    "John Doe",
				"phone_number": "+62234567890",
				"Password":     "password",
			},
			mockOutput:          1,
			expectedCode:        http.StatusCreated,
			expectedError:       true,
			isInputValidated:    true,
			isPhoneNumberUnique: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isInputValidated {
				mockRepository.EXPECT().CheckPhoneNumber(gomock.Any(), gomock.Any()).Return(int64(tc.mockOutput), nil)
				if tc.isPhoneNumberUnique {
					mockRepository.EXPECT().RegisterUser(gomock.Any()).Return(nil)
				}
			}
			reqBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			// Create a new Echo instance and handle the request
			e := echo.New()
			c := e.NewContext(req, rec)

			err := server.RegisterTheUser(c)

			if tc.expectedError {
				assert.Error(t, err)
				assert.NotEqual(t, tc.expectedCode, rec.Code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCode, rec.Code)
			}
		})
	}
}

func TestLoginUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := repository.NewMockRepositoryInterface(ctrl)

	server := &Server{
		Repository: mockRepository,
	}
	hashedPassword, _ := utils.HashPassword("password")
	tests := []struct {
		name            string
		requestBody     interface{}
		mockOutput      *repository.User
		expectedCode    int
		expectedError   bool
		isInputValidate bool
		isProceedLogin  bool
	}{
		{
			name: "Successful Login",
			requestBody: map[string]interface{}{
				"phone_number": "+62234567890",
				"Password":     "password",
			},
			mockOutput: &repository.User{
				UserID:                   "mockUserID",
				PhoneNumber:              "+62234567890",
				Password:                 hashedPassword, // Assuming this is the hashed password stored in the repository
				SuccessfullLoginAttempts: 0,
				CreatedAt:                time.Now(),
				UpdatedAt:                time.Now(),
			},
			expectedCode:    http.StatusOK,
			expectedError:   false,
			isInputValidate: true,
			isProceedLogin:  true,
		},
		{
			name: "Phone Number Empty",
			requestBody: map[string]interface{}{
				"phone_number": "",
				"Password":     "password",
			},
			mockOutput: &repository.User{
				UserID:                   "mockUserID",
				PhoneNumber:              "+62234567890",
				Password:                 hashedPassword, // Assuming this is the hashed password stored in the repository
				SuccessfullLoginAttempts: 0,
				CreatedAt:                time.Now(),
				UpdatedAt:                time.Now(),
			},
			expectedCode:    http.StatusOK,
			expectedError:   true,
			isInputValidate: false,
			isProceedLogin:  false,
		},
		{
			name: "Password Empty",
			requestBody: map[string]interface{}{
				"phone_number": "+62234567890",
				"Password":     "",
			},
			mockOutput: &repository.User{
				UserID:                   "mockUserID",
				PhoneNumber:              "+62234567890",
				Password:                 hashedPassword, // Assuming this is the hashed password stored in the repository
				SuccessfullLoginAttempts: 0,
				CreatedAt:                time.Now(),
				UpdatedAt:                time.Now(),
			},
			expectedCode:    http.StatusOK,
			expectedError:   true,
			isInputValidate: false,
			isProceedLogin:  false,
		},
		{
			name: "Invalid Password",
			requestBody: map[string]interface{}{
				"phone_number": "+62234567890",
				"Password":     "asewrewr",
			},
			mockOutput: &repository.User{
				UserID:                   "mockUserID",
				PhoneNumber:              "+62234567890",
				Password:                 hashedPassword, // Assuming this is the hashed password stored in the repository
				SuccessfullLoginAttempts: 0,
				CreatedAt:                time.Now(),
				UpdatedAt:                time.Now(),
			},
			expectedCode:    http.StatusOK,
			expectedError:   true,
			isInputValidate: true,
			isProceedLogin:  false,
		},
		// Add more test cases as needed
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isInputValidate {
				mockRepository.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(tc.mockOutput, nil)
				if tc.isProceedLogin {
					mockRepository.EXPECT().UpdateLoginUser(gomock.Any(), gomock.Any()).Return(nil)
				}
			}
			reqBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Create a new Echo instance and handle the request
			e := echo.New()
			c := e.NewContext(req, rec)

			err := server.LoginUser(c)

			if tc.expectedError {
				assert.Error(t, err)
				//assert.NotEqual(t, tc.expectedCode, rec.Code)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCode, rec.Code)
			}
		})
	}
}

func TestGetProfile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepository := repository.NewMockRepositoryInterface(ctrl)
	os.Setenv("JWT_SECRET", "verysecret")
	server := &Server{
		Repository: mockRepository,
	}
	jwtToken, _, _ := utils.GenerateJWTToken("mockUserID", "verysecret")
	tests := []struct {
		name           string
		authHeader     string
		mockInput      string
		mockOutput     *repository.User
		expectedCode   int
		expectedResult map[string]map[string]string
		expectedError  string
	}{
		{
			name:       "Successful Profile Retrieval",
			authHeader: "Bearer " + jwtToken,
			mockInput:  "mockUserID",
			mockOutput: &repository.User{
				UserID:      "mockUserID",
				FullName:    "John Doe",
				PhoneNumber: "1234567890",
			},
			expectedCode: http.StatusOK,
			expectedResult: map[string]map[string]string{
				"data": {
					"full_name":    "John Doe",
					"phone_number": "1234567890",
				},
			},
			expectedError: "",
		},
		// Add more test cases as needed
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepository.EXPECT().GetUserByUserId(gomock.Any(), tc.mockInput).Return(tc.mockOutput, nil)

			req := httptest.NewRequest(http.MethodGet, "/user", nil)
			req.Header.Set("Authorization", tc.authHeader)
			rec := httptest.NewRecorder()

			// Create a new Echo instance and handle the request
			e := echo.New()
			c := e.NewContext(req, rec)

			err := server.GetProfile(c)

			if tc.expectedError != "" {
				assert.Error(t, err)
				//assert.Contains(t, err.Error(), tc.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCode, rec.Code)
				//assert.JSONEq(t, utils.ToJSON(tc.expectedResult), rec.Body.String())
			}
		})
	}
}

func TestUpdateProfile(t *testing.T) {
	// Set up mock repository and server
	// Initialize your server and mock repository
	os.Setenv("JWT_SECRET", "verysecret")
	jwtToken, _, _ := utils.GenerateJWTToken("mockUserID", "verysecret")

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockRepository := repository.NewMockRepositoryInterface(ctrl)
	server := &Server{
		Repository: mockRepository,
	}
	tests := []struct {
		name                 string
		requestBody          map[string]string
		userID               string // User ID extracted from JWT token
		existingPhone        string // Existing phone number in the database
		mockOutput           *repository.User
		mockCheckPhoneNumber int // Return value of CheckPhoneNumber
		expectedCode         int
		expectedMsg          string
		expectedError        bool
		isInputValidate      bool
		isProceedUpdate      bool
	}{
		{
			name: "Successful Profile Update",
			requestBody: map[string]string{
				"full_name":    "John Doe",
				"phone_number": "+621234567890",
			},
			userID:        "mockUserID",
			existingPhone: "+621234567890",
			mockOutput: &repository.User{
				UserID:      "mockUserID",
				FullName:    "Old Name",
				PhoneNumber: "+621234567890",
			},
			mockCheckPhoneNumber: 0,
			expectedCode:         http.StatusAccepted,
			expectedMsg:          "User Profile Successfully Updated!",
			expectedError:        false,
			isInputValidate:      true,
			isProceedUpdate:      true,
		},
		{
			name: "Empty Full Name",
			requestBody: map[string]string{
				"full_name":    "",
				"phone_number": "+621234567890",
			},
			userID:        "mockUserID",
			existingPhone: "+621234567890",
			mockOutput: &repository.User{
				UserID:      "mockUserID",
				FullName:    "Old Name",
				PhoneNumber: "+621234567890",
			},
			mockCheckPhoneNumber: 0,
			expectedCode:         http.StatusAccepted,
			expectedMsg:          "User Profile Successfully Updated!",
			expectedError:        true,
			isInputValidate:      false,
			isProceedUpdate:      false,
		},
		{
			name: "Empty Phone Number",
			requestBody: map[string]string{
				"full_name":    "John Doe",
				"phone_number": "",
			},
			userID:        "mockUserID",
			existingPhone: "+621234567890",
			mockOutput: &repository.User{
				UserID:      "mockUserID",
				FullName:    "Old Name",
				PhoneNumber: "+621234567890",
			},
			mockCheckPhoneNumber: 0,
			expectedCode:         http.StatusAccepted,
			expectedMsg:          "User Profile Successfully Updated!",
			expectedError:        true,
			isInputValidate:      false,
			isProceedUpdate:      false,
		},
		{
			name: "Invalid Phone Number",
			requestBody: map[string]string{
				"full_name":    "John Doe",
				"phone_number": "21234567890",
			},
			userID:        "mockUserID",
			existingPhone: "+621234567890",
			mockOutput: &repository.User{
				UserID:      "mockUserID",
				FullName:    "Old Name",
				PhoneNumber: "+621234567890",
			},
			mockCheckPhoneNumber: 0,
			expectedCode:         http.StatusAccepted,
			expectedMsg:          "User Profile Successfully Updated!",
			expectedError:        true,
			isInputValidate:      false,
			isProceedUpdate:      false,
		},
		{
			name: "Phone Number Existed",
			requestBody: map[string]string{
				"full_name":    "John Doe",
				"phone_number": "+621234567890",
			},
			userID:        "mockUserID",
			existingPhone: "+621234567890",
			mockOutput: &repository.User{
				UserID:      "mockUserID",
				FullName:    "Old Name",
				PhoneNumber: "+621234567893",
			},
			mockCheckPhoneNumber: 1,
			expectedCode:         http.StatusAccepted,
			expectedMsg:          "User Profile Successfully Updated!",
			expectedError:        true,
			isInputValidate:      true,
			isProceedUpdate:      false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.isInputValidate {
				mockRepository.EXPECT().GetUserByUserId(gomock.Any(), tc.userID).Return(tc.mockOutput, nil)
				mockRepository.EXPECT().CheckPhoneNumber(gomock.Any(), tc.requestBody["phone_number"]).Return(int64(tc.mockCheckPhoneNumber), nil)
				if tc.isProceedUpdate {
					mockRepository.EXPECT().UpdateUserProfile(gomock.Any(), gomock.Any()).Return(nil)
				}
			}
			reqBody, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("error marshaling request body: %v", err)
			}

			req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewReader(reqBody))
			req.Header.Set("Authorization", "Bearer "+jwtToken)
			rec := httptest.NewRecorder()

			e := echo.New()
			c := e.NewContext(req, rec)

			err = server.UpdateProfile(c)

			if !tc.expectedError {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCode, rec.Code)

			} else {
				assert.Error(t, err)
				assert.NotEqual(t, tc.expectedCode, rec.Code)
			}
		})
	}
}
