package buff_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
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

func TestListBuffsForStream(t *testing.T) {
	sentinelUUID := uuid.New()

	var tests = []struct {
		name                 string
		requestParams        map[string]string
		storeResponse        []model.Buff
		storeError           error
		expectResponseCode   int
		expectResponseData   interface{}
		expectStoreNotCalled bool
	}{
		{
			name: "happy path returns ok",
			requestParams: map[string]string{
				"uuid": sentinelUUID.String(),
			},
			storeResponse: []model.Buff{
				{
					ID:       model.BuffID(sentinelUUID),
					Stream:   model.VideoStreamID(sentinelUUID),
					Question: "what's the answer to life, the universe, and everything?",
					Answers: []model.Answer{
						{ID: model.AnswerID(sentinelUUID), Text: "42", Correct: true},
						{ID: model.AnswerID(sentinelUUID), Text: "43", Correct: false},
						{ID: model.AnswerID(sentinelUUID), Text: "44", Correct: false},
					},
				},
			},
			storeError:         nil,
			expectResponseCode: http.StatusOK,
			expectResponseData: []types.Buff{
				{
					UUID:             sentinelUUID.String(),
					VideoStreamUUID:  sentinelUUID.String(),
					Question:         "what's the answer to life, the universe, and everything?",
					CorrectAnswer:    "42",
					IncorrectAnswers: []string{"43", "44"},
				},
			},
		},
		{
			name: "returns empty list on store no found",
			requestParams: map[string]string{
				"uuid": sentinelUUID.String(),
			},
			storeResponse:      nil,
			storeError:         model.ErrNotFound,
			expectResponseCode: http.StatusOK,
			expectResponseData: []types.Buff{},
		},
		{
			name: "returns bad request on missformated uuid",
			requestParams: map[string]string{
				"uuid": "not_a_valid_uuid",
			},
			storeResponse:        nil,
			storeError:           nil,
			expectResponseCode:   http.StatusBadRequest,
			expectResponseData:   errors.New("invalid UUID length: 16"),
			expectStoreNotCalled: true,
		},
		{
			name: "returns internal error on unexpected store error",
			requestParams: map[string]string{
				"uuid": sentinelUUID.String(),
			},
			storeResponse:      nil,
			storeError:         errors.New("the world exploded"),
			expectResponseCode: http.StatusInternalServerError,
			expectResponseData: errors.New("the world exploded"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			// Setup the mock store object to return the data configured in the test fixture
			testingStore := testmodel.NewModelMock()
			testingStore.On("ListBuffForStream", mock.Anything, mock.Anything, mock.Anything).Return(tt.storeResponse, tt.storeError)

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
			handler := buff.NewListForStreamHandler(testingStore)
			handler.ServeCodec(codec, nil, req)

			// assert that the handler returns the expected data
			codec.AssertCalled(t, "Respond", mock.Anything, nil, tt.expectResponseCode, tt.expectResponseData)

			// assert that the handler responded only once
			codec.AssertNumberOfCalls(t, "Respond", 1)

			// If the handler needs to use the store, assert it made the right call
			if tt.expectStoreNotCalled {
				// assert that no calls to the store were made
				testingStore.AssertNotCalled(t, "ListBuffForStream", mock.Anything, mock.Anything, mock.Anything)
			} else {
				// assert that the store was called with the correct uuid
				testingStore.AssertCalled(t, "ListBuffForStream", mock.Anything, 0, 0)
			}
		})
	}
}

