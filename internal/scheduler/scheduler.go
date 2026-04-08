package scheduler

import "github.com/robfig/cron/v3"

type Scheduler struct {
	cron                 *cron.Cron
	syncExercisesService *SyncExercisesService
}

func NewScheduler(syncExercisesService *SyncExercisesService) *Scheduler {
	c := cron.New(
		cron.WithSeconds(),
		cron.WithChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
		),
	)

	return &Scheduler{
		cron:                 c,
		syncExercisesService: syncExercisesService,
	}
}

func (s *Scheduler) Start() error {
	if err := s.registerJobs(); err != nil {
		return err
	}

	s.cron.Start()
	return nil
}

func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
}

func (s *Scheduler) registerJobs() error {
	_, err := s.cron.AddFunc("0 */1 * * * *", func() {
		_ = s.syncExercisesService.Run()
	})
	return err
}
