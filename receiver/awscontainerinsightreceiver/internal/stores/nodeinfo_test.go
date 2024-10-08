// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package stores

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	v1 "k8s.io/api/core/v1"
)

func TestSetGetCPUCapacity(t *testing.T) {
	nodeInfo := newNodeInfo("testNode1", &mockNodeInfoProvider{}, zap.NewNop())
	nodeInfo.setCPUCapacity(int(4))
	assert.Equal(t, uint64(4), nodeInfo.getCPUCapacity())

	nodeInfo.setCPUCapacity(int32(2))
	assert.Equal(t, uint64(2), nodeInfo.getCPUCapacity())

	nodeInfo.setCPUCapacity(int64(4))
	assert.Equal(t, uint64(4), nodeInfo.getCPUCapacity())

	nodeInfo.setCPUCapacity(uint(2))
	assert.Equal(t, uint64(2), nodeInfo.getCPUCapacity())

	nodeInfo.setCPUCapacity(uint32(4))
	assert.Equal(t, uint64(4), nodeInfo.getCPUCapacity())

	nodeInfo.setCPUCapacity(uint64(2))
	assert.Equal(t, uint64(2), nodeInfo.getCPUCapacity())

	// with invalid type
	nodeInfo.setCPUCapacity("2")
	assert.Equal(t, uint64(0), nodeInfo.getCPUCapacity())

	// with negative value
	nodeInfo.setCPUCapacity(int64(-2))
	assert.Equal(t, uint64(0), nodeInfo.getCPUCapacity())
	nodeInfo.setCPUCapacity(int(-3))
	assert.Equal(t, uint64(0), nodeInfo.getCPUCapacity())
	nodeInfo.setCPUCapacity(int32(-4))
	assert.Equal(t, uint64(0), nodeInfo.getCPUCapacity())
}

func TestSetGetMemCapacity(t *testing.T) {
	nodeInfo := newNodeInfo("testNode1", &mockNodeInfoProvider{}, zap.NewNop())
	nodeInfo.setMemCapacity(int(2048))
	assert.Equal(t, uint64(2048), nodeInfo.getMemCapacity())

	nodeInfo.setMemCapacity(int32(1024))
	assert.Equal(t, uint64(1024), nodeInfo.getMemCapacity())

	nodeInfo.setMemCapacity(int64(2048))
	assert.Equal(t, uint64(2048), nodeInfo.getMemCapacity())

	nodeInfo.setMemCapacity(uint(1024))
	assert.Equal(t, uint64(1024), nodeInfo.getMemCapacity())

	nodeInfo.setMemCapacity(uint32(2048))
	assert.Equal(t, uint64(2048), nodeInfo.getMemCapacity())

	nodeInfo.setMemCapacity(uint64(1024))
	assert.Equal(t, uint64(1024), nodeInfo.getMemCapacity())

	// with invalid type
	nodeInfo.setMemCapacity("2")
	assert.Equal(t, uint64(0), nodeInfo.getMemCapacity())

	// with negative value
	nodeInfo.setMemCapacity(int64(-2))
	assert.Equal(t, uint64(0), nodeInfo.getMemCapacity())
}

func TestGetNodeStatusCapacityPods(t *testing.T) {
	nodeInfo := newNodeInfo("testNode1", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCapacityPods, valid := nodeInfo.getNodeStatusCapacityPods()
	assert.True(t, valid)
	assert.Equal(t, uint64(5), nodeStatusCapacityPods)

	nodeInfo = newNodeInfo("testNodeNonExistent", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCapacityPods, valid = nodeInfo.getNodeStatusCapacityPods()
	assert.False(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCapacityPods)
}

func TestGetNodeStatusAllocatablePods(t *testing.T) {
	nodeInfo := newNodeInfo("testNode1", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusAllocatablePods, valid := nodeInfo.getNodeStatusAllocatablePods()
	assert.True(t, valid)
	assert.Equal(t, uint64(15), nodeStatusAllocatablePods)

	nodeInfo = newNodeInfo("testNodeNonExistent", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusAllocatablePods, valid = nodeInfo.getNodeStatusAllocatablePods()
	assert.False(t, valid)
	assert.Equal(t, uint64(0), nodeStatusAllocatablePods)
}

func TestGetNodeStatusCapacityGPUs(t *testing.T) {
	nodeInfo := newNodeInfo("testNode1", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCapacityGPUs, valid := nodeInfo.getNodeStatusCapacityGPUs()
	assert.True(t, valid)
	assert.Equal(t, uint64(20), nodeStatusCapacityGPUs)

	nodeInfo = newNodeInfo("testNode2", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCapacityGPUs, valid = nodeInfo.getNodeStatusCapacityGPUs()
	assert.True(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCapacityGPUs)
}

func TestGetNodeStatusCondition(t *testing.T) {
	nodeInfo := newNodeInfo("testNode1", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCondition, valid := nodeInfo.getNodeStatusCondition(v1.NodeReady)
	assert.True(t, valid)
	assert.Equal(t, uint64(1), nodeStatusCondition)
	nodeStatusCondition, valid = nodeInfo.getNodeStatusCondition(v1.NodeDiskPressure)
	assert.True(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)
	nodeStatusCondition, valid = nodeInfo.getNodeStatusCondition(v1.NodeMemoryPressure)
	assert.True(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)
	nodeStatusCondition, valid = nodeInfo.getNodeStatusCondition(v1.NodePIDPressure)
	assert.True(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)
	nodeStatusCondition, valid = nodeInfo.getNodeStatusCondition(v1.NodeNetworkUnavailable)
	assert.True(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)

	nodeInfo = newNodeInfo("testNode2", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCondition, valid = nodeInfo.getNodeStatusCondition(v1.NodeNetworkUnavailable)
	assert.False(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)

	nodeInfo = newNodeInfo("testNodeNonExistent", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCondition, valid = nodeInfo.getNodeStatusCondition(v1.NodeReady)
	assert.False(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)
}

func TestGetNodeConditionUnknown(t *testing.T) {
	nodeInfo := newNodeInfo("testNode1", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCondition, valid := nodeInfo.getNodeConditionUnknown()
	assert.True(t, valid)
	assert.Equal(t, uint64(1), nodeStatusCondition)

	nodeInfo = newNodeInfo("testNode2", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCondition, valid = nodeInfo.getNodeConditionUnknown()
	assert.True(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)

	nodeInfo = newNodeInfo("testNodeNonExistent", &mockNodeInfoProvider{}, zap.NewNop())
	nodeStatusCondition, valid = nodeInfo.getNodeStatusCondition(v1.NodeReady)
	assert.False(t, valid)
	assert.Equal(t, uint64(0), nodeStatusCondition)
}
