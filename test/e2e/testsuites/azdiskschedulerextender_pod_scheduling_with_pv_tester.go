/*
Copyright 2020 The Kubernetes Authors.

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

package testsuites

import (
	"github.com/onsi/ginkgo/v2"
	v1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	"sigs.k8s.io/azuredisk-csi-driver/test/e2e/driver"
	"sigs.k8s.io/azuredisk-csi-driver/test/resources"
	nodeutil "sigs.k8s.io/azuredisk-csi-driver/test/utils/node"
)

// AzDiskSchedulerExtenderPodSchedulingWithPVTest will provision required PV(s), PVC(s) and Pod(s)
// Pod with PV should successfully be scheduled in a cluster with AzDriverNode and AzVolumeAttachment resources
type AzDiskSchedulerExtenderPodSchedulingWithPVTest struct {
	CSIDriver              driver.DynamicPVTestDriver
	Pod                    resources.PodDetails
	StorageClassParameters map[string]string
}

func (t *AzDiskSchedulerExtenderPodSchedulingWithPVTest) Run(client clientset.Interface, namespace *v1.Namespace, schedulerName string) {
	tpod := resources.NewTestPod(client, namespace, t.Pod.Cmd, schedulerName, t.Pod.IsWindows, t.Pod.WinServerVer)
	volume := t.Pod.Volumes[0]
	tpvc, pvcCleanup := volume.SetupDynamicPersistentVolumeClaim(client, namespace, t.CSIDriver, t.StorageClassParameters)
	for i := range pvcCleanup {
		i := i
		defer pvcCleanup[i]()
	}
	volumeName := volume.VolumeMount.NameGenerate + "1"
	tpod.SetupVolume(tpvc.PersistentVolumeClaim, volumeName, volume.VolumeMount.MountPathGenerate+"1", volume.VolumeMount.ReadOnly)

	// Get the list of available nodes for scheduling the pod
	nodeNames := nodeutil.ListAgentNodeNames(client, t.Pod.IsWindows)
	if len(nodeNames) < 1 {
		ginkgo.Skip("need at least 1 agent node to verify the test case. Current agent node count is %d", len(nodeNames))
	}

	ginkgo.By("deploying the pod")
	tpod.Create()
	defer tpod.Cleanup()
	ginkgo.By("checking that the pod's command exits with no error")
	tpod.WaitForSuccess()
}
