package main

import (
	"flag"
	"log"
	"net"

	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
	scheduler "scheduler/sched"
)

var (
	master  = flag.String("master", "IP:PORT", "master address")
	address = flag.String("address", "localip", "address")
	image   = flag.String("docker-image", "", "docker image addr")
)

func main() {
	flag.Parse()
	addr, err := net.LookupIP(*address)
	if err != nil {
		log.Printf("parse IP get error", err)
	}

	if len(addr) < 1 {
		log.Printf("fail to Parse IP %v", address)
	}
	bindingAddress := addr[0]

	log.Printf("starting my framework")
	driver, err := sched.NewMesosSchedulerDriver(sched.DriverConfig{
		Scheduler: scheduler.NewMyScheduler(*image, 0.1, 32.0),
		Framework: &mesos.FrameworkInfo{
			Name: proto.String("april"),
			User: proto.String("root"),
		},
		Master:         *master,
		BindingAddress: bindingAddress,
		BindingPort:    bindingPort,
	})

	if err != nil {
		log.Printf("create scheduler error", err.Error())
		return
	}

	log.Printf("complete myscheduler create")

	if stat, err := driver.Run(); err != nil {
		log.Printf("Framwork stopped with status %s and error: %s\n", stat.String(), err.Error())
	}

	log.Printf("framework terminating")

}
