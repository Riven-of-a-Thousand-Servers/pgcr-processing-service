package repository

import (
	"database/sql"
	"fmt"

	"pgcr-processing-service/internal/model"
)

type WeaponRepository struct {
	Conn *sql.DB
}

func (r *WeaponRepository) AddWeapon(tx *sql.Tx, entity model.WeaponEntity) (*model.WeaponEntity, error) {
	_, err := tx.Exec(`
    INSERT INTO weapon (weapon_hash, weapon_icon, weapon_name, weapon_damage_type, weapon_equipment_slot)
    VALUES ($1, $2, $3, $4, $5)`,
		entity.WeaponHash, entity.WeaponIcon, entity.WeaponName, entity.WeaponDamageType, entity.WeaponEquipmentSlot)

	if err != nil {
		return nil, fmt.Errorf("Error inserting weapon with hash [%d] and name [%s] into table: %v", entity.WeaponHash, entity.WeaponName, err)
	}

	return &entity, nil
}
