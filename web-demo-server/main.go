package main

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	cronmgr "github.com/michaelzx/cron-mgr"
)

var jobMater *cronmgr.JobMgr

func main() {

	logInfo := make(chan string)
	logError := make(chan error)

	go func() {
		for {
			select {
			case s := <-logInfo:
				fmt.Println("info---->", s)
			case e := <-logError:
				fmt.Println("error---->", e)
				fmt.Printf("%+v", e)
			}
		}
	}()
	jobMater = cronmgr.NewJobMgr(logInfo, logError)

	runWebServer()
}

func runWebServer() {
	app := fiber.New(fiber.Config{})
	app.Get("/list", jobList)
	app.Get("/create/repeat", createRepeat)
	app.Get("/create/success", createSuccess)
	app.Get("/create/panic", createPanic)
	app.Get("/create/error", createError)
	app.Get("/remove/:id", remove)
	err := app.Listen(":8888")
	if err != nil {
		panic(err)
	}
}