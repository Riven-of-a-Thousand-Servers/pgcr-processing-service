package service

import (
	"context"
	"database/sql"
	"pgcr-processing-service/internal/model"

	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
	"github.com/stretchr/testify/mock"
)

type MockPlayerRepository struct {
	mock.Mock
}

type MockRaidRepository struct {
	mock.Mock
}

type MockInstanceActivityRepository struct {
	mock.Mock
}

type MockWeaponRepository struct {
	mock.Mock
}

type MockRedisService struct {
	mock.Mock
}

type MockActivityWeaponStatsRepository struct {
	mock.Mock
}

func (m *MockActivityWeaponStatsRepository) AddInstanceWeaponStats(tx *sql.Tx, entity model.InstanceWeaponStats) (*model.InstanceWeaponStats, error) {
	args := m.Called(tx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*model.InstanceWeaponStats), args.Error(1)
	}
}

func (m *MockRedisService) GetManifestEntity(ctx context.Context, hash string) (*types.ManifestObject, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*types.ManifestObject), args.Error(1)
}

func (m *MockPlayerRepository) AddPlayer(tx *sql.Tx, entity model.PlayerEntity) (*model.PlayerEntity, error) {
	args := m.Called(tx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*model.PlayerEntity), args.Error(1)
	}
}

func (m *MockRaidRepository) AddRaidInfo(tx *sql.Tx, entity model.RaidEntity) (*model.RaidEntity, error) {
	args := m.Called(tx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*model.RaidEntity), args.Error(1)
	}
}

func (m *MockInstanceActivityRepository) AddInstanceActivity(tx *sql.Tx, entity model.InstanceActivityEntity) (*model.InstanceActivityEntity, error) {
	args := m.Called(tx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.InstanceActivityEntity), args.Error(1)
}

func (m *MockWeaponRepository) AddWeapon(tx *sql.Tx, entity model.WeaponEntity) (*model.WeaponEntity, error) {
	args := m.Called(tx, entity)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	} else {
		return args.Get(0).(*model.WeaponEntity), args.Error(1)
	}
}
