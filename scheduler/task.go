package scheduler

import (
	"my_framework/types"
	"strconv"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/golang/glog"
	mesos "github.com/mesos/mesos-go/mesosproto"
	"github.com/mesos/mesos-go/mesosutil"
)

func NewMyTask(req *types.TaskRequest, taskTotal int, taskFrontEnd string) *types.MyTask {
	taskBackEnd := strconv.FormatInt(time.Now().Unix(), 8)
	return &types.MyTask{
		ID:       taskFrontEnd + "-" + taskBackEnd + "-" + strconv.Itoa(taskTotal),
		TaskCpu:  req.Cpu,
		TaskMem:  req.Mem,
		Cmd:      req.Cmd,
		Env:      req.Env,
		Image:    req.Image,
		Name:     req.Name,
		Count:    req.Count,
		Hostname: req.Hostname,
	}
}

func CreateTaskInfo(offer *mesos.Offer, task *types.MyTask) (*types.MyTask, *mesos.TaskInfo) {
	//	task.Hostname = *offer.Hostname
	glog.Infof("task : %+v\n", task)
	task.SlaveId = offer.SlaveId.GetValue()
	task.FrameworkId = offer.FrameworkId.GetValue()

	taskId := &mesos.TaskID{
		Value: proto.String(task.ID),
	}

	var env []*mesos.Environment_Variable
	var commandInfo *mesos.CommandInfo
	if len(task.Env) > 0 {
		for k, v := range task.Env {
			env = append(env, &mesos.Environment_Variable{
				Name:  proto.String(k),
				Value: proto.String(v),
			})
		}

		commandInfo = &mesos.CommandInfo{
			Environment: &mesos.Environment{
				Variables: env,
			},
		}
	} else {
		commandInfo = &mesos.CommandInfo{}
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
				Image: proto.String(task.Image),
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
