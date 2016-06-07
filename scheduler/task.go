package scheduler

import (
	"github.com/gogo/protobuf/proto"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/mesosutil"
	"my_framework/types"
	"strconv"
	"time"
)

func NewMyTask(req *types.TaskRequest, taskTotal int, taskFrontEnd string) *types.MyTask {
	taskBackEnd := strconv.FormatInt(time.Now().Unix(), 8)
	return &types.MyTask{
		ID:       taskFrontEnd + taskBackEnd + strconv.Itoa(taskTotal),
		TaskCpu:  req.Cpu,
		TaskMem:  req.Mem,
		Cmd:      req.Cmd,
		Env:      req.Env,
		Image:    req.Image,
		Name:     req.Name,
		Count:    req.Count,
		Hostname: req.Hostname,
		Status:   mesos.TaskState_TASK_STAGING.String(),
	}
}

func CreateTaskInfo(offer *mesos.Offer, task *types.MyTask) (*types.MyTask, *mesos.TaskInfo) {
	//	task.Hostname = *offer.Hostname
	task.SlaveId = offer.SlaveId.GetValue()
	task.FrameworkId = offer.FrameworkId.GetValue()

	taskId := &mesos.TaskID{
		Value: proto.String(task.ID),
	}

	var env []*mesos.Environment_Variable
	for k, v := range task.Env {
		env = append(env, &mesos.Environment_Variable{
			Name:  proto.String(k),
			Value: proto.String(v),
		})
	}

	commandInfo := &mesos.CommandInfo{
		Environment: &mesos.Environment{
			Variables: env,
		},
	}

	if task.Cmd == "" {
		commandInfo.Shell = proto.Bool(false)
	} else {
		commandInfo.Value = &task.Cmd
	}

	taskInfo := &mesos.TaskInfo{
		Name:    proto.String(task.Name),
		TaskId:  taskId,
		SlaveId: offer.SlaveId,
		Container: &mesos.ContainerInfo{
			Type: mesos.ContainerInfo_DOCKER.Enum(),
			Docker: &mesos.ContainerInfo_DockerInfo{
				Image: proto.String(task.Name),
			},
		},
		Command: commandInfo,
		Resources: []*mesos.Resource{
			mesosutil.NewScalarResource("cpus", task.TaskCpu),
			mesosutil.NewScalarResource("mem", task.TaskMem),
		},
	}
	return task, taskInfo
}
