// The MIT License
//
// Copyright (c) 2020 Temporal Technologies Inc.  All rights reserved.
//
// Copyright (c) 2020 Uber Technologies, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package persistence

import (
	"context"

	commonpb "go.temporal.io/api/common/v1"
	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/api/serviceerror"

	"go.temporal.io/server/common/log"
	"go.temporal.io/server/common/quotas"
)

var (
	// ErrPersistenceLimitExceeded is the error indicating QPS limit reached.
	ErrPersistenceLimitExceeded = serviceerror.NewResourceExhausted(enumspb.RESOURCE_EXHAUSTED_CAUSE_SYSTEM_OVERLOADED, "Persistence Max QPS Reached.")
)

type (
	shardRateLimitedPersistenceClient struct {
		rateLimiter quotas.RateLimiter
		persistence ShardManager
		logger      log.Logger
	}

	executionRateLimitedPersistenceClient struct {
		rateLimiter quotas.RateLimiter
		persistence ExecutionManager
		logger      log.Logger
	}

	taskRateLimitedPersistenceClient struct {
		rateLimiter quotas.RateLimiter
		persistence TaskManager
		logger      log.Logger
	}

	metadataRateLimitedPersistenceClient struct {
		rateLimiter quotas.RateLimiter
		persistence MetadataManager
		logger      log.Logger
	}

	clusterMetadataRateLimitedPersistenceClient struct {
		rateLimiter quotas.RateLimiter
		persistence ClusterMetadataManager
		logger      log.Logger
	}

	queueRateLimitedPersistenceClient struct {
		rateLimiter quotas.RateLimiter
		persistence Queue
		logger      log.Logger
	}
)

var _ ShardManager = (*shardRateLimitedPersistenceClient)(nil)
var _ ExecutionManager = (*executionRateLimitedPersistenceClient)(nil)
var _ TaskManager = (*taskRateLimitedPersistenceClient)(nil)
var _ MetadataManager = (*metadataRateLimitedPersistenceClient)(nil)
var _ ClusterMetadataManager = (*clusterMetadataRateLimitedPersistenceClient)(nil)
var _ Queue = (*queueRateLimitedPersistenceClient)(nil)

// NewShardPersistenceRateLimitedClient creates a client to manage shards
func NewShardPersistenceRateLimitedClient(persistence ShardManager, rateLimiter quotas.RateLimiter, logger log.Logger) ShardManager {
	return &shardRateLimitedPersistenceClient{
		persistence: persistence,
		rateLimiter: rateLimiter,
		logger:      logger,
	}
}

// NewExecutionPersistenceRateLimitedClient creates a client to manage executions
func NewExecutionPersistenceRateLimitedClient(persistence ExecutionManager, rateLimiter quotas.RateLimiter, logger log.Logger) ExecutionManager {
	return &executionRateLimitedPersistenceClient{
		persistence: persistence,
		rateLimiter: rateLimiter,
		logger:      logger,
	}
}

// NewTaskPersistenceRateLimitedClient creates a client to manage tasks
func NewTaskPersistenceRateLimitedClient(persistence TaskManager, rateLimiter quotas.RateLimiter, logger log.Logger) TaskManager {
	return &taskRateLimitedPersistenceClient{
		persistence: persistence,
		rateLimiter: rateLimiter,
		logger:      logger,
	}
}

// NewMetadataPersistenceRateLimitedClient creates a MetadataManager client to manage metadata
func NewMetadataPersistenceRateLimitedClient(persistence MetadataManager, rateLimiter quotas.RateLimiter, logger log.Logger) MetadataManager {
	return &metadataRateLimitedPersistenceClient{
		persistence: persistence,
		rateLimiter: rateLimiter,
		logger:      logger,
	}
}

// NewClusterMetadataPersistenceRateLimitedClient creates a MetadataManager client to manage metadata
func NewClusterMetadataPersistenceRateLimitedClient(persistence ClusterMetadataManager, rateLimiter quotas.RateLimiter, logger log.Logger) ClusterMetadataManager {
	return &clusterMetadataRateLimitedPersistenceClient{
		persistence: persistence,
		rateLimiter: rateLimiter,
		logger:      logger,
	}
}

// NewQueuePersistenceRateLimitedClient creates a client to manage queue
func NewQueuePersistenceRateLimitedClient(persistence Queue, rateLimiter quotas.RateLimiter, logger log.Logger) Queue {
	return &queueRateLimitedPersistenceClient{
		persistence: persistence,
		rateLimiter: rateLimiter,
		logger:      logger,
	}
}

