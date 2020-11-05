# cron mgr

```golang
type IJobMgr interface {
	AddOnceJob(desc string, nextTime time.Time, runFunc jobFunc) (*Job, error)
	AddJob(desc string, spec string, runFunc jobFunc) (*Job, error)
	DelJob(id string)
	GetJob(id string) *Job
	GetJobList() []*Job
}
type IJob interface {
	GetEntity() cron.Entry
	OnFail(f jobFailFunc)
	OnSuccess(f jobSuccessFunc)
	AddBeforeFunc(f jobFunc)
	AddAfterFunc(f jobFunc)
	Run()
}
```

# Usage
