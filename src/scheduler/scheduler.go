package scheduler

import (
	"flag"
	"log"
	"net"

	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/mesosutil"
	sched "github.com/mesos/mesos-go/scheduler"
)

type Myscheduler struct {
	imageName   string
	executor    *mesos.ExecutorInfo
	taskRunning []string
	taskPending []string
	taskTotal   int
	cpuPerTask  float64
	memPerTask  float64
}

func NewMyScheduler(imageName string, cpus, mem float64) *Myscheduler {
	return &Myscheduler{
		imageName:  imageName,
		executor:   &mesos.ExecutorInfo{},
		taskTotal:  0,
		cpuPerTask: cpus,
		memPerTask: mem,
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
		remainingCpus := getOfferCpu(offer)
		remainingMem := getOfferMem(off)

		var tasks []*mesos.TaskInfo
		for mysched.cpuPerTask <= remainingCpus &&
			mysched.memPerTask <= remainingMem &&
			len(mysched.taskRunning)+len(mysched.taskPending) < mysched.taskTotal {
			taskId := &mesos.TaskID{
				Value: proto.String( /*id*/),
			}

			taskPending = append(taskPending, taskId.GetValue())

			//set dockerInfo and containerinfo
			docker := &mesos.DockerInfo{
				Image: proto.String(mysched.imageName),
			}

			container := &mesos.ContainerInfo{
				Type:   mesos.DOCKER.enum(),
				Docker: docker,
			}

			command := &mesos.CommandInfo{
				Shell: proto.String("false"),
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
			RefuseSeconds: proto.float64(1),
		}
		driver.LaunchTasks(offer.Id, tasks, filter)
	}
}

func (mysched *Myscheduler) OfferRescinded(_ sched.SchedulerDriver, offerId *mesos.OfferID) {
	log.Printf("offer %s rescind", offerId)
}

func (mysched *Myscheduler) StatusUpdate(_ sched.SchedulerDriver, status *mesos.TaskStatus) {
	taskStatus := status.GetState()

	if taskStatus == mesos.TaskState_Task_KILLED ||
		taskStatus == mesos.TaskState_Task_LOST ||
		taskStatus == mesos.TaskState_Task_FAILED {
		log.Infof("Abort Task %v in state %v with message %v\n", status.TaskId.GetValue(), status.State.String(), status.GetMessage())
	}

	switch taskStatus {
	case mesos.TaskState_Task_RUNNING:
	case mesos.TaskState_Task_FINISHED:
		mysched.taskFinished++
		log.Infof("task %v finished\n", status.TaskId.GetValue())
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