func (p *shardRateLimitedPersistenceClient) GetName() string {
	return p.persistence.GetName()
}

func (p *shardRateLimitedPersistenceClient) GetOrCreateShard(request *GetOrCreateShardRequest) (*GetOrCreateShardResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	response, err := p.persistence.GetOrCreateShard(request)
	return response, err
}

func (p *shardRateLimitedPersistenceClient) UpdateShard(request *UpdateShardRequest) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	err := p.persistence.UpdateShard(request)
	return err
}

func (p *shardRateLimitedPersistenceClient) Close() {
	p.persistence.Close()
}

func (p *executionRateLimitedPersistenceClient) GetName() string {
	return p.persistence.GetName()
}

func (p *executionRateLimitedPersistenceClient) CreateWorkflowExecution(
	ctx context.Context,
	request *CreateWorkflowExecutionRequest,
) (*CreateWorkflowExecutionResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	response, err := p.persistence.CreateWorkflowExecution(ctx, request)
	return response, err
}

func (p *executionRateLimitedPersistenceClient) GetWorkflowExecution(
	ctx context.Context,
	request *GetWorkflowExecutionRequest,
) (*GetWorkflowExecutionResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	response, err := p.persistence.GetWorkflowExecution(ctx, request)
	return response, err
}

func (p *executionRateLimitedPersistenceClient) SetWorkflowExecution(
	ctx context.Context,
	request *SetWorkflowExecutionRequest,
) (*SetWorkflowExecutionResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	response, err := p.persistence.SetWorkflowExecution(ctx, request)
	return response, err
}

func (p *executionRateLimitedPersistenceClient) UpdateWorkflowExecution(
	ctx context.Context,
	request *UpdateWorkflowExecutionRequest,
) (*UpdateWorkflowExecutionResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	resp, err := p.persistence.UpdateWorkflowExecution(ctx, request)
	return resp, err
}

func (p *executionRateLimitedPersistenceClient) ConflictResolveWorkflowExecution(
	ctx context.Context,
	request *ConflictResolveWorkflowExecutionRequest,
) (*ConflictResolveWorkflowExecutionResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	response, err := p.persistence.ConflictResolveWorkflowExecution(ctx, request)
	return response, err
}

func (p *executionRateLimitedPersistenceClient) DeleteWorkflowExecution(
	ctx context.Context,
	request *DeleteWorkflowExecutionRequest,
) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	err := p.persistence.DeleteWorkflowExecution(ctx, request)
	return err
}

func (p *executionRateLimitedPersistenceClient) DeleteCurrentWorkflowExecution(
	ctx context.Context,
	request *DeleteCurrentWorkflowExecutionRequest,
) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	err := p.persistence.DeleteCurrentWorkflowExecution(ctx, request)
	return err
}

func (p *executionRateLimitedPersistenceClient) GetCurrentExecution(
	ctx context.Context,
	request *GetCurrentExecutionRequest,
) (*GetCurrentExecutionResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	response, err := p.persistence.GetCurrentExecution(ctx, request)
	return response, err
}

func (p *executionRateLimitedPersistenceClient) ListConcreteExecutions(
	ctx context.Context,
	request *ListConcreteExecutionsRequest,
) (*ListConcreteExecutionsResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	response, err := p.persistence.ListConcreteExecutions(ctx, request)
	return response, err
}

func (p *executionRateLimitedPersistenceClient) AddHistoryTasks(
	ctx context.Context,
	request *AddHistoryTasksRequest,
) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	err := p.persistence.AddHistoryTasks(ctx, request)
	return err
}

func (p *executionRateLimitedPersistenceClient) GetHistoryTask(
	ctx context.Context,
	request *GetHistoryTaskRequest,
) (*GetHistoryTaskResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	response, err := p.persistence.GetHistoryTask(ctx, request)
	return response, err
}

func (p *executionRateLimitedPersistenceClient) GetHistoryTasks(
	ctx context.Context,
	request *GetHistoryTasksRequest,
) (*GetHistoryTasksResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	response, err := p.persistence.GetHistoryTasks(ctx, request)
	return response, err
}

func (p *executionRateLimitedPersistenceClient) CompleteHistoryTask(
	ctx context.Context,
	request *CompleteHistoryTaskRequest,
) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	err := p.persistence.CompleteHistoryTask(ctx, request)
	return err
}

func (p *executionRateLimitedPersistenceClient) RangeCompleteHistoryTasks(
	ctx context.Context,
	request *RangeCompleteHistoryTasksRequest,
) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	err := p.persistence.RangeCompleteHistoryTasks(ctx, request)
	return err
}

