package git

/*
#include <git2.h>
*/
import "C"

type ResetType int

const (
	ResetSoft  ResetType = C.GIT_RESET_SOFT
	ResetMixed ResetType = C.GIT_RESET_MIXED
	ResetHard  ResetType = C.GIT_RESET_HARD
)

func (r *Repository) ResetToCommit(commit *Commit, resetType ResetType, opts *CheckoutOptions) error {
	var err error
	cOpts := populateCheckoutOptions(&C.git_checkout_options{}, opts, &err)
	defer freeCheckoutOptions(cOpts)

	ret := C.git_reset(r.ptr, commit.ptr, C.git_reset_t(resetType), cOpts)
	if ret == C.int(ErrorCodeUser) && err != nil {
		return err
	}
	if ret < 0 {
		return MakeFastGitError(ret)
	}
	return nil
}

func (r *Repository) ResetDefaultToCommit(commit *Commit, pathspecs []string) error {
	cpathspecs := C.git_strarray{}
	cpathspecs.count = C.size_t(len(pathspecs))
	cpathspecs.strings = makeCStringsFromStrings(pathspecs)
	defer freeStrarray(&cpathspecs)

	ret := C.git_reset_default(r.ptr, commit.ptr, &cpathspecs)

	if ret < 0 {
		return MakeFastGitError(ret)
	}
	return nil
}
