package git

/*
#include <git2.h>
*/
import "C"
import (
	"runtime"
)

type CherrypickOptions struct {
	Mainline        uint
	MergeOptions    MergeOptions
	CheckoutOptions CheckoutOptions
}

func cherrypickOptionsFromC(c *C.git_cherrypick_options) CherrypickOptions {
	opts := CherrypickOptions{
		Mainline:        uint(c.mainline),
		MergeOptions:    mergeOptionsFromC(&c.merge_opts),
		CheckoutOptions: checkoutOptionsFromC(&c.checkout_opts),
	}
	return opts
}

func populateCherrypickOptions(copts *C.git_cherrypick_options, opts *CherrypickOptions, errorTarget *error) *C.git_cherrypick_options {
	C.git_cherrypick_options_init(copts, C.GIT_CHERRYPICK_OPTIONS_VERSION)
	if opts == nil {
		return nil
	}
	copts.mainline = C.uint(opts.Mainline)
	populateMergeOptions(&copts.merge_opts, &opts.MergeOptions)
	populateCheckoutOptions(&copts.checkout_opts, &opts.CheckoutOptions, errorTarget)
	return copts
}

func freeCherrypickOpts(copts *C.git_cherrypick_options) {
	if copts == nil {
		return
	}
	freeMergeOptions(&copts.merge_opts)
	freeCheckoutOptions(&copts.checkout_opts)
}

func DefaultCherrypickOptions() (CherrypickOptions, error) {
	c := C.git_cherrypick_options{}

	ecode := C.git_cherrypick_options_init(&c, C.GIT_CHERRYPICK_OPTIONS_VERSION)
	if ecode < 0 {
		return CherrypickOptions{}, MakeFastGitError(ecode)
	}
	defer freeCherrypickOpts(&c)
	return cherrypickOptionsFromC(&c), nil
}

func (v *Repository) Cherrypick(commit *Commit, opts CherrypickOptions) error {
	var err error
	cOpts := populateCherrypickOptions(&C.git_cherrypick_options{}, &opts, &err)
	defer freeCherrypickOpts(cOpts)

	ret := C.git_cherrypick(v.ptr, commit.cast_ptr, cOpts)
	runtime.KeepAlive(v)
	runtime.KeepAlive(commit)
	if ret == C.int(ErrorCodeUser) && err != nil {
		return err
	}
	if ret < 0 {
		return MakeFastGitError(ret)
	}
	return nil
}

func (r *Repository) CherrypickCommit(pick, our *Commit, opts CherrypickOptions) (*Index, error) {
	cOpts := populateMergeOptions(&C.git_merge_options{}, &opts.MergeOptions)
	defer freeMergeOptions(cOpts)

	var ptr *C.git_index
	ret := C.git_cherrypick_commit(&ptr, r.ptr, pick.cast_ptr, our.cast_ptr, C.uint(opts.Mainline), cOpts)
	runtime.KeepAlive(pick)
	runtime.KeepAlive(our)
	if ret < 0 {
		return nil, MakeFastGitError(ret)
	}
	return newIndexFromC(ptr, r), nil
}
