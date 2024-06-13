// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package chi

import (
	"strings"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	api "github.com/altinity/clickhouse-operator/pkg/apis/clickhouse.altinity.com/v1"
	"github.com/altinity/clickhouse-operator/pkg/util"
)

// ActionPlan is an action plan with a list of differences between two CHIs
type ActionPlan struct {
	old *api.ClickHouseInstallation
	new *api.ClickHouseInstallation

	specDiff  util.DiffResult
	specEqual bool

	labelsDiff  util.DiffResult
	labelsEqual bool

	deletionTimestampDiff  util.DiffResult
	deletionTimestampEqual bool

	finalizersDiff  util.DiffResult
	finalizersEqual bool

	attributesDiff  util.DiffResult
	attributesEqual bool
}

// NewActionPlan makes new ActionPlan out of two CHIs
func NewActionPlan(old, new *api.ClickHouseInstallation) *ActionPlan {
	ap := &ActionPlan{
		old: old,
		new: new,
	}

	if (old != nil) && (new != nil) {
		ap.specDiff, ap.specEqual = util.DeepDiff(ap.old.Spec, ap.new.Spec)
		ap.labelsDiff, ap.labelsEqual = util.DeepDiff(ap.old.Labels, ap.new.Labels)
		ap.deletionTimestampEqual = ap.timestampEqual(ap.old.DeletionTimestamp, ap.new.DeletionTimestamp)
		ap.deletionTimestampDiff, _ = util.DeepDiff(ap.old.DeletionTimestamp, ap.new.DeletionTimestamp)
		ap.finalizersDiff, ap.finalizersEqual = util.DeepDiff(ap.old.Finalizers, ap.new.Finalizers)
		ap.attributesDiff, ap.attributesEqual = util.DeepDiff(ap.old.EnsureRuntime().GetAttributes(), ap.new.EnsureRuntime().GetAttributes())
	} else if old == nil {
		ap.specDiff, ap.specEqual = util.DeepDiff(nil, ap.new.Spec)
		ap.labelsDiff, ap.labelsEqual = util.DeepDiff(nil, ap.new.Labels)
		ap.deletionTimestampEqual = ap.timestampEqual(nil, ap.new.DeletionTimestamp)
		ap.deletionTimestampDiff, _ = util.DeepDiff(nil, ap.new.DeletionTimestamp)
		ap.finalizersDiff, ap.finalizersEqual = util.DeepDiff(nil, ap.new.Finalizers)
		ap.attributesDiff, ap.attributesEqual = util.DeepDiff(nil, ap.new.EnsureRuntime().GetAttributes())
	} else if new == nil {
		ap.specDiff, ap.specEqual = util.DeepDiff(ap.old.Spec, nil)
		ap.labelsDiff, ap.labelsEqual = util.DeepDiff(ap.old.Labels, nil)
		ap.deletionTimestampEqual = ap.timestampEqual(ap.old.DeletionTimestamp, nil)
		ap.deletionTimestampDiff, _ = util.DeepDiff(ap.old.DeletionTimestamp, nil)
		ap.finalizersDiff, ap.finalizersEqual = util.DeepDiff(ap.old.Finalizers, nil)
		ap.attributesDiff, ap.attributesEqual = util.DeepDiff(ap.old.EnsureRuntime().GetAttributes(), nil)
	} else {
		// Both are nil
		ap.specDiff = util.DiffResult{}
		ap.specEqual = true

		ap.labelsDiff = util.DiffResult{}
		ap.labelsEqual = true

		ap.deletionTimestampDiff = util.DiffResult{}
		ap.deletionTimestampEqual = true

		ap.finalizersDiff = util.DiffResult{}
		ap.finalizersEqual = true

		ap.attributesDiff = util.DiffResult{}
		ap.attributesEqual = true
	}

	ap.excludePaths()

	return ap
}

func (ap *ActionPlan) timestampEqual(old, new *meta.Time) bool {
	switch {
	case (old == nil) && (new == nil):
		// Both are useless
		return true
	case (old == nil) && (new != nil):
		// Timestamp assigned
		return false
	case (old != nil) && (new == nil):
		// Timestamp unassigned
		return false
	}
	return old.Equal(new)
}

