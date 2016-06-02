package scheduler

import (
	"github.com/gogo/protobuf/proto"
	log "github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
	"my_framework/types"
	"net"
	"os"
)

func RunScheduler(mysched *Myscheduler, cfg *types.Config) {
	bindingAddress := net.ParseIP(cfg.Address)
	if bindingAddress == nil {
		log.Errorln("Error parse IP from ", cfg.Address)
		os.Exit(-1)
	}
	driver, err := sched.NewMesosSchedulerDriver(sched.DriverConfig{
		Scheduler: mysched,
		Framework: &mesos.FrameworkInfo{
			Name: proto.String(cfg.Name),
			User: proto.String(cfg.User),
		},
		Master:         cfg.Master,
		BindingAddress: bindingAddress,
		BindingPort:    cfg.Port,
	})

	if err != nil {
		log.Errorln("Create Scheduler ", err)
		os.Exit(-1)
	}

	//go func() {
	//	<-mysched.shutdown
	//	driver.Stop(false)
	//}()

	if stat, err := driver.Run(); err != nil {
		log.Errorf("Framework stop with Status %s and err %s\n", stat.String(), err.Error())
	}
}
