// Copyright (c) 2021 Terminus, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package addon

import (
	"reflect"
	"testing"

	"bou.ke/monkey"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"github.com/erda-project/erda/apistructs"
	"github.com/erda-project/erda/bundle"
	"github.com/erda-project/erda/internal/tools/orchestrator/dbclient"
	"github.com/erda-project/erda/pkg/database/dbengine"
	"github.com/erda-project/erda/pkg/kms/kmstypes"
	"github.com/erda-project/erda/pkg/parser/diceyml"
)

func TestSetlabelsFromOptions(t *testing.T) {
	type args struct {
		labels map[string]string
		opts   map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "test org label",
			args: args{
				labels: map[string]string{},
				opts: map[string]string{
					apistructs.LabelOrgName: "erda",
				},
			},
			want: map[string]string{
				apistructs.EnvDiceOrgName: "erda",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetlabelsFromOptions(tt.args.opts, tt.args.labels)
			if !reflect.DeepEqual(tt.args.labels, tt.want) {
				t.Errorf("SetlabelsFromOptions() = %v, want %v", tt.args.labels, tt.want)
			}
		})
	}
}

func TestBuildRocketMQOperaotrServiceItem(t *testing.T) {
	params := &apistructs.AddonHandlerCreateItem{
		Options: map[string]string{},
		Plan:    "basic",
	}
	addonIns := &dbclient.AddonInstance{
		ID: "1",
	}
	addonSpec := &apistructs.AddonExtension{
		Name: "rocketmq",
		Plan: map[string]apistructs.AddonPlanItem{
			"basic": {
				InsideMoudle: map[string]apistructs.AddonPlanItem{
					"rocketmq-namesrv": {
						CPU:   1,
						Mem:   2048,
						Nodes: 1,
					},
					"rocketmq-broker": {
						CPU:   1,
						Mem:   2048,
						Nodes: 2,
					},
					"rocketmq-console": {
						CPU:   0.5,
						Mem:   1024,
						Nodes: 1,
					},
				},
			},
		},
	}
	addonDice := &diceyml.Object{
		Services: diceyml.Services{
			"rocketmq-namesrv": {
				Labels: map[string]string{},
				Envs: map[string]string{
					"ADDON_TYPE": "rocketmq",
				},
			},
			"rocketmq-broker": {
				Labels: map[string]string{},
				Envs: map[string]string{
					"ADDON_TYPE": "rocketmq",
				},
			},
			"rocketmq-console": {
				Labels: map[string]string{},
				Envs: map[string]string{
					"ADDON_TYPE": "rocketmq",
				},
			},
		},
	}
	a := &Addon{}

	err := a.BuildRocketMQOperaotrServiceItem(params, addonIns, addonSpec, addonDice, nil, "5.0.0")
	assert.NoError(t, err)
	assert.Equal(t, 2, addonDice.Services["rocketmq-broker"].Deployments.Replicas)
}

func TestBuildRedisServiceItem(t *testing.T) {
	params := &apistructs.AddonHandlerCreateItem{
		Options: map[string]string{},
		Plan:    "basic",
	}
	addonIns := &dbclient.AddonInstance{
		ID: "1",
	}
	addonSpec := &apistructs.AddonExtension{
		Name: "redis",
		Plan: map[string]apistructs.AddonPlanItem{
			"basic": {
				CPU:   0.5,
				Mem:   256,
				Nodes: 1,
			},
		},
	}
	addonDice := &diceyml.Object{
		Services: diceyml.Services{
			"redis-master": {
				Labels: map[string]string{},
				Envs: map[string]string{
					"ADDON_TYPE":     "redis",
					"ADDON_GROUP_ID": "redis",
				},
			},
			"redis-slave": {
				Labels: map[string]string{},
				Envs: map[string]string{
					"ADDON_TYPE":     "redis",
					"ADDON_GROUP_ID": "redis",
				},
			},
			"redis-sentinel": {
				Labels: map[string]string{},
				Envs: map[string]string{
					"ADDON_TYPE":     "redis",
					"ADDON_GROUP_ID": "redis",
				},
			},
		},
	}
	bdl := bundle.New()

	pm1 := monkey.PatchInstanceMethod(reflect.TypeOf(bdl), "KMSEncrypt", func(_ *bundle.Bundle, req apistructs.KMSEncryptRequest) (*kmstypes.EncryptResponse, error) {
		return &kmstypes.EncryptResponse{KeyID: "1", CiphertextBase64: "xxx"}, nil
	})
	defer pm1.Unpatch()

	db := &dbclient.DBClient{
		DBEngine: &dbengine.DBEngine{
			DB: &gorm.DB{},
		},
	}
	pm2 := monkey.PatchInstanceMethod(reflect.TypeOf(db), "CreateAddonInstanceExtra", func(_ *dbclient.DBClient, addonInstanceExtra *dbclient.AddonInstanceExtra) error {
		return nil
	})
	defer pm2.Unpatch()

	pm3 := monkey.PatchInstanceMethod(reflect.TypeOf(&gorm.DB{}), "Create", func(_ *gorm.DB, value interface{}) *gorm.DB {
		return &gorm.DB{}
	})
	defer pm3.Unpatch()

	a := &Addon{db: db, bdl: bdl}
	err := a.BuildRedisServiceItem(params, addonIns, addonSpec, addonDice)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(addonDice.Services[apistructs.RedisMasterNamePrefix].SideCars))
}

func TestBuildEsServiceItem(t *testing.T) {
	params := &apistructs.AddonHandlerCreateItem{
		Options: map[string]string{},
		Plan:    "basic",
	}
	addonIns := &dbclient.AddonInstance{
		ID: "1",
	}
	addonSpec := &apistructs.AddonExtension{
		Name: "elasticsearch",
		Plan: map[string]apistructs.AddonPlanItem{
			"basic": {
				CPU:   0.5,
				Mem:   256,
				Nodes: 1,
			},
		},
	}
	addonDice := &diceyml.Object{
		Services: diceyml.Services{
			"elasticsearch": {
				Labels: map[string]string{},
				Envs: map[string]string{
					"ADDON_TYPE": "elasticsearch",
				},
			},
		},
	}

	a := &Addon{}
	err := a.BuildEsServiceItem(params, addonIns, addonSpec, addonDice, &apistructs.ClusterInfoData{})
	assert.NoError(t, err)
	assert.Equal(t, &[]int64{1000}[0], addonDice.Services["elasticsearch-1"].K8SSnippet.Container.SecurityContext.RunAsUser)
}
