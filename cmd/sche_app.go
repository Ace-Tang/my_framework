package main

import (
	"flag"
	"log"
	"net"

	"github.com/gogo/protobuf/proto"
	//mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
	"my_framework/route"
	scheduler "my_framework/scheduler"
	"my_framework/store"
	"my_framework/types"
)

var (
	master  = flag.String("master", "IP:PORT", "master address")
	address = flag.String("address", "localip", "address")
	port    = flag.Uint("port", 10000, "framwork binding port")
	webui   = flag.String("webui", "localhost:8080", "webui address")
	name    = flag.String("name", "april", "framework name")
	user    = flag.String("user", "root", "who can use framework")
	db      = flag.String("db", "root@ace/mesos", "location of db")
)

func main() {
	flag.Parse()
	//addr, err := net.LookupIP(*address)
	//if err != nil {
	//	log.Printf("parse IP get error", err)
	//}

	//if len(addr) < 1 {
	//	log.Printf("fail to Parse IP %v", address)
	//}
	//bindingAddress := addr[0]

	//log.Printf("starting my framework")
	//driver, err := sched.NewMesosSchedulerDriver(sched.DriverConfig{
	//	Scheduler: scheduler.NewMyScheduler(*image, 0.1, 32.0),
	//	Framework: &mesos.FrameworkInfo{
	//		Name: proto.String("april"),
	//		User: proto.String("root"),
	//	},
	//	Master:         *master,
	//	BindingAddress: bindingAddress,
	//	BindingPort:    uint16(*port),
	//})

	//if err != nil {
	//	log.Printf("create scheduler error", err.Error())
	//	return
	//}

	//log.Printf("complete myscheduler create")

	//if stat, err := driver.Run(); err != nil {
	//	log.Printf("Framework stopped with status %s and error: %s\n", stat.String(), err.Error())
	//}

	//log.Printf("framework terminating")

	storage := store.NewStorage(*db)
	err := storage.Open()
	if err != nil {
		log.Errorln("open mysql ", err)
		os.Exit(-1)
	}
	log.Infoln("open database on ", *db)
	defer storage.Close()

	cfg := &types.Config{
		Master:  *master,
		Address: *address,
		Port:    uint16(*port),
		Webui:   *webui,
		Name:    *name,
		User:    *user,
		Db:      storage,
	}

	mysched := scheduler.NewMyScheduler()
	go func() {
		scheduler.RunScheduler(mysched, cfg)
		//close other
	}()

	//catch interrupt
	//go func() {
	//	c := make(chan os.Signal, 1)
	//	signal.Notify(c, os.Interrupt, os.Kill)
	//	s := <-c
	//	if s != os.Interrupt || s != os.Kill {
	//		return
	//	}
	//	scheduler.Stop()
	//}()

	hander := route.NewHander(*webui)
	routeM := route.NewRouteManage(storage)

	routeM.ContainerRegister(hander.container)
	hander.HanderServer()

}
