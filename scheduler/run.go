package scheduler

import (
	"my_framework/types"
	"net"
	"os"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
)

func RunScheduler(mysched *Myscheduler, cfg *types.Config) {
	bindingAddress := net.ParseIP(cfg.Address)
	if bindingAddress == nil {
		glog.Infoln("Error parse IP from ", cfg.Address)
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
		glog.Infoln("Create Scheduler ", err)
		os.Exit(-1)
	}

	go func() {
		<-mysched.shutdown
		driver.Stop(true)
	}()

	if stat, err := driver.Run(); err != nil {
		glog.Infoln("Framework stop with Status %s and err %s\n", stat.String(), err.Error())
	}
}
