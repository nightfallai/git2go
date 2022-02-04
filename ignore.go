package git

/*
#include <git2.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

func (v *Repository) AddIgnoreRule(rules string) error {
	crules := C.CString(rules)
	defer C.free(unsafe.Pointer(crules))

	ret := C.git_ignore_add_rule(v.ptr, crules)
	runtime.KeepAlive(v)
	if ret < 0 {
		return MakeFastGitError(ret)
	}
	return nil
}

func (v *Repository) ClearInternalIgnoreRules() error {
	ret := C.git_ignore_clear_internal_rules(v.ptr)
	runtime.KeepAlive(v)
	if ret < 0 {
		return MakeFastGitError(ret)
	}
	return nil
}

func (v *Repository) IsPathIgnored(path string) (bool, error) {
	var ignored C.int

	cpath := C.CString(path)
	defer C.free(unsafe.Pointer(cpath))

	ret := C.git_ignore_path_is_ignored(&ignored, v.ptr, cpath)
	runtime.KeepAlive(v)
	if ret < 0 {
		return false, MakeFastGitError(ret)
	}
	return ignored == 1, nil
}