func (p *executionRateLimitedPersistenceClient) PutReplicationTaskToDLQ(
	ctx context.Context,
	request *PutReplicationTaskToDLQRequest,
) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	return p.persistence.PutReplicationTaskToDLQ(ctx, request)
}

func (p *executionRateLimitedPersistenceClient) GetReplicationTasksFromDLQ(
	ctx context.Context,
	request *GetReplicationTasksFromDLQRequest,
) (*GetHistoryTasksResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	return p.persistence.GetReplicationTasksFromDLQ(ctx, request)
}

func (p *executionRateLimitedPersistenceClient) DeleteReplicationTaskFromDLQ(
	ctx context.Context,
	request *DeleteReplicationTaskFromDLQRequest,
) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	return p.persistence.DeleteReplicationTaskFromDLQ(ctx, request)
}

func (p *executionRateLimitedPersistenceClient) RangeDeleteReplicationTaskFromDLQ(
	ctx context.Context,
	request *RangeDeleteReplicationTaskFromDLQRequest,
) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	return p.persistence.RangeDeleteReplicationTaskFromDLQ(ctx, request)
}

func (p *executionRateLimitedPersistenceClient) Close() {
	p.persistence.Close()
}

func (p *taskRateLimitedPersistenceClient) GetName() string {
	return p.persistence.GetName()
}

func (p *taskRateLimitedPersistenceClient) CreateTasks(request *CreateTasksRequest) (*CreateTasksResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	response, err := p.persistence.CreateTasks(request)
	return response, err
}

func (p *taskRateLimitedPersistenceClient) GetTasks(request *GetTasksRequest) (*GetTasksResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	response, err := p.persistence.GetTasks(request)
	return response, err
}

func (p *taskRateLimitedPersistenceClient) CompleteTask(request *CompleteTaskRequest) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	err := p.persistence.CompleteTask(request)
	return err
}

func (p *taskRateLimitedPersistenceClient) CompleteTasksLessThan(request *CompleteTasksLessThanRequest) (int, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return 0, ErrPersistenceLimitExceeded
	}
	return p.persistence.CompleteTasksLessThan(request)
}

func (p *taskRateLimitedPersistenceClient) CreateTaskQueue(request *CreateTaskQueueRequest) (*CreateTaskQueueResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	return p.persistence.CreateTaskQueue(request)
}

func (p *taskRateLimitedPersistenceClient) UpdateTaskQueue(request *UpdateTaskQueueRequest) (*UpdateTaskQueueResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	return p.persistence.UpdateTaskQueue(request)
}

func (p *taskRateLimitedPersistenceClient) GetTaskQueue(request *GetTaskQueueRequest) (*GetTaskQueueResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	return p.persistence.GetTaskQueue(request)
}

func (p *taskRateLimitedPersistenceClient) ListTaskQueue(request *ListTaskQueueRequest) (*ListTaskQueueResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	return p.persistence.ListTaskQueue(request)
}

func (p *taskRateLimitedPersistenceClient) DeleteTaskQueue(request *DeleteTaskQueueRequest) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}
	return p.persistence.DeleteTaskQueue(request)
}

func (p *taskRateLimitedPersistenceClient) Close() {
	p.persistence.Close()
}

func (p *metadataRateLimitedPersistenceClient) GetName() string {
	return p.persistence.GetName()
}

func (p *metadataRateLimitedPersistenceClient) CreateNamespace(request *CreateNamespaceRequest) (*CreateNamespaceResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	response, err := p.persistence.CreateNamespace(request)
	return response, err
}

func (p *metadataRateLimitedPersistenceClient) GetNamespace(request *GetNamespaceRequest) (*GetNamespaceResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	response, err := p.persistence.GetNamespace(request)
	return response, err
}

func (p *metadataRateLimitedPersistenceClient) UpdateNamespace(request *UpdateNamespaceRequest) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	err := p.persistence.UpdateNamespace(request)
	return err
}

func (p *metadataRateLimitedPersistenceClient) RenameNamespace(request *RenameNamespaceRequest) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	err := p.persistence.RenameNamespace(request)
	return err
}

func (p *metadataRateLimitedPersistenceClient) DeleteNamespace(request *DeleteNamespaceRequest) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	err := p.persistence.DeleteNamespace(request)
	return err
}

func (p *metadataRateLimitedPersistenceClient) DeleteNamespaceByName(request *DeleteNamespaceByNameRequest) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	err := p.persistence.DeleteNamespaceByName(request)
	return err
}

