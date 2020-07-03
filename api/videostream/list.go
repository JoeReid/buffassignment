package videostream

import (
	"net/http"

	"github.com/JoeReid/apiutils"
	"github.com/JoeReid/apiutils/tracer"
	"github.com/JoeReid/buffassignment/api/types"
	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/opentracing/opentracing-go"
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
	sp, ctx := opentracing.StartSpanFromContext(r.Context(), "list stream handler")
	defer sp.Finish()

	tracer.Log(sp, "read pagination values from request")
	count, skip, err := apiutils.Paginate(r, apiutils.DefaultCount(10), apiutils.MaxCount(10))
	if err != nil {
		tracer.SetError(sp, err)
		c.Respond(ctx, w, http.StatusBadRequest, err)
		return
	}

	tracer.Log(sp, "list video streams from store")
	streams, err := s.store.ListVideoStream(ctx, count*skip, count)
	if err != nil {
		if err == model.ErrNotFound {
			tracer.Log(sp, "streams not found")
			c.Respond(ctx, w, http.StatusOK, []types.VideoStream{})
			return
		}

		tracer.Log(sp, "unexpected store error")
		tracer.SetError(sp, err)
		c.Respond(ctx, w, http.StatusInternalServerError, err)
		return
	}

	tracer.Log(sp, "return streams")
	c.Respond(ctx, w, http.StatusOK, types.NewVideoStreams(streams))
}
