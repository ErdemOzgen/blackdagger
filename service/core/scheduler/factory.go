package scheduler

import (
	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/engine"
	"github.com/ErdemOzgen/blackdagger/service/core/scheduler/job"
	"github.com/ErdemOzgen/blackdagger/service/core/scheduler/scheduler"
	"time"
)

type jobFactory struct {
	Command       string
	WorkDir       string
	EngineFactory engine.Factory
}

func (jf jobFactory) NewJob(dag *dag.DAG, next time.Time) scheduler.Job {
	return &job.Job{
		DAG:           dag,
		Command:       jf.Command,
		WorkDir:       jf.WorkDir,
		Next:          next,
		EngineFactory: jf.EngineFactory,
	}
}