// excludePaths - sanitize diff - do not pay attention to changes in some paths, such as
// ObjectMeta.ResourceVersion
func (ap *ActionPlan) excludePaths() {
	if len(ap.specDiff.Modified) == 0 {
		return
	}

	excludePaths := make([]string, 0)
	// Walk over all .diff.Modified paths and find .ObjectMeta.ResourceVersion path
	for path := range ap.specDiff.Modified {
		pathParts := strings.Split(path, ".")
		for i := range pathParts {
			pathNodeCurr := pathParts[i]
			pathNodePrev := ""
			if i > 0 {
				// We have prev node
				pathNodePrev = pathParts[i-1]
			}

			if ap.isExcludedPath(pathNodePrev, pathNodeCurr) {
				// This path should be excluded from Modified
				excludePaths = append(excludePaths, path)
				break
			}
		}
	}

	// Exclude paths from diff.Modified
	for _, path := range excludePaths {
		delete(ap.specDiff.Modified, path)
	}
}

// isExcludedPath checks whether path is excluded
func (ap *ActionPlan) isExcludedPath(prev, cur string) bool {
	if ((prev == "ObjectMeta") && (cur == ".ResourceVersion")) ||
		((prev == ".ObjectMeta") && (cur == ".ResourceVersion")) {
		return true
	}

	if ((prev == "Status") && (cur == "Status")) ||
		((prev == ".Status") && (cur == ".Status")) {
		return true
	}

	return false
}

// HasActionsToDo checks whether there are any actions to do - meaning changes between states to reconcile
func (ap *ActionPlan) HasActionsToDo() bool {
	if ap.specEqual && ap.labelsEqual && ap.deletionTimestampEqual && ap.finalizersEqual && ap.attributesEqual {
		// All is equal - no actions to do
		return false
	}

	// Something is not equal

	if len(ap.specDiff.Added)+len(ap.specDiff.Removed)+len(ap.specDiff.Modified) > 0 {
		// Spec section has some modifications
		return true
	}

	if len(ap.labelsDiff.Added)+len(ap.labelsDiff.Removed)+len(ap.labelsDiff.Modified) > 0 {
		// Labels section has some modifications
		return true
	}

	return !ap.deletionTimestampEqual || !ap.finalizersEqual || !ap.attributesEqual
}

// String stringifies ActionPlan
func (ap *ActionPlan) String() string {
	if !ap.HasActionsToDo() {
		return ""
	}

	str := ""

	if len(ap.specDiff.Added) > 0 {
		// Something added
		str += util.MessageDiffItemString("added spec items", "none", "", ap.specDiff.Added)
	}

	if len(ap.specDiff.Removed) > 0 {
		// Something removed
		str += util.MessageDiffItemString("removed spec items", "none", "", ap.specDiff.Removed)
	}

	if len(ap.specDiff.Modified) > 0 {
		// Something modified
		str += util.MessageDiffItemString("modified spec items", "none", "", ap.specDiff.Modified)
	}

	if len(ap.labelsDiff.Added) > 0 {
		// Something added
		str += "added labels\n"
	}

	if len(ap.labelsDiff.Removed) > 0 {
		// Something removed
		str += "removed labels\n"
	}

	if len(ap.labelsDiff.Modified) > 0 {
		// Something modified
		str += "modified labels\n"
	}

	if !ap.deletionTimestampEqual {
		str += "modified deletion timestamp:\n"
		str += util.MessageDiffItemString("modified deletion timestamp", "none", ".metadata.deletionTimestamp", ap.deletionTimestampDiff.Modified)
	}

	if !ap.finalizersEqual {
		str += "modified finalizer:\n"
		str += util.MessageDiffItemString("modified finalizers", "none", ".metadata.finalizers", ap.finalizersDiff.Modified)
	}

	return str
}

// GetNewHostsNum - total number of hosts to be achieved
func (ap *ActionPlan) GetNewHostsNum() int {
	return ap.new.HostsCount()
}