func TestListBuffs(t *testing.T) {
	sentinelUUID := uuid.New()

	var tests = []struct {
		name                 string
		requestURLValues     map[string]string
		storeResponse        []model.Buff
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
			storeResponse: []model.Buff{
				{
					ID:       model.BuffID(sentinelUUID),
					Stream:   model.VideoStreamID(sentinelUUID),
					Question: "what's the answer to life, the universe, and everything?",
					Answers: []model.Answer{
						{ID: model.AnswerID(sentinelUUID), Text: "42", Correct: true},
						{ID: model.AnswerID(sentinelUUID), Text: "43", Correct: false},
						{ID: model.AnswerID(sentinelUUID), Text: "44", Correct: false},
					},
				},
				{
					ID:       model.BuffID(sentinelUUID),
					Stream:   model.VideoStreamID(sentinelUUID),
					Question: "what's the answer to life, the universe, and everything?",
					Answers: []model.Answer{
						{ID: model.AnswerID(sentinelUUID), Text: "42", Correct: true},
						{ID: model.AnswerID(sentinelUUID), Text: "43", Correct: false},
						{ID: model.AnswerID(sentinelUUID), Text: "44", Correct: false},
					},
				},
				{
					ID:       model.BuffID(sentinelUUID),
					Stream:   model.VideoStreamID(sentinelUUID),
					Question: "what's the answer to life, the universe, and everything?",
					Answers: []model.Answer{
						{ID: model.AnswerID(sentinelUUID), Text: "42", Correct: true},
						{ID: model.AnswerID(sentinelUUID), Text: "43", Correct: false},
						{ID: model.AnswerID(sentinelUUID), Text: "44", Correct: false},
					},
				},
				{
					ID:       model.BuffID(sentinelUUID),
					Stream:   model.VideoStreamID(sentinelUUID),
					Question: "what's the answer to life, the universe, and everything?",
					Answers: []model.Answer{
						{ID: model.AnswerID(sentinelUUID), Text: "42", Correct: true},
						{ID: model.AnswerID(sentinelUUID), Text: "43", Correct: false},
						{ID: model.AnswerID(sentinelUUID), Text: "44", Correct: false},
					},
				},
				{
					ID:       model.BuffID(sentinelUUID),
					Stream:   model.VideoStreamID(sentinelUUID),
					Question: "what's the answer to life, the universe, and everything?",
					Answers: []model.Answer{
						{ID: model.AnswerID(sentinelUUID), Text: "42", Correct: true},
						{ID: model.AnswerID(sentinelUUID), Text: "43", Correct: false},
						{ID: model.AnswerID(sentinelUUID), Text: "44", Correct: false},
					},
				},
				{
					ID:       model.BuffID(sentinelUUID),
					Stream:   model.VideoStreamID(sentinelUUID),
					Question: "what's the answer to life, the universe, and everything?",
					Answers: []model.Answer{
						{ID: model.AnswerID(sentinelUUID), Text: "42", Correct: true},
						{ID: model.AnswerID(sentinelUUID), Text: "43", Correct: false},
						{ID: model.AnswerID(sentinelUUID), Text: "44", Correct: false},
					},
				},
				{
					ID:       model.BuffID(sentinelUUID),
					Stream:   model.VideoStreamID(sentinelUUID),
					Question: "what's the answer to life, the universe, and everything?",
					Answers: []model.Answer{
						{ID: model.AnswerID(sentinelUUID), Text: "42", Correct: true},
						{ID: model.AnswerID(sentinelUUID), Text: "43", Correct: false},
						{ID: model.AnswerID(sentinelUUID), Text: "44", Correct: false},
					},
				},
				{
					ID:       model.BuffID(sentinelUUID),
					Stream:   model.VideoStreamID(sentinelUUID),
					Question: "what's the answer to life, the universe, and everything?",
					Answers: []model.Answer{
						{ID: model.AnswerID(sentinelUUID), Text: "42", Correct: true},
						{ID: model.AnswerID(sentinelUUID), Text: "43", Correct: false},
						{ID: model.AnswerID(sentinelUUID), Text: "44", Correct: false},
					},
				},
				{
					ID:       model.BuffID(sentinelUUID),
					Stream:   model.VideoStreamID(sentinelUUID),
					Question: "what's the answer to life, the universe, and everything?",
					Answers: []model.Answer{
						{ID: model.AnswerID(sentinelUUID), Text: "42", Correct: true},
						{ID: model.AnswerID(sentinelUUID), Text: "43", Correct: false},
						{ID: model.AnswerID(sentinelUUID), Text: "44", Correct: false},
					},
				},
				{
					ID:       model.BuffID(sentinelUUID),
					Stream:   model.VideoStreamID(sentinelUUID),
					Question: "what's the answer to life, the universe, and everything?",
					Answers: []model.Answer{
						{ID: model.AnswerID(sentinelUUID), Text: "42", Correct: true},
						{ID: model.AnswerID(sentinelUUID), Text: "43", Correct: false},
						{ID: model.AnswerID(sentinelUUID), Text: "44", Correct: false},
					},
				},
			},
			storeError:         nil,
			expectOffset:       0,
			expectLimit:        10,
			expectResponseCode: http.StatusOK,
			expectResponseData: []types.Buff{
				{
					UUID:             sentinelUUID.String(),
					VideoStreamUUID:  sentinelUUID.String(),
					Question:         "what's the answer to life, the universe, and everything?",
					CorrectAnswer:    "42",
					IncorrectAnswers: []string{"43", "44"},
				},
				{
					UUID:             sentinelUUID.String(),
					VideoStreamUUID:  sentinelUUID.String(),
					Question:         "what's the answer to life, the universe, and everything?",
					CorrectAnswer:    "42",
					IncorrectAnswers: []string{"43", "44"},
				},
				{
					UUID:             sentinelUUID.String(),
					VideoStreamUUID:  sentinelUUID.String(),
					Question:         "what's the answer to life, the universe, and everything?",
					CorrectAnswer:    "42",
					IncorrectAnswers: []string{"43", "44"},
				},
				{
					UUID:             sentinelUUID.String(),
					VideoStreamUUID:  sentinelUUID.String(),
					Question:         "what's the answer to life, the universe, and everything?",
					CorrectAnswer:    "42",
					IncorrectAnswers: []string{"43", "44"},
				},
				{
					UUID:             sentinelUUID.String(),
					VideoStreamUUID:  sentinelUUID.String(),
					Question:         "what's the answer to life, the universe, and everything?",
					CorrectAnswer:    "42",
					IncorrectAnswers: []string{"43", "44"},
				},
				{
					UUID:             sentinelUUID.String(),
					VideoStreamUUID:  sentinelUUID.String(),
					Question:         "what's the answer to life, the universe, and everything?",
					CorrectAnswer:    "42",
					IncorrectAnswers: []string{"43", "44"},
				},
				{
					UUID:             sentinelUUID.String(),
					VideoStreamUUID:  sentinelUUID.String(),
					Question:         "what's the answer to life, the universe, and everything?",
					CorrectAnswer:    "42",
					IncorrectAnswers: []string{"43", "44"},
				},
				{
					UUID:             sentinelUUID.String(),
					VideoStreamUUID:  sentinelUUID.String(),
					Question:         "what's the answer to life, the universe, and everything?",
					CorrectAnswer:    "42",
					IncorrectAnswers: []string{"43", "44"},
				},
				{
					UUID:             sentinelUUID.String(),
					VideoStreamUUID:  sentinelUUID.String(),
					Question:         "what's the answer to life, the universe, and everything?",
					CorrectAnswer:    "42",
					IncorrectAnswers: []string{"43", "44"},
				},
				{
					UUID:             sentinelUUID.String(),
					VideoStreamUUID:  sentinelUUID.String(),
					Question:         "what's the answer to life, the universe, and everything?",
					CorrectAnswer:    "42",
					IncorrectAnswers: []string{"43", "44"},
				},
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
			expectResponseData: []types.Buff{},
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
			storeResponse: []model.Buff{
				{
					ID:       model.BuffID(sentinelUUID),
					Stream:   model.VideoStreamID(sentinelUUID),
					Question: "question 7?",
					Answers: []model.Answer{
						{ID: model.AnswerID(sentinelUUID), Text: "answer 19", Correct: true},
						{ID: model.AnswerID(sentinelUUID), Text: "answer 20", Correct: false},
						{ID: model.AnswerID(sentinelUUID), Text: "answer 21", Correct: false},
					},
				},
				{
					ID:       model.BuffID(sentinelUUID),
					Stream:   model.VideoStreamID(sentinelUUID),
					Question: "question 8?",
					Answers: []model.Answer{
						{ID: model.AnswerID(sentinelUUID), Text: "answer 22", Correct: true},
						{ID: model.AnswerID(sentinelUUID), Text: "answer 23", Correct: false},
						{ID: model.AnswerID(sentinelUUID), Text: "answer 24", Correct: false},
					},
				},
				{
					ID:       model.BuffID(sentinelUUID),
					Stream:   model.VideoStreamID(sentinelUUID),
					Question: "question 9?",
					Answers: []model.Answer{
						{ID: model.AnswerID(sentinelUUID), Text: "answer 25", Correct: true},
						{ID: model.AnswerID(sentinelUUID), Text: "answer 26", Correct: false},
						{ID: model.AnswerID(sentinelUUID), Text: "answer 27", Correct: false},
					},
				},
			},
			storeError:         nil,
			expectOffset:       6,
			expectLimit:        3,
			expectResponseCode: http.StatusOK,
			expectResponseData: []types.Buff{
				{
					UUID:             sentinelUUID.String(),
					VideoStreamUUID:  sentinelUUID.String(),
					Question:         "question 7?",
					CorrectAnswer:    "answer 19",
					IncorrectAnswers: []string{"answer 20", "answer 21"},
				},
				{
					UUID:             sentinelUUID.String(),
					VideoStreamUUID:  sentinelUUID.String(),
					Question:         "question 8?",
					CorrectAnswer:    "answer 22",
					IncorrectAnswers: []string{"answer 23", "answer 24"},
				},
				{
					UUID:             sentinelUUID.String(),
					VideoStreamUUID:  sentinelUUID.String(),
					Question:         "question 9?",
					CorrectAnswer:    "answer 25",
					IncorrectAnswers: []string{"answer 26", "answer 27"},
				},
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
			testingStore.On("ListBuff", mock.Anything, mock.Anything, mock.Anything).Return(tt.storeResponse, tt.storeError)

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
			handler := buff.NewListHandler(testingStore)
			handler.ServeCodec(codec, nil, req)

			// assert that the handler returns the expected data
			codec.AssertCalled(t, "Respond", mock.Anything, nil, tt.expectResponseCode, tt.expectResponseData)

			// assert that the handler responded only once
			codec.AssertNumberOfCalls(t, "Respond", 1)

			// If the handler needs to use the store, assert it made the right call
			if tt.expectStoreNotCalled {
				// assert that no calls to the store were made
				testingStore.AssertNotCalled(t, "ListBuff", mock.Anything, mock.Anything)
			} else {
				// assert that the store was called with the correct uuid
				testingStore.AssertCalled(t, "ListBuff", tt.expectOffset, tt.expectLimit)
			}
		})
	}
}
