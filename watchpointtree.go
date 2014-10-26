package notify

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/rjeczalik/fs"
)

// Point TODO
type Point struct {
	Name   string
	Parent map[string]interface{}
}

// WatchPointTree TODO
type WatchPointTree struct {
	// FS TODO
	FS fs.Filesystem

	cwd  Point
	root map[string]interface{}
}

// NewWatchPointTree TODO
func NewWatchPointTree() *WatchPointTree {
	return &WatchPointTree{
		FS:   fs.Default,
		root: make(map[string]interface{}),
	}
}

// WalkPointFunc TODO
type WalkPointFunc func(pt Point, last bool) error

// WalkPoint TODO
func (w WatchPointTree) WalkPoint(p string, fn WalkPointFunc) (err error) {
	parent, i := w.begin(p)
	for j := 0; ; {
		if j = strings.IndexRune(p[i:], os.PathSeparator); j == -1 {
			break
		}
		pt := Point{Name: p[i : i+j], Parent: parent}
		if err = fn(pt, false); err != nil {
			return
		}
		// TODO(rjeczalik): handle edge case where parent[pt.Name] is a file
		cd, ok := parent[pt.Name].(map[string]interface{})
		if !ok {
			cd = make(map[string]interface{})
			parent[pt.Name] = cd
		}
		i += j + 1
		parent = cd
	}
	if i < len(p) {
		err = fn(Point{Name: p[i:], Parent: parent}, true)
	}
	return
}

func (w WatchPointTree) begin(p string) (d map[string]interface{}, n int) {
	if len(w.cwd.Name) != 0 {
		if i := strings.Index(p, w.cwd.Name); i != -1 && p[i] == os.PathSeparator {
			return w.cwd.Parent, i
		}
	}
	vol := filepath.VolumeName(p)
	n = len(vol)
	if n == 0 {
		return w.root, 1
	}
	d, ok := w.root[vol].(map[string]interface{})
	if !ok {
		d = make(map[string]interface{})
		w.root[vol] = d
	}
	return d, n
}

func (w WatchPointTree) Register(p string, c chan<- EventInfo, e Event) error {
	if _, err := w.FS.Stat(p); err != nil {
		return err
	}
	return nil
}