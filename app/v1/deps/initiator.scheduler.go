package deps

import "dot/infrastructure/v1/scheduler/cron"

func (d *dependency) initScheduler() {
	d.initCronScheduler()
}

func (d *dependency) initCronScheduler() {
	d.base.Schlr = cron.New()
}