func (p *metadataRateLimitedPersistenceClient) ListNamespaces(request *ListNamespacesRequest) (*ListNamespacesResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	response, err := p.persistence.ListNamespaces(request)
	return response, err
}

func (p *metadataRateLimitedPersistenceClient) GetMetadata() (*GetMetadataResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	response, err := p.persistence.GetMetadata()
	return response, err
}

func (p *metadataRateLimitedPersistenceClient) Close() {
	p.persistence.Close()
}

// AppendHistoryNodes add a node to history node table
func (p *executionRateLimitedPersistenceClient) AppendHistoryNodes(
	ctx context.Context,
	request *AppendHistoryNodesRequest,
) (*AppendHistoryNodesResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	return p.persistence.AppendHistoryNodes(ctx, request)
}

// ReadHistoryBranch returns history node data for a branch
func (p *executionRateLimitedPersistenceClient) ReadHistoryBranch(
	ctx context.Context,
	request *ReadHistoryBranchRequest,
) (*ReadHistoryBranchResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	response, err := p.persistence.ReadHistoryBranch(ctx, request)
	return response, err
}

// ReadHistoryBranch returns history node data for a branch
func (p *executionRateLimitedPersistenceClient) ReadHistoryBranchReverse(
	ctx context.Context,
	request *ReadHistoryBranchReverseRequest,
) (*ReadHistoryBranchReverseResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	response, err := p.persistence.ReadHistoryBranchReverse(ctx, request)
	return response, err
}

// ReadHistoryBranchByBatch returns history node data for a branch
func (p *executionRateLimitedPersistenceClient) ReadHistoryBranchByBatch(
	ctx context.Context,
	request *ReadHistoryBranchRequest,
) (*ReadHistoryBranchByBatchResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	response, err := p.persistence.ReadHistoryBranchByBatch(ctx, request)
	return response, err
}

// ReadHistoryBranchByBatch returns history node data for a branch
func (p *executionRateLimitedPersistenceClient) ReadRawHistoryBranch(
	ctx context.Context,
	request *ReadHistoryBranchRequest,
) (*ReadRawHistoryBranchResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	response, err := p.persistence.ReadRawHistoryBranch(ctx, request)
	return response, err
}

// ForkHistoryBranch forks a new branch from a old branch
func (p *executionRateLimitedPersistenceClient) ForkHistoryBranch(
	ctx context.Context,
	request *ForkHistoryBranchRequest,
) (*ForkHistoryBranchResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	response, err := p.persistence.ForkHistoryBranch(ctx, request)
	return response, err
}

// DeleteHistoryBranch removes a branch
func (p *executionRateLimitedPersistenceClient) DeleteHistoryBranch(
	ctx context.Context,
	request *DeleteHistoryBranchRequest,
) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}
	err := p.persistence.DeleteHistoryBranch(ctx, request)
	return err
}

// TrimHistoryBranch trims a branch
func (p *executionRateLimitedPersistenceClient) TrimHistoryBranch(
	ctx context.Context,
	request *TrimHistoryBranchRequest,
) (*TrimHistoryBranchResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	resp, err := p.persistence.TrimHistoryBranch(ctx, request)
	return resp, err
}

// GetHistoryTree returns all branch information of a tree
func (p *executionRateLimitedPersistenceClient) GetHistoryTree(
	ctx context.Context,
	request *GetHistoryTreeRequest,
) (*GetHistoryTreeResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	response, err := p.persistence.GetHistoryTree(ctx, request)
	return response, err
}

func (p *executionRateLimitedPersistenceClient) GetAllHistoryTreeBranches(
	ctx context.Context,
	request *GetAllHistoryTreeBranchesRequest,
) (*GetAllHistoryTreeBranchesResponse, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	response, err := p.persistence.GetAllHistoryTreeBranches(ctx, request)
	return response, err
}

func (p *queueRateLimitedPersistenceClient) EnqueueMessage(blob commonpb.DataBlob) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	return p.persistence.EnqueueMessage(blob)
}

func (p *queueRateLimitedPersistenceClient) ReadMessages(lastMessageID int64, maxCount int) ([]*QueueMessage, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	return p.persistence.ReadMessages(lastMessageID, maxCount)
}

func (p *queueRateLimitedPersistenceClient) UpdateAckLevel(metadata *InternalQueueMetadata) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	return p.persistence.UpdateAckLevel(metadata)
}

func (p *queueRateLimitedPersistenceClient) GetAckLevels() (*InternalQueueMetadata, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	return p.persistence.GetAckLevels()
}

