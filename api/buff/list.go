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

// NewListHandler returns a new instance of the list action of
// the buff API using the given store instance.
//
// The store is provided as an argument for easy dependency injection in tests
// I.E: using the testing mock store rather than a full DB for API testing
func NewListHandler(store model.BuffStore) apiutils.Handler {
	return &buffList{store}
}

// NewListForStreamHandler returns a new instance of the list for stream action of
// the buff API using the given store instance.
//
// The store is provided as an argument for easy dependency injection in tests
// I.E: using the testing mock store rather than a full DB for API testing
func NewListForStreamHandler(store model.BuffStore) apiutils.Handler {
	return &buffListForStream{store}
}

// buffList implements the apiutils.Handler interface to provide the
// list portion of the buff API
type buffList struct {
	store model.BuffStore
}

// ServeCodec serves the API using the apiutils.Handler pattern
// This allows the business logic to live here, and the encoding to live separate from it
// This also makes testing easier, as there is a test codec that allows us to peek at the output
// in a testing context.
func (b *buffList) ServeCodec(c apiutils.Codec, w http.ResponseWriter, r *http.Request) {
	sp, ctx := opentracing.StartSpanFromContext(r.Context(), "list buff handler")
	defer sp.Finish()

	tracer.Log(sp, "read pagination values from request")
	count, skip, err := apiutils.Paginate(r, apiutils.DefaultCount(10), apiutils.MaxCount(10))
	if err != nil {
		tracer.SetError(sp, err)
		c.Respond(ctx, w, http.StatusBadRequest, err)
		return
	}

	tracer.Log(sp, "list buffs from store")
	buffs, err := b.store.ListBuff(ctx, count*skip, count)
	if err != nil {
		if err == model.ErrNotFound {
			tracer.Log(sp, "buffs not found")
			c.Respond(ctx, w, http.StatusOK, []types.Buff{})
			return
		}

		tracer.Log(sp, "unexpected store error")
		tracer.SetError(sp, err)
		c.Respond(ctx, w, http.StatusInternalServerError, err)
		return
	}

	tracer.Log(sp, "return buffs")
	c.Respond(ctx, w, http.StatusOK, types.NewBuffs(buffs))
}

// buffListForStream implements the apiutils.Handler interface to provide the
// list for stream portion of the buff API
type buffListForStream struct {
	store model.BuffStore
}

// ServeCodec serves the API using the apiutils.Handler pattern
// This allows the business logic to live here, and the encoding to live separate from it
// This also makes testing easier, as there is a test codec that allows us to peek at the output
// in a testing context.
func (b *buffListForStream) ServeCodec(c apiutils.Codec, w http.ResponseWriter, r *http.Request) {
	sp, ctx := opentracing.StartSpanFromContext(r.Context(), "list buff for stream handler")
	defer sp.Finish()

	tracer.Log(sp, "get stream uuid from url params")
	vID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		tracer.SetError(sp, err)
		c.Respond(ctx, w, http.StatusBadRequest, err)
		return
	}

	// Assume the list is short, we can add pagination later if needed
	tracer.Logf(sp, "list buffs from store, streamid=%v", vID)
	buffs, err := b.store.ListBuffForStream(ctx, model.VideoStreamID(vID), 0, 0)
	if err != nil {
		if err == model.ErrNotFound {
			tracer.Log(sp, "buffs not found")
			c.Respond(ctx, w, http.StatusOK, []types.Buff{})
			return
		}

		tracer.Log(sp, "unexpected store error")
		tracer.SetError(sp, err)
		c.Respond(ctx, w, http.StatusInternalServerError, err)
		return
	}

	tracer.Log(sp, "return buffs")
	c.Respond(ctx, w, http.StatusOK, types.NewBuffs(buffs))
}
