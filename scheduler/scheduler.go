package scheduler

import (
	"strconv"
	"time"

	store "my_framework/store"
	"my_framework/types"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
)

var (
	filter = &mesos.Filters{
		RefuseSeconds: proto.Float64(10),
	}
)

type Myscheduler struct {
	mesosDriver  sched.SchedulerDriver
	taskSet      map[string]*types.MyTask // map[TaskId]Hostname
	taskTotal    int
	taskFrontEnd string
	start        chan string
	shutdown     chan struct{}
	db           *store.Storage
}

func NewMyScheduler(db *store.Storage) *Myscheduler {
	return &Myscheduler{
		taskTotal:    0,
		taskFrontEnd: strconv.FormatInt(time.Now().Unix(), 32),
		start:        make(chan string, 1000),
		db:           db,
		shutdown:     make(chan struct{}),
	}
}

func (mysched *Myscheduler) Registered(driver sched.SchedulerDriver, frameworkId *mesos.FrameworkID, masterInfo *mesos.MasterInfo) {
	glog.Infof("scheduler register mesos with framework id %s, master id %s", frameworkId.GetValue(), masterInfo.GetId())
	mysched.mesosDriver = driver
}

func (mysched *Myscheduler) Reregistered(driver sched.SchedulerDriver, masterInfo *mesos.MasterInfo) {
	glog.Infoln("scheduler reregister mesos ", masterInfo)
	mysched.mesosDriver = driver
}

func (mysched *Myscheduler) Disconnected(driver sched.SchedulerDriver) {
	glog.Infoln("scheduler Disconnect with master")
}

func (mysched *Myscheduler) ResourceOffers(driver sched.SchedulerDriver, offers []*mesos.Offer) {
	glog.Infof("receive %d offers\n", len(offers))

loop:
	time.Sleep(20 * time.Second)
	for len(offers) > 0 {

		var tasks []*mesos.TaskInfo

		select {
		case <-mysched.shutdown:
			glog.Infof("stop framework ....\n")
			break loop
		case tid := <-mysched.start:
			task, err := mysched.db.GetTask(tid)
			if err != nil {
				glog.Infof("unable to find task %v in db\n", tid)
			}
			glog.Infof("try to launch task %s on slave %s", tid, task.Hostname)

			var ok int = 0
			var offer *mesos.Offer
			for _, offer = range offers {
				remainingCpus := getOfferRes("cpus", offer)
				remainingMem := getOfferRes("mem", offer)
				glog.Infof("receiced offer %v with cpus %v, memory %v\n", offer.GetHostname(), remainingCpus, remainingMem)

				if *(offer.Hostname) == task.Hostname &&
					int(remainingCpus) > int(task.TaskCpu)*task.Count &&
					int(remainingMem) > int(task.TaskMem)*task.Count {
					ok = 1
					break
				}
			}

			if ok == 0 {
				glog.Infof("task %s cannot launing on slave %s\n", task.ID, task.Hostname)
				go func() {
					mysched.start <- tid
				}()
				break loop
			}

			t, taskInfo := CreateTaskInfo(offer, task)
			glog.Infof("%d task %s pending\n", t.Count, t.ID)
			mysched.db.UpdateTask(t.ID, t.SlaveId, t.FrameworkId)
			//add in store, queue
			for i := 0; i < task.Count; i++ {
				tasks = append(tasks, taskInfo)
			}
			glog.Infof("Launching task %s for offer %v\n", task.ID, offer.Id.GetValue())
			driver.LaunchTasks([]*mesos.OfferID{offer.Id}, tasks, filter)
			mysched.taskTotal++

			continue
		default:
			break loop
		}

	}
	glog.Infoln("no task run, decline offers")
	for _, offer := range offers {
		driver.DeclineOffer(offer.Id, filter)
	}
}

func (mysched *Myscheduler) OfferRescinded(_ sched.SchedulerDriver, offerId *mesos.OfferID) {
	glog.Infof("offer %s rescind", offerId)
}

func (mysched *Myscheduler) StatusUpdate(_ sched.SchedulerDriver, status *mesos.TaskStatus) {
	taskStatus := status.GetState()
	id := status.TaskId.GetValue()

	t, exist := mysched.taskSet[id]
	if !exist {
		glog.Infof("task %v not in memory\n", id)
		return
	}
	t.Status = converTaskStatus(status.GetState())

	state := status.State.String()
	if taskStatus == mesos.TaskState_TASK_KILLED ||
		taskStatus == mesos.TaskState_TASK_LOST ||
		taskStatus == mesos.TaskState_TASK_FAILED {
		glog.Infof("Abort Task %v in state %v with message %v\n", status.TaskId.GetValue(), state, status.GetMessage())
	}

	switch taskStatus {
	case mesos.TaskState_TASK_RUNNING:
		//mysched.taskPending--
		//mysched.taskRunning++
		glog.Infof("task %v running\n", status.TaskId.GetValue())
	case mesos.TaskState_TASK_FINISHED:
		//mysched.taskTotal++
		//mysched.taskRunning--
		glog.Infof("task %v finished\n", status.TaskId.GetValue())
	}
	mysched.db.UpdateTaskStatus(id, state)
}

func (mysched *Myscheduler) FrameworkMessage(_ sched.SchedulerDriver, executorId *mesos.ExecutorID, slaveId *mesos.SlaveID, data string) {
	glog.Infoln("get framework info")
	glog.Infof("using executor %s, in slave %s\n", executorId, slaveId)
}

func (mysched *Myscheduler) SlaveLost(_ sched.SchedulerDriver, slaveId *mesos.SlaveID) {
	glog.Infof("slave is lost with ID %s", slaveId)
}

func (mysched *Myscheduler) ExecutorLost(_ sched.SchedulerDriver, executorId *mesos.ExecutorID, slaveId *mesos.SlaveID, status int) {
	glog.Infof("executor is lost with id %s, in slave slaveId", executorId, slaveId)
}

func (mysched *Myscheduler) Error(_ sched.SchedulerDriver, message string) {
	glog.Infoln("Get error, ", message)
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
	glog.Infof("receive task %s , add %s to queue\n", task.Name, task.ID)

	select {
	case mysched.start <- task.ID:
		mysched.db.PutTask(task)
		//store,queue
		//return task.ID, nil
	}
}

func (mysched *Myscheduler) KillTask(taskId string) {
	t, exist := mysched.taskSet[taskId]
	if !exist {
		glog.Infof("task %s not exist in memory\n", taskId)
		return
	}

	if t.Status != ContainerTaskState_RUNNING {
		glog.Infof("task %s not running, do not need kill\n", taskId)
		return
	}
	_, err := mysched.mesosDriver.KillTask(mesosutil.NewTaskID(taskId))
	if err != nil {
		glog.Infof("kill task %s error %v\n", taskId, err)
	}
	glog.Infof("kill task %s successful\n", taskId)
}

func (mysched *Myscheduler) Stop() {
	close(mysched.shutdown)
}
