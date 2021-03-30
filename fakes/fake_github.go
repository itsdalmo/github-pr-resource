// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/shurcooL/githubv4"
	resource "github.com/telia-oss/github-pr-resource"
)

type FakeGithub struct {
	DeletePreviousCommentsStub        func(string) error
	deletePreviousCommentsMutex       sync.RWMutex
	deletePreviousCommentsArgsForCall []struct {
		arg1 string
	}
	deletePreviousCommentsReturns struct {
		result1 error
	}
	deletePreviousCommentsReturnsOnCall map[int]struct {
		result1 error
	}
	GetChangedFilesStub        func(string, string) ([]resource.ChangedFileObject, error)
	getChangedFilesMutex       sync.RWMutex
	getChangedFilesArgsForCall []struct {
		arg1 string
		arg2 string
	}
	getChangedFilesReturns struct {
		result1 []resource.ChangedFileObject
		result2 error
	}
	getChangedFilesReturnsOnCall map[int]struct {
		result1 []resource.ChangedFileObject
		result2 error
	}
	GetPullRequestStub        func(string, string) (*resource.PullRequest, error)
	getPullRequestMutex       sync.RWMutex
	getPullRequestArgsForCall []struct {
		arg1 string
		arg2 string
	}
	getPullRequestReturns struct {
		result1 *resource.PullRequest
		result2 error
	}
	getPullRequestReturnsOnCall map[int]struct {
		result1 *resource.PullRequest
		result2 error
	}
	ListModifiedFilesStub        func(int) ([]string, error)
	listModifiedFilesMutex       sync.RWMutex
	listModifiedFilesArgsForCall []struct {
		arg1 int
	}
	listModifiedFilesReturns struct {
		result1 []string
		result2 error
	}
	listModifiedFilesReturnsOnCall map[int]struct {
		result1 []string
		result2 error
	}
	ListPullRequestsStub        func([]githubv4.PullRequestState, string) ([]*resource.PullRequest, error)
	listPullRequestsMutex       sync.RWMutex
	listPullRequestsArgsForCall []struct {
		arg1 []githubv4.PullRequestState
		arg2 string
	}
	listPullRequestsReturns struct {
		result1 []*resource.PullRequest
		result2 error
	}
	listPullRequestsReturnsOnCall map[int]struct {
		result1 []*resource.PullRequest
		result2 error
	}
	PostCommentStub        func(string, string) error
	postCommentMutex       sync.RWMutex
	postCommentArgsForCall []struct {
		arg1 string
		arg2 string
	}
	postCommentReturns struct {
		result1 error
	}
	postCommentReturnsOnCall map[int]struct {
		result1 error
	}
	UpdateCommitStatusStub        func(string, string, string, string, string, string) error
	updateCommitStatusMutex       sync.RWMutex
	updateCommitStatusArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 string
		arg5 string
		arg6 string
	}
	updateCommitStatusReturns struct {
		result1 error
	}
	updateCommitStatusReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeGithub) DeletePreviousComments(arg1 string) error {
	fake.deletePreviousCommentsMutex.Lock()
	ret, specificReturn := fake.deletePreviousCommentsReturnsOnCall[len(fake.deletePreviousCommentsArgsForCall)]
	fake.deletePreviousCommentsArgsForCall = append(fake.deletePreviousCommentsArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("DeletePreviousComments", []interface{}{arg1})
	fake.deletePreviousCommentsMutex.Unlock()
	if fake.DeletePreviousCommentsStub != nil {
		return fake.DeletePreviousCommentsStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.deletePreviousCommentsReturns
	return fakeReturns.result1
}

func (fake *FakeGithub) DeletePreviousCommentsCallCount() int {
	fake.deletePreviousCommentsMutex.RLock()
	defer fake.deletePreviousCommentsMutex.RUnlock()
	return len(fake.deletePreviousCommentsArgsForCall)
}

func (fake *FakeGithub) DeletePreviousCommentsCalls(stub func(string) error) {
	fake.deletePreviousCommentsMutex.Lock()
	defer fake.deletePreviousCommentsMutex.Unlock()
	fake.DeletePreviousCommentsStub = stub
}

func (fake *FakeGithub) DeletePreviousCommentsArgsForCall(i int) string {
	fake.deletePreviousCommentsMutex.RLock()
	defer fake.deletePreviousCommentsMutex.RUnlock()
	argsForCall := fake.deletePreviousCommentsArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeGithub) DeletePreviousCommentsReturns(result1 error) {
	fake.deletePreviousCommentsMutex.Lock()
	defer fake.deletePreviousCommentsMutex.Unlock()
	fake.DeletePreviousCommentsStub = nil
	fake.deletePreviousCommentsReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGithub) DeletePreviousCommentsReturnsOnCall(i int, result1 error) {
	fake.deletePreviousCommentsMutex.Lock()
	defer fake.deletePreviousCommentsMutex.Unlock()
	fake.DeletePreviousCommentsStub = nil
	if fake.deletePreviousCommentsReturnsOnCall == nil {
		fake.deletePreviousCommentsReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deletePreviousCommentsReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGithub) GetChangedFiles(arg1 string, arg2 string) ([]resource.ChangedFileObject, error) {
	fake.getChangedFilesMutex.Lock()
	ret, specificReturn := fake.getChangedFilesReturnsOnCall[len(fake.getChangedFilesArgsForCall)]
	fake.getChangedFilesArgsForCall = append(fake.getChangedFilesArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("GetChangedFiles", []interface{}{arg1, arg2})
	fake.getChangedFilesMutex.Unlock()
	if fake.GetChangedFilesStub != nil {
		return fake.GetChangedFilesStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getChangedFilesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeGithub) GetChangedFilesCallCount() int {
	fake.getChangedFilesMutex.RLock()
	defer fake.getChangedFilesMutex.RUnlock()
	return len(fake.getChangedFilesArgsForCall)
}

func (fake *FakeGithub) GetChangedFilesCalls(stub func(string, string) ([]resource.ChangedFileObject, error)) {
	fake.getChangedFilesMutex.Lock()
	defer fake.getChangedFilesMutex.Unlock()
	fake.GetChangedFilesStub = stub
}

func (fake *FakeGithub) GetChangedFilesArgsForCall(i int) (string, string) {
	fake.getChangedFilesMutex.RLock()
	defer fake.getChangedFilesMutex.RUnlock()
	argsForCall := fake.getChangedFilesArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeGithub) GetChangedFilesReturns(result1 []resource.ChangedFileObject, result2 error) {
	fake.getChangedFilesMutex.Lock()
	defer fake.getChangedFilesMutex.Unlock()
	fake.GetChangedFilesStub = nil
	fake.getChangedFilesReturns = struct {
		result1 []resource.ChangedFileObject
		result2 error
	}{result1, result2}
}

func (fake *FakeGithub) GetChangedFilesReturnsOnCall(i int, result1 []resource.ChangedFileObject, result2 error) {
	fake.getChangedFilesMutex.Lock()
	defer fake.getChangedFilesMutex.Unlock()
	fake.GetChangedFilesStub = nil
	if fake.getChangedFilesReturnsOnCall == nil {
		fake.getChangedFilesReturnsOnCall = make(map[int]struct {
			result1 []resource.ChangedFileObject
			result2 error
		})
	}
	fake.getChangedFilesReturnsOnCall[i] = struct {
		result1 []resource.ChangedFileObject
		result2 error
	}{result1, result2}
}

func (fake *FakeGithub) GetPullRequest(arg1 string, arg2 string) (*resource.PullRequest, error) {
	fake.getPullRequestMutex.Lock()
	ret, specificReturn := fake.getPullRequestReturnsOnCall[len(fake.getPullRequestArgsForCall)]
	fake.getPullRequestArgsForCall = append(fake.getPullRequestArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("GetPullRequest", []interface{}{arg1, arg2})
	fake.getPullRequestMutex.Unlock()
	if fake.GetPullRequestStub != nil {
		return fake.GetPullRequestStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.getPullRequestReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeGithub) GetPullRequestCallCount() int {
	fake.getPullRequestMutex.RLock()
	defer fake.getPullRequestMutex.RUnlock()
	return len(fake.getPullRequestArgsForCall)
}

func (fake *FakeGithub) GetPullRequestCalls(stub func(string, string) (*resource.PullRequest, error)) {
	fake.getPullRequestMutex.Lock()
	defer fake.getPullRequestMutex.Unlock()
	fake.GetPullRequestStub = stub
}

func (fake *FakeGithub) GetPullRequestArgsForCall(i int) (string, string) {
	fake.getPullRequestMutex.RLock()
	defer fake.getPullRequestMutex.RUnlock()
	argsForCall := fake.getPullRequestArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeGithub) GetPullRequestReturns(result1 *resource.PullRequest, result2 error) {
	fake.getPullRequestMutex.Lock()
	defer fake.getPullRequestMutex.Unlock()
	fake.GetPullRequestStub = nil
	fake.getPullRequestReturns = struct {
		result1 *resource.PullRequest
		result2 error
	}{result1, result2}
}

func (fake *FakeGithub) GetPullRequestReturnsOnCall(i int, result1 *resource.PullRequest, result2 error) {
	fake.getPullRequestMutex.Lock()
	defer fake.getPullRequestMutex.Unlock()
	fake.GetPullRequestStub = nil
	if fake.getPullRequestReturnsOnCall == nil {
		fake.getPullRequestReturnsOnCall = make(map[int]struct {
			result1 *resource.PullRequest
			result2 error
		})
	}
	fake.getPullRequestReturnsOnCall[i] = struct {
		result1 *resource.PullRequest
		result2 error
	}{result1, result2}
}

func (fake *FakeGithub) ListModifiedFiles(arg1 int) ([]string, error) {
	fake.listModifiedFilesMutex.Lock()
	ret, specificReturn := fake.listModifiedFilesReturnsOnCall[len(fake.listModifiedFilesArgsForCall)]
	fake.listModifiedFilesArgsForCall = append(fake.listModifiedFilesArgsForCall, struct {
		arg1 int
	}{arg1})
	fake.recordInvocation("ListModifiedFiles", []interface{}{arg1})
	fake.listModifiedFilesMutex.Unlock()
	if fake.ListModifiedFilesStub != nil {
		return fake.ListModifiedFilesStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listModifiedFilesReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeGithub) ListModifiedFilesCallCount() int {
	fake.listModifiedFilesMutex.RLock()
	defer fake.listModifiedFilesMutex.RUnlock()
	return len(fake.listModifiedFilesArgsForCall)
}

func (fake *FakeGithub) ListModifiedFilesCalls(stub func(int) ([]string, error)) {
	fake.listModifiedFilesMutex.Lock()
	defer fake.listModifiedFilesMutex.Unlock()
	fake.ListModifiedFilesStub = stub
}

func (fake *FakeGithub) ListModifiedFilesArgsForCall(i int) int {
	fake.listModifiedFilesMutex.RLock()
	defer fake.listModifiedFilesMutex.RUnlock()
	argsForCall := fake.listModifiedFilesArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeGithub) ListModifiedFilesReturns(result1 []string, result2 error) {
	fake.listModifiedFilesMutex.Lock()
	defer fake.listModifiedFilesMutex.Unlock()
	fake.ListModifiedFilesStub = nil
	fake.listModifiedFilesReturns = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeGithub) ListModifiedFilesReturnsOnCall(i int, result1 []string, result2 error) {
	fake.listModifiedFilesMutex.Lock()
	defer fake.listModifiedFilesMutex.Unlock()
	fake.ListModifiedFilesStub = nil
	if fake.listModifiedFilesReturnsOnCall == nil {
		fake.listModifiedFilesReturnsOnCall = make(map[int]struct {
			result1 []string
			result2 error
		})
	}
	fake.listModifiedFilesReturnsOnCall[i] = struct {
		result1 []string
		result2 error
	}{result1, result2}
}

func (fake *FakeGithub) ListPullRequests(arg1 []githubv4.PullRequestState, arg2 string) ([]*resource.PullRequest, error) {
	var arg1Copy []githubv4.PullRequestState
	if arg1 != nil {
		arg1Copy = make([]githubv4.PullRequestState, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.listPullRequestsMutex.Lock()
	ret, specificReturn := fake.listPullRequestsReturnsOnCall[len(fake.listPullRequestsArgsForCall)]
	fake.listPullRequestsArgsForCall = append(fake.listPullRequestsArgsForCall, struct {
		arg1 []githubv4.PullRequestState
		arg2 string
	}{arg1Copy, arg2})
	fake.recordInvocation("ListPullRequests", []interface{}{arg1Copy, arg2})
	fake.listPullRequestsMutex.Unlock()
	if fake.ListPullRequestsStub != nil {
		return fake.ListPullRequestsStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listPullRequestsReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeGithub) ListPullRequestsCallCount() int {
	fake.listPullRequestsMutex.RLock()
	defer fake.listPullRequestsMutex.RUnlock()
	return len(fake.listPullRequestsArgsForCall)
}

func (fake *FakeGithub) ListPullRequestsCalls(stub func([]githubv4.PullRequestState, string) ([]*resource.PullRequest, error)) {
	fake.listPullRequestsMutex.Lock()
	defer fake.listPullRequestsMutex.Unlock()
	fake.ListPullRequestsStub = stub
}

func (fake *FakeGithub) ListPullRequestsArgsForCall(i int) ([]githubv4.PullRequestState, string) {
	fake.listPullRequestsMutex.RLock()
	defer fake.listPullRequestsMutex.RUnlock()
	argsForCall := fake.listPullRequestsArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeGithub) ListPullRequestsReturns(result1 []*resource.PullRequest, result2 error) {
	fake.listPullRequestsMutex.Lock()
	defer fake.listPullRequestsMutex.Unlock()
	fake.ListPullRequestsStub = nil
	fake.listPullRequestsReturns = struct {
		result1 []*resource.PullRequest
		result2 error
	}{result1, result2}
}

func (fake *FakeGithub) ListPullRequestsReturnsOnCall(i int, result1 []*resource.PullRequest, result2 error) {
	fake.listPullRequestsMutex.Lock()
	defer fake.listPullRequestsMutex.Unlock()
	fake.ListPullRequestsStub = nil
	if fake.listPullRequestsReturnsOnCall == nil {
		fake.listPullRequestsReturnsOnCall = make(map[int]struct {
			result1 []*resource.PullRequest
			result2 error
		})
	}
	fake.listPullRequestsReturnsOnCall[i] = struct {
		result1 []*resource.PullRequest
		result2 error
	}{result1, result2}
}

func (fake *FakeGithub) PostComment(arg1 string, arg2 string) error {
	fake.postCommentMutex.Lock()
	ret, specificReturn := fake.postCommentReturnsOnCall[len(fake.postCommentArgsForCall)]
	fake.postCommentArgsForCall = append(fake.postCommentArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("PostComment", []interface{}{arg1, arg2})
	fake.postCommentMutex.Unlock()
	if fake.PostCommentStub != nil {
		return fake.PostCommentStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.postCommentReturns
	return fakeReturns.result1
}

func (fake *FakeGithub) PostCommentCallCount() int {
	fake.postCommentMutex.RLock()
	defer fake.postCommentMutex.RUnlock()
	return len(fake.postCommentArgsForCall)
}

func (fake *FakeGithub) PostCommentCalls(stub func(string, string) error) {
	fake.postCommentMutex.Lock()
	defer fake.postCommentMutex.Unlock()
	fake.PostCommentStub = stub
}

func (fake *FakeGithub) PostCommentArgsForCall(i int) (string, string) {
	fake.postCommentMutex.RLock()
	defer fake.postCommentMutex.RUnlock()
	argsForCall := fake.postCommentArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeGithub) PostCommentReturns(result1 error) {
	fake.postCommentMutex.Lock()
	defer fake.postCommentMutex.Unlock()
	fake.PostCommentStub = nil
	fake.postCommentReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGithub) PostCommentReturnsOnCall(i int, result1 error) {
	fake.postCommentMutex.Lock()
	defer fake.postCommentMutex.Unlock()
	fake.PostCommentStub = nil
	if fake.postCommentReturnsOnCall == nil {
		fake.postCommentReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.postCommentReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGithub) UpdateCommitStatus(arg1 string, arg2 string, arg3 string, arg4 string, arg5 string, arg6 string) error {
	fake.updateCommitStatusMutex.Lock()
	ret, specificReturn := fake.updateCommitStatusReturnsOnCall[len(fake.updateCommitStatusArgsForCall)]
	fake.updateCommitStatusArgsForCall = append(fake.updateCommitStatusArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
		arg4 string
		arg5 string
		arg6 string
	}{arg1, arg2, arg3, arg4, arg5, arg6})
	fake.recordInvocation("UpdateCommitStatus", []interface{}{arg1, arg2, arg3, arg4, arg5, arg6})
	fake.updateCommitStatusMutex.Unlock()
	if fake.UpdateCommitStatusStub != nil {
		return fake.UpdateCommitStatusStub(arg1, arg2, arg3, arg4, arg5, arg6)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.updateCommitStatusReturns
	return fakeReturns.result1
}

func (fake *FakeGithub) UpdateCommitStatusCallCount() int {
	fake.updateCommitStatusMutex.RLock()
	defer fake.updateCommitStatusMutex.RUnlock()
	return len(fake.updateCommitStatusArgsForCall)
}

func (fake *FakeGithub) UpdateCommitStatusCalls(stub func(string, string, string, string, string, string) error) {
	fake.updateCommitStatusMutex.Lock()
	defer fake.updateCommitStatusMutex.Unlock()
	fake.UpdateCommitStatusStub = stub
}

func (fake *FakeGithub) UpdateCommitStatusArgsForCall(i int) (string, string, string, string, string, string) {
	fake.updateCommitStatusMutex.RLock()
	defer fake.updateCommitStatusMutex.RUnlock()
	argsForCall := fake.updateCommitStatusArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5, argsForCall.arg6
}

func (fake *FakeGithub) UpdateCommitStatusReturns(result1 error) {
	fake.updateCommitStatusMutex.Lock()
	defer fake.updateCommitStatusMutex.Unlock()
	fake.UpdateCommitStatusStub = nil
	fake.updateCommitStatusReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGithub) UpdateCommitStatusReturnsOnCall(i int, result1 error) {
	fake.updateCommitStatusMutex.Lock()
	defer fake.updateCommitStatusMutex.Unlock()
	fake.UpdateCommitStatusStub = nil
	if fake.updateCommitStatusReturnsOnCall == nil {
		fake.updateCommitStatusReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.updateCommitStatusReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGithub) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.deletePreviousCommentsMutex.RLock()
	defer fake.deletePreviousCommentsMutex.RUnlock()
	fake.getChangedFilesMutex.RLock()
	defer fake.getChangedFilesMutex.RUnlock()
	fake.getPullRequestMutex.RLock()
	defer fake.getPullRequestMutex.RUnlock()
	fake.listModifiedFilesMutex.RLock()
	defer fake.listModifiedFilesMutex.RUnlock()
	fake.listPullRequestsMutex.RLock()
	defer fake.listPullRequestsMutex.RUnlock()
	fake.postCommentMutex.RLock()
	defer fake.postCommentMutex.RUnlock()
	fake.updateCommitStatusMutex.RLock()
	defer fake.updateCommitStatusMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeGithub) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ resource.Github = new(FakeGithub)
