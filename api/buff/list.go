package buff

import (
	"net/http"

	"github.com/JoeReid/apiutils"
	"github.com/JoeReid/buffassignment/api/types"
	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func NewListHandler(store model.BuffStore) apiutils.Handler {
	return &buffList{store}
}

func NewListForStreamHandler(store model.BuffStore) apiutils.Handler {
	return &buffListForStream{store}
}

type buffList struct {
	store model.BuffStore
}

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

type buffListForStream struct {
	store model.BuffStore
}

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
