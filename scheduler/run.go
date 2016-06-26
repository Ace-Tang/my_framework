package scheduler

import (
	"github.com/gogo/protobuf/proto"
	_ "github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
	"log"
	"my_framework/types"
	"net"
	"os"
)

func RunScheduler(mysched *Myscheduler, cfg *types.Config) {
	bindingAddress := net.ParseIP(cfg.Address)
	if bindingAddress == nil {
		log.Println("Error parse IP from ", cfg.Address)
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
		log.Println("Create Scheduler ", err)
		os.Exit(-1)
	}

	go func() {
		<-mysched.shutdown
		driver.Stop(true)
	}()

	if stat, err := driver.Run(); err != nil {
		log.Println("Framework stop with Status %s and err %s\n", stat.String(), err.Error())
	}
}
