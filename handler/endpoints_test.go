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
		name          string
		requestBody   interface{}
		mockInput     *repository.User
		mockOutput    int
		expectedCode  int
		expectedError string
	}{
		{
			name: "Successful Registration",
			requestBody: map[string]interface{}{
				"full_name":    "John Doe",
				"phone_number": "+62234567890",
				"Password":     "password",
			},
			mockInput: &repository.User{
				UserID:                   "mockUserID",
				FullName:                 "John Doe",
				PhoneNumber:              "+62234567890",
				Password:                 "hashedPassword",
				SuccessfullLoginAttempts: 0,
				CreatedAt:                time.Now(),
				UpdatedAt:                time.Now(),
			},
			mockOutput:    0,
			expectedCode:  http.StatusCreated,
			expectedError: "",
		},
		// Add more test cases as needed
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepository.EXPECT().CheckPhoneNumber(gomock.Any(), gomock.Any()).Return(int64(tc.mockOutput), nil)
			mockRepository.EXPECT().RegisterUser(gomock.Any()).Return(nil)

			reqBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()
			// Create a new Echo instance and handle the request
			e := echo.New()
			c := e.NewContext(req, rec)

			err := server.RegisterTheUser(c)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
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
		name          string
		requestBody   interface{}
		mockInput     *repository.User
		mockOutput    *repository.User
		expectedCode  int
		expectedError string
	}{
		{
			name: "Successful Login",
			requestBody: map[string]interface{}{
				"phone_number": "+62234567890",
				"Password":     "password",
			},
			mockInput: &repository.User{
				UserID:                   "mockUserID",
				PhoneNumber:              "+62234567890",
				Password:                 hashedPassword, // Assuming this is the hashed password stored in the repository
				SuccessfullLoginAttempts: 0,
				CreatedAt:                time.Now(),
				UpdatedAt:                time.Now(),
			},
			mockOutput: &repository.User{
				UserID:                   "mockUserID",
				PhoneNumber:              "+62234567890",
				Password:                 hashedPassword, // Assuming this is the hashed password stored in the repository
				SuccessfullLoginAttempts: 0,
				CreatedAt:                time.Now(),
				UpdatedAt:                time.Now(),
			},
			expectedCode:  http.StatusOK,
			expectedError: "",
		},
		// Add more test cases as needed
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepository.EXPECT().CheckUser(gomock.Any(), gomock.Any()).Return(tc.mockOutput, nil)
			mockRepository.EXPECT().UpdateLoginUser(gomock.Any(), gomock.Any()).Return(nil)

			reqBody, _ := json.Marshal(tc.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Create a new Echo instance and handle the request
			e := echo.New()
			c := e.NewContext(req, rec)

			err := server.LoginUser(c)

			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
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
		name          string
		requestBody   map[string]string
		userID        string // User ID extracted from JWT token
		existingPhone string // Existing phone number in the database
		mockOutput    *repository.User
		mockCheckUser int // Return value of CheckPhoneNumber
		expectedCode  int
		expectedMsg   string
		expectedError bool
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
			mockCheckUser: 0,
			expectedCode:  http.StatusAccepted,
			expectedMsg:   "User Profile Successfully Updated!",
			expectedError: false,
		},
		// Add more test cases for validation errors, authorization errors, etc.
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mockRepository.EXPECT().GetUserByUserId(gomock.Any(), tc.userID).Return(tc.mockOutput, nil)
			mockRepository.EXPECT().CheckPhoneNumber(gomock.Any(), tc.requestBody["phone_number"]).Return(int64(0), nil)
			mockRepository.EXPECT().UpdateUserProfile(gomock.Any(), gomock.Any()).Return(nil)
			// Marshal the request body
			reqBody, err := json.Marshal(tc.requestBody)
			if err != nil {
				t.Fatalf("error marshaling request body: %v", err)
			}

			// Create a new request with the test token and request body
			req := httptest.NewRequest(http.MethodPut, "/user", bytes.NewReader(reqBody))
			req.Header.Set("Authorization", "Bearer "+jwtToken)
			rec := httptest.NewRecorder()

			// Mock repository setup: GetUserByUserId, CheckPhoneNumber, UpdateUserProfile
			// Set expectations for repository calls based on test case

			// Create a new Echo instance and handle the request
			e := echo.New()
			c := e.NewContext(req, rec)

			// Call the function under test
			err = server.UpdateProfile(c)

			// Perform assertions on the response and error, if any
			if !tc.expectedError {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedCode, rec.Code)

			} else {
				assert.Error(t, err)
				// Handle other expected error cases
				// Assert error type and message
			}
		})
	}
}
