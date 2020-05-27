package models

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID               uuid.UUID  `json:"id"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"-"`
	Username         string     `json:"username"`
	RealName         string     `json:"real_name"`
	Email            string     `json:"email"`
	Password         string     `json:"-"`
	PasswordStrength int        `json:"password_strength"`
	LastLoginOn      time.Time  `json:"last_login_on"`
	LastLoginFrom    net.IP     `json:"last_login_from"`
}

func (u *User) Conflicts() bool {
	var count int
	if err := session.Model(&User{}).Where("username = ?", u.Username).Count(&count).Error; err != nil {
		return true
	}
	if count > 0 {
		return true
	}
	return false
}

func (u *User) Validate() error {
	u.Username = strings.ToLower(u.Username)
	if matched, _ := regexp.MatchString(`^[\w\-\.]{1,63}$`, u.Username); !matched {
		return errors.New("username is not valid")
	}
	if len(u.Email) == 0 {
		return errors.New("missing email")
	}
	u.PasswordStrength = passwordStrength(u.Password)
	if u.PasswordStrength < viper.GetInt("minimumPasswordStrength") {
		return errors.New("password is not complex enough")
	}
	return nil
}

func (u *User) EagerLoad() error {
	return nil
}

func (u *User) Save() error {
	return session.Save(u).Error
}

func (u *User) Update(items map[string]interface{}) error {
	if password, ok := items["password"]; ok {
		items["password_strength"] = passwordStrength(password.(string))
	}
	if err := session.Model(&u).Updates(items).Error; err != nil {
		return err
	}

	item, err := GetUser(u.ID)
	*u = *item
	return err
}

func (u *User) Delete() error {
	if u.ID == uuid.Nil {
		return fmt.Errorf("missing Primary Key")
	}
	return session.Delete(u).Error
}

func NewUser(u *User) error {
	if u.ID == uuid.Nil {
		var err error
		if u.ID, err = uuid.NewRandom(); err != nil {
			return err
		}
	}
	if u.Conflicts() {
		return fmt.Errorf("the item already exists")
	}
	if err := u.Validate(); err != nil {
		return err
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	if err := session.Create(&u).Error; err != nil {
		return err
	}
	for _, r := range viper.GetStringSlice("defaultRoles") {
		roleID, err := uuid.Parse(r)
		if err != nil {
			return err
		}
		p := Permission{
			UserID:     u.ID,
			ResourceID: WildCardResource,
			RoleID:     roleID,
		}
		if err := NewPermission(&p); err != nil {
			return err
		}
	}
	return nil
}

func ListUsers(filters map[string]interface{}) ([]User, error) {
	var items []User
	if err := session.Where(filters).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func SearchUsers(query string) ([]User, error) {
	var items []User
	if err := session.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(query)+"%").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func GetUser(id uuid.UUID) (*User, error) {
	var item User
	err := session.Where("id = ?", id).First(&item).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func GetUserByName(username string) (*User, error) {
	var item User
	err := session.Where("username = ?", username).First(&item).Error
	if gorm.IsRecordNotFoundError(err) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func passwordStrength(password string) int {
	var complexity int
	checks := []*regexp.Regexp{
		regexp.MustCompile(`.{8,}`),        // minimum 8 chars
		regexp.MustCompile(`.{12,}`),       // minimum 12 chars
		regexp.MustCompile(`[a-z]`),        // contains lower-case letters
		regexp.MustCompile(`[A-Z]`),        // contains upper-case letters
		regexp.MustCompile(`[0-9]`),        // contains numbers
		regexp.MustCompile(`[^0-9a-zA-Z]`), // contains special characters
	}
	for _, check := range checks {
		if check.MatchString(password) {
			complexity++
		}
	}
	return complexity
}
