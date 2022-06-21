package main_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func TestSingleLabelRebuild(t *testing.T) {

	envMap, mapErr := getEnvMap()
	if mapErr != nil {
		t.Fatalf("%v", mapErr)
	}
	hosts := []string{"ncn-w001"}
	var rebuildResponse RebuildResponse

	err := rebuildHosts(envMap["REBUILD_URL"], hosts, &rebuildResponse)

	if err != nil {
		t.Fatalf("could not rebuild hosts: %v", err.Error())
	}

	// wait for workflow to start
	// DI: let me know if theres a smarter way to handle this
	time.Sleep(5 * time.Second)

	var getResponse GetResponse
	getRebuildStatus(envMap["STATUS_URL"], "target-ncns="+hosts[0], &getResponse)

	for getResponse[0].Status.Phase == "Running" {
		// make get request to check status
		// TODO: handle error that this returns
		getRebuildStatus(envMap["STATUS_URL"], "target-ncns="+hosts[0], &getResponse)

		// DI: Let me know if you would like me to sleep here or just request as many times as possible?
		// time.Sleep(2 * time.Second)

	}

	// TODO: Fail here in more cases
	if getResponse[0].Status.Phase != "Succeeded" {
		t.Fatalf("Expected phase to be Succeeded but got: %v", getResponse[0].Status.Phase)

	}

}

func rebuildHosts(url string, hosts []string, target interface{}) error {

	hoststostring := ``

	for i := 0; i < len(hosts); i++ {
		if i == len(hosts)-1 {
			hoststostring += `"` + hosts[i] + `"`
		} else {
			hoststostring += `"` + hosts[i] + `",`
		}
	}

	requestBody := strings.NewReader(`
		{
		"dryRun": true,
		"hosts": [
		  ` + hoststostring + `
		]
	  	}`)

	response, err := http.Post(url, "application/json", requestBody)
	defer response.Body.Close()
	if err != nil {
		return errors.New("could not complete POST request: " + err.Error())
	}

	// content, _ := ioutil.ReadAll(response.Body)

	return json.NewDecoder(response.Body).Decode(target)

}

func getEnvMap() (map[string]string, error) {

	err := godotenv.Load(".env")
	if err != nil {
		return nil, errors.New("could not load .env file")
	}
	envMap, mapErr := godotenv.Read(".env")
	if mapErr != nil {
		return nil, errors.New("could not read .env file")
	}
	return envMap, nil
}

func getRebuildStatus(url string, label string, target interface{}) error {

	url += "?labelSelector=" + label

	response, err := http.Get(url)
	defer response.Body.Close()
	if err != nil {
		return err
	}

	return json.NewDecoder(response.Body).Decode(target)
}

type RebuildResponse struct {
	Name       string   `json:"name"`
	TargetNcns []string `json:"targetNcns"`
}

