//go:build integration
// +build integration

/*
 *  MIT License
 *
 *  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a
 *  copy of this software and associated documentation files (the "Software"),
 *  to deal in the Software without restriction, including without limitation
 *  the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *  and/or sell copies of the Software, and to permit persons to whom the
 *  Software is furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included
 *  in all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 *  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 *  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 *  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 *  OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package services_iuf

import (
	mocks "github.com/Cray-HPE/cray-nls/src/api/mocks/services"
	services_shared "github.com/Cray-HPE/cray-nls/src/api/services/shared"
	"github.com/golang/mock/gomock"
	"github.com/spf13/viper"
	"k8s.io/client-go/kubernetes/fake"
	"log"
	"os"
	"regexp"
	"testing"

	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/Cray-HPE/cray-nls/src/utils"
	"github.com/alecthomas/assert"
	"github.com/google/uuid"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

// NOTE: This test is not run as part of the normal test suite. It is run as part of the integration test suite.
// it is mainly used for local development and debugging.
// it reads from .env so you can use this file for debugging argo templates
func TestNcnRolloutWorkflowGen(t *testing.T) {
	activityName, _, iufSvc := ncnRolloutTestSetup(t)

	t.Run("it can generate management rollout workflow", func(t *testing.T) {
		session := iuf.Session{
			Name:     uuid.NewString(),
			Products: []iuf.Product{},
			Workflows: []iuf.SessionWorkflow{
				{
					Id:  "1",
					Url: "1",
				},
			},
			CurrentStage: "management-nodes-rollout",
			CurrentState: iuf.SessionStateInProgress,
			InputParameters: iuf.InputParameters{
				Stages: []string{"management-nodes-rollout"},
			},
			ActivityRef: activityName,
		}

		workflow, err, _ := iufSvc.CreateIufWorkflow(session)
		assert.NoError(t, err)
		assert.Equal(t, "ncn-m001", workflow.Spec.NodeSelector["kubernetes.io/hostname"])
		data, _ := yaml.Marshal(&workflow)
		t.Logf("%s", string(data))
	})

}

func ncnRolloutTestSetup(t *testing.T) (string, string, IufService) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	re := regexp.MustCompile(`^(.*` + "cray-nls" + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))
	viper.SetConfigFile(string(rootPath) + "/.env")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error while reading config file %s", err)
	}
	var env utils.Env
	viper.Unmarshal(&env)

	name := uuid.NewString()
	activityData := `{"name":"fullstack-230217","input_parameters":{"media_dir":"/fullstack-230217","site_parameters":"","limit_management_nodes":["Management_Worker"],"limit_managed_nodes":["Compute"],"managed_rollout_strategy":"stage","concurrent_management_rollout_percentage":20,"media_host":"ncn-m001","concurrency":0,"bootprep_config_managed":"","bootprep_config_management":"","stages":["process-media"],"force":false},"site_parameters":{"global":{},"products":{"HFP-firmware":{"version":"23.1.2"},"analytics":{"version":"1.4.20","working_branch":"integration-1.4.20"},"cocn":{"version":"2.0.0-20230306163418-f05bc66"},"cos":{"version":"2.5.97","working_branch":"integration-2.5.97"},"cpe":{"version":"23.2.2","working_branch":"integration-23.2.2"},"cray-sdu-rda":{"version":"2.2.0-beta15"},"csm":{"version":"1.4.0-alpha.29"},"csm-diags":{"version":"1.4.18"},"default":{"network_type":"mellanox","note":"johnn-","suffix":"-fullstack-230217","wlm":"future use","working_branch":"integration-{{version_x_y_z}}"},"hfp":{"version":"22.10.0"},"hpc-csm-software-recipe":{"version":"23.4.0-alpha.1"},"pbs":{"version":"1.2.9-20230209172608.28ca7d3","working_branch":"integration-1.2.9"},"recipe":{"version":"23.4.0-alpha.1"},"sat":{"version":"2.5.14"},"sle-os-backports-15-sp3":{"version":""},"sle-os-backports-15-sp4":{"version":""},"sle-os-backports-sle-15-sp4-x86_64":{"version":"23.2.7-20230106"},"sle-os-products-15-sp3":{"version":""},"sle-os-products-15-sp4":{"version":""},"sle-os-products-15-sp4-x86_64":{"version":"23.2.7-20230106"},"sle-os-updates-15-sp3":{"version":""},"sle-os-updates-15-sp4":{"version":""},"sle-os-updates-15-sp4-x86_64":{"version":"23.2.7-20230106"},"sleupdates-15-SP4-23.2.7-20230106-supplement":{"version":"23.2.7-20230106"},"slingshot":{"version":"2.0.2-941"},"slingshot-host-software":{"version":"2.0.2-97-csm-1.4","working_branch":"integration-2.0.2"},"slurm":{"version":"1.2.9-20230209172621.c445c79","working_branch":"integration-1.2.9"},"sma":{"version":"1.8.6"},"uan":{"version":"2.6.0-alpha.5","working_branch":"integration-2.6.0"},"wlm-slurm":{"version":"1.2.9"}}},"operation_outputs":{"stage_params":{"pre-install-check":{"preflight-checks-for-services":{"preflight-checks(0)":{"output":""},"prom-metrics":{"diff-time-value":"21"}}},"prepare-images":{"prepare-managed-images":{"prom-metrics":{"diff-time-value":"4539"},"sat-bootprep-run":{"script_stdout":"{\n    \"images\": [\n        {\n            \"name\": \"lnet-johnn-cray-shasta-compute-sles15sp4.noarch-2.5.30-fullstack-230217\",\n            \"preconfigured_image_id\": null,\n            \"final_image_id\": \"6ff313ff-00ad-47d9-be40-9a33687cde30\",\n            \"configuration\": \"johnn-lnet-23.4.0-alpha.1-fullstack-230217\",\n            \"configuration_group_names\": [\n                \"Application\",\n                \"Application_LNETRouter\"\n            ]\n        },\n        {\n            \"name\": \"compute-johnn-cray-shasta-compute-sles15sp4.noarch-2.5.30-fullstack-230217\",\n            \"preconfigured_image_id\": null,\n            \"final_image_id\": \"c60f94da-8312-47ce-9aac-2bec27186742\",\n            \"configuration\": \"johnn-compute-23.4.0-alpha.1-fullstack-230217\",\n            \"configuration_group_names\": [\n                \"Compute\"\n            ]\n        },\n        {\n            \"name\": \"uan-johnn-cray-shasta-compute-sles15sp4.noarch-2.5.30-fullstack-230217\",\n            \"preconfigured_image_id\": null,\n            \"final_image_id\": \"676df716-ccc9-4fc8-a0f7-4bdd199ba2d5\",\n            \"configuration\": \"johnn-uan-23.4.0-alpha.1-fullstack-230217\",\n            \"configuration_group_names\": [\n                \"Application\",\n                \"Application_UAN\"\n            ]\n        },\n        {\n            \"name\": \"johnn-cray-shasta-compute-sles15sp4.noarch-2.5.30-fullstack-230217\",\n            \"preconfigured_image_id\": \"32c55e03-23d5-4de1-8a69-8c4130c7202f\",\n            \"final_image_id\": \"32c55e03-23d5-4de1-8a69-8c4130c7202f\",\n            \"configuration\": null,\n            \"configuration_group_names\": null\n        }\n    ],\n    \"session_templates\": [\n        {\n            \"name\": \"johnn-compute-23.4.0-alpha.1-fullstack-230217\",\n            \"configuration\": \"johnn-compute-23.4.0-alpha.1-fullstack-230217\"\n        },\n        {\n            \"name\": \"johnn-uan-23.4.0-alpha.1-fullstack-230217\",\n            \"configuration\": \"johnn-uan-23.4.0-alpha.1-fullstack-230217\"\n        },\n        {\n            \"name\": \"johnn-lnet-23.4.0-alpha.1-fullstack-230217\",\n            \"configuration\": \"johnn-lnet-23.4.0-alpha.1-fullstack-230217\"\n        }\n    ]\n}"}},"prepare-management-images":{"prom-metrics":{"diff-time-value":"3283"},"sat-bootprep-run":{"script_stdout":"{\n    \"images\": [\n        {\n            \"name\": \"johnn-storage-secure-storage-ceph-0.4.35-x86_64.squashfs-fullstack-230217\",\n            \"preconfigured_image_id\": null,\n            \"final_image_id\": \"0a1349c6-5803-474f-8afa-4d23ba4072ba\",\n            \"configuration\": \"johnn-management-23.4.0-alpha.1-fullstack-230217\",\n            \"configuration_group_names\": [\n                \"Management_Storage\"\n            ]\n        },\n        {\n            \"name\": \"johnn-worker-secure-kubernetes-0.4.35-x86_64.squashfs-fullstack-230217\",\n            \"preconfigured_image_id\": null,\n            \"final_image_id\": \"981d10ef-1b2d-4d17-8070-1ab788e25802\",\n            \"configuration\": \"johnn-management-23.4.0-alpha.1-fullstack-230217\",\n            \"configuration_group_names\": [\n                \"Management_Worker\"\n            ]\n        },\n        {\n            \"name\": \"johnn-master-secure-kubernetes-0.4.35-x86_64.squashfs-fullstack-230217\",\n            \"preconfigured_image_id\": null,\n            \"final_image_id\": \"2dad54d1-43d6-4131-9ac7-597c2a935811\",\n            \"configuration\": \"johnn-management-23.4.0-alpha.1-fullstack-230217\",\n            \"configuration_group_names\": [\n                \"Management_Master\"\n            ]\n        }\n    ]\n}"}}},"process-media":{"products":{"HFP-firmware-23-1-2":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/HFP-firmware-23.1.2"},"analytics-1-4-20":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/analytics-1.4.20"},"cos-2-5-97":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/cos-2.5.97"},"cpe-23-2-2":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/cpe-23.02-sles15-sp4"},"cray-sdu-rda-2-2-0-beta15":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/cray-sdu-rda-2.2.0-beta15"},"csm-diags-1-4-18":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/csm-diags-1.4.18"},"pbs-1-2-9-20230209172608-28ca7d3":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/wlm-pbs-1.2.9-20230209172608.28ca7d3"},"sat-2-5-14":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/sat-2.5.14"},"sle-os-backports-sle-15-sp4-x86_64-23-2-7-20230106":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/SUSE-Backports-SLE-15-SP4-x86_64-23.2.7-20230106"},"sle-os-products-15-sp4-x86_64-23-2-7-20230106":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/SUSE-Products-15-SP4-x86_64-23.2.7-20230106"},"sle-os-updates-15-sp4-x86_64-23-2-7-20230106":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/SUSE-Updates-15-SP4-x86_64-23.2.7-20230106"},"slingshot-2-0-2-941":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/slingshot-2.0.2-941"},"slingshot-host-software-2-0-2-97-cos-2-5":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/slingshot-host-software-2.0.2-97-cos-2.5"},"slingshot-host-software-2-0-2-97-csm-1-4":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/slingshot-host-software-2.0.2-97-csm-1.4"},"slurm-1-2-9-20230209172621-c445c79":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/wlm-slurm-1.2.9-20230209172621.c445c79"},"sma-1-8-6":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/sma-1.8.6"},"uan-2-6-0-alpha-5":{"parent_directory":"/etc/cray/upgrade/csm/fullstack-230217/uan-2.6.0-alpha.5"}}},"update-cfs-config":{"update-managed-cfs-config":{"prom-metrics":{"diff-time-value":"38"},"sat-bootprep-run":{"script_stdout":"{\n    \"configurations\": [\n        {\n            \"name\": \"johnn-compute-23.4.0-alpha.1-fullstack-230217\"\n        },\n        {\n            \"name\": \"johnn-lnet-23.4.0-alpha.1-fullstack-230217\"\n        },\n        {\n            \"name\": \"johnn-uan-23.4.0-alpha.1-fullstack-230217\"\n        }\n    ]\n}"}},"update-management-cfs-config":{"prom-metrics":{"diff-time-value":"57"},"sat-bootprep-run":{"script_stdout":"{\n    \"configurations\": [\n        {\n            \"name\": \"johnn-management-23.4.0-alpha.1-fullstack-230217\"\n        }\n    ]\n}"}}}}},"products":[{"name":"sle-os-updates-15-sp4-x86_64","version":"23.2.7-20230106","original_location":"/etc/cray/upgrade/csm/fullstack-230217/SUSE-Updates-15-SP4-x86_64-23.2.7-20230106","validated":true,"manifest":"{\"content\":{\"docker\":[{\"path\":\"docker/cray\"}],\"nexus_blob_stores\":{\"yaml_path\":\"nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"nexus-repositories.yaml\"},\"rpms\":[{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Basesystem-15-SP4-x86_64-Debug-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Basesystem-15-SP4-x86_64-Debug-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Basesystem-15-SP4-x86_64-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Basesystem-15-SP4-x86_64-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Containers-15-SP4-x86_64-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Containers-15-SP4-x86_64-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Desktop-Applications-15-SP4-x86_64-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Desktop-Applications-15-SP4-x86_64-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Development-Tools-15-SP4-x86_64-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Development-Tools-15-SP4-x86_64-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-HPC-15-SP4-x86_64-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-HPC-15-SP4-x86_64-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Legacy-15-SP4-x86_64-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Legacy-15-SP4-x86_64-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Packagehub-Subpackages-15-SP4-x86_64-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Packagehub-Subpackages-15-SP4-x86_64-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Public-Cloud-15-SP4-x86_64-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Public-Cloud-15-SP4-x86_64-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Python3-15-SP4-x86_64-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Python3-15-SP4-x86_64-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Server-Applications-15-SP4-x86_64-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Server-Applications-15-SP4-x86_64-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Transactional-Server-15-SP4-x86_64-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Transactional-Server-15-SP4-x86_64-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Web-Scripting-15-SP4-x86_64-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Web-Scripting-15-SP4-x86_64-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Product-HPC-15-SP4-x86_64-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Product-HPC-15-SP4-x86_64-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Product-SLES-15-SP4-x86_64-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Product-SLES-15-SP4-x86_64-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Product-WE-15-SP4-x86_64-Updates\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Product-WE-15-SP4-x86_64-Updates\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-PTF-15-SP4-x86_64-PTF\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-PTF-15-SP4-x86_64-PTF\",\"repository_type\":\"raw\"}]},\"description\":\"SUSE Linux Enterprise (SLE).\\n\",\"iuf_version\":\"^0.6.0\",\"name\":\"sle-os-updates-15-sp4-x86_64\",\"version\":\"23.2.7-20230106\"}"},{"name":"slingshot","version":"2.0.2-941","original_location":"/etc/cray/upgrade/csm/fullstack-230217/slingshot-2.0.2-941","validated":true,"manifest":"{\"content\":{\"docker\":[{\"path\":\"docker\"}],\"helm\":[{\"path\":\"helm\"}],\"loftsman\":[{\"deploy\":true,\"path\":\"manifests/slingshot.yaml\",\"use_manifestgen\":true}]},\"description\":\"Slingshot Fabric Manager\\n\",\"hooks\":{\"pre_install_check\":{\"pre\":{\"script_path\":\"hooks/pre-install-check.sh\"}}},\"iuf_version\":\"^0.1.0\",\"name\":\"slingshot\",\"version\":\"2.0.2-941\"}"},{"name":"cpe","version":"23.2.2","original_location":"/etc/cray/upgrade/csm/fullstack-230217/cpe-23.02-sles15-sp4","validated":true,"manifest":"{\"content\":{\"ims\":{\"content_dirs\":[\"ims/recipes/x86_64\"]},\"nexus_blob_stores\":{\"yaml_path\":\"nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"nexus-repositories.yaml\"},\"rpms\":[{\"path\":\"rpms/cray-sles15-sp4-cn\",\"repository_name\":\"cpe-23.02-sles15-sp4\",\"repository_type\":\"raw\"}],\"s3\":[{\"bucket\":\"boot-images\",\"key\":\"PE/CPE-base.x86_64-23.02.squashfs\",\"path\":\"squashfs/cpe-base-sles15sp4.x86_64-23.02.squashfs\"},{\"bucket\":\"boot-images\",\"key\":\"PE/CPE-amd.x86_64-23.02.squashfs\",\"path\":\"squashfs/cpe-amd-sles15sp4.x86_64-23.02.squashfs\"},{\"bucket\":\"boot-images\",\"key\":\"PE/CPE-aocc.x86_64-23.02.squashfs\",\"path\":\"squashfs/cpe-aocc-sles15sp4.x86_64-23.02.squashfs\"},{\"bucket\":\"boot-images\",\"key\":\"PE/CPE-intel.x86_64-23.02.squashfs\",\"path\":\"squashfs/cpe-intel-sles15sp4.x86_64-23.02.squashfs\"},{\"bucket\":\"boot-images\",\"key\":\"PE/CPE-nvidia.x86_64-23.02.squashfs\",\"path\":\"squashfs/cpe-nvidia-sles15sp4.x86_64-23.02.squashfs\"}],\"vcs\":{\"path\":\"vcs\"}},\"description\":\"The Cray Programming Environment (CPE).\\n\",\"iuf_version\":\"^0.1.0\",\"name\":\"cpe\",\"version\":\"23.2.2\"}"},{"name":"sle-os-products-15-sp4-x86_64","version":"23.2.7-20230106","original_location":"/etc/cray/upgrade/csm/fullstack-230217/SUSE-Products-15-SP4-x86_64-23.2.7-20230106","validated":true,"manifest":"{\"content\":{\"docker\":[{\"path\":\"docker/cray\"}],\"nexus_blob_stores\":{\"yaml_path\":\"nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"nexus-repositories.yaml\"},\"rpms\":[{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Basesystem-15-SP4-x86_64-Debug-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Basesystem-15-SP4-x86_64-Debug-Pool\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Basesystem-15-SP4-x86_64-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Basesystem-15-SP4-x86_64-Pool\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Containers-15-SP4-x86_64-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Containers-15-SP4-x86_64-Pool\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Desktop-Applications-15-SP4-x86_64-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Desktop-Applications-15-SP4-x86_64-Pool\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Development-Tools-15-SP4-x86_64-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Development-Tools-15-SP4-x86_64-Pool\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-HPC-15-SP4-x86_64-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-HPC-15-SP4-x86_64-Pool\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Legacy-15-SP4-x86_64-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Legacy-15-SP4-x86_64-Pool\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Packagehub-Subpackages-15-SP4-x86_64-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Packagehub-Subpackages-15-SP4-x86_64-Pool\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Public-Cloud-15-SP4-x86_64-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Public-Cloud-15-SP4-x86_64-Pool\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Python3-15-SP4-x86_64-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Python3-15-SP4-x86_64-Pool\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Server-Applications-15-SP4-x86_64-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Server-Applications-15-SP4-x86_64-Pool\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Transactional-Server-15-SP4-x86_64-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Transactional-Server-15-SP4-x86_64-Pool\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Module-Web-Scripting-15-SP4-x86_64-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Module-Web-Scripting-15-SP4-x86_64-Pool\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Product-HPC-15-SP4-x86_64-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Product-HPC-15-SP4-x86_64-Pool\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Product-SLES-15-SP4-x86_64-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Product-SLES-15-SP4-x86_64-Pool\",\"repository_type\":\"raw\"},{\"path\":\"repos/SUSE-23.2.7-20230106-SLE-Product-WE-15-SP4-x86_64-Pool\",\"repository_name\":\"SUSE-23.2.7-20230106-SLE-Product-WE-15-SP4-x86_64-Pool\",\"repository_type\":\"raw\"}]},\"description\":\"SUSE Linux Enterprise (SLE).\\n\",\"iuf_version\":\"^0.6.0\",\"name\":\"sle-os-products-15-sp4-x86_64\",\"version\":\"23.2.7-20230106\"}"},{"name":"slurm","version":"1.2.9-20230209172621.c445c79","original_location":"/etc/cray/upgrade/csm/fullstack-230217/wlm-slurm-1.2.9-20230209172621.c445c79","validated":true,"manifest":"{\"content\":{\"docker\":[{\"path\":\"docker/dtr.dev.cray.com\"}],\"helm\":[{\"path\":\"helm\"}],\"loftsman\":[{\"deploy\":true,\"path\":\"manifests/slurm.yaml\",\"use_manifestgen\":true},{\"deploy\":true,\"path\":\"manifests/slurm-operator.yaml\",\"use_manifestgen\":true}],\"nexus_blob_stores\":{\"yaml_path\":\"nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"nexus-repositories.yaml\"},\"rpms\":[{\"path\":\"rpms/cray-sles15-sp1-cn\",\"repository_name\":\"wlm-slurm-1.2.9-20230209172621.c445c79-sle-15sp1-compute\",\"repository_type\":\"yum\"},{\"path\":\"rpms/cray-sles15-sp2-cn\",\"repository_name\":\"wlm-slurm-1.2.9-20230209172621.c445c79-sle-15sp2-compute\",\"repository_type\":\"yum\"},{\"path\":\"rpms/cray-sles15-sp3-cn\",\"repository_name\":\"wlm-slurm-1.2.9-20230209172621.c445c79-sle-15sp3-compute\",\"repository_type\":\"yum\"},{\"path\":\"rpms/cray-sles15-sp4-cn\",\"repository_name\":\"wlm-slurm-1.2.9-20230209172621.c445c79-sle-15sp4-compute\",\"repository_type\":\"yum\"}],\"vcs\":{\"path\":\"vcs\"}},\"hooks\":{\"deploy_product\":{\"pre\":{\"script_path\":\"hooks/pre-deploy.sh\"}},\"post_install_service_check\":{\"post\":{\"script_path\":\"hooks/validate.sh\"}}},\"iuf_version\":\"^0.6.0\",\"name\":\"slurm\",\"version\":\"1.2.9-20230209172621.c445c79\"}"},{"name":"slingshot-host-software","version":"2.0.2-97-cos-2.5","original_location":"/etc/cray/upgrade/csm/fullstack-230217/slingshot-host-software-2.0.2-97-cos-2.5","validated":true,"manifest":"{\"content\":{\"nexus_blob_stores\":{\"yaml_path\":\"nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"nexus-repositories.yaml\"},\"rpms\":[{\"path\":\"rpms/mellanox/sle15-sp4/cn\",\"repository_name\":\"slingshot-host-software-2.0.2-97-cos-2.5-sle15-sp4-cn-mellanox\",\"repository_type\":\"raw\"},{\"path\":\"rpms/cassini/sle15-sp4/cn\",\"repository_name\":\"slingshot-host-software-2.0.2-97-cos-2.5-sle15-sp4-cn-cassini\",\"repository_type\":\"raw\"}],\"vcs\":{\"path\":\"vcs\"}},\"description\":\"Slingshot Host Software (SHS)\\n\",\"iuf_version\":\"^0.1.0\",\"name\":\"slingshot-host-software\",\"version\":\"2.0.2-97-cos-2.5\"}"},{"name":"cos","version":"2.5.97","original_location":"/etc/cray/upgrade/csm/fullstack-230217/cos-2.5.97","validated":true,"manifest":"{\"content\":{\"docker\":[{\"path\":\"docker\"}],\"helm\":[{\"path\":\"helm\"}],\"ims\":{\"content_dirs\":[\"ims/recipes/x86_64\"]},\"loftsman\":[{\"deploy\":true,\"path\":\"manifests/cos-services.yaml\",\"use_manifestgen\":true}],\"nexus_blob_stores\":{\"yaml_path\":\"nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"nexus-repositories.yaml\"},\"rpms\":[{\"path\":\"rpms/cray-sles15-sp4-ncn\",\"repository_name\":\"cos-2.5.97-sle-15sp4\",\"repository_type\":\"raw\"},{\"path\":\"rpms/cray-sles15-sp4-net-ncn\",\"repository_name\":\"cos-2.5.97-net-sle-15sp4-shs-2.0\",\"repository_type\":\"raw\"},{\"path\":\"rpms/cray-sles15-sp4-cn\",\"repository_name\":\"cos-2.5.97-sle-15sp4-compute\",\"repository_type\":\"raw\"},{\"path\":\"rpms/cray-sles15-sp4-net-cn\",\"repository_name\":\"cos-2.5.97-net-sle-15sp4-compute-shs-2.0\",\"repository_type\":\"raw\"}],\"vcs\":{\"path\":\"vcs\"}},\"description\":\"The Cray Operating System (COS).\\n\",\"hooks\":{\"deliver_product\":{\"post\":{\"script_path\":\"hooks/deliver_product-posthook.sh\"}},\"post_install_service_check\":{\"pre\":{\"script_path\":\"hooks/post_install_service_check-prehook.sh\"}}},\"iuf_version\":\"^0.5.0\",\"name\":\"cos\",\"version\":\"2.5.97\"}"},{"name":"sle-os-backports-sle-15-sp4-x86_64","version":"23.2.7-20230106","original_location":"/etc/cray/upgrade/csm/fullstack-230217/SUSE-Backports-SLE-15-SP4-x86_64-23.2.7-20230106","validated":true,"manifest":"{\"content\":{\"docker\":[{\"path\":\"docker/cray\"}],\"nexus_blob_stores\":{\"yaml_path\":\"nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"nexus-repositories.yaml\"},\"rpms\":[{\"path\":\"repos/SUSE-23.2.7-20230106-Backports-SLE-15-SP4-x86_64\",\"repository_name\":\"SUSE-23.2.7-20230106-Backports-SLE-15-SP4-x86_64\",\"repository_type\":\"raw\"}]},\"description\":\"SUSE Linux Enterprise (SLE).\\n\",\"iuf_version\":\"^0.6.0\",\"name\":\"sle-os-backports-sle-15-sp4-x86_64\",\"version\":\"23.2.7-20230106\"}"},{"name":"uan","version":"2.6.0-alpha.5","original_location":"/etc/cray/upgrade/csm/fullstack-230217/uan-2.6.0-alpha.5","validated":true,"manifest":"{\"content\":{\"docker\":[{\"path\":\"docker\"}],\"helm\":[{\"path\":\"helm\"}],\"ims\":{\"images\":[{\"initrd\":{\"md5sum\":\"366813ef725e16496f70166291d9a344\",\"path\":\"initrd.img-0.2.1.xz\"},\"kernel\":{\"md5sum\":\"baa029f3f2bd6c014ab6084f88a6b547\",\"path\":\"5.3.18-150300.59.87-default-0.2.1.kernel\"},\"path\":\"images/application/cray-application-sles15sp3.x86_64-0.2.1\",\"rootfs\":{\"md5sum\":\"0843e55e84556278b36ba4d68c53882e\",\"path\":\"application-0.2.1.squashfs\"}}]},\"loftsman\":[{\"deploy\":true,\"path\":\"manifests/uan.yaml\",\"use_manifestgen\":true}],\"nexus_blob_stores\":{\"yaml_path\":\"nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"nexus-repositories.yaml\"},\"rpms\":[{\"path\":\"rpms/sle-15sp4\",\"repository_name\":\"uan-2.6.0-sle-15sp4\",\"repository_type\":\"raw\"},{\"path\":\"third-party\",\"repository_name\":\"uan-2.6.0-third-party\",\"repository_type\":\"raw\"}],\"vcs\":{\"path\":\"vcs\"}},\"description\":\"User Access Nodes (UAN).\\n\",\"iuf_version\":\"^0.1.0\",\"name\":\"uan\",\"version\":\"2.6.0-alpha.5\"}"},{"name":"sat","version":"2.5.14","original_location":"/etc/cray/upgrade/csm/fullstack-230217/sat-2.5.14","validated":true,"manifest":"{\"content\":{\"docker\":[{\"path\":\"docker\"}],\"nexus_blob_stores\":{\"yaml_path\":\"nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"nexus-repositories.yaml\"},\"rpms\":[{\"path\":\"rpms/sat-sle-15sp2\",\"repository_name\":\"sat-2.5.14-sle-15sp2\",\"repository_type\":\"raw\"}],\"vcs\":{\"path\":\"ansible\"}},\"description\":\"The System Admin Toolkit (SAT)\\n\",\"iuf_version\":\"^0.1.0\",\"name\":\"sat\",\"version\":\"2.5.14\"}"},{"name":"pbs","version":"1.2.9-20230209172608.28ca7d3","original_location":"/etc/cray/upgrade/csm/fullstack-230217/wlm-pbs-1.2.9-20230209172608.28ca7d3","validated":true,"manifest":"{\"content\":{\"docker\":[{\"path\":\"docker/dtr.dev.cray.com\"}],\"helm\":[{\"path\":\"helm\"}],\"loftsman\":[{\"deploy\":true,\"path\":\"manifests/pbs.yaml\",\"use_manifestgen\":true}],\"nexus_blob_stores\":{\"yaml_path\":\"nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"nexus-repositories.yaml\"},\"rpms\":[{\"path\":\"rpms/cray-sles15-sp1-cn\",\"repository_name\":\"wlm-pbs-1.2.9-20230209172608.28ca7d3-sle-15sp1-compute\",\"repository_type\":\"yum\"},{\"path\":\"rpms/cray-sles15-sp2-cn\",\"repository_name\":\"wlm-pbs-1.2.9-20230209172608.28ca7d3-sle-15sp2-compute\",\"repository_type\":\"yum\"},{\"path\":\"rpms/cray-sles15-sp3-cn\",\"repository_name\":\"wlm-pbs-1.2.9-20230209172608.28ca7d3-sle-15sp3-compute\",\"repository_type\":\"yum\"},{\"path\":\"rpms/cray-sles15-sp4-cn\",\"repository_name\":\"wlm-pbs-1.2.9-20230209172608.28ca7d3-sle-15sp4-compute\",\"repository_type\":\"yum\"}],\"vcs\":{\"path\":\"vcs\"}},\"hooks\":{\"deploy_product\":{\"pre\":{\"script_path\":\"hooks/pre-deploy.sh\"}},\"post_install_service_check\":{\"post\":{\"script_path\":\"hooks/validate.sh\"}}},\"iuf_version\":\"^0.6.0\",\"name\":\"pbs\",\"version\":\"1.2.9-20230209172608.28ca7d3\"}"},{"name":"sma","version":"1.8.6","original_location":"/etc/cray/upgrade/csm/fullstack-230217/sma-1.8.6","validated":true,"manifest":"{\"content\":{\"docker\":[{\"path\":\"docker/dtr.dev.cray.com\"}],\"helm\":[{\"path\":\"helm\"}],\"loftsman\":[{\"deploy\":true,\"path\":\"manifests/sma.yaml\",\"use_manifestgen\":true}],\"nexus_blob_stores\":{\"yaml_path\":\"nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"nexus-repositories.yaml\"},\"rpms\":[{\"path\":\"rpms/sle15_sp4_cn\",\"repository_name\":\"sma-1.8.6-sle-15sp4-compute\",\"repository_type\":\"raw\"},{\"path\":\"rpms/sle15_sp4_ncn\",\"repository_name\":\"sma-1.8.6-sle-15sp4\",\"repository_type\":\"raw\"}],\"vcs\":{\"path\":\"vcs\"}},\"description\":\"The System Monitoring Application (SMA)\\n\",\"hooks\":{\"deliver_product\":{\"pre\":{\"execution_context\":\"master_host\",\"script_path\":\"hooks/pre-deliver-product.sh\"}}},\"iuf_version\":\"^0.6.0\",\"name\":\"sma\",\"version\":\"1.8.6\"}"},{"name":"cray-sdu-rda","version":"2.2.0-beta15","original_location":"/etc/cray/upgrade/csm/fullstack-230217/cray-sdu-rda-2.2.0-beta15","validated":true,"manifest":"{\"content\":{\"docker\":[{\"path\":\"oci_img\"}],\"nexus_blob_stores\":{\"yaml_path\":\"iuf/nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"iuf/nexus-repositories.yaml\"},\"rpms\":[{\"path\":\"rpms/cray-sdu-rda\",\"repository_name\":\"cray-sdu-rda-2.2.0-beta15\",\"repository_type\":\"raw\"}]},\"description\":\"This is the manifest for Cray System Management (CSM)) tarfile that delivers the System Diagnostic Utility (SDU) and Remote Device Access (RDA) software via an Open Container Initiative (OCI) image. Additionally the container interfacing RPM by the name of cray-sdu-rda is also delivered in the tarfile.\",\"hooks\":{\"deliver_product\":{\"pre\":{\"script_path\":\"iuf/hooks/sdu_pre_deliver.sh\"}}},\"iuf_version\":\"^0.0.0\",\"name\":\"cray-sdu-rda\",\"version\":\"2.2.0-beta15\"}"},{"name":"slingshot-host-software","version":"2.0.2-97-csm-1.4","original_location":"/etc/cray/upgrade/csm/fullstack-230217/slingshot-host-software-2.0.2-97-csm-1.4","validated":true,"manifest":"{\"content\":{\"nexus_blob_stores\":{\"yaml_path\":\"nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"nexus-repositories.yaml\"},\"rpms\":[{\"path\":\"rpms/mellanox/sle15-sp4/ncn\",\"repository_name\":\"slingshot-host-software-2.0.2-97-csm-1.4-sle15-sp4-ncn-mellanox\",\"repository_type\":\"raw\"},{\"path\":\"rpms/cassini/sle15-sp4/ncn\",\"repository_name\":\"slingshot-host-software-2.0.2-97-csm-1.4-sle15-sp4-ncn-cassini\",\"repository_type\":\"raw\"}],\"vcs\":{\"path\":\"vcs\"}},\"description\":\"Slingshot Host Software (SHS)\\n\",\"iuf_version\":\"^0.1.0\",\"name\":\"slingshot-host-software\",\"version\":\"2.0.2-97-csm-1.4\"}"},{"name":"csm-diags","version":"1.4.18","original_location":"/etc/cray/upgrade/csm/fullstack-230217/csm-diags-1.4.18","validated":true,"manifest":"{\"content\":{\"docker\":[{\"path\":\"docker/arti.hpc.amslabs.hpecorp.net\"},{\"path\":\"docker/arti.hpc.amslabs.hpecorp.net/csm-docker-remote/stable/\"}],\"helm\":[{\"path\":\"helm\"}],\"loftsman\":[{\"deploy\":true,\"path\":\"manifests/csm-diags.yaml\",\"use_manifestgen\":true}],\"nexus_blob_stores\":{\"yaml_path\":\"nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"nexus-repositories.yaml\"},\"rpms\":[{\"path\":\"rpms/csm-diags\",\"repository_name\":\"csm-diags-1.4.18\",\"repository_type\":\"raw\"}],\"vcs\":{\"path\":\"vcs\"}},\"description\":\"CSM Diagnostics (CSM-DIAGS).\\n\",\"hooks\":{\"post_install_service_check\":{\"post\":{\"script_path\":\"post-install-check.py\"}},\"pre_install_check\":{\"pre\":{\"script_path\":\"pre-install-checks.py\"}}},\"iuf_version\":\"^0.1.0\",\"name\":\"csm-diags\",\"version\":\"1.4.18\"}"},{"name":"HFP-firmware","version":"23.1.2","original_location":"/etc/cray/upgrade/csm/fullstack-230217/HFP-firmware-23.1.2","validated":true,"manifest":"{\"content\":{\"nexus_blob_stores\":{\"yaml_path\":\"nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"nexus-repositories.yaml\"},\"rpms\":[{\"path\":\"rpm\",\"repository_name\":\"HFP-firmware-23.1.2\",\"repository_type\":\"raw\"}]},\"description\":\"HPC Firmware Pack (HFP).\\n\",\"hooks\":{\"deliver_product\":{\"post\":{\"execution_context\":\"master_host\",\"script_path\":\"post-deliver-product.sh\"}},\"pre_install_check\":{\"pre\":{\"execution_context\":\"master_host\",\"script_path\":\"preinstall.sh\"}}},\"iuf_version\":\"^0.4.0\",\"name\":\"HFP-firmware\",\"version\":\"23.1.2\"}"},{"name":"analytics","version":"1.4.20","original_location":"/etc/cray/upgrade/csm/fullstack-230217/analytics-1.4.20","validated":true,"manifest":"{\"content\":{\"nexus_blob_stores\":{\"yaml_path\":\"nexus-blobstores.yaml\"},\"nexus_repositories\":{\"yaml_path\":\"nexus-repositories.yaml\"},\"s3\":[{\"bucket\":\"boot-images\",\"key\":\"Analytics/Cray-Analytics.x86_64-1.4.20.squashfs\",\"path\":\"squashfs/Cray-Analytics.x86_64-1.4.20.squashfs\"}],\"vcs\":{\"path\":\"vcs\"}},\"hooks\":{\"deliver_product\":{\"post\":{\"script_path\":\"hooks/deliver_product-posthook.sh\"}}},\"iuf_version\":\"^0.6.0\",\"name\":\"analytics\",\"version\":\"1.4.20\"}"}],"activity_state":"wait_for_admin"}`
	configmap := v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: DEFAULT_NAMESPACE,
			Labels: map[string]string{
				"type": LABEL_ACTIVITY,
			},
		},
		Data: map[string]string{LABEL_ACTIVITY: activityData},
	}
	fakeClient := fake.NewSimpleClientset(&configmap)

	mockTokenValue := "{{workflow.parameters.auth_token}}"
	keycloakServiceMock := mocks.NewMockKeycloakService(ctrl)
	keycloakServiceMock.EXPECT().NewKeycloakAccessToken().Return(mockTokenValue, nil).AnyTimes()

	iufSvc := NewIufService(utils.GetLogger(), services_shared.NewArgoService(env), services_shared.K8sService{
		Client: fakeClient,
	}, keycloakServiceMock, env)
	return name, mockTokenValue, iufSvc
}
