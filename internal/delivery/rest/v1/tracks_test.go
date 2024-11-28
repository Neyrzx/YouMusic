package v1_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	v1 "github.com/neyrzx/youmusic/internal/delivery/rest/v1"
	"github.com/neyrzx/youmusic/internal/domain/errors"
	"github.com/neyrzx/youmusic/mocks/internal_/delivery/rest/v1/mocks"
	"github.com/neyrzx/youmusic/pkg/validator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		requestPath      string
		requestBody      string
		requestMethod    string
		expectedStatus   int
		expectedResponse string
		expectedCalls    func(mockService *mocks.MockTracksService)
	}{
		{
			name:        "success creation",
			requestPath: "/api/v1/tracks",
			requestBody: `{
				"group": "Muse",
				"song": "Supermassive Black Hole"
			}`,
			requestMethod:    http.MethodPost,
			expectedStatus:   http.StatusCreated,
			expectedResponse: ``,
			expectedCalls: func(mockService *mocks.MockTracksService) {
				mockService.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Once()
			},
		},
		{
			name:           "request body malformed",
			requestPath:    "/api/v1/tracks",
			requestBody:    `{"`,
			requestMethod:  http.MethodPost,
			expectedStatus: http.StatusBadRequest,
			expectedResponse: `
			{
				"message":"request body malformed"
			}`,
			expectedCalls: func(mockService *mocks.MockTracksService) {
				mockService.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Unset()
			},
		},
		{
			name:           "validation error - empty request",
			requestPath:    "/api/v1/tracks",
			requestBody:    ``,
			requestMethod:  http.MethodPost,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedResponse: `
			{
				"message":{
					"TracksCreateRequest.Group":"field is required",
					"TracksCreateRequest.Song":"field is required"
				}
			}`,
			expectedCalls: func(mockService *mocks.MockTracksService) {
				mockService.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Unset()
			},
		},
		{
			name:        "validation error - group only",
			requestPath: "/api/v1/tracks",
			requestBody: `
			{
				"group": "Muse"
			}`,
			requestMethod:  http.MethodPost,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedResponse: `
			{
				"message": {
					"TracksCreateRequest.Song": "field is required"
				}
			}`,
			expectedCalls: func(mockService *mocks.MockTracksService) {
				mockService.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Unset()
			},
		},
		{
			name:        "validation error - song only",
			requestPath: "/api/v1/tracks",
			requestBody: `
			{
				"song": "track title"
			}`,
			requestMethod:  http.MethodPost,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedResponse: `
			{
				"message": {
					"TracksCreateRequest.Group": "field is required"
				}
			}`,
			expectedCalls: func(mockService *mocks.MockTracksService) {
				mockService.EXPECT().Create(mock.Anything, mock.Anything).Return(nil).Unset()
			},
		},
		{
			name:        "track already exists",
			requestPath: "/api/v1/tracks",
			requestBody: `
			{
				"group": "Muse",
				"song": "Supermassive Black Hole"
			}`,
			expectedStatus: http.StatusBadRequest,
			expectedResponse: `
			{
				"message": "track already exists"
			}`,
			expectedCalls: func(mockService *mocks.MockTracksService) {
				mockService.EXPECT().Create(mock.Anything, mock.Anything).Return(errors.ErrTrackAlreadyExists).Once()
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			e := echo.New()
			e.Validator = validator.NewValidator()
			tracksGroup := e.Group("/api/v1/tracks")

			rec, req := buildRequest(t, http.MethodPost, test.requestPath, strings.NewReader(test.requestBody))
			mockService := mocks.NewMockTracksService(t)
			test.expectedCalls(mockService)

			handler := v1.NewTracksHandlers(tracksGroup, mockService)
			handler.Create(e.NewContext(req, rec))

			assert.Equal(t, test.expectedStatus, rec.Code)
			if test.expectedResponse != "" {
				assert.JSONEq(t, test.expectedResponse, rec.Body.String())
			}
		})
	}
}

func buildRequest(
	t *testing.T, method string, target string, body io.Reader,
) (*httptest.ResponseRecorder, *http.Request) {
	t.Helper()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(method, target, body)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	return rec, req
}