func (p *queueRateLimitedPersistenceClient) DeleteMessagesBefore(messageID int64) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	return p.persistence.DeleteMessagesBefore(messageID)
}

func (p *queueRateLimitedPersistenceClient) EnqueueMessageToDLQ(blob commonpb.DataBlob) (int64, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return EmptyQueueMessageID, ErrPersistenceLimitExceeded
	}

	return p.persistence.EnqueueMessageToDLQ(blob)
}

func (p *queueRateLimitedPersistenceClient) ReadMessagesFromDLQ(firstMessageID int64, lastMessageID int64, pageSize int, pageToken []byte) ([]*QueueMessage, []byte, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, nil, ErrPersistenceLimitExceeded
	}

	return p.persistence.ReadMessagesFromDLQ(firstMessageID, lastMessageID, pageSize, pageToken)
}

func (p *queueRateLimitedPersistenceClient) RangeDeleteMessagesFromDLQ(firstMessageID int64, lastMessageID int64) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	return p.persistence.RangeDeleteMessagesFromDLQ(firstMessageID, lastMessageID)
}
func (p *queueRateLimitedPersistenceClient) UpdateDLQAckLevel(metadata *InternalQueueMetadata) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	return p.persistence.UpdateDLQAckLevel(metadata)
}

func (p *queueRateLimitedPersistenceClient) GetDLQAckLevels() (*InternalQueueMetadata, error) {
	if ok := p.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}

	return p.persistence.GetDLQAckLevels()
}

func (p *queueRateLimitedPersistenceClient) DeleteMessageFromDLQ(messageID int64) error {
	if ok := p.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}

	return p.persistence.DeleteMessageFromDLQ(messageID)
}

func (p *queueRateLimitedPersistenceClient) Close() {
	p.persistence.Close()
}

func (p *queueRateLimitedPersistenceClient) Init(blob *commonpb.DataBlob) error {
	return p.persistence.Init(blob)
}

func (c *clusterMetadataRateLimitedPersistenceClient) Close() {
	c.persistence.Close()
}

func (c *clusterMetadataRateLimitedPersistenceClient) GetName() string {
	return c.persistence.GetName()
}

func (c *clusterMetadataRateLimitedPersistenceClient) GetClusterMembers(request *GetClusterMembersRequest) (*GetClusterMembersResponse, error) {
	if ok := c.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	return c.persistence.GetClusterMembers(request)
}

func (c *clusterMetadataRateLimitedPersistenceClient) UpsertClusterMembership(request *UpsertClusterMembershipRequest) error {
	if ok := c.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}
	return c.persistence.UpsertClusterMembership(request)
}

func (c *clusterMetadataRateLimitedPersistenceClient) PruneClusterMembership(request *PruneClusterMembershipRequest) error {
	if ok := c.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}
	return c.persistence.PruneClusterMembership(request)
}

func (c *clusterMetadataRateLimitedPersistenceClient) ListClusterMetadata(request *ListClusterMetadataRequest) (*ListClusterMetadataResponse, error) {
	if ok := c.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	return c.persistence.ListClusterMetadata(request)
}

func (c *clusterMetadataRateLimitedPersistenceClient) GetCurrentClusterMetadata() (*GetClusterMetadataResponse, error) {
	if ok := c.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	return c.persistence.GetCurrentClusterMetadata()
}

func (c *clusterMetadataRateLimitedPersistenceClient) GetClusterMetadata(request *GetClusterMetadataRequest) (*GetClusterMetadataResponse, error) {
	if ok := c.rateLimiter.Allow(); !ok {
		return nil, ErrPersistenceLimitExceeded
	}
	return c.persistence.GetClusterMetadata(request)
}

func (c *clusterMetadataRateLimitedPersistenceClient) SaveClusterMetadata(request *SaveClusterMetadataRequest) (bool, error) {
	if ok := c.rateLimiter.Allow(); !ok {
		return false, ErrPersistenceLimitExceeded
	}
	return c.persistence.SaveClusterMetadata(request)
}

func (c *clusterMetadataRateLimitedPersistenceClient) DeleteClusterMetadata(request *DeleteClusterMetadataRequest) error {
	if ok := c.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}
	return c.persistence.DeleteClusterMetadata(request)
}

func (c *metadataRateLimitedPersistenceClient) InitializeSystemNamespaces(currentClusterName string) error {
	if ok := c.rateLimiter.Allow(); !ok {
		return ErrPersistenceLimitExceeded
	}
	return c.persistence.InitializeSystemNamespaces(currentClusterName)
}
