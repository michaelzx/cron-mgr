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
## basic

```go
	jobMater = cronmgr.NewJobMgr(&cronmgr.JobMgrOption{})
    
```
## add repeat job

```go
	job, err := jobMater.AddJob("test-job-repeat", "*/5 * * * * ?", func(thisJob *cronmgr.Job) error {
		fmt.Println(time.Now(), thisJob.ID, "run")
		return nil
	})
	if err != nil {
		panic(err)
	}
	job.OnSuccess(func(thisJob *cronmgr.Job) {
		// ...
	})
	job.OnFail(func(thisJob *cronmgr.Job, jobErr error) {
		// ...
	})
```