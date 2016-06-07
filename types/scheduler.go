package types

type Scheduler interface {
	ScheduleTask(req *TaskRequest)
}
