package controllers

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	// short name for revel
	r "github.com/revel/revel"
	// YOUR APP NAME
	"database/sql"
	"github.com/mweiss/lang-ex-app-server/app/models"
)

// type: revel controller with `*gorm.DB`
// c.Txn will keep `Gdb *gorm.DB`
type GormController struct {
	*r.Controller
	Txn  *gorm.DB
	User *models.User
}

const authenticationHeader = "AuthenticationToken"

// it can be used for jobs
var Gdb *gorm.DB

// init db
func InitDB() {
	var err error

	// open db
	Gdb, err = gorm.Open("mysql", r.Config.StringDefault("db.path", ""))
	if err != nil {
		r.ERROR.Println("FATAL", err)
		panic(err)
	}

	Gdb.AutoMigrate(&models.TestEntity{},
		&models.User{},
		&models.UserLanguage{},
		&models.Badge{},
		&models.Post{},
		&models.PostEdit{},
		&models.PostCorrection{},
		&models.UserAuthentication{})

	// uniquie index if need
	//Gdb.Model(&models.User{}).AddUniqueIndex("idx_user_name", "name")
}

// Initialize the user based on the authentication token.  If the authentication header
// is invalid, then the user is not initialized.
func (c *GormController) InitUser() r.Result {

	authenticationToken := c.Request.Header.Get(authenticationHeader)

	if authenticationToken != "" {
		var userAuthentication *models.UserAuthentication
		c.Txn.Where("login_token = ? and deletion_date is null", authenticationToken).First(&userAuthentication)
		if userAuthentication != nil {
			var user *models.User
			c.Txn.Where("id = ?", userAuthentication.UserId).First(&user)
			c.User = user
		}
	}

	return nil
}

// transactions

// This method fills the c.Txn before each transaction
func (c *GormController) Begin() r.Result {
	txn := Gdb.Begin()
	if txn.Error != nil {
		panic(txn.Error)
	}
	c.Txn = txn
	return nil
}

// This method clears the c.Txn after each transaction
func (c *GormController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	c.Txn.Commit()
	if err := c.Txn.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

// This method clears the c.Txn after each transaction, too
func (c *GormController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	c.Txn.Rollback()
	if err := c.Txn.Error; err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
