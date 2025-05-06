package local

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/ErdemOzgen/blackdagger/internal/dag"
	"github.com/ErdemOzgen/blackdagger/internal/persistence"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/filecache"
	"github.com/ErdemOzgen/blackdagger/internal/persistence/grep"
	"github.com/ErdemOzgen/blackdagger/internal/util"
)

type dagStoreImpl struct {
	dir       string
	metaCache *filecache.Cache[*dag.DAG]
}

type NewDAGStoreArgs struct {
	Dir string
}

func NewDAGStore(args *NewDAGStoreArgs) persistence.DAGStore {
	dagStore := &dagStoreImpl{
		dir:       args.Dir,
		metaCache: filecache.New[*dag.DAG](0, time.Hour*24),
	}
	dagStore.metaCache.StartEviction()
	return dagStore
}

func (d *dagStoreImpl) GetMetadata(name string) (*dag.DAG, error) {
	loc, err := d.resolve(name)
	if err != nil {
		return nil, err
	}

	return d.metaCache.LoadLatest(loc, func() (*dag.DAG, error) {
		return dag.LoadMetadata(loc)
	})
}

func (d *dagStoreImpl) GetDetails(name string) (*dag.DAG, error) {
	loc, err := d.resolve(name)
	if err != nil {
		return nil, err
	}

	dat, err := dag.LoadWithoutEval(loc)
	if err != nil {
		return nil, err
	}
	return dat, nil
}

func (d *dagStoreImpl) GetSpec(name string) (string, error) {
	loc, err := d.fileLocation(name)
	if err != nil {
		return "", err
	}
	dat, err := os.ReadFile(loc)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}

// TODO: use 0600 // nolint: gosec
const defaultPerm os.FileMode = 0744

var errDOGFileNotExist = errors.New("the DAG file does not exist")

func (d *dagStoreImpl) UpdateSpec(name string, spec []byte) error {
	// Validate the new spec
	_, err := dag.LoadYAML(spec)
	if err != nil {
		return err
	}

	loc, err := d.resolve(name)
	if err != nil {
		return err
	}

	// Check if the file exists
	if !exists(loc) {
		return fmt.Errorf("%w: %s", errDOGFileNotExist, loc)
	}

	// Write the updated spec to the file
	err = os.WriteFile(loc, spec, defaultPerm)
	if err != nil {
		return err
	}

	// Invalidate the cache
	d.metaCache.Invalidate(loc)
	return nil
}

var errDAGFileAlreadyExists = errors.New("the DAG file already exists")

func (d *dagStoreImpl) Create(name string, spec []byte) (string, error) {
	if err := d.ensureDirExist(); err != nil {
		return "", err
	}

	name += ".yaml"
	loc := filepath.Join(d.dir, name)
	fmt.Printf("Creating DAG file at %s...\n", loc)

	if exists(loc) {
		return "", fmt.Errorf("%w: %s", errDAGFileAlreadyExists, loc)
	}
	// nolint: gosec
	fmt.Printf("Creating DAG file at %s...\n", loc)
	return name, os.WriteFile(loc, spec, 0644)
}

func (d *dagStoreImpl) Delete(name string) error {
	loc, err := d.resolve(name)
	if err != nil {
		return err
	}

	err = os.Remove(loc)
	if err != nil {
		return err
	}

	d.metaCache.Invalidate(loc)
	return nil
}

func exists(file string) bool {
	_, err := os.Stat(file)
	return !os.IsNotExist(err)
}

func (d *dagStoreImpl) fileLocation(name string) (string, error) {
	if filepath.IsAbs(name) {
		// If the name is already an absolute path, return it as-is
		return name, nil
	}
	loc, err := d.resolve(name)
	if err != nil {
		return "", err
	}
	fmt.Printf("file location: %s\n", loc)
	return util.AddYamlExtension(loc), nil
}

func (d *dagStoreImpl) ensureDirExist() error {
	if !exists(d.dir) {
		if err := os.MkdirAll(d.dir, 0755); err != nil {
			return err
		}
	}
	return nil
}

