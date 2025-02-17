// Copyright 2022 The ILLA Authors.
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

package repository

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type KVState struct {
	ID        int       `json:"id" 		   gorm:"column:id;type:bigserial"`
	StateType int       `json:"state_type" gorm:"column:state_type;type:bigint"`
	AppRefID  int       `json:"app_ref_id" gorm:"column:app_ref_id;type:bigint"`
	Version   int       `json:"version"    gorm:"column:version;type:bigint"`
	Key       string    `json:"key" 	   gorm:"column:key;type:text"`
	Value     string    `json:"value" 	   gorm:"column:value;type:jsonb"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;type:timestamp"`
	CreatedBy int       `json:"created_by" gorm:"column:created_by;type:bigint"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at;type:timestamp"`
	UpdatedBy int       `json:"updated_by" gorm:"column:updated_by;type:bigint"`
}

type KVStateRepository interface {
	Create(kvstate *KVState) error
	Delete(kvstateID int) error
	Update(kvstate *KVState) error
	RetrieveByID(kvstateID int) (*KVState, error)
	RetrieveKVStatesByVersion(versionID int) ([]*KVState, error)
	RetrieveKVStatesByKey(key string) ([]*KVState, error)
	RetrieveKVStatesByApp(apprefid int, statetype int, version int) ([]*KVState, error)
	RetrieveEditVersionByAppAndKey(apprefid int, statetype int, key string) (*KVState, error)
	RetrieveAllTypeKVStatesByApp(apprefid int, version int) ([]*KVState, error)
	DeleteAllTypeKVStatesByApp(apprefid int) error
	DeleteAllKVStatesByAppVersionAndType(apprefid int, version int, stateType int) error
}

type KVStateRepositoryImpl struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

func NewKVStateRepositoryImpl(logger *zap.SugaredLogger, db *gorm.DB) *KVStateRepositoryImpl {
	return &KVStateRepositoryImpl{
		logger: logger,
		db:     db,
	}
}

func (impl *KVStateRepositoryImpl) Create(kvstate *KVState) error {
	if err := impl.db.Create(kvstate).Error; err != nil {
		return err
	}
	return nil
}

func (impl *KVStateRepositoryImpl) Delete(kvstateID int) error {
	if err := impl.db.Delete(&KVState{}, kvstateID).Error; err != nil {
		return err
	}
	return nil
}

func (impl *KVStateRepositoryImpl) Update(kvstate *KVState) error {
	if err := impl.db.Model(kvstate).UpdateColumns(KVState{
		ID:        kvstate.ID,
		StateType: kvstate.StateType,
		AppRefID:  kvstate.AppRefID,
		Version:   kvstate.Version,
		Key:       kvstate.Key,
		Value:     kvstate.Value,
		UpdatedAt: kvstate.UpdatedAt,
		UpdatedBy: kvstate.UpdatedBy,
	}).Error; err != nil {
		return err
	}
	return nil
}

func (impl *KVStateRepositoryImpl) RetrieveByID(kvstateID int) (*KVState, error) {
	kvstate := &KVState{}
	if err := impl.db.First(kvstate, kvstateID).Error; err != nil {
		return &KVState{}, err
	}
	return kvstate, nil
}

func (impl *KVStateRepositoryImpl) RetrieveKVStatesByVersion(version int) ([]*KVState, error) {
	var kvstates []*KVState
	if err := impl.db.Where("version = ?", version).Find(&kvstates).Error; err != nil {
		return nil, err
	}
	return kvstates, nil
}

func (impl *KVStateRepositoryImpl) RetrieveKVStatesByKey(key string) ([]*KVState, error) {
	var kvstates []*KVState
	if err := impl.db.Where("key = ?", key).Find(&kvstates).Error; err != nil {
		return nil, err
	}
	return kvstates, nil
}

func (impl *KVStateRepositoryImpl) RetrieveKVStatesByApp(apprefid int, statetype int, version int) ([]*KVState, error) {
	var kvstates []*KVState
	if err := impl.db.Where("app_ref_id = ? AND state_type = ? AND version = ?", apprefid, statetype, version).Find(&kvstates).Error; err != nil {
		return nil, err
	}
	return kvstates, nil
}

func (impl *KVStateRepositoryImpl) RetrieveEditVersionByAppAndKey(apprefid int, statetype int, key string) (*KVState, error) {
	var kvstate *KVState
	if err := impl.db.Where("app_ref_id = ? AND state_type = ? AND version = ? AND key = ?", apprefid, statetype, APP_EDIT_VERSION, key).First(&kvstate).Error; err != nil {
		return nil, err
	}
	return kvstate, nil
}

func (impl *KVStateRepositoryImpl) RetrieveAllTypeKVStatesByApp(apprefid int, version int) ([]*KVState, error) {
	var kvstates []*KVState
	if err := impl.db.Where("app_ref_id = ? AND version = ?", apprefid, version).Find(&kvstates).Error; err != nil {
		return nil, err
	}
	return kvstates, nil
}

func (impl *KVStateRepositoryImpl) DeleteAllTypeKVStatesByApp(apprefid int) error {
	if err := impl.db.Where("app_ref_id = ?", apprefid).Delete(&KVState{}).Error; err != nil {
		return err
	}
	return nil
}

func (impl *KVStateRepositoryImpl) DeleteAllKVStatesByAppVersionAndType(apprefid int, version int, stateType int) error {
	if err := impl.db.Where("app_ref_id = ? AND version = ? AND state_type = ?", apprefid, version, stateType).Delete(&KVState{}).Error; err != nil {
		return err
	}
	return nil
}
