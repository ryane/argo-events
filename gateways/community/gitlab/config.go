/*
Copyright 2018 BlackRock, Inc.

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

package gitlab

import (
	"github.com/argoproj/argo-events/gateways/common"
	gwcommon "github.com/argoproj/argo-events/gateways/common"
	"github.com/ghodss/yaml"
	"github.com/rs/zerolog"
	"github.com/xanzy/go-gitlab"
	"k8s.io/client-go/kubernetes"
)

// GitlabEventSourceExecutor implements ConfigExecutor
type GitlabEventSourceExecutor struct {
	Log zerolog.Logger
	// Clientset is kubernetes client
	Clientset kubernetes.Interface
	// Namespace where gateway is deployed
	Namespace string
}

// RouteConfig contains the configuration information for a route
type RouteConfig struct {
	route     *gwcommon.Route
	clientset kubernetes.Interface
	client    *gitlab.Client
	hook      *gitlab.ProjectHook
	namespace string
	ges       *gitlabEventSource
}

// gitlabEventSource contains information to setup a gitlab project integration
type gitlabEventSource struct {
	// Webhook Id
	Id int `json:"id"`
	// Webhook
	Hook *common.Webhook `json:"hook"`
	// ProjectId is the id of project for which integration needs to setup
	ProjectId string `json:"projectId"`
	// Event is a gitlab event to listen to.
	// Refer https://github.com/xanzy/go-gitlab/blob/bf34eca5d13a9f4c3f501d8a97b8ac226d55e4d9/projects.go#L794.
	Event string `json:"event"`
	// AccessToken is reference to k8 secret which holds the gitlab api access information
	AccessToken *GitlabSecret `json:"accessToken"`
	// EnableSSLVerification to enable ssl verification
	EnableSSLVerification bool `json:"enableSSLVerification"`
	// GitlabBaseURL is the base URL for API requests to a custom endpoint
	GitlabBaseURL string `json:"gitlabBaseUrl"`
}

// GitlabSecret contains information of k8 secret which holds the gitlab api access information
type GitlabSecret struct {
	// Key within the K8 secret for access token
	Key string
	// Name of K8 secret containing access token info
	Name string
}

// cred stores the api access token
type cred struct {
	// token is gitlab api access token
	token string
}

// parseEventSource parses an event sources of gateway
func parseEventSource(config string) (interface{}, error) {
	var g *gitlabEventSource
	err := yaml.Unmarshal([]byte(config), &g)
	if err != nil {
		return nil, err
	}
	return g, err
}