func (d *dagStoreImpl) searchName(fileName string, searchText *string) bool {
	if searchText == nil {
		return true
	}

	fileName = strings.TrimSuffix(fileName, path.Ext(fileName))
	fileName = strings.ToLower(fileName)

	*searchText = strings.ToLower(*searchText)

	return strings.Contains(fileName, *searchText)
}

func (d *dagStoreImpl) searchDescription(description string, searchDescription *string) bool {
	if searchDescription == nil {
		return true
	}

	return strings.Contains(description, *searchDescription)
}

func (d *dagStoreImpl) searchTags(tags []string, searchTag *string) bool {
	if searchTag == nil {
		return true
	}

	for _, tag := range tags {
		if tag == *searchTag {
			return true
		}
	}

	return false
}

func (d *dagStoreImpl) getTagList(tagSet map[string]struct{}) []string {
	tagList := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		tagList = append(tagList, tag)
	}
	return tagList
}
func (d *dagStoreImpl) ListPagination(params persistence.DAGListPaginationArgs) (*persistence.DagListPaginationResult, error) {
	var (
		dagList    = make([]*dag.DAG, 0)
		errList    = make([]string, 0)
		count      int
		currentDag *dag.DAG
	)

	err := filepath.WalkDir(d.dir, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			errList = append(errList, fmt.Sprintf("error accessing %s: %s", path, err))
			return nil
		}

		// Skip directories
		if dir.IsDir() {
			return nil
		}

		// Process only files with valid extensions
		if checkExtension(dir.Name()) {
			currentDag, err = d.GetMetadata(path)
			if err != nil {
				errList = append(errList, fmt.Sprintf("reading %s failed: %s", path, err))
				return nil
			}

			if !d.searchName(dir.Name(), params.Name) && !d.searchDescription(currentDag.Description, params.Name) || currentDag == nil || !d.searchTags(currentDag.Tags, params.Tag) {
				return nil
			}

			count++
			if count > (params.Page-1)*params.Limit && len(dagList) < params.Limit {
				dagList = append(dagList, currentDag)
			}
		}

		return nil
	})

	if err != nil {
		return &persistence.DagListPaginationResult{
			DagList:   dagList,
			Count:     count,
			ErrorList: append(errList, err.Error()),
		}, err
	}

	return &persistence.DagListPaginationResult{
		DagList:   dagList,
		Count:     count,
		ErrorList: errList,
	}, nil
}

func (d *dagStoreImpl) List() (ret []*dag.DAG, errs []string, err error) {
	if err = d.ensureDirExist(); err != nil {
		errs = append(errs, err.Error())
		return
	}

	err = filepath.WalkDir(d.dir, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			errs = append(errs, fmt.Sprintf("error accessing %s: %s", path, err))
			return nil
		}

		if dir.IsDir() {
			return nil
		}

		if checkExtension(dir.Name()) {
			dat, err := d.GetMetadata(path)
			if err == nil {
				ret = append(ret, dat)
			} else {
				errs = append(errs, fmt.Sprintf("reading %s failed: %s", path, err))
			}
		}

		return nil
	})

	if err != nil {
		errs = append(errs, fmt.Sprintf("error walking directory: %s", err))
	}

	return ret, errs, nil
}

var extensions = []string{".yaml", ".yml"}

func checkExtension(file string) bool {
	ext := filepath.Ext(file)
	for _, e := range extensions {
		if e == ext {
			return true
		}
	}
	return false
}

