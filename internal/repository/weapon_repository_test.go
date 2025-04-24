package repository

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Riven-of-a-Thousand-Servers/rivenbot-commons/pkg/types"
	"github.com/stretchr/testify/assert"
	"pgcr-processing-service/internal/model"
)

func TestAddWeapon_Success(t *testing.T) {
	// given: a weapon types
	weapon := model.WeaponEntity{
		WeaponHash:          3211806999,
		WeaponIcon:          "/some/route/here/",
		WeaponName:          "Izanagi's Burden",
		WeaponDamageType:    types.KINETIC,
		WeaponEquipmentSlot: types.PRIMARY,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error stablishing stub connection to database")
	}

	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	weaponRepository := WeaponRepositoryImpl{
		Conn: db,
	}

	mock.ExpectExec("INSERT INTO weapon").
		WithArgs(weapon.WeaponHash, weapon.WeaponIcon, weapon.WeaponName, weapon.WeaponDamageType, weapon.WeaponEquipmentSlot).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// when: AddWeapon is called
	result, err := weaponRepository.AddWeapon(tx, weapon)

	if err != nil {
		t.Fatalf("Error was not expected, but got one: %v", err)
	}

	// then: the result is correct
	assert := assert.New(t)
	assert.Equal(weapon.WeaponHash, result.WeaponHash, "Weapon hashes should match")
	assert.Equal(weapon.WeaponIcon, result.WeaponIcon, "Weapon icon route should match")
	assert.Equal(weapon.WeaponName, result.WeaponName, "Weapon name should match")
	assert.Equal(weapon.WeaponDamageType, result.WeaponDamageType, "Weapon damage type should match")
	assert.Equal(weapon.WeaponEquipmentSlot, result.WeaponEquipmentSlot, "Weapon equipment slot should match")

	// and: All database expecations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Mock expectations were not met: %v", err)
	}
}

func TestAddWeapon_ErrorOnWeaponInsert(t *testing.T) {
	// given: a weapon types
	weapon := model.WeaponEntity{
		WeaponHash:          3211806999,
		WeaponIcon:          "/some/route/here/",
		WeaponName:          "Izanagi's Burden",
		WeaponDamageType:    types.KINETIC,
		WeaponEquipmentSlot: types.PRIMARY,
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Error stablishing stub connection to database")
	}

	defer db.Close()

	mock.ExpectBegin()
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Error beginning transaction: %v", err)
	}

	weaponRepository := WeaponRepositoryImpl{
		Conn: db,
	}

	mock.ExpectExec("INSERT INTO weapon").
		WithArgs(weapon.WeaponHash, weapon.WeaponIcon, weapon.WeaponName, weapon.WeaponDamageType, weapon.WeaponEquipmentSlot).
		WillReturnError(fmt.Errorf("Something happened while inserting weapon in DB"))

	// when: AddWeapon is called
	_, err = weaponRepository.AddWeapon(tx, weapon)

	// then: An error is expected to happen
	if err == nil {
		t.Fatalf("Expected error, found none")
	}

	// and: All database expecations are met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Mock expectations were not met: %v", err)
	}
}
