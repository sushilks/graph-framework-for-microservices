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

package kdefault

import (
	"path/filepath"

	"gitlab.eng.vmware.com/nsx-allspark_users/nexus-sdk/kubebuilder.git/pkg/machinery"
)

var _ machinery.Template = &ManagerConfigPatch{}

// ManagerConfigPatch scaffolds a ManagerConfigPatch for a Resource
type ManagerConfigPatch struct {
	machinery.TemplateMixin
}

// SetTemplateDefaults implements input.Template
func (f *ManagerConfigPatch) SetTemplateDefaults() error {
	if f.Path == "" {
		f.Path = filepath.Join("config", "default", "manager_config_patch.yaml")
	}

	f.TemplateBody = managerConfigPatchTemplate

	return nil
}

const managerConfigPatchTemplate = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: controller-manager
  namespace: system
spec:
  template:
    spec:
      containers:
      - name: manager
        args:
        - "--config=controller_manager_config.yaml"
        volumeMounts:
        - name: manager-config
          mountPath: /controller_manager_config.yaml
          subPath: controller_manager_config.yaml
      volumes:
      - name: manager-config
        configMap:
          name: manager-config
`