func (d *dagStoreImpl) Grep(pattern string) (ret []*persistence.GrepResult, errs []string, err error) {
	if err = d.ensureDirExist(); err != nil {
		errs = append(errs, fmt.Sprintf("failed to create DAGs directory %s", d.dir))
		return
	}

	err = filepath.WalkDir(d.dir, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			errs = append(errs, fmt.Sprintf("error accessing %s: %s", path, err))
			return nil
		}

		// Skip directories
		if dir.IsDir() {
			return nil
		}

		// Process only files with valid extensions
		if checkExtension(dir.Name()) {
			file := path
			dat, err := os.ReadFile(file)
			if err != nil {
				errs = append(errs, fmt.Sprintf("read DAG file %s failed: %s", file, err))
				return nil
			}

			m, err := grep.Grep(dat, fmt.Sprintf("(?i)%s", pattern), &grep.Options{
				IsRegexp: true,
				Before:   2,
				After:    2,
			})
			if err != nil {
				errs = append(errs, fmt.Sprintf("grep %s failed: %s", file, err))
				return nil
			}

			dg, err := dag.LoadMetadata(file)
			if err != nil {
				errs = append(errs, fmt.Sprintf("check %s failed: %s", file, err))
				return nil
			}

			ret = append(ret, &persistence.GrepResult{
				Name:    strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)),
				DAG:     dg,
				Matches: m,
			})
		}

		return nil
	})

	return ret, errs, err
}

func (d *dagStoreImpl) Rename(oldID, newID string) error {
	// Resolve the old file location
	oldLoc, err := d.resolve(oldID)
	if err != nil {
		return err
	}
	fmt.Printf("Old location: %s\n", oldLoc)

	oldExt := filepath.Ext(oldLoc)
	dir := filepath.Dir(oldLoc)
	newID = newID + oldExt
	newLoc := filepath.Join(dir, newID)

	fmt.Printf("New location: %s\n", newLoc)

	// Rename the file
	err = os.Rename(oldLoc, newLoc)
	if err != nil {
		return fmt.Errorf("failed to rename %s to %s: %w", oldLoc, newLoc, err)
	}

	return nil
}

func (d *dagStoreImpl) Find(name string) (*dag.DAG, error) {
	file, err := d.resolve(name)
	if err != nil {
		return nil, err
	}
	return dag.LoadWithoutEval(file)
}

func (d *dagStoreImpl) resolve(name string) (string, error) {
	// Check if the name is an absolute or relative file path
	if strings.Contains(name, string(filepath.Separator)) {
		if util.FileExists(name) {
			return name, nil
		}
		return "", fmt.Errorf("workflow %s not found", name)
	}

	// Search recursively under d.dir
	var foundPath string
	err := filepath.WalkDir(d.dir, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if dir.IsDir() {
			return nil
		}

		// Check if the file matches the name (with or without extension)
		if filepath.Base(path) == name || strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)) == name {
			foundPath = path
			return filepath.SkipDir // Stop further traversal
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if foundPath == "" {
		return "", fmt.Errorf("workflow %s not found", name)
	}

	return foundPath, nil
}

// find finds the sub workflow file with the given name.
func find(name string) (string, error) {
	ext := path.Ext(name)
	if ext == "" {
		// try all supported extensions
		for _, ext := range dag.Exts {
			if util.FileExists(name + ext) {
				return filepath.Abs(name + ext)
			}
		}
	} else if util.FileExists(name) {
		// the name has an extension
		return filepath.Abs(name)
	}
	return "", fmt.Errorf("sub workflow %s not found", name)
}
func (d *dagStoreImpl) TagList() ([]string, []string, error) {
	var (
		errList    = make([]string, 0)
		tagSet     = make(map[string]struct{})
		currentDag *dag.DAG
	)

	err := filepath.WalkDir(d.dir, func(path string, dir fs.DirEntry, err error) error {
		if err != nil {
			errList = append(errList, fmt.Sprintf("error accessing %s: %s", path, err))
			return nil
		}

		if dir.IsDir() {
			return nil
		}

		if checkExtension(dir.Name()) {
			currentDag, err = d.GetMetadata(path)
			if err != nil {
				errList = append(errList, fmt.Sprintf("reading %s failed: %s", path, err))
				return nil
			}

			if currentDag == nil {
				return nil
			}

			for _, tag := range currentDag.Tags {
				tagSet[tag] = struct{}{}
			}
		}

		return nil
	})

	if err != nil {
		return nil, append(errList, err.Error()), err
	}

	return d.getTagList(tagSet), errList, nil
}
