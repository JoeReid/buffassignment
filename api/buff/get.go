package buff

import (
	"net/http"

	"github.com/JoeReid/apiutils"
	"github.com/JoeReid/apiutils/tracer"
	"github.com/JoeReid/buffassignment/api/types"
	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

// NewGetHandler returns a new instance of the get action of
// the buff API using the given store instance.
//
// The store is provided as an argument for easy dependency injection in tests
// I.E: using the testing mock store rather than a full DB for API testing
func NewGetHandler(store model.BuffStore) apiutils.Handler {
	return &buffGet{store}
}

// buffGet implements the apiutils.Handler interface to provide the
// get portion of the buff API
type buffGet struct {
	store model.BuffStore
}

// ServeCodec serves the API using the apiutils.Handler pattern
// This allows the business logic to live here, and the encoding to live separate from it
// This also makes testing easier, as there is a test codec that allows us to peek at the output
// in a testing context.
func (b *buffGet) ServeCodec(c apiutils.Codec, w http.ResponseWriter, r *http.Request) {
	sp, ctx := opentracing.StartSpanFromContext(r.Context(), "get buff handler")
	defer sp.Finish()

	tracer.Log(sp, "get uuid from url params")
	bID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		tracer.SetError(sp, err)
		c.Respond(ctx, w, http.StatusBadRequest, err)
		return
	}

	tracer.Logf(sp, "get buff: %v", bID)
	buff, err := b.store.GetBuff(ctx, model.BuffID(bID))
	if err != nil {
		if err == model.ErrNotFound {
			tracer.Log(sp, "buff not found")
			c.Respond(ctx, w, http.StatusNotFound, err)
			return
		}

		tracer.Log(sp, "unexpected store error")
		tracer.SetError(sp, err)
		c.Respond(ctx, w, http.StatusInternalServerError, err)
		return
	}

	tracer.Log(sp, "return buff")
	c.Respond(ctx, w, http.StatusOK, types.NewBuff(*buff))
}
