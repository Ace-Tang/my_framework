package scheduler

import (
	"log"
	"strconv"

	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
)

type Myscheduler struct {
	imageName   string
	executor    *mesos.ExecutorInfo
	taskRunning int
	taskPending int
	taskTotal   int
	cpuPerTask  float64
	memPerTask  float64
}

func NewMyScheduler(imageName string, cpus, mem float64) *Myscheduler {
	return &Myscheduler{
		imageName:   imageName,
		executor:    &mesos.ExecutorInfo{},
		taskRunning: 0,
		taskPending: 0,
		taskTotal:   0,
		cpuPerTask:  cpus,
		memPerTask:  mem,
	}
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
		remainingCpus := 0.0
		remainingMem := 0.0
		cpuRes := mesosutil.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
			return res.GetName() == "cpus"
		})
		for _, cpus := range cpuRes {
			remainingCpus += cpus.GetScalar().GetValue()
		}
		memRes := mesosutil.FilterResources(offer.Resources, func(res *mesos.Resource) bool {
			return res.GetName() == "mem"
		})
		for _, mem := range memRes {
			remainingMem += mem.GetScalar().GetValue()
		}

		var tasks []*mesos.TaskInfo
		for mysched.cpuPerTask <= remainingCpus &&
			mysched.memPerTask <= remainingMem &&
			mysched.taskRunning+mysched.taskPending < mysched.taskTotal {
			taskId := &mesos.TaskID{
				Value: proto.String(strconv.Itoa(mysched.taskRunning)),
			}

			mysched.taskPending++

			//set dockerInfo and containerinfo
			docker := &mesos.ContainerInfo_DockerInfo{
				Image: proto.String(mysched.imageName),
			}

			dockerTye := mesos.ContainerInfo_DOCKER
			container := &mesos.ContainerInfo{
				Type:   &dockerTye,
				Docker: docker,
			}

			command := &mesos.CommandInfo{
				Shell: proto.Bool(false),
			}

			log.Printf("task %s pending\n", taskId.GetValue())
			task := &mesos.TaskInfo{
				Name:     proto.String("ace-task" + taskId.GetValue()),
				TaskId:   taskId,
				SlaveId:  offer.SlaveId,
				Executor: mysched.executor,
				Resources: []*mesos.Resource{
					mesosutil.NewScalarResource("cpus", mysched.cpuPerTask),
					mesosutil.NewScalarResource("mem", mysched.memPerTask),
				},
				Container: container,
				Command:   command,
			}

			tasks = append(tasks, task)
			remainingCpus -= mysched.cpuPerTask
			remainingMem -= mysched.memPerTask

		}
		log.Printf("Launching %d task for offer %v\n", len(tasks), offer.Id.GetValue())
		filter := &mesos.Filters{
			RefuseSeconds: proto.Float64(1),
		}
		driver.LaunchTasks([]*mesos.OfferID{offer.Id}, tasks, filter)
	}
}

func (mysched *Myscheduler) OfferRescinded(_ sched.SchedulerDriver, offerId *mesos.OfferID) {
	log.Printf("offer %s rescind", offerId)
}

func (mysched *Myscheduler) StatusUpdate(_ sched.SchedulerDriver, status *mesos.TaskStatus) {
	taskStatus := status.GetState()

	if taskStatus == mesos.TaskState_TASK_KILLED ||
		taskStatus == mesos.TaskState_TASK_LOST ||
		taskStatus == mesos.TaskState_TASK_FAILED {
		log.Printf("Abort Task %v in state %v with message %v\n", status.TaskId.GetValue(), status.State.String(), status.GetMessage())
	}

	switch taskStatus {
	case mesos.TaskState_TASK_RUNNING:
		mysched.taskPending--
		mysched.taskRunning++
	case mesos.TaskState_TASK_FINISHED:
		mysched.taskTotal++
		log.Printf("task %v finished\n", status.TaskId.GetValue())
	}
}

func (mysched *Myscheduler) FrameworkMessage(_ sched.SchedulerDriver, executorId *mesos.ExecutorID, slaveId *mesos.SlaveID, data string) {
	log.Printf("get framework info\n")
	log.Printf("using executor %s, in slave %s\n", executorId, slaveId)
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

func getOfferCpu() {

}
