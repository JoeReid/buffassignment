package buff

import (
	"net/http"

	"github.com/JoeReid/apiutils"
	"github.com/JoeReid/buffassignment/api/types"
	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
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
	count, skip, err := apiutils.Paginate(r, apiutils.DefaultCount(10), apiutils.MaxCount(10))
	if err != nil {
		c.Respond(r.Context(), w, http.StatusBadRequest, err)
		return
	}

	buffs, err := b.store.ListBuff(count*skip, count)
	if err != nil {
		if err == model.ErrNotFound {
			c.Respond(r.Context(), w, http.StatusOK, []types.Buff{})
			return
		}
		c.Respond(r.Context(), w, http.StatusInternalServerError, err)
		return
	}
	c.Respond(r.Context(), w, http.StatusOK, types.NewBuffs(buffs))
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
	vID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		c.Respond(r.Context(), w, http.StatusBadRequest, err)
		return
	}

	// Assume the list is short, we can add pagination later if needed
	buffs, err := b.store.ListBuffForStream(model.VideoStreamID(vID), 0, 0)
	if err != nil {
		if err == model.ErrNotFound {
			c.Respond(r.Context(), w, http.StatusOK, []types.Buff{})
			return
		}
		c.Respond(r.Context(), w, http.StatusInternalServerError, err)
		return
	}
	c.Respond(r.Context(), w, http.StatusOK, types.NewBuffs(buffs))
}
