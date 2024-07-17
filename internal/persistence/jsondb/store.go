package jsondb

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/persistence/filecache"

	"github.com/ErdemOzgen/blackdagger/internal/config"
	"github.com/ErdemOzgen/blackdagger/internal/persistence"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/model"

	"github.com/ErdemOzgen/blackdagger/internal/utils"
)

// Store is the interface to store dags status in local.
// It stores status in JSON format in a directory as per each dagFile.
// Multiple JSON data can be stored in a single file and each data
// is separated by newline.
// When a data is updated, it appends a new line to the file.
// Only the latest data in a single file can be read.
// When Compact is called, it removes old data.
// Compact must be called only once per file.
type Store struct {
	dir     string
	dagsDir string
	writer  *writer
	cache   *filecache.Cache[*model.Status]
}

var (
	errRequestIdNotFound  = errors.New("requestId not found")
	errCreateNewDirectory = errors.New("failed to create new directory")
	errDAGFileEmpty       = errors.New("dagFile is empty")
)

const (
	defaultCacheSize = 300
)

// New creates a new Store with default configuration.
func New(dir, dagsDir string) *Store {
	// dagsDir is used to calculate the directory that is compatible with the old version.
	s := &Store{
		dir:     dir,
		dagsDir: dagsDir,
		cache:   filecache.New[*model.Status](defaultCacheSize, time.Hour*3),
	}
	s.cache.StartEviction()
	return s
}

func (store *Store) Update(dagFile, requestId string, s *model.Status) error {
	f, err := store.FindByRequestId(dagFile, requestId)
	if err != nil {
		return err
	}
	w := &writer{target: f.File}
	if err := w.open(); err != nil {
		return err
	}
	defer func() {
		store.cache.Invalidate(f.File)
		_ = w.close()
	}()
	return w.write(s)
}

func (store *Store) Open(dagFile string, t time.Time, requestId string) error {
	writer, _, err := store.newWriter(dagFile, t, requestId)
	if err != nil {
		return err
	}
	if err := writer.open(); err != nil {
		return err
	}
	store.writer = writer
	return nil
}

func (store *Store) Write(s *model.Status) error {
	return store.writer.write(s)
}

func (store *Store) Close() error {
	if store.writer == nil {
		return nil
	}
	defer func() {
		_ = store.writer.close()
		store.writer = nil
	}()
	if err := store.Compact(store.writer.dagFile, store.writer.target); err != nil {
		return err
	}
	store.cache.Invalidate(store.writer.target)
	return store.writer.close()
}

func ParseFile(file string) (*model.Status, error) {
	f, err := os.Open(file)
	if err != nil {
		log.Printf("failed to open file. err: %v", err)
		return nil, err
	}
	defer func() {
		_ = f.Close()
	}()
	var offset int64
	var ret *model.Status
	for {
		line, err := readLineFrom(f, offset)
		if err == io.EOF {
			if ret == nil {
				return nil, err
			}
			return ret, nil
		} else if err != nil {
			return nil, err
		}
		offset += int64(len(line)) + 1 // +1 for newline
		if len(line) > 0 {
			var m *model.Status
			m, err = model.StatusFromJson(string(line))
			if err == nil {
				ret = m
				continue
			}
		}
	}
}

// NewWriter creates a new writer for a status.
func (store *Store) newWriter(dagFile string, t time.Time, requestId string) (*writer, string, error) {
	f, err := store.newFile(dagFile, t, requestId)
	if err != nil {
		return nil, "", err
	}
	w := &writer{target: f, dagFile: dagFile}
	return w, f, nil
}

// ReadStatusRecent returns recent n status
func (store *Store) ReadStatusRecent(dagFile string, n int) []*model.StatusFile {
	var ret []*model.StatusFile
	files := store.latest(store.pattern(dagFile)+"*.dat", n)
	for _, file := range files {
		status, err := store.cache.LoadLatest(file, func() (*model.Status, error) {
			return ParseFile(file)
		})
		if err != nil {
			continue
		}
		ret = append(ret, &model.StatusFile{
			File:   file,
			Status: status,
		})
	}
	return ret
}

// ReadStatusToday returns a list of status files.
func (store *Store) ReadStatusToday(dagFile string) (*model.Status, error) {
	// TODO: let's fix below not to use config here
	readLatestStatus := config.Get().LatestStatusToday
	file, err := store.latestToday(dagFile, time.Now(), readLatestStatus)
	if err != nil {
		return nil, err
	}
	return store.cache.LoadLatest(file, func() (*model.Status, error) {
		return ParseFile(file)
	})
}

