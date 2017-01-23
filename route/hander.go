package route

import (
	"net/http"

	restful "github.com/emicklei/go-restful"
	"github.com/golang/glog"
)

type Hander struct {
	Addr      string
	Container *restful.Container
}

func NewHander(addr string) *Hander {
	return &Hander{
		Addr:      addr,
		Container: restful.NewContainer(),
	}
}

func (h *Hander) HanderServer() {
	glog.Infoln("start to listen on ", h.Addr)
	server := &http.Server{
		Addr:    h.Addr,
		Handler: h.Container,
	}
	glog.Fatalln(server.ListenAndServe())
}
