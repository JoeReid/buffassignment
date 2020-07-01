package videostream

import (
	"net/http"

	"github.com/JoeReid/apiutils"
	"github.com/JoeReid/buffassignment/api/types"
	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// NewGetHandler returns a new instance of the get action of
// the videostream API using the given store instance.
//
// The store is provided as an argument for easy dependency injection in tests
// I.E: using the testing mock store rather than a full DB for API testing
func NewGetHandler(store model.VideoStreamStore) apiutils.Handler {
	return &streamGet{store}
}

// streamGet implements the apiutils.Handler interface to provide the
// get portion of the videostream API
type streamGet struct {
	store model.VideoStreamStore
}

// ServeCodec serves the API using the apiutils.Handler pattern
// This allows the business logic to live here, and the encoding to live separate from it
// This also makes testing easier, as there is a test codec that allows us to peek at the output
// in a testing context.
func (s *streamGet) ServeCodec(c apiutils.Codec, w http.ResponseWriter, r *http.Request) {
	vID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		c.Respond(r.Context(), w, http.StatusBadRequest, err)
		return
	}

	stream, err := s.store.GetVideoStream(model.VideoStreamID(vID))
	if err != nil {
		if err == model.ErrNotFound {
			c.Respond(r.Context(), w, http.StatusNotFound, err)
			return
		}
		c.Respond(r.Context(), w, http.StatusInternalServerError, err)
		return
	}
	c.Respond(r.Context(), w, http.StatusOK, types.NewVideoStream(*stream))
}
