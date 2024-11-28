package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/neyrzx/youmusic/internal/domain/services"
	"github.com/neyrzx/youmusic/internal/dtos"
	"github.com/neyrzx/youmusic/mocks/internal_/domain/services/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Вынести в слой репозитория
var (
	ErrTrackAlreadyExists     = errors.New("track already exists")
	ErrTrackRequestInfoFailed = errors.New("failed to get info about the track from external api gateway")
	ErrTrackCreateDBError     = errors.New("failed to save track into db")
)

func TestTracksServiceCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name                     string
		trackDTO                 dtos.TrackCreateDTO
		trackInfoResultDTO       dtos.TrackInfoResultDTO
		expectedRepoCalls        func(mockRepo *mocks.MockTracksRepository)
		expectedInfoGatewayCalls func(mockInfoGateway *mocks.MockTracksInfoGateway, trackInfoResult dtos.TrackInfoResultDTO)
		expectedRepoError        error
		expectedInfoGatewayError error
	}{
		{
			"case: failed to repo.Exists: already exists",
			dtos.TrackCreateDTO{Title: "Track title", Artist: "Muse"},
			dtos.TrackInfoResultDTO{},
			func(mockRepo *mocks.MockTracksRepository) {
				mockRepo.EXPECT().Exists(mock.Anything, mock.Anything).
					Return(ErrTrackAlreadyExists).
					Once()
			},
			func(mockInfoGateway *mocks.MockTracksInfoGateway, trackInfoResult dtos.TrackInfoResultDTO) {
				mockInfoGateway.EXPECT().Info(mock.Anything, mock.Anything).
					Return(trackInfoResult, nil).
					Unset()
			},
			ErrTrackAlreadyExists,
			nil,
		},
		{
			"case: failed to infoGateway.Info: fetch error",
			dtos.TrackCreateDTO{Title: "Track title", Artist: "Muse"},
			dtos.TrackInfoResultDTO{},
			func(mockRepo *mocks.MockTracksRepository) {
				mockRepo.EXPECT().Exists(mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			func(mockInfoGateway *mocks.MockTracksInfoGateway, trackInfoResult dtos.TrackInfoResultDTO) {
				mockInfoGateway.EXPECT().Info(mock.Anything, mock.Anything).
					Return(trackInfoResult, ErrTrackRequestInfoFailed).
					Once()
			},
			nil,
			ErrTrackRequestInfoFailed,
		},
		{
			"case: failed to repo.Create: db error",
			dtos.TrackCreateDTO{Title: "Track title", Artist: "Muse"},
			dtos.TrackInfoResultDTO{
				ReleaseDate: time.Now(),
				Text:        "Verse1\n\nVerse2",
				Link:        "https://cdn.youmusic.com/dijw0-wqre3-rewr-rewr/)",
			},
			func(mockRepo *mocks.MockTracksRepository) {
				mockRepo.EXPECT().Exists(mock.Anything, mock.Anything).
					Return(nil).
					Once()
				mockRepo.EXPECT().Create(mock.Anything, mock.Anything).
					Return(ErrTrackCreateDBError).
					Once()
			},
			func(mockInfoGateway *mocks.MockTracksInfoGateway, trackInfoResult dtos.TrackInfoResultDTO) {
				mockInfoGateway.EXPECT().Info(mock.Anything, mock.Anything).
					Return(trackInfoResult, nil).
					Once()
			},
			nil,
			ErrTrackCreateDBError,
		},
		{
			"case: success created",
			dtos.TrackCreateDTO{Title: "Track title", Artist: "Muse"},
			dtos.TrackInfoResultDTO{
				ReleaseDate: time.Now(),
				Text:        "Verse1\n\nVerse2",
				Link:        "https://cdn.youmusic.com/dijw0-wqre3-rewr-rewr/)",
			},
			func(mockRepo *mocks.MockTracksRepository) {
				mockRepo.EXPECT().Exists(mock.Anything, mock.Anything).
					Return(nil).
					Once()
				mockRepo.EXPECT().Create(mock.Anything, mock.Anything).
					Return(nil).
					Once()
			},
			func(mockInfoGateway *mocks.MockTracksInfoGateway, trackInfoResult dtos.TrackInfoResultDTO) {
				mockInfoGateway.EXPECT().Info(mock.Anything, mock.Anything).
					Return(trackInfoResult, nil).
					Once()
			},
			nil,
			nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			mockRespo := mocks.NewMockTracksRepository(t)
			mockInfoGateway := mocks.NewMockTracksInfoGateway(t)

			test.expectedRepoCalls(mockRespo)
			test.expectedInfoGatewayCalls(mockInfoGateway, test.trackInfoResultDTO)

			service := services.NewTracksService(mockRespo, mockInfoGateway)
			actualError := service.Create(context.Background(), test.trackDTO)

			switch {
			case test.expectedRepoError != nil:
				require.ErrorIs(t, actualError, test.expectedRepoError)
			case test.expectedInfoGatewayError != nil:
				require.ErrorIs(t, actualError, test.expectedInfoGatewayError)
			default:
				require.NoError(t, actualError)
			}
		})
	}
}
