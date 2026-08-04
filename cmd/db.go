package cmd

import (
	"fmt"
	"log"
	"stock-tracker/constant"
	"stock-tracker/cryption"
	"stock-tracker/db"
	"stock-tracker/model"
	"stock-tracker/utils"
	"time"

	"gorm.io/gorm"
)

func dbUp() {
	log.Println("Running database migrations...")
	dbUpTable()
	dbUpView()
}

func dbUpTable() {
	var err error

	conn := db.GetPostgresConnection()

	err = conn.Migrator().AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().AutoMigrate(&model.Item{})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().AutoMigrate(&model.Itemlog{})
	if err != nil {
		panic(err)
	}
}

func dbUpView() {
	var err error

	conn := db.GetPostgresConnection()

	err = conn.Migrator().DropView(model.VIEW_USER)
	if err != nil {
		panic(err)
	}
	vUser := conn.Model(&model.User{}).Unscoped().
		Select("users.*, '' as photo_url, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = users.create_by").
		Joins("left join users u2 on u2.id = users.update_by")
	err = conn.Migrator().CreateView(model.VIEW_USER, gorm.ViewOption{
		Replace: true,
		Query:   vUser,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_ITEM)
	if err != nil {
		panic(err)
	}
	vItem := conn.Model(&model.Item{}).Unscoped().
		Select("items.*, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join users u1 on u1.id = items.create_by").
		Joins("left join users u2 on u2.id = items.update_by")
	err = conn.Migrator().CreateView(model.VIEW_ITEM, gorm.ViewOption{
		Replace: true,
		Query:   vItem,
	})
	if err != nil {
		panic(err)
	}

	err = conn.Migrator().DropView(model.VIEW_ITEMLOG)
	if err != nil {
		panic(err)
	}
	vItemlog := conn.Model(&model.Itemlog{}).Unscoped().
		Select("itemlogs.*, items.name as item_name, u1.fullname as create_name, u2.fullname as update_name").
		Joins("left join items items on items.id = itemlogs.item_id").
		Joins("left join users u1 on u1.id = itemlogs.create_by").
		Joins("left join users u2 on u2.id = itemlogs.update_by")
	err = conn.Migrator().CreateView(model.VIEW_ITEMLOG, gorm.ViewOption{
		Replace: true,
		Query:   vItemlog,
	})
	if err != nil {
		panic(err)
	}
}

func dbDown() {
	log.Println("Reverting database migrations...")
	var err error

	conn := db.GetPostgresConnection()

	err = conn.Exec("DROP SCHEMA public CASCADE").Error
	if err != nil {
		panic(err)
	}

	err = conn.Exec("CREATE SCHEMA public").Error
	if err != nil {
		panic(err)
	}

	err = conn.Exec("GRANT ALL ON SCHEMA public TO postgres").Error
	if err != nil {
		panic(err)
	}

	err = conn.Exec("GRANT ALL ON SCHEMA public TO public").Error
	if err != nil {
		panic(err)
	}
}

func dbSeed() {
	conn := db.GetPostgresConnection()

	tx := conn.Begin()

	adminID := utils.GetUniqueID()
	userID := utils.GetUniqueID()

	now := time.Now()

	password, err := cryption.EncryptAES64("123456")
	if err != nil {
		panic(err)
	}

	users := []model.User{
		{
			ID:                adminID,
			Role:              constant.RoleAdmin,
			Email:             "jihanlugas2@gmail.com",
			Username:          "jihanlugas",
			PhoneNumber:       utils.FormatPhoneTo62("6287770333043"),
			Fullname:          "Jihan Lugas",
			Address:           "Jl. Gunung Sahari No. 10, Jakarta Pusat",
			Passwd:            password,
			PassVersion:       1,
			IsActive:          true,
			AccountVerifiedDt: &now,
			CreateBy:          adminID,
			UpdateBy:          adminID,
		},
		{
			ID:                userID,
			Role:              constant.RoleUser,
			Email:             "user1@gmail.com",
			Username:          "user1",
			PhoneNumber:       utils.FormatPhoneTo62("6287770333043"),
			Fullname:          "User 1",
			Address:           "Jl. Ahmad Yani No. 10, Jakarta Timur",
			Passwd:            password,
			PassVersion:       1,
			IsActive:          true,
			AccountVerifiedDt: &now,
			CreateBy:          adminID,
			UpdateBy:          adminID,
		},
	}
	tx.Create(&users)

	items := []model.Item{
		{
			ID:       utils.GetUniqueID(),
			Name:     "Marinir",
			Notes:    "Kalender Marinir 2027",
			Stock:    0,
			Sent:     0,
			CreateBy: adminID,
			UpdateBy: adminID,
		},
		{
			ID:       utils.GetUniqueID(),
			Name:     "AU",
			Notes:    "Kalender AU 2027",
			Stock:    0,
			Sent:     0,
			CreateBy: adminID,
			UpdateBy: adminID,
		},
	}
	tx.Create(&items)

	err = tx.Commit().Error
	if err != nil {
		panic(fmt.Errorf("failed to commit transaction: %w", err))
	}

	log.Println("Seeding the database with initial data end")
}

func dbReset() {
	dbDown()
	dbUp()
	dbSeed()
}
