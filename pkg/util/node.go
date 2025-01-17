/*
 Copyright 2022 The Koordinator Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package util

import (
	"fmt"
	"path"

	corev1 "k8s.io/api/core/v1"

	slov1alpha1 "github.com/koordinator-sh/koordinator/apis/slo/v1alpha1"
	"github.com/koordinator-sh/koordinator/pkg/util/system"
)

// GetNodeKey returns a generated key with given node
func GetNodeKey(node *corev1.Node) string {
	return fmt.Sprintf("%v/%v", node.GetNamespace(), node.GetName())
}

// GetNodeMetricKey returns a generated key with given nodeMetric
func GetNodeMetricKey(nodeMetric *slov1alpha1.NodeMetric) string {
	return fmt.Sprintf("%v/%v", nodeMetric.GetNamespace(), nodeMetric.GetName())
}

// @output like kubepods.slice/kubepods-besteffort.slice/
func GetKubeQosRelativePath(qosClass corev1.PodQOSClass) string {
	return GetPodCgroupDirWithKube(system.CgroupPathFormatter.QOSDirFn(qosClass))
}

// GetRootCgroupCPUSetDir gets the cpuset parent directory of the specified podQos' root cgroup
// @output /sys/fs/cgroup/cpuset/kubepods.slice/kubepods-besteffort.slice
func GetRootCgroupCPUSetDir(qosClass corev1.PodQOSClass) string {
	rootCgroupParentDir := GetKubeQosRelativePath(qosClass)
	return path.Join(system.Conf.CgroupRootDir, system.CgroupCPUSetDir, rootCgroupParentDir)
}

// GetRootCgroupCurCPUSet gets the current cpuset of the specified podQos' root cgroup
func GetRootCgroupCurCPUSet(qosClass corev1.PodQOSClass) ([]int32, error) {
	rawContent, err := system.CgroupFileRead(GetKubeQosRelativePath(qosClass), system.CPUSet)
	if err != nil {
		return nil, err
	}

	return ParseCPUSetStr(rawContent)
}
