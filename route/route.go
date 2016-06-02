package route

import (
	"github.com/emicklei/go-restful"
	_ "github.com/emicklei/go-restful/swagger"
	log "github.com/golang/glog"
	sched "my_framework/scheduler"
	"my_framework/store" //has mysql
	"my_framework/types"
	"net/http"
)

type RouteManage struct {
	db *store.Storage
}

func NewRouteManage(db *store.Storage) *RouteManage {
	return &RouteManage{
		db: db,
	}
}

func (r *RouteManage) ContainerRegister(h Hander) {
	r.RegisterTask(h.container)
	r.RegisterSlave(h.container)
}

func (r *RouteManage) RegisterTask(c *restful.Container) {
	ws := new(restful.WebService)

	ws.
		Path("/task").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("").To(r.addTask).
		Operation("addTask").
		Reads(types.TaskRequest{}))

	ws.Route(ws.GET("").To(r.lsTask).
		Operation("lsTask").
		Param(ws.PathParameter("hostname", "identifier of the slave node")).DataType("string"))

	ws.Route(ws.GET("/{name}").To(r.descTask).
		Operation("descTask").
		Param(ws.PathParameter("name", "identifier of task")).DataType("string"))

	ws.Route(ws.DELETE("/{id}").To(r.rmTask).
		Operation("rmTask").
		Param(PathParameter("id", "identifier of task")).DataType("string"))

	c.Add(ws)
}

func (r *RouteManage) addTask(req *restful.Request, resp *restful.Response) {
	taskReq := types.TaskRequest{}
	err := req.ReadEntity(&taskReq)
	if err != nil {
		resp.AddHeader("Context-Type", "text/plain")
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	sched.ScheduleTask(&taskReq)

}

type ls_tasks struct {
	Tasks []*types.MyTask `json:"tasks"`
}

func (r *RouteManage) lsTask(req *restful.Request, resp *restful.Response) {
	hostname := req.QueryParameter("hostname")

	ts := &types.SlaveNode{}
	if len(hostname) == 0 {
		ts = r.db.ListTask(hostname)

	} else {
		ts = r.db.ListAllTask()
	}
	err := resp.WriteAsJson(&ls_tasks{
		Tasks: ts,
	})

	if err != nil {
		resp.AddHeader("Context-Type", "text/plain")
		resp.WriteErrorString(http, StatusInternalServerError, err.Error())
		return
	}
}

func (r *RouteManage) descTask(req *restful.Request, resp *restful.Response) {
	name := req.QueryParameter("name")
	t := r.db.DescTask(name)

	err := resp.WriteAsJson(t)

	if err != nil {
		resp.AddHeader("Context-Type", "text/plain")
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	//flush
}

func (r *RouteManage) rmTask(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	r.db.RmTask(id)
	//tell scheduler
}

func (r *RouteManage) RegisterSlave(c *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/slaves").
		Consumes(rrestful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("").To(s.addSlaves).
		Operation("addSlaves").
		Reads(add_slave{}))

	ws.Route(ws.GET("").To(s.lsSlaves).
		Operation("lsSlaves").
		Param(ws.PathParameter("hostname", "identifier of the slave node").DataType("string")))

	ws.Route(ws.DELETE("/{hostname}").To(s.rmSlave).
		Operation("rmSlave").
		Param(ws.PathParameter("hostname", "identifier of the slave node").DataType("string")))

	c.Add(ws)
}

func (r *RouteManage) addSlaves(req *restful.Request, resp *restful.Response) {
	args := types.SlaveNode{}
	err := req.ReadEntity(&args)
	if err != nil {
		resp.AddHeader("Context-Type", "text/plain")
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	log.Infoln("add node to my_framework ", args)
	r.db.AddNode(&args)
}

type ls_slaves struct {
	Slaves []*types.SlaveNode `json:"slaves"`
}

func (r *RouteManage) lsSlaves(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("hostname")
	slaves := []*types.SlaveNode{}
	if len(id) == 0 {
		slaves = r.db.LsNode()
	} else {
		node := r.db.DescNode(id)
		slaves = append(slaves, node)
	}

	err := resp.WriteAsJson(&ls_slaves{
		Slaves: slaves,
	})

	if err != nil {
		resp.AddHeader("Context-type", "text/plain")
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
}

func (r *RouteManage) rmSlave(req *restful.Request, resp *restful.Response) {
	id := req.PathParameter("hostname")
	r.db.RmNode(id)
}
