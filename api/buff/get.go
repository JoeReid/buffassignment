package buff

import (
	"net/http"

	"github.com/JoeReid/apiutils"
	"github.com/JoeReid/buffassignment/api/types"
	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func NewGetHandler(store model.BuffStore) apiutils.Handler {
	return &buffGet{store}
}

type buffGet struct {
	store model.BuffStore
}

func (b *buffGet) ServeCodec(c apiutils.Codec, w http.ResponseWriter, r *http.Request) {
	bID, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		c.Respond(r.Context(), w, http.StatusBadRequest, err)
		return
	}

	buff, err := b.store.GetBuff(model.BuffID(bID))
	if err != nil {
		if err == model.ErrNotFound {
			c.Respond(r.Context(), w, http.StatusNotFound, err)
			return
		}
		c.Respond(r.Context(), w, http.StatusInternalServerError, err)
		return
	}
	c.Respond(r.Context(), w, http.StatusOK, types.NewBuff(*buff))
}