type GetResponse []struct {
	Name  string `json:"name"`
	Label struct {
		NodeType                   string `json:"node-type"`
		TargetNcns                 string `json:"target-ncns"`
		Type                       string `json:"type"`
		WorkflowsArgoprojIoCreator string `json:"workflows.argoproj.io/creator"`
		WorkflowsArgoprojIoPhase   string `json:"workflows.argoproj.io/phase"`
	} `json:"label"`
	Status struct {
		Phase      string      `json:"phase"`
		StartedAt  time.Time   `json:"startedAt"`
		FinishedAt interface{} `json:"finishedAt"`
		Progress   string      `json:"progress"`
		Nodes      struct {
			NcnLifecycleRebuildX9K8F struct {
				ID            string      `json:"id"`
				Name          string      `json:"name"`
				DisplayName   string      `json:"displayName"`
				Type          string      `json:"type"`
				TemplateName  string      `json:"templateName"`
				TemplateScope string      `json:"templateScope"`
				Phase         string      `json:"phase"`
				StartedAt     time.Time   `json:"startedAt"`
				FinishedAt    interface{} `json:"finishedAt"`
				Progress      string      `json:"progress"`
				Children      []string    `json:"children"`
			} `json:"ncn-lifecycle-rebuild-x9k8f"`
			NcnLifecycleRebuildX9K8F104190284 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children []string `json:"children"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-104190284"`
			NcnLifecycleRebuildX9K8F1092205481 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children     []string `json:"children"`
				HostNodeName string   `json:"hostNodeName"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-1092205481"`
			NcnLifecycleRebuildX9K8F131465828 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children     []string `json:"children"`
				HostNodeName string   `json:"hostNodeName"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-131465828"`
			NcnLifecycleRebuildX9K8F1857508145 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children []string `json:"children"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-1857508145"`
			NcnLifecycleRebuildX9K8F1958437438 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children     []string `json:"children"`
				HostNodeName string   `json:"hostNodeName"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-1958437438"`
			NcnLifecycleRebuildX9K8F2020474418 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children []string `json:"children"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-2020474418"`
			NcnLifecycleRebuildX9K8F2033783197 struct {
				ID                string    `json:"id"`
				Name              string    `json:"name"`
				DisplayName       string    `json:"displayName"`
				Type              string    `json:"type"`
				TemplateName      string    `json:"templateName"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Children      []string `json:"children"`
				OutboundNodes []string `json:"outboundNodes"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-2033783197"`
			NcnLifecycleRebuildX9K8F2115364659 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children []string `json:"children"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-2115364659"`
			NcnLifecycleRebuildX9K8F2480791111 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children     []string `json:"children"`
				HostNodeName string   `json:"hostNodeName"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-2480791111"`
			NcnLifecycleRebuildX9K8F2535232093 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children []string `json:"children"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-2535232093"`
			NcnLifecycleRebuildX9K8F2600770821 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children []string `json:"children"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-2600770821"`
			NcnLifecycleRebuildX9K8F2624410248 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children     []string `json:"children"`
				HostNodeName string   `json:"hostNodeName"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-2624410248"`
			NcnLifecycleRebuildX9K8F3086706440 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					Result   time.Time `json:"result"`
					ExitCode string    `json:"exitCode"`
				} `json:"outputs"`
				Children     []string `json:"children"`
				HostNodeName string   `json:"hostNodeName"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-3086706440"`
			NcnLifecycleRebuildX9K8F3145789770 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope string      `json:"templateScope"`
				Phase         string      `json:"phase"`
				BoundaryID    string      `json:"boundaryID"`
				StartedAt     time.Time   `json:"startedAt"`
				FinishedAt    interface{} `json:"finishedAt"`
				Progress      string      `json:"progress"`
				Inputs        struct {
					Parameters []struct {
						Name  string    `json:"name"`
						Value time.Time `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Children []string `json:"children"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-3145789770"`
			NcnLifecycleRebuildX9K8F3344334626 struct {
				ID            string      `json:"id"`
				Name          string      `json:"name"`
				DisplayName   string      `json:"displayName"`
				Type          string      `json:"type"`
				TemplateName  string      `json:"templateName"`
				TemplateScope string      `json:"templateScope"`
				Phase         string      `json:"phase"`
				BoundaryID    string      `json:"boundaryID"`
				StartedAt     time.Time   `json:"startedAt"`
				FinishedAt    interface{} `json:"finishedAt"`
				Progress      string      `json:"progress"`
				Inputs        struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Children []string `json:"children"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-3344334626"`
			NcnLifecycleRebuildX9K8F3495968945 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					Result   time.Time `json:"result"`
					ExitCode string    `json:"exitCode"`
				} `json:"outputs"`
				Children []string `json:"children"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-3495968945"`
			NcnLifecycleRebuildX9K8F3521380144 struct {
				ID                string    `json:"id"`
				Name              string    `json:"name"`
				DisplayName       string    `json:"displayName"`
				Type              string    `json:"type"`
				TemplateName      string    `json:"templateName"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children     []string `json:"children"`
				HostNodeName string   `json:"hostNodeName"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-3521380144"`
			NcnLifecycleRebuildX9K8F3617667273 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children []string `json:"children"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-3617667273"`
			NcnLifecycleRebuildX9K8F3781305857 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children []string `json:"children"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-3781305857"`
			NcnLifecycleRebuildX9K8F3787273452 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children     []string `json:"children"`
				HostNodeName string   `json:"hostNodeName"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-3787273452"`
			NcnLifecycleRebuildX9K8F3871122776 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children     []string `json:"children"`
				HostNodeName string   `json:"hostNodeName"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-3871122776"`
			NcnLifecycleRebuildX9K8F3983236595 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children     []string `json:"children"`
				HostNodeName string   `json:"hostNodeName"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-3983236595"`
			NcnLifecycleRebuildX9K8F4212224497 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope string      `json:"templateScope"`
				Phase         string      `json:"phase"`
				BoundaryID    string      `json:"boundaryID"`
				StartedAt     time.Time   `json:"startedAt"`
				FinishedAt    interface{} `json:"finishedAt"`
				Progress      string      `json:"progress"`
				Inputs        struct {
					Parameters []struct {
						Name  string    `json:"name"`
						Value time.Time `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-4212224497"`
			NcnLifecycleRebuildX9K8F523633377 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children     []string `json:"children"`
				HostNodeName string   `json:"hostNodeName"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-523633377"`
			NcnLifecycleRebuildX9K8F750204496 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children []string `json:"children"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-750204496"`
			NcnLifecycleRebuildX9K8F900712624 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children     []string `json:"children"`
				HostNodeName string   `json:"hostNodeName"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-900712624"`
			NcnLifecycleRebuildX9K8F91206074 struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				Type        string `json:"type"`
				TemplateRef struct {
					Name     string `json:"name"`
					Template string `json:"template"`
				} `json:"templateRef"`
				TemplateScope     string    `json:"templateScope"`
				Phase             string    `json:"phase"`
				BoundaryID        string    `json:"boundaryID"`
				StartedAt         time.Time `json:"startedAt"`
				FinishedAt        time.Time `json:"finishedAt"`
				Progress          string    `json:"progress"`
				ResourcesDuration struct {
					CPU    int `json:"cpu"`
					Memory int `json:"memory"`
				} `json:"resourcesDuration"`
				Inputs struct {
					Parameters []struct {
						Name  string `json:"name"`
						Value string `json:"value"`
					} `json:"parameters"`
				} `json:"inputs"`
				Outputs struct {
					Artifacts []struct {
						Name string `json:"name"`
						S3   struct {
							Key string `json:"key"`
						} `json:"s3"`
					} `json:"artifacts"`
					ExitCode string `json:"exitCode"`
				} `json:"outputs"`
				Children []string `json:"children"`
			} `json:"ncn-lifecycle-rebuild-x9k8f-91206074"`
		} `json:"nodes"`
		Conditions []struct {
			Type   string `json:"type"`
			Status string `json:"status"`
		} `json:"conditions"`
		ResourcesDuration struct {
			CPU    int `json:"cpu"`
			Memory int `json:"memory"`
		} `json:"resourcesDuration"`
	} `json:"status"`
}
