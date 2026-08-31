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

	marinirID := utils.GetUniqueID()
	auID := utils.GetUniqueID()

	items := []model.Item{
		{
			ID:       marinirID,
			Name:     "Marinir 2027",
			Notes:    "Kalender Marinir 2027",
			Stock:    15650,
			Sent:     25850,
			CreateBy: adminID,
			CreateDt: time.Now().AddDate(1, 0, -30),
			UpdateBy: adminID,
		},
		{
			ID:       auID,
			Name:     "AU 2027",
			Notes:    "Kalender AU 2027",
			Stock:    12000,
			Sent:     22500,
			CreateBy: adminID,
			CreateDt: time.Now().AddDate(1, 0, -30),
			UpdateBy: adminID,
		},
	}
	tx.Create(&items)

	itemlogs := []model.Itemlog{
		// =========================================================
		// MARINIR - Penerimaan (STOCK)
		// =========================================================
		{
			ItemID:   marinirID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 1",
			Quantity: 5000,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -30),
		},
		{
			ItemID:   marinirID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 2",
			Quantity: 5050,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -28),
		},
		{
			ItemID:   marinirID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 3",
			Quantity: 3000,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -25),
		},
		{
			ItemID:   marinirID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 4",
			Quantity: 2500,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -22),
		},
		{
			ItemID:   marinirID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 5",
			Quantity: 4500,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -19),
		},
		{
			ItemID:   marinirID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 6",
			Quantity: 3500,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -16),
		},
		{
			ItemID:   marinirID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 7",
			Quantity: 5000,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -13),
		},
		{
			ItemID:   marinirID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 8",
			Quantity: 2750,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -10),
		},
		{
			ItemID:   marinirID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 9",
			Quantity: 4000,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -7),
		},
		{
			ItemID:   marinirID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 10",
			Quantity: 3200,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -4),
		},
		{
			ItemID:   marinirID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 11",
			Quantity: 3000,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -2),
		},

		// =========================================================
		// MARINIR - Pengiriman (SENT)
		// =========================================================
		{
			ItemID:   marinirID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko A",
			Quantity: 2000,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -27),
		},
		{
			ItemID:   marinirID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko B",
			Quantity: 3000,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -24),
		},
		{
			ItemID:   marinirID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko C",
			Quantity: 2500,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -21),
		},
		{
			ItemID:   marinirID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko D",
			Quantity: 1300,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -18),
		},
		{
			ItemID:   marinirID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko E",
			Quantity: 2800,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -15),
		},
		{
			ItemID:   marinirID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko F",
			Quantity: 1750,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -12),
		},
		{
			ItemID:   marinirID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko G",
			Quantity: 2200,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -10),
		},
		{
			ItemID:   marinirID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko H",
			Quantity: 1900,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -8),
		},
		{
			ItemID:   marinirID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko I",
			Quantity: 2600,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -6),
		},
		{
			ItemID:   marinirID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko J",
			Quantity: 2100,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -4),
		},
		{
			ItemID:   marinirID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko K",
			Quantity: 2400,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -2),
		},
		{
			ItemID:   marinirID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko L",
			Quantity: 1300,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now,
		},

		// =========================================================
		// AU - Penerimaan (STOCK)
		// =========================================================
		{
			ItemID:   auID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 1",
			Quantity: 3000,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -30),
		},
		{
			ItemID:   auID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 2",
			Quantity: 3200,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -27),
		},
		{
			ItemID:   auID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 3",
			Quantity: 2500,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -24),
		},
		{
			ItemID:   auID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 4",
			Quantity: 4000,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -21),
		},
		{
			ItemID:   auID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 5",
			Quantity: 3500,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -18),
		},
		{
			ItemID:   auID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 6",
			Quantity: 2800,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -15),
		},
		{
			ItemID:   auID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 7",
			Quantity: 4500,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -12),
		},
		{
			ItemID:   auID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 8",
			Quantity: 3000,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -9),
		},
		{
			ItemID:   auID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 9",
			Quantity: 3500,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -6),
		},
		{
			ItemID:   auID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 10",
			Quantity: 2500,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -3),
		},
		{
			ItemID:   auID,
			Type:     "STOCK",
			Notes:    "Penerimaan barang dari supplier - batch 11",
			Quantity: 2000,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now,
		},

		// =========================================================
		// AU - Pengiriman (SENT)
		// =========================================================
		{
			ItemID:   auID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko A",
			Quantity: 1800,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -28),
		},
		{
			ItemID:   auID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko B",
			Quantity: 2200,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -25),
		},
		{
			ItemID:   auID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko C",
			Quantity: 2000,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -22),
		},
		{
			ItemID:   auID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko D",
			Quantity: 1500,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -19),
		},
		{
			ItemID:   auID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko E",
			Quantity: 2500,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -16),
		},
		{
			ItemID:   auID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko F",
			Quantity: 1500,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -13),
		},
		{
			ItemID:   auID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko G",
			Quantity: 1800,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -11),
		},
		{
			ItemID:   auID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko H",
			Quantity: 2300,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -9),
		},
		{
			ItemID:   auID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko I",
			Quantity: 1700,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -7),
		},
		{
			ItemID:   auID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko J",
			Quantity: 2100,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -5),
		},
		{
			ItemID:   auID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko K",
			Quantity: 1600,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now.AddDate(0, 0, -3),
		},
		{
			ItemID:   auID,
			Type:     "SENT",
			Notes:    "Pengiriman ke toko L",
			Quantity: 1500,
			CreateBy: adminID,
			UpdateBy: adminID,
			CreateDt: now,
		},
	}

	tx.Create(&itemlogs)

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
