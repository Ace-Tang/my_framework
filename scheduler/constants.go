package scheduler

import "github.com/mesos/mesos-go/mesosproto"

const (
	ContainerTaskState_PENDING  int = 0
	ContainerTaskState_RUNNING  int = 1
	ContainerTaskState_FINISHED int = 2
	ContainerTaskState_FAILED   int = 3
	ContainerTaskState_KILL     int = 4
	ContainerTaskState_LOST     int = 5
	ContainerTaskState_STAGING  int = 6
	ContainerTaskState_ERROR    int = 7 //define cannot be launched
)

var ContainerState_name = map[int]string{
	0: "ContainerTaskState_PENDING",
	1: "ContainerTaskState_RUNNING",
	2: "ContainerTaskState_FINISHED",
	3: "ContainerTaskState_FAILED",
	4: "ContainerTaskState_KILL",
	5: "ContainerTaskState_LOST",
	6: "ContainerTaskState_ERROR",
}

func getContainerState(st int) string {
	return ContainerState_name[st]
}

func converTaskStatus(stat mesosproto.TaskState) int {
	mesos2local := []int{
		ContainerTaskState_PENDING,
		ContainerTaskState_RUNNING,
		ContainerTaskState_FINISHED,
		ContainerTaskState_FAILED,
		ContainerTaskState_KILL,
		ContainerTaskState_LOST,
		ContainerTaskState_STAGING,
		ContainerTaskState_ERROR,
	}

	return mesos2local[stat]
}
