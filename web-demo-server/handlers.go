package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	cronmgr "github.com/michaelzx/cron-mgr"
	"github.com/pkg/errors"
	"time"
)

func jobList(ctx *fiber.Ctx) error {
	return ctx.JSON(jobMater.GetJobList())
}

func createRepeat(ctx *fiber.Ctx) error {
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
	return ctx.JSON(job)
}

func createSuccess(ctx *fiber.Ctx) error {
	execTime := time.Now().Add(time.Duration(10) * time.Second)
	job, err := jobMater.AddOnceJob("test-job-success", execTime, func(thisJob *cronmgr.Job) error {
		fmt.Println(thisJob.ID, "run")
		return nil
	})
	if err != nil {
		panic(err)
	}
	// optional
	job.OnSuccess(func(thisJob *cronmgr.Job) {
		fmt.Println(thisJob.ID, "success")
	})
	// optional
	job.OnFail(func(thisJob *cronmgr.Job, jobErr error) {
		fmt.Println(thisJob.ID, jobErr)
	})
	return ctx.JSON(job)
}
func createPanic(ctx *fiber.Ctx) error {
	execTime := time.Now().Add(time.Duration(10) * time.Second)
	job, err := jobMater.AddOnceJob("test-job-panic", execTime, func(thisJob *cronmgr.Job) error {
		fmt.Println(thisJob.ID, "run")
		panic(errors.New(thisJob.ID + ":panic"))
		return nil
	})
	if err != nil {
		panic(err)
	}
	return ctx.JSON(job)
}
func createError(ctx *fiber.Ctx) error {
	execTime := time.Now().Add(time.Duration(10) * time.Second)
	job, err := jobMater.AddOnceJob("test-job-panic", execTime, func(thisJob *cronmgr.Job) error {
		fmt.Println(thisJob.ID, "run")
		return errors.New(thisJob.ID + ":return error")
	})
	if err != nil {
		panic(err)
	}
	return ctx.JSON(job)
}

func remove(ctx *fiber.Ctx) error {

	id := ctx.Params("id")
	jobMater.DelJob(id)
	return ctx.Send(nil)
}
