package route

import (
	"github.com/emicklei/go-restful"
	"net/http"
)

type Hander struct {
	addr      string
	container *restful.Container
}

func NewHander(addr string) *Hander {
	return &Hander{
		addr:      addr,
		container: restful.NewContainer(),
	}
}

func (h *Hander) HanderServer() {
	log.Infoln("start to listen on ", addr)
	server := &http.Server{
		Addr:    h.addr,
		Handler: h.container,
	}
	log.Fatal(server.ListenAndServe())
}
