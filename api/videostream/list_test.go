package videostream_test

import (
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/JoeReid/apiutils/testingcodec"
	"github.com/JoeReid/buffassignment/api/types"
	"github.com/JoeReid/buffassignment/api/videostream"
	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/JoeReid/buffassignment/internal/model/testmodel"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListVideoStreams(t *testing.T) {
	sentinelUUID := uuid.New()
	sentinelTime := time.Now()

	var tests = []struct {
		name                 string
		requestURLValues     map[string]string
		storeResponse        []model.VideoStream
		storeError           error
		expectOffset         int
		expectLimit          int
		expectResponseCode   int
		expectResponseData   interface{}
		expectStoreNotCalled bool
	}{
		{
			name:             "list with defaults returns 10 items",
			requestURLValues: nil,
			storeResponse: []model.VideoStream{
				{ID: model.VideoStreamID(sentinelUUID), Title: "1", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{ID: model.VideoStreamID(sentinelUUID), Title: "2", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{ID: model.VideoStreamID(sentinelUUID), Title: "3", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{ID: model.VideoStreamID(sentinelUUID), Title: "4", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{ID: model.VideoStreamID(sentinelUUID), Title: "5", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{ID: model.VideoStreamID(sentinelUUID), Title: "6", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{ID: model.VideoStreamID(sentinelUUID), Title: "7", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{ID: model.VideoStreamID(sentinelUUID), Title: "8", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{ID: model.VideoStreamID(sentinelUUID), Title: "9", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{ID: model.VideoStreamID(sentinelUUID), Title: "10", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
			},
			storeError:         nil,
			expectOffset:       0,
			expectLimit:        10,
			expectResponseCode: http.StatusOK,
			expectResponseData: []types.VideoStream{
				{UUID: sentinelUUID.String(), Title: "1", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{UUID: sentinelUUID.String(), Title: "2", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{UUID: sentinelUUID.String(), Title: "3", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{UUID: sentinelUUID.String(), Title: "4", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{UUID: sentinelUUID.String(), Title: "5", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{UUID: sentinelUUID.String(), Title: "6", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{UUID: sentinelUUID.String(), Title: "7", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{UUID: sentinelUUID.String(), Title: "8", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{UUID: sentinelUUID.String(), Title: "9", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{UUID: sentinelUUID.String(), Title: "10", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
			},
		},
		{
			name:               "returns empty list on not found error",
			requestURLValues:   nil,
			storeResponse:      nil,
			storeError:         model.ErrNotFound,
			expectOffset:       0,
			expectLimit:        10,
			expectResponseCode: http.StatusOK,
			expectResponseData: []types.VideoStream{},
		},
		{
			name:                 "returns error on requested data size > 10",
			requestURLValues:     map[string]string{"count": "11"},
			expectStoreNotCalled: true,
			expectResponseCode:   http.StatusBadRequest,
			expectResponseData:   errors.New("paginate error: count 11 must be < 10"),
		},
		{
			name:             "custom pagination values correctly computed",
			requestURLValues: map[string]string{"count": "3", "skip": "2"},
			storeResponse: []model.VideoStream{
				{ID: model.VideoStreamID(sentinelUUID), Title: "7", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{ID: model.VideoStreamID(sentinelUUID), Title: "8", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{ID: model.VideoStreamID(sentinelUUID), Title: "9", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
			},
			storeError:         nil,
			expectOffset:       6,
			expectLimit:        3,
			expectResponseCode: http.StatusOK,
			expectResponseData: []types.VideoStream{
				{UUID: sentinelUUID.String(), Title: "7", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{UUID: sentinelUUID.String(), Title: "8", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
				{UUID: sentinelUUID.String(), Title: "9", CreatedAt: sentinelTime, UpdatedAt: sentinelTime},
			},
		},
		{
			name:               "returns internal error on unexpected store error",
			storeResponse:      nil,
			storeError:         errors.New("the world exploded"),
			expectOffset:       0,
			expectLimit:        10,
			expectResponseCode: http.StatusInternalServerError,
			expectResponseData: errors.New("the world exploded"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Setup the mock store object to return the data configured in the test fixture
			testingStore := testmodel.NewModelMock()
			testingStore.On("ListVideoStream", mock.Anything, mock.Anything, mock.Anything).Return(tt.storeResponse, tt.storeError)

			// Build the request to the spec of the test fixture
			req, err := http.NewRequest("GET", "", nil)
			require.NoError(t, err, "failed to build request for test")
			vals := url.Values{}
			for k, v := range tt.requestURLValues {
				vals.Add(k, v)
			}
			req.URL.RawQuery = vals.Encode()

			// Use the testing codec to assert handler behaviour
			codec := testingcodec.New()
			codec.On("Respond", mock.Anything, nil, mock.Anything, mock.Anything).Return()

			// Create the handler under test, and execute it
			handler := videostream.NewListHandler(testingStore)
			handler.ServeCodec(codec, nil, req)

			// assert that the handler returns the expected data
			codec.AssertCalled(t, "Respond", mock.Anything, nil, tt.expectResponseCode, tt.expectResponseData)

			// assert that the handler responded only once
			codec.AssertNumberOfCalls(t, "Respond", 1)

			// If the handler needs to use the store, assert it made the right call
			if tt.expectStoreNotCalled {
				// assert that no calls to the store were made
				testingStore.AssertNotCalled(t, "ListVideoStream", mock.Anything, mock.Anything)
			} else {
				// assert that the store was called with the correct uuid
				testingStore.AssertCalled(t, "ListVideoStream", tt.expectOffset, tt.expectLimit)
			}
		})
	}
}
