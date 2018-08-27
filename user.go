package panel

import (
	"time"

	"github.com/jinzhu/gorm"
)

//User model
type User struct {
	gorm.Model
	Mail     string
	Password string
	IsAdmin  bool
	Expire   time.Time
}
