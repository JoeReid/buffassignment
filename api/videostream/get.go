package videostream

import (
	"net/http"

	"github.com/JoeReid/apiutils"
	"github.com/JoeReid/buffassignment/api/types"
	"github.com/JoeReid/buffassignment/internal/model"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func NewGetHandler(store model.VideoStreamStore) apiutils.Handler {
	return &streamGet{store}
}

type streamGet struct {
	store model.VideoStreamStore
}

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
