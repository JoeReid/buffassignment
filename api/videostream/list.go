package videostream

import (
	"net/http"

	"github.com/JoeReid/apiutils"
	"github.com/JoeReid/buffassignment/api/types"
	"github.com/JoeReid/buffassignment/internal/model"
)

func NewListHandler(store model.VideoStreamStore) apiutils.Handler {
	return &streamList{store}
}

type streamList struct {
	store model.VideoStreamStore
}

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
