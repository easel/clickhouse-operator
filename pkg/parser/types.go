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

package parser

import (
	chiv1 "github.com/altinity/clickhouse-operator/pkg/apis/clickhouse.altinity.com/v1"
	apps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
)

// ObjectKind defines k8s objects list kind
type ObjectKind uint8

// ObjectsMap defines map of a generated k8s objects
type ObjectsMap map[ObjectKind]interface{}

// ConfigMapList defines a list of the ConfigMap objects
type ConfigMapList []*corev1.ConfigMap

// StatefulSetList defines a list of the StatefulSet objects
type StatefulSetList []*apps.StatefulSet

// ServiceList defines a list of the Service objects
type ServiceList []*corev1.Service

type generatorOptions struct {

	// deploymentNumber contains number of each deployment required to satisfy clusters' infrastructure
	// deploymentNumber[fingerprint] = int number of such deployments required
	deploymentNumber NamedNumber

	// fullDeploymentIDToDeployment[fullDeploymentID] = &replica.Deployment
	fullDeploymentIDToDeployment map[string]*chiv1.ChiDeployment

	// fullDeploymentIDToMacrosData[fullDeploymentID] = macrosDataShardDescriptionList
	fullDeploymentIDToMacrosData map[string]macrosDataShardDescriptionList

	// commonConfigSections maps section name to section XML config
	commonConfigSections map[string]string

	// commonUsersConfigSections maps section name to section XML config
	commonUsersConfigSections map[string]string
}

// macrosDataShardDescriptionList is a list of all shards descriptions
type macrosDataShardDescriptionList []*macrosDataShardDescription

// macrosDataShardDescription is used in generating macros config file
// and describes which cluster this shard belongs to and
// index (1-based, human-friendly and XML-config usable) of this shard within cluster
// Used to build this:
//<yandex>
//    <macros>
//        <installation>example-02</installation>
//        <2shard-1repl>example-02-2shard-1repl</02-2shard-1repl>
//        <2shard-1repl-shard>1 [OWN UNIQUE SHARD ID]</02-2shard-1repl-shard>
//        <replica>1eb454-1 [OWN UNIQUE ID]</replica>
//    </macros>
//</yandex>
type macrosDataShardDescription struct {
	clusterName string
	index       int
}

// NamedNumber maps Deployment fingerprint to its usage count
type NamedNumber map[string]int

// mergeAndReplaceWithBiggerValues combines NamedNumber object with another one
// and replaces local values with another one's values in case another's value is bigger
func (d NamedNumber) mergeAndReplaceWithBiggerValues(another NamedNumber) {

	// Loop over another struct and bring in new OR bigger values
	for key, value := range another {
		_, ok := d[key]

		if !ok {
			// No such key - new key/value pair just - include/merge it in
			d[key] = value
		} else if value > d[key] {
			// Have such key, but "another"'s value is bigger, overwrite local value
			d[key] = value
		}
	}
}

// volumeClaimTemplatesIndex maps volume claim template name - which is .spec.templates.volumeClaimTemplates.name
// to "Volume Claim Template"-like structure (simplified).
// Used to provide dictionary/index for templates
type volumeClaimTemplatesIndex map[string]*volumeClaimTemplatesIndexData

type volumeClaimTemplatesIndexData struct {
	useDefaultName        bool
	persistentVolumeClaim *corev1.PersistentVolumeClaim
}

// podTemplatesIndex maps pod template name - which is .spec.templates.podTemplates.name
// to "Pod Template"-like structure (simplified).
// Used to provide dictionary/index for templates
type podTemplatesIndex map[string]*podTemplatesIndexData

type podTemplatesIndexData struct {
	containers []corev1.Container
	volumes    []corev1.Volume
}
