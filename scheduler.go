package main

import (
	"github.com/go-co-op/gocron/v2"
	"time"
)

type Scheduler struct {
	underlyingScheduler gocron.Scheduler
}

func NewScheduler() *Scheduler {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}
	return &Scheduler{underlyingScheduler: scheduler}
}

func (s *Scheduler) Shutdown() {
	s.underlyingScheduler.Shutdown()
}

func (s *Scheduler) ScheduleTaskAndStart(function any, parameters ...any) error {
	_, err := s.underlyingScheduler.NewJob(
		gocron.WeeklyJob(1,
			gocron.NewWeekdays(time.Wednesday),
			gocron.NewAtTimes(gocron.NewAtTime(10, 0, 0))),
		gocron.NewTask(
			function,
			parameters...,
		),
	)
	/*_, err := s.underlyingScheduler.NewJob(
		gocron.DurationJob(5*time.Second),
		gocron.NewTask(
			function,
			parameters...,
		),
	)*/
	if err != nil {
		return err
	}

	// start the scheduler
	s.underlyingScheduler.Start()

	return nil
}
