package gateways_test

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/neyrzx/youmusic/internal/config"
	"github.com/neyrzx/youmusic/internal/domain/entities"
	"github.com/neyrzx/youmusic/internal/gateways"
	"github.com/neyrzx/youmusic/mocks/internal_/gateways/mocks"
	"github.com/neyrzx/youmusic/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var errGatewayClientBadRequest = errors.New("status bad request")

func TestInfo(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		trackInfoDTO    entities.TrackInfo
		gatewayConfig   config.GatewayHTTPClient
		gatewayResponse io.ReadCloser
		gatewayError    error
		expectedResult  entities.TrackInfoResult
		expectedError   error
	}{
		{
			name: "success response",
			trackInfoDTO: entities.TrackInfo{
				Song:  "title",
				Group: "group",
			},
			gatewayConfig: config.GatewayHTTPClient{},
			gatewayResponse: newResponse(`
			{
				"releaseDate":"16.07.2006",
				"text":"Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight  ",
				"link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw"
			}`),
			gatewayError: nil,
			expectedResult: entities.TrackInfoResult{
				ReleaseDate: mustTimeParse("16.07.2006"),
				Text:        "Ooh baby, don't you know I suffer?\nOoh baby, can you hear me moan?\nYou caught me under false pretenses\nHow long before you let me go?\n\nOoh\nYou set my soul alight\nOoh\nYou set my soul alight  ",
				Link:        "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
			},
			expectedError: nil,
		},
		{
			name: "gateway client error",
			trackInfoDTO: entities.TrackInfo{
				Song:  "title",
				Group: "group",
			},
			gatewayConfig:   config.GatewayHTTPClient{},
			gatewayResponse: newResponse(``),
			gatewayError:    errGatewayClientBadRequest,
			expectedResult:  entities.TrackInfoResult{},
			expectedError:   errGatewayClientBadRequest,
		},
		{
			name: "gateway client error #2",
			trackInfoDTO: entities.TrackInfo{
				Song:  "title",
				Group: "group",
			},
			gatewayConfig:   config.GatewayHTTPClient{},
			gatewayResponse: newResponse(``),
			gatewayError:    errGatewayClientBadRequest,
			expectedResult:  entities.TrackInfoResult{},
			expectedError:   errGatewayClientBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mockInfoGatewayCalls := mocks.NewMockClient(t)
			mockInfoGatewayCalls.EXPECT().Get(mock.Anything, mock.Anything).
				Return(test.gatewayResponse, test.gatewayError).
				Once()

			gateway := gateways.NewMusicInfoGateway(mockInfoGatewayCalls, test.gatewayConfig)
			actualResult, actualErr := gateway.Info(context.Background(), test.trackInfoDTO)

			require.ErrorIs(t, actualErr, test.expectedError)
			assert.EqualValues(t, test.expectedResult, actualResult)
		})
	}
}

func newResponse(body string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(body))
}

func mustTimeParse(value string) time.Time {
	t, err := time.Parse(utils.ReleaseDateLayout, value)
	if err != nil {
		panic(fmt.Errorf("musicinfo_test: failed to time.Parse(%s, %s): %w", utils.ReleaseDateLayout, value, err))
	}
	return t
}
