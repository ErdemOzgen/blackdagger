package entry_reader

import (
	"github.com/ErdemOzgen/blackdagger/internal/engine"
	"github.com/ErdemOzgen/blackdagger/internal/logger"
	"github.com/ErdemOzgen/blackdagger/internal/logger/tag"
	"github.com/ErdemOzgen/blackdagger/service/core/scheduler/filenotify"
	"github.com/ErdemOzgen/blackdagger/service/core/scheduler/scheduler"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/utils"
	"github.com/fsnotify/fsnotify"
)

type JobFactory interface {
	NewJob(dag *dag.DAG, next time.Time) scheduler.Job
}

type Params struct {
	DagsDir       string
	JobFactory    JobFactory
	Logger        logger.Logger
	EngineFactory engine.Factory
}

type EntryReader struct {
	dagsDir       string
	dagsLock      sync.Mutex
	dags          map[string]*dag.DAG
	jf            JobFactory
	logger        logger.Logger
	engineFactory engine.Factory
}

func New(params Params) *EntryReader {
	er := &EntryReader{
		dagsDir:       params.DagsDir,
		dagsLock:      sync.Mutex{},
		dags:          map[string]*dag.DAG{},
		jf:            params.JobFactory,
		logger:        params.Logger,
		engineFactory: params.EngineFactory,
	}
	if err := er.initDags(); err != nil {
		er.logger.Error("failed to init entry_reader dags", tag.Error(err))
	}
	go er.watchDags()
	return er
}

func (er *EntryReader) Read(now time.Time) ([]*scheduler.Entry, error) {
	var entries []*scheduler.Entry
	er.dagsLock.Lock()
	defer er.dagsLock.Unlock()

	f := func(d *dag.DAG, s []*dag.Schedule, e scheduler.Type) {
		for _, ss := range s {
			next := ss.Parsed.Next(now)
			entries = append(entries, &scheduler.Entry{
				Next: ss.Parsed.Next(now),
				// TODO: fix this
				Job:       er.jf.NewJob(d, next),
				EntryType: e,
				Logger:    er.logger,
			})
		}
	}

	e := er.engineFactory.Create()
	for _, d := range er.dags {
		if e.IsSuspended(d.Name) {
			continue
		}
		f(d, d.Schedule, scheduler.Start)
		f(d, d.StopSchedule, scheduler.Stop)
		f(d, d.RestartSchedule, scheduler.Restart)
	}

	return entries, nil
}

func (er *EntryReader) initDags() error {
	er.dagsLock.Lock()
	defer er.dagsLock.Unlock()
	cl := dag.Loader{}
	fis, err := os.ReadDir(er.dagsDir)
	if err != nil {
		return err
	}
	var fileNames []string
	for _, fi := range fis {
		if utils.MatchExtension(fi.Name(), dag.EXTENSIONS) {
			dag, err := cl.LoadMetadata(filepath.Join(er.dagsDir, fi.Name()))
			if err != nil {
				er.logger.Error("failed to read DAG cfg", tag.Error(err))
				continue
			}
			er.dags[fi.Name()] = dag
			fileNames = append(fileNames, fi.Name())
		}
	}
	er.logger.Info("init backend dags", "files", strings.Join(fileNames, ","))
	return nil
}

func (er *EntryReader) watchDags() {
	cl := dag.Loader{}
	watcher, err := filenotify.New(time.Minute)
	if err != nil {
		er.logger.Error("failed to init file watcher", tag.Error(err))
		return
	}
	defer func() {
		_ = watcher.Close()
	}()
	_ = watcher.Add(er.dagsDir)
	for {
		select {
		case event, ok := <-watcher.Events():
			if !ok {
				return
			}
			if !utils.MatchExtension(event.Name, dag.EXTENSIONS) {
				continue
			}
			er.dagsLock.Lock()
			if event.Op == fsnotify.Create || event.Op == fsnotify.Write {
				dag, err := cl.LoadMetadata(filepath.Join(er.dagsDir, filepath.Base(event.Name)))
				if err != nil {
					er.logger.Error("failed to read DAG cfg", tag.Error(err))
				} else {
					er.dags[filepath.Base(event.Name)] = dag
					er.logger.Info("reload DAG entry_reader", "file", event.Name)
				}
			}
			if event.Op == fsnotify.Rename || event.Op == fsnotify.Remove {
				delete(er.dags, filepath.Base(event.Name))
				er.logger.Info("remove DAG entry_reader", "file", event.Name)
			}
			er.dagsLock.Unlock()
		case err, ok := <-watcher.Errors():
			if !ok {
				return
			}
			er.logger.Error("watch entry_reader DAGs error", tag.Error(err))
		}
	}

}
