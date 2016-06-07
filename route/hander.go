package route

import (
	"github.com/emicklei/go-restful"
	_ "github.com/golang/glog"
	"log"
	"net/http"
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
	log.Println("start to listen on ", h.Addr)
	server := &http.Server{
		Addr:    h.Addr,
		Handler: h.Container,
	}
	log.Println(server.ListenAndServe())
}
