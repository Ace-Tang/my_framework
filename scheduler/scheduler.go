package scheduler

import (
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
	_ "github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
	"log"
	store "my_framework/store"
	"my_framework/types"
)

var (
	filter = &mesos.Filters{
		RefuseSeconds: proto.Float64(1),
	}
)

type Myscheduler struct {
	taskTotal    int
	taskFrontEnd string
	start        chan string
	//	shutdown chan
	db *store.Storage
}

func NewMyScheduler(db *store.Storage) *Myscheduler {
	return &Myscheduler{
		taskTotal:    0,
		taskFrontEnd: strconv.FormatInt(time.Now().Unix(), 32),
		start:        make(chan string, 1000),
		db:           db,
		//shutdown:
	}
}

func (mysched *Myscheduler) Registered(_ sched.SchedulerDriver, frameworkId *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	log.Printf("scheduler register mesos with framework id %s, master id %s", frameworkId.GetValue(), masterInfo.GetId())
}

func (mysched *Myscheduler) Reregistered(_ sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	log.Println("scheduler reregister mesos ", masterInfo)
}

func (mysched *Myscheduler) Disconnected(driver sched.SchedulerDriver) {
	log.Println("scheduler Disconnect with master")
}

func (mysched *Myscheduler) ResourceOffers(driver sched.SchedulerDriver, offers []*mesos.Offer) {
	log.Printf("receive %d offers\n", len(offers))

loop:
	for len(offers) > 0 {

		var tasks []*mesos.TaskInfo

		select {
		case tid := <-mysched.start:
			task, err := mysched.db.GetTask(tid)
			if err != nil {
				log.Printf("unable to find task %v in db\n", tid)
			}
			log.Printf("try to launching task %s on slave %s", tid, task.Hostname)

			var ok int = 0
			var offer *mesos.Offer
			for _, offer = range offers {
				remainingCpus := getOfferRes("cpus", offer)
				remainingMem := getOfferRes("mem", offer)
				log.Printf("receiced offer %v with cpus %v, memory %v\n", offer.Id.GetValue(), remainingCpus, remainingMem)

				if *(offer.Hostname) == task.Hostname &&
					int(remainingCpus) > int(task.TaskCpu)*task.Count &&
					int(remainingMem) > int(task.TaskMem)*task.Count {
					ok = 1
					break
				}
			}

			if ok == 0 {
				log.Printf("task %s cannot launing on slave %s\n", task.ID, task.Hostname)
				go func() {
					mysched.start <- tid
				}()
				break loop
			}

			t, taskInfo := CreateTaskInfo(offer, task)
			log.Printf("%d task %s pending\n", t.Count, t.ID)
			mysched.db.UpdateTask(t.ID, t.SlaveId, t.FrameworkId)
			//add in store, queue
			for i := 0; i < task.Count; i++ {
				tasks = append(tasks, taskInfo)
			}
			log.Printf("Launching task %s for offer %v\n", task.ID, offer.Id.GetValue())
			driver.LaunchTasks([]*mesos.OfferID{offer.Id}, tasks, filter)
			mysched.taskTotal++

			continue
		default:
			break loop
		}

	}
	log.Println("no task run, decline offers")
	for _, offer := range offers {
		driver.DeclineOffer(offer.Id, filter)
	}
}

func (mysched *Myscheduler) OfferRescinded(_ sched.SchedulerDriver, offerId *mesos.OfferID) {
	log.Printf("offer %s rescind", offerId)
}

func (mysched *Myscheduler) StatusUpdate(_ sched.SchedulerDriver, status *mesos.TaskStatus) {
	taskStatus := status.GetState()
	id := status.TaskId.GetValue()

	_, err := mysched.db.GetTask(id)
	if err != nil {
		log.Printf("unable to find task %v in db\n", id)
	}

	state := status.State.String()
	if taskStatus == mesos.TaskState_TASK_KILLED ||
		taskStatus == mesos.TaskState_TASK_LOST ||
		taskStatus == mesos.TaskState_TASK_FAILED {
		log.Printf("Abort Task %v in state %v with message %v\n", status.TaskId.GetValue(), state, status.GetMessage())
	}

	switch taskStatus {
	case mesos.TaskState_TASK_RUNNING:
		//mysched.taskPending--
		//mysched.taskRunning++
		log.Printf("task %v running\n", status.TaskId.GetValue())
	case mesos.TaskState_TASK_FINISHED:
		//mysched.taskTotal++
		//mysched.taskRunning--
		log.Printf("task %v finished\n", status.TaskId.GetValue())
	}
	mysched.db.UpdateTaskStatus(id, state)
}

func (mysched *Myscheduler) FrameworkMessage(_ sched.SchedulerDriver, executorId *mesos.ExecutorID, slaveId *mesos.SlaveID, data string) {
	log.Println("get framework info")
	log.Printf("using executor %s, in slave %s\n", executorId, slaveId)
}

func (mysched *Myscheduler) SlaveLost(_ sched.SchedulerDriver, slaveId *mesos.SlaveID) {
	log.Printf("slave is lost with ID %s", slaveId)
}

func (mysched *Myscheduler) ExecutorLost(_ sched.SchedulerDriver, executorId *mesos.ExecutorID, slaveId *mesos.SlaveID, status int) {
	log.Printf("executor is lost with id %s, in slave slaveId", executorId, slaveId)
}

func (mysched *Myscheduler) Error(_ sched.SchedulerDriver, message string) {
	log.Println("Get error, ", message)
}

func getOfferRes(name string, offer *mesos.Offer) float64 {
	ress := mesosutil.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
		return res.GetName() == name
	})

	val := 0.0
	for _, res := range ress {
		val += res.GetScalar().GetValue()
	}

	return val
}

func (mysched *Myscheduler) ScheduleTask(req *types.TaskRequest) {
	task := NewMyTask(req, mysched.taskTotal, mysched.taskFrontEnd)
	log.Printf("receive task %s , add %s to queue\n", task.Name, task.ID)

	select {
	case mysched.start <- task.ID:
		mysched.db.PutTask(task)
		//store,queue
		//return task.ID, nil
	}
}

func (mysched *Myscheduler) Stop() {

}
