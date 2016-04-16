package main

import (
	"flag"
	"log"
	"net"

	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	sched "github.com/mesos/mesos-go/scheduler"
)

type Myscheduler struct {
}

func newMyScheduler() *Myscheduler {
	return &Myscheduler{}
}

func (mysched *Myscheduler) Registered(_ sched.SchedulerDriver, frameworkId *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	log.Printf("scheduler register mesos with framework id %s, master id %s", frameworkId.GetValue(), masterInfo.GetId())
}

func (mysched *Myscheduler) Reregistered(_ sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	log.Printf("scheduler reregister mesos ", masterInfo)
}

func (mysched *Myscheduler) Disconnected(driver sched.SchedulerDriver) {
	log.Printf("scheduler Disconnect with master")
}

func (mysched *Myscheduler) ResourceOffers(driver sched.SchedulerDriver, offers []*mesos.Offer) {
	log.Printf("receiced %d offers", len(offers))
	for _, offer := range offers {
		driver.DeclineOffer(offer.GetId(), nil)
		log.Printf("decline offer %s", offer.GetId().GetValue())
	}
}

func (mysched *Myscheduler) OfferRescinded(_ sched.SchedulerDriver, offerId *mesos.OfferID) {
	log.Printf("offer %s rescind", offerId)
}

func (mysched *Myscheduler) StatusUpdate(_ sched.SchedulerDriver, status *mesos.TaskStatus) {
	if status.GetState() == mesos.TaskState_TASK_FAILED ||
		status.GetState() == mesos.TaskState_TASK_KILLED ||
		status.GetState() == mesos.TaskState_TASK_LOST {
		log.Printf("Task  go error\n")
	}

}

func (mysched *Myscheduler) FrameworkMessage(_ sched.SchedulerDriver, executorId *mesos.ExecutorID, slaveId *mesos.SlaveID, data string) {
	log.Printf("get framework info")
	log.Printf("using executor %s, in slave %s", executorId, slaveId)
}

func (mysched *Myscheduler) SlaveLost(_ sched.SchedulerDriver, slaveId *mesos.SlaveID) {
	log.Printf("slave is lost with ID %s", slaveId)
}

func (mysched *Myscheduler) ExecutorLost(_ sched.SchedulerDriver, executorId *mesos.ExecutorID, slaveId *mesos.SlaveID, status int) {
	log.Printf("executor is lost with id %s, in slave slaveId", executorId, slaveId)
}

func (mysched *Myscheduler) Error(_ sched.SchedulerDriver, message string) {
	log.Printf("Get error, ", message)
}

var (
	master  = flag.String("master", "IP:PORT", "master address")
	address = flag.String("address", "localip", "address")
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
		Scheduler: newMyScheduler(),
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