// GetRemovedHostsNum - how many hosts would be removed
func (ap *ActionPlan) GetRemovedHostsNum() int {
	var count int
	ap.WalkRemoved(
		func(cluster *api.Cluster) {
			count += cluster.HostsCount()
		},
		func(shard *api.ChiShard) {
			count += shard.HostsCount()
		},
		func(host *api.ChiHost) {
			count++
		},
	)
	return count
}

// WalkRemoved walk removed cluster items
func (ap *ActionPlan) WalkRemoved(
	clusterFunc func(cluster *api.Cluster),
	shardFunc func(shard *api.ChiShard),
	hostFunc func(host *api.ChiHost),
) {
	// TODO refactor to map[string]object handling, instead of slice
	for path := range ap.specDiff.Removed {
		switch ap.specDiff.Removed[path].(type) {
		case api.Cluster:
			cluster := ap.specDiff.Removed[path].(api.Cluster)
			clusterFunc(&cluster)
		case api.ChiShard:
			shard := ap.specDiff.Removed[path].(api.ChiShard)
			shardFunc(&shard)
		case api.ChiHost:
			host := ap.specDiff.Removed[path].(api.ChiHost)
			hostFunc(&host)
		case *api.Cluster:
			cluster := ap.specDiff.Removed[path].(*api.Cluster)
			clusterFunc(cluster)
		case *api.ChiShard:
			shard := ap.specDiff.Removed[path].(*api.ChiShard)
			shardFunc(shard)
		case *api.ChiHost:
			host := ap.specDiff.Removed[path].(*api.ChiHost)
			hostFunc(host)
		}
	}
}

// WalkAdded walk added cluster items
func (ap *ActionPlan) WalkAdded(
	clusterFunc func(cluster *api.Cluster),
	shardFunc func(shard *api.ChiShard),
	hostFunc func(host *api.ChiHost),
) {
	// TODO refactor to map[string]object handling, instead of slice
	for path := range ap.specDiff.Added {
		switch ap.specDiff.Added[path].(type) {
		case api.Cluster:
			cluster := ap.specDiff.Added[path].(api.Cluster)
			clusterFunc(&cluster)
		case api.ChiShard:
			shard := ap.specDiff.Added[path].(api.ChiShard)
			shardFunc(&shard)
		case api.ChiHost:
			host := ap.specDiff.Added[path].(api.ChiHost)
			hostFunc(&host)
		case *api.Cluster:
			cluster := ap.specDiff.Added[path].(*api.Cluster)
			clusterFunc(cluster)
		case *api.ChiShard:
			shard := ap.specDiff.Added[path].(*api.ChiShard)
			shardFunc(shard)
		case *api.ChiHost:
			host := ap.specDiff.Added[path].(*api.ChiHost)
			hostFunc(host)
		}
	}
}

// WalkModified walk modified cluster items
func (ap *ActionPlan) WalkModified(
	clusterFunc func(cluster *api.Cluster),
	shardFunc func(shard *api.ChiShard),
	hostFunc func(host *api.ChiHost),
) {
	// TODO refactor to map[string]object handling, instead of slice
	for path := range ap.specDiff.Modified {
		switch ap.specDiff.Modified[path].(type) {
		case api.Cluster:
			cluster := ap.specDiff.Modified[path].(api.Cluster)
			clusterFunc(&cluster)
		case api.ChiShard:
			shard := ap.specDiff.Modified[path].(api.ChiShard)
			shardFunc(&shard)
		case api.ChiHost:
			host := ap.specDiff.Modified[path].(api.ChiHost)
			hostFunc(&host)
		case *api.Cluster:
			cluster := ap.specDiff.Modified[path].(*api.Cluster)
			clusterFunc(cluster)
		case *api.ChiShard:
			shard := ap.specDiff.Modified[path].(*api.ChiShard)
			shardFunc(shard)
		case *api.ChiHost:
			host := ap.specDiff.Modified[path].(*api.ChiHost)
			hostFunc(host)
		}
	}
}