// FindByRequestId finds a status file by requestId.
func (store *Store) FindByRequestId(dagFile string, requestId string) (*model.StatusFile, error) {
	if requestId == "" {
		return nil, errRequestIdNotFound
	}
	pattern := store.pattern(dagFile) + "*.dat"
	matches, err := filepath.Glob(pattern)
	if len(matches) > 0 || err == nil {
		sort.Slice(matches, func(i, j int) bool {
			return strings.Compare(matches[i], matches[j]) >= 0
		})
		for _, f := range matches {
			status, err := ParseFile(f)
			if err != nil {
				log.Printf("parsing failed %s : %s", f, err)
				continue
			}
			if status != nil && status.RequestId == requestId {
				return &model.StatusFile{
					File:   f,
					Status: status,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("%w : %s", persistence.ErrRequestIdNotFound, requestId)
}

// RemoveAll removes all files in a directory.
func (store *Store) RemoveAll(dagFile string) error {
	return store.RemoveOld(dagFile, 0)
}

// RemoveOld removes old files.
func (store *Store) RemoveOld(dagFile string, retentionDays int) error {
	pattern := store.pattern(dagFile) + "*.dat"
	var lastErr error
	if retentionDays >= 0 {
		matches, _ := filepath.Glob(pattern)
		ot := time.Now().AddDate(0, 0, -1*retentionDays)
		for _, m := range matches {
			info, err := os.Stat(m)
			if err == nil {
				if info.ModTime().Before(ot) {
					lastErr = os.Remove(m)
				}
			}
		}
	}
	return lastErr
}

// Compact creates a new file with only the latest data and removes old data.
func (store *Store) Compact(_, original string) error {
	status, err := ParseFile(original)
	if err != nil {
		return err
	}

	newFile := fmt.Sprintf("%s_c.dat",
		strings.TrimSuffix(filepath.Base(original), path.Ext(original)))
	f := path.Join(filepath.Dir(original), newFile)
	w := &writer{target: f}
	if err := w.open(); err != nil {
		return err
	}
	defer func() {
		_ = w.close()
	}()

	if err := w.write(status); err != nil {
		if err := os.Remove(f); err != nil {
			log.Printf("failed to remove %s : %s", f, err.Error())
		}
		return err
	}

	return os.Remove(original)
}

func (store *Store) normalizeInternalName(name string) string {
	a := strings.TrimSuffix(name, ".yaml")
	a = strings.TrimSuffix(a, ".yml")
	a = path.Join(store.dagsDir, a)
	return fmt.Sprintf("%s.yaml", a)
}

func (store *Store) exists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func (store *Store) Rename(oldName, newName string) error {
	// This is needed to ensure backward compatibility.
	on := store.normalizeInternalName(oldName)
	nn := store.normalizeInternalName(newName)

	oldDir := store.directory(on, prefix(on))
	newDir := store.directory(nn, prefix(nn))
	if !store.exists(oldDir) {
		// Nothing to do
		return nil
	}
	if !store.exists(newDir) {
		if err := os.MkdirAll(newDir, 0755); err != nil {
			return fmt.Errorf("%w: %s : %s", errCreateNewDirectory, newDir, err.Error())
		}
	}
	matches, err := filepath.Glob(store.pattern(on) + "*.dat")
	if err != nil {
		return err
	}
	oldPattern := path.Base(store.pattern(on))
	newPattern := path.Base(store.pattern(nn))
	for _, m := range matches {
		base := path.Base(m)
		f := strings.Replace(base, oldPattern, newPattern, 1)
		_ = os.Rename(m, path.Join(newDir, f))
	}
	if files, _ := os.ReadDir(oldDir); len(files) == 0 {
		_ = os.Remove(oldDir)
	}
	return nil
}

func (store *Store) directory(name string, prefix string) string {
	h := md5.New()
	_, _ = h.Write([]byte(name))
	v := hex.EncodeToString(h.Sum(nil))
	return filepath.Join(store.dir, fmt.Sprintf("%s-%s", prefix, v))
}

func (store *Store) newFile(dagFile string, t time.Time, requestId string) (string, error) {
	if dagFile == "" {
		return "", errDAGFileEmpty
	}
	fileName := fmt.Sprintf("%s.%s.%s.dat", store.pattern(dagFile), t.Format("20060102.15:04:05.000"), utils.TruncString(requestId, 8))
	return fileName, nil
}

func (store *Store) pattern(dagFile string) string {
	p := prefix(dagFile)
	dir := store.directory(dagFile, p)
	return filepath.Join(dir, p)
}

func (store *Store) latestToday(dagFile string, day time.Time, latestStatusToday bool) (string, error) {
	var ret []string
	pattern := ""
	if latestStatusToday {
		pattern = fmt.Sprintf("%s.%s*.*.dat", store.pattern(dagFile), day.Format("20060102"))
	} else {
		pattern = fmt.Sprintf("%s.*.*.dat", store.pattern(dagFile))
	}
	matches, err := filepath.Glob(pattern)
	if err != nil || len(matches) == 0 {
		return "", persistence.ErrNoStatusDataToday
	}
	ret = filterLatest(matches, 1)

	if len(ret) == 0 {
		return "", persistence.ErrNoStatusData
	}
	return ret[0], err
}

func (store *Store) latest(pattern string, n int) []string {
	matches, err := filepath.Glob(pattern)
	var ret = []string{}
	if err == nil || len(matches) >= 0 {
		ret = filterLatest(matches, n)
	}
	return ret
}

var rTimestamp = regexp.MustCompile(`2\d{7}.\d{2}:\d{2}:\d{2}`)

func filterLatest(files []string, n int) []string {
	if len(files) == 0 {
		return []string{}
	}
	sort.Slice(files, func(i, j int) bool {
		t1 := timestamp(files[i])
		t2 := timestamp(files[j])
		return t1 > t2
	})
	ret := make([]string, 0, n)
	for i := 0; i < n && i < len(files); i++ {
		ret = append(ret, files[i])
	}
	return ret
}

func timestamp(file string) string {
	return rTimestamp.FindString(file)
}

func readLineFrom(f *os.File, offset int64) ([]byte, error) {
	if _, err := f.Seek(offset, io.SeekStart); err != nil {
		return nil, err
	}
	r := bufio.NewReader(f)
	var ret []byte
	for {
		b, isPrefix, err := r.ReadLine()
		if err == io.EOF {
			return ret, err
		} else if err != nil {
			log.Printf("read line failed. %s", err)
			return nil, err
		}
		if err == nil {
			ret = append(ret, b...)
			if !isPrefix {
				break
			}
		}
	}
	return ret, nil
}

func prefix(dagFile string) string {
	return strings.TrimSuffix(
		filepath.Base(dagFile),
		path.Ext(dagFile),
	)
}
