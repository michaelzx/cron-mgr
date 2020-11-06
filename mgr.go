package cronmgr

import (
	"errors"
	"github.com/robfig/cron/v3"
	"sort"
	"sync"
	"time"
)

type IJobMgr interface {
	AddOnceJob(desc string, nextTime time.Time, runFunc JobFunc) (*Job, error)
	AddJob(desc string, spec string, runFunc JobFunc) (*Job, error)
	DelJob(id string)
	GetJob(id string) *Job
	GetJobList() []*Job
}
type IJob interface {
	GetEntity() cron.Entry
	OnFail(f JobFailFunc)
	OnSuccess(f JobSuccessFunc)
	AddBeforeFunc(f JobFunc)
	AddAfterFunc(f JobFunc)
	Run()
}
type JobMgr struct {
	cron     *cron.Cron
	jobMap   map[string]*Job
	jobMapRW sync.RWMutex
	logInfo  chan string
	logError chan error
}
type JobMgrOption struct {
	LogInfo  chan string
	LogError chan error
}

func NewJobMgr(opt *JobMgrOption) *JobMgr {
	var options []cron.Option
	options = append(options, cron.WithSeconds())
	c := cron.New(options...)
	c.Start()
	return &JobMgr{
		cron:     c,
		jobMap:   make(map[string]*Job),
		logInfo:  opt.LogInfo,
		logError: opt.LogError,
	}
}

func (mgr *JobMgr) DelJob(id string) {
	job, exists := mgr.jobMap[id]
	if !exists {
		return
	}
	mgr.jobMapRW.Lock()
	defer mgr.jobMapRW.Unlock()
	delete(mgr.jobMap, job.ID)
	mgr.cron.Remove(job.EntityID)
}

func (mgr *JobMgr) AddOnceJob(desc string, nextTime time.Time, runFunc JobFunc) (*Job, error) {
	if nextTime.Before(time.Now()) {
		return nil, errors.New("job already expired")
	}
	job := NewJob(mgr, desc, runFunc)
	spec := nextTime.Format("05 04 15 02 01 ?")
	id, err := mgr.cron.AddJob(spec, job)
	if err != nil {
		return nil, err
	}
	entity := mgr.cron.Entry(id)
	job.EntityID = entity.ID
	job.NextTime = entity.Next
	job.AddAfterFunc(func(thisJob *Job) error {
		mgr.DelJob(thisJob.ID)
		return nil
	})
	// add to mgr
	mgr.jobMapRW.Lock()
	defer mgr.jobMapRW.Unlock()
	mgr.jobMap[job.ID] = job
	return job, nil
}
func (mgr *JobMgr) AddJob(desc string, spec string, runFunc JobFunc) (*Job, error) {
	job := NewJob(mgr, desc, runFunc)
	id, err := mgr.cron.AddJob(spec, job)
	if err != nil {
		return nil, err
	}
	entity := mgr.cron.Entry(id)
	job.EntityID = entity.ID
	job.NextTime = entity.Next
	job.AddAfterFunc(func(thisJob *Job) error {
		thisJob.NextTime = thisJob.GetEntity().Next
		return nil
	})
	// add to mgr
	mgr.jobMapRW.Lock()
	defer mgr.jobMapRW.Unlock()
	mgr.jobMap[job.ID] = job
	return job, nil
}

func (mgr *JobMgr) GetJobList() []*Job {
	var list []*Job
	mgr.jobMapRW.RLock()
	defer mgr.jobMapRW.RUnlock()
	for _, j := range mgr.jobMap {
		list = append(list, j)
	}
	sort.Sort(byTime(list))
	return list
}
func (mgr *JobMgr) GetJob(id string) *Job {
	mgr.jobMapRW.RLock()
	defer mgr.jobMapRW.RUnlock()
	if j, ok := mgr.jobMap[id]; ok {
		return j
	}
	return nil
}
