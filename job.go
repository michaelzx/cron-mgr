package cronmgr

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	uuid "github.com/satori/go.uuid"
	"time"
)

type JobStatus uint8

const (
	JobStatusWait JobStatus = iota
	JobStatusRunning
)

// 执行顺序：
// before
// ...before
// run
// success/fail
// ...after
// after
type JobFunc func(thisJob *Job) error
type JobFailFunc func(thisJob *Job, jobErr error)
type JobSuccessFunc func(thisJob *Job)
type Job struct {
	EntityID       cron.EntryID
	ID             string    // uuid
	Desc           string    // 任务描述
	Status         JobStatus // 状态：1=待执行、2、正在执行
	NextTime       time.Time
	beforeFuncList []JobFunc      `json:"-"` // 前置逻辑
	runFunc        JobFunc        `json:"-"` // 运行逻辑
	successFunc    JobSuccessFunc `json:"-"`
	failFunc       JobFailFunc    `json:"-"`
	afterFuncList  []JobFunc      `json:"-"` // 后置逻辑
	mgr            *JobMgr
}

func NewJob(mgr *JobMgr, desc string, runFunc JobFunc) *Job {
	return &Job{
		ID:      uuid.NewV4().String(),
		Desc:    desc,
		Status:  JobStatusWait,
		runFunc: runFunc,
		mgr:     mgr,
	}
}

func (j *Job) GetEntity() cron.Entry {
	return j.mgr.cron.Entry(j.EntityID)
}
func (j *Job) AddBeforeFunc(f JobFunc) {
	if j.beforeFuncList == nil {
		j.beforeFuncList = make([]JobFunc, 0, 0)
	}
	j.beforeFuncList = append(j.beforeFuncList, f)
}
func (j *Job) OnFail(f JobFailFunc) {
	j.failFunc = f
}
func (j *Job) OnSuccess(f JobSuccessFunc) {
	j.successFunc = f
}
func (j *Job) AddAfterFunc(f JobFunc) {
	if j.afterFuncList == nil {
		j.afterFuncList = make([]JobFunc, 0, 0)
	}
	j.afterFuncList = append(j.afterFuncList, f)
}
func (j *Job) Run() {
	if j.beforeFuncList != nil && len(j.beforeFuncList) > 0 {
		for _, f := range j.beforeFuncList {
			err := f(j)
			if err != nil {
				j.failFunc(j, err) // handle before error
				return
			}
		}
	}
	j.run()
	if j.afterFuncList != nil && len(j.afterFuncList) > 0 {
		for _, f := range j.afterFuncList {
			_ = f(j) // ignore afterFunc error
			// TODO what about panic??
		}
	}
}
func (j *Job) run() {
	defer func() {
		if err := recover(); err != nil {
			if e, ok := err.(error); ok {
				// fmt.Println("recover is error")
				if j.failFunc != nil {
					// fmt.Println("failFunc is not nil")
					j.failFunc(j, e)
				} else {
					// fmt.Println("failFunc is nil")
					if j.mgr.logError != nil {
						// fmt.Println("LogError is not nil")
						j.mgr.logError <- errors.Wrap(e, "job fail without failFunc")
					} else {
						// fmt.Println("LogError is nil")
						// fmt.Println("job fail without failFunc", e)
					}
				}
			} else {
				fmt.Println("!!!!!!!!!error", err)
			}
		}
	}()
	j.Status = JobStatusRunning
	err := j.runFunc(j)
	if err != nil {
		j.failFunc(j, err)
	} else {
		if j.successFunc != nil {
			j.successFunc(j)
		}
	}
}
