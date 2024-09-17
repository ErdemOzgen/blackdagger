package scheduler

import (
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/engine"
	"github.com/ErdemOzgen/blackdagger/service/scheduler/job"
	"github.com/ErdemOzgen/blackdagger/service/scheduler/scheduler"
)

type jobFactory struct {
	Executable    string
	WorkDir       string
	EngineFactory engine.Factory
}

func (jf jobFactory) NewJob(d *dag.DAG, next time.Time) scheduler.Job {
	return &job.Job{
		DAG:           d,
		Executable:    jf.Executable,
		WorkDir:       jf.WorkDir,
		Next:          next,
		EngineFactory: jf.EngineFactory,
	}
}
