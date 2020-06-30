package buff_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/JoeReid/apiutils/testingcodec"
	"github.com/JoeReid/buffassignment/api/buff"
	"github.com/JoeReid/buffassignment/api/types"
	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/JoeReid/buffassignment/internal/model/testmodel"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetBuffs(t *testing.T) {
	sentinelUUID := uuid.New()

	var tests = []struct {
		name                 string
		requestParams        map[string]string
		storeResponse        *model.Buff
		storeError           error
		expectResponseCode   int
		expectResponseData   interface{}
		expectStoreNotCalled bool
	}{
		{
			name: "returns correct data on happy path",
			requestParams: map[string]string{
				"uuid": sentinelUUID.String(),
			},
			storeResponse: &model.Buff{
				ID:       model.BuffID(sentinelUUID),
				Stream:   model.VideoStreamID(sentinelUUID),
				Question: "what's the answer to life, the universe, and everything?",
				Answers: []model.Answer{
					{ID: model.AnswerID(sentinelUUID), Text: "42", Correct: true},
					{ID: model.AnswerID(sentinelUUID), Text: "43", Correct: false},
					{ID: model.AnswerID(sentinelUUID), Text: "44", Correct: false},
				},
			},
			storeError: nil,
			expectResponseData: types.Buff{
				UUID:             sentinelUUID.String(),
				VideoStreamUUID:  sentinelUUID.String(),
				Question:         "what's the answer to life, the universe, and everything?",
				CorrectAnswer:    "42",
				IncorrectAnswers: []string{"43", "44"},
			},
			expectResponseCode: http.StatusOK,
		},
		{
			name: "returns not found on store not found error",
			requestParams: map[string]string{
				"uuid": sentinelUUID.String(),
			},
			storeResponse:      nil,
			storeError:         model.ErrNotFound,
			expectResponseData: model.ErrNotFound,
			expectResponseCode: http.StatusNotFound,
		},
		{
			name: "returns bad request on missformated uuid",
			requestParams: map[string]string{
				"uuid": "not_a_valid_uuid",
			},
			storeResponse:        nil,
			storeError:           nil,
			expectResponseData:   errors.New("invalid UUID length: 16"),
			expectResponseCode:   http.StatusBadRequest,
			expectStoreNotCalled: true,
		},
		{
			name: "returns internal error on unexpected store error",
			requestParams: map[string]string{
				"uuid": sentinelUUID.String(),
			},
			storeResponse:      nil,
			storeError:         errors.New("the world exploded"),
			expectResponseData: errors.New("the world exploded"),
			expectResponseCode: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			// Setup the mock store object to return the data configured in the test fixture
			testingStore := testmodel.NewModelMock()
			testingStore.On("GetBuff", mock.Anything).Return(tt.storeResponse, tt.storeError)

			// Build the request to the spec of the test fixture
			rctx := chi.NewRouteContext()
			for k, v := range tt.requestParams {
				rctx.URLParams.Add(k, v)
			}
			req, err := http.NewRequest("GET", "", nil)
			require.NoError(t, err, "failed to build request for test")
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

			// Use the testing codec to assert handler behaviour
			codec := testingcodec.New()
			codec.On("Respond", mock.Anything, nil, mock.Anything, mock.Anything).Return()

			// Create the handler under test, and execute it
			handler := buff.NewGetHandler(testingStore)
			handler.ServeCodec(codec, nil, req)

			// assert that the handler returns the expected data
			codec.AssertCalled(t, "Respond", mock.Anything, nil, tt.expectResponseCode, tt.expectResponseData)

			// assert that the handler responded only once
			codec.AssertNumberOfCalls(t, "Respond", 1)

			// If the handler needs to use the store, assert it made the right call
			if tt.expectStoreNotCalled {
				// assert that no calls to the store were made
				testingStore.AssertNotCalled(t, "GetBuff", mock.Anything)
			} else {
				// assert that the store was called with the correct uuid
				testingStore.AssertCalled(t, "GetBuff", model.BuffID(sentinelUUID))
			}
		})
	}
}
