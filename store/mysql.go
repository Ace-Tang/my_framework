package store

import (
	"database/sql"
	_ "errors"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang/glog"
	"log"
	"my_framework/types"
)

type Storage struct {
	db   *sql.DB
	addr string
}

func NewStorage(addr string) *Storage {
	return &Storage{
		db:   nil,
		addr: addr,
	}
}

func (s *Storage) Open() error {
	if s.db != nil {
		return nil
	}

	var err error
	s.db, err = sql.Open("mysql", s.addr)
	log.Println("in open , ", s.db)
	if err != nil {
		s.db = nil
	}
	return err
}

func (s *Storage) Close() error {
	if s.db != nil {
		err := s.db.Close()
		return err
	}
	return nil
}

func (s *Storage) AddNode(n *types.SlaveNode) {
	sql := "insert into `slave_info` (`hostname`,`attachment`) values (?,?);"
	_, err := s.db.Exec(sql, n.Hostname, n.Attachment)

	//stmt, err := s.db.Prepare(`INSERT slave_info (hostname, attachment) values (?,?)`)
	//_, err = stmt.Exec(hostname, attachment)

	if err != nil {
		log.Println("insert to table slave_info : ", err)
	}
}

func (s *Storage) LsNode() []*types.SlaveNode {
	sql := "select * from `slave_info`;"

	rows, err := s.db.Query(sql)
	if err != nil {
		log.Println("select from db : ", err)
	}
	nodes := []*types.SlaveNode{}
	for rows.Next() {
		node := &types.SlaveNode{}
		rows.Scan(&node.Hostname, &node.Attachment)
		nodes = append(nodes, node)
	}

	return nodes
}

func (s *Storage) DescNode(hostname string) *types.SlaveNode {
	sql := "select `hostname`, `attachment` from `slave_info` where `hostname`=?;"

	rows, err := s.db.Query(sql, hostname)
	if err != nil {
		log.Println("select slave_info from db ", err)
		return nil
	}

	node := &types.SlaveNode{}
	for rows.Next() {
		rows.Scan(&node.Hostname, &node.Attachment)
	}

	return node
}

func (s *Storage) RmNode(hostname string) {
	sql := "delete from `slave_info` where `hostname`=?;"
	_, err := s.db.Exec(sql, hostname)
	if err != nil {
		log.Println("delete from db : ", err)
	}
}

func (s *Storage) PutTask(t *types.MyTask) {
	//env will be array k:v
	sql := "insert into `task_info` (`task_cpu`, `task_mem`, `id`, `cmd`, `env`, `image`, `hostname`, `name`, `status`, `count`) values (?,?,?,?,?,?,?,?,?,?,);"

	_, err := s.db.Exec(sql, t.TaskCpu, t.TaskMem, t.ID, t.Cmd, t.Image, t.SlaveId, t.Hostname, t.Name, t.FrameworkId, t.Status, t.Count)
	if err != nil {
		log.Println("insert to table slave_info : ", err)
	}
}

func (s *Storage) UpdateTask(id, slave_id, framework_id string) {
	//sql := "update `task_info` set `slave_id`=" + slave_id + " and `framework_id`=" + framework_id + " where `id`=" + id + ";"
}

func (s *Storage) UpdateTaskStatus(id, status string) {
	//	sql := "update `task_info` set `status` =" + status + " where `id`=" + id + ";"
}

func (s *Storage) GetTask(id string) (*types.MyTask, error) {
	sql := "select `id`, `hostname`, `task_cpu`, `task_mem`, `count` from `task_info` where `id`=?;"
	rows, err := s.db.Query(sql, id)
	if err != nil {
		log.Println("GetTask from db ", err)
		return nil, err
	}
	t := &types.MyTask{}
	for rows.Next() {
		rows.Scan(&t.ID, &t.Hostname, &t.TaskCpu, &t.TaskMem, &t.Count)
	}

	return t, nil
}

func (s *Storage) ListAllTask() []*types.MyTask {
	sql := "select `id`, `name`, `hostname` from `task_info`;"
	rows, err := s.db.Query(sql)

	if err != nil {
		log.Println("select from task_info ", err)
		return nil
	}
	tasks := []*types.MyTask{}
	for rows.Next() {
		t := &types.MyTask{}
		rows.Scan(&t.ID, &t.Name, &t.Hostname)
		tasks = append(tasks, t)
	}
	return tasks
}

func (s *Storage) ListTask(hostname string) []*types.MyTask {
	sql := "select `id`, `name`, `hostname` from `task_info` where `hostname`=?;"
	rows, err := s.db.Query(sql, hostname)

	if err != nil {
		log.Println("select from task_info ", err)
		return nil
	}

	tasks := []*types.MyTask{}
	for rows.Next() {
		t := &types.MyTask{}
		rows.Scan(&t.ID, &t.Name, &t.Hostname)
		tasks = append(tasks, t)
	}
	return tasks
}

func (s *Storage) DescTask(name string) *types.MyTask {
	sql := "select `task_cpu`, `task_mem`, `cmd`, `env`, `image`, `slave_id`, `hostname`, `id`, `framework_id`, `status`, `count` from `task_info` where `name`=?;"

	rows, err := s.db.Query(sql, name)
	if err != nil {
		log.Println("select from db : ", err)
	}
	ts := &types.MyTask{}
	for rows.Next() {
		rows.Scan(&ts.TaskCpu, &ts.TaskMem, &ts.Cmd, &ts.Env, &ts.Image, &ts.SlaveId, &ts.Hostname, &ts.ID, &ts.FrameworkId, &ts.Status, &ts.Count)
	}

	return ts
}

func (s *Storage) RmTask(id string) {
	sql := "delete from `task_info` where `id`=?;"

	_, err := s.db.Exec(sql, id)
	if err != nil {
		log.Printf("delete task %s :%v\n", id, err)
	}
}
