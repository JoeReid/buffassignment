package videostream

import (
	"net/http"

	"github.com/JoeReid/apiutils"
	"github.com/JoeReid/buffassignment/api/types"
	"github.com/JoeReid/buffassignment/internal/model"
)

// NewListHandler returns a new instance of the list action of
// the videostream API using the given store instance.
//
// The store is provided as an argument for easy dependency injection in tests
// I.E: using the testing mock store rather than a full DB for API testing
func NewListHandler(store model.VideoStreamStore) apiutils.Handler {
	return &streamList{store}
}

// streamGet implements the apiutils.Handler interface to provide the
// list portion of the videostream API
type streamList struct {
	store model.VideoStreamStore
}

// ServeCodec serves the API using the apiutils.Handler pattern
// This allows the business logic to live here, and the encoding to live separate from it
// This also makes testing easier, as there is a test codec that allows us to peek at the output
// in a testing context.
func (s *streamList) ServeCodec(c apiutils.Codec, w http.ResponseWriter, r *http.Request) {
	count, skip, err := apiutils.Paginate(r, apiutils.DefaultCount(10), apiutils.MaxCount(10))
	if err != nil {
		c.Respond(r.Context(), w, http.StatusBadRequest, err)
		return
	}

	streams, err := s.store.ListVideoStream(count*skip, count)
	if err != nil {
		if err == model.ErrNotFound {
			c.Respond(r.Context(), w, http.StatusOK, []types.VideoStream{})
			return
		}
		c.Respond(r.Context(), w, http.StatusInternalServerError, err)
		return
	}
	c.Respond(r.Context(), w, http.StatusOK, types.NewVideoStreams(streams))
}
