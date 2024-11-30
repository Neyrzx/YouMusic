package v1_test

// func TestRetrieveCases(t *testing.T) {
// 	t.Parallel()

// 	tests := []struct {
// 		name             string
// 		requestPath      string
// 		requestBody      string
// 		requestMethod    string
// 		expectedStatus   int
// 		expectedResponse *v1.TracksRetrieveResponse
// 		expectedCalls    func(mockService *mocks.MockTracksService, ret v1.TracksRetrieveResponse)
// 	}{
// 		{
// 			name:           "retrieve track",
// 			requestPath:    "/api/v1/tracks/1",
// 			expectedStatus: http.StatusOK,
// 			expectedResponse: &v1.TracksRetrieveResponse{
// 				Artist:   "artist",
// 				Track:    "track",
// 				Lyric:    []string{"verse 1", "verse 2"},
// 				Link:     "https://y.be/cdsahfw3fiuiejfc",
// 				Released: time.Now(),
// 			},
// 			expectedCalls: func(mockService *mocks.MockTracksService, ret v1.TracksRetrieveResponse) {
// 				mockService.EXPECT().GetByID(mock.Anything, mock.Anything).
// 					Return(dtos.Track{
// 						Artist: ret.Artist,
// 						Track:  ret.Track,
// 					}, nil).
// 					Once()
// 			},
// 		},
// 	}

// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			t.Parallel()

// 			e := echo.New()
// 			e.Validator = validator.NewValidator()
// 			tracksGroup := e.Group("/api/v1/tracks")

// 			rec, req := buildRequest(t, http.MethodGet, test.requestPath, strings.NewReader(test.requestBody))
// 			mockService := mocks.NewMockTracksService(t)
// 			test.expectedCalls(mockService, test.expectedResponse)

// 			handler := v1.NewTracksHandlers(tracksGroup, mockService)
// 			handler.Retrieve(e.NewContext(req, rec))

// 			assert.Equal(t, test.expectedStatus, rec.Code)
// 			if test.expectedResponse != nil {
// 				assert.JSONEq(t, test.expectedResponse, rec.Body.String())
// 			}
// 		})
// 	}
// }
