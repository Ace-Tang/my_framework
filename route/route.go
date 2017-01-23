package route

import (
	"log"
	sched "my_framework/scheduler"
	"my_framework/store" //has mysql
	"my_framework/types"
	"net/http"

	restful "github.com/emicklei/go-restful"
	_ "github.com/emicklei/go-restful/swagger"
	_ "github.com/golang/glog"
)

type RouteManage struct {
	db    *store.Storage
	sched *sched.Myscheduler //Scheduler.go already has this package, so we use interface
}

func NewRouteManage(db *store.Storage, mysched *sched.Myscheduler) *RouteManage {
	return &RouteManage{
		db:    db,
		sched: mysched,
	}
}

func (r *RouteManage) ContainerRegister(h *Hander) {
	r.RegisterTask(h.Container)
	r.RegisterSlave(h.Container)
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
		Param(ws.PathParameter("hostname", "identifier of the slave node").DataType("string")))

	ws.Route(ws.GET("/{name}").To(r.descTask).
		Operation("descTask").
		Param(ws.PathParameter("name", "identifier of task").DataType("string")))

	ws.Route(ws.DELETE("/{id}").To(r.rmTask).
		Operation("rmTask").
		Param(ws.PathParameter("id", "identifier of task").DataType("string")))

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

	r.sched.ScheduleTask(&taskReq)

}

type ls_tasks struct {
	Tasks []*types.MyTask `json:"tasks"`
}

func (r *RouteManage) lsTask(req *restful.Request, resp *restful.Response) {
	hostname := req.QueryParameter("hostname")

	ts := []*types.MyTask{}
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
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
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
	r.sched.KillTask(id)
	r.db.RmTask(id)
	//tell scheduler
}

func (r *RouteManage) RegisterSlave(c *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/slaves").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	ws.Route(ws.POST("").To(r.addSlaves).
		Operation("addSlaves").
		Reads(types.SlaveNode{}))

	ws.Route(ws.GET("").To(r.lsSlaves).
		Operation("lsSlaves").
		Param(ws.PathParameter("hostname", "identifier of the slave node").DataType("string")))

	ws.Route(ws.DELETE("/{hostname}").To(r.rmSlave).
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

	log.Println("add node to my_framework ", args)
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
