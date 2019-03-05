package my

import (
	"time"

	"github.com/nishrit/aero/db"
)

type IDKey struct {
	ID uint `sql:"auto_increment;not null;primary_key" json:"id" insert:"no" update:"no"`
}

type User struct {
	Name        string     `sql:"TYPE:varchar(128);not null;" json:"name" insert:"must" fako:"full_name"`
	Email       string     `sql:"TYPE:varchar(128);not null;" json:"email" insert:"must" update:"no" unique:"idx_email_uniq" fako:"unique_email"`
	Password    string     `sql:"TYPE:varchar(256);not null;" json:"-" insert:"must" update:"no" fako:"simple_password"`
	Mobile      string     `sql:"TYPE:varchar(16);" json:"mobile" fako:"phone"`
	Phone       string     `sql:"TYPE:varchar(16);" json:"phone" fako:"phone"`
	Active      uint8      `sql:"TYPE:tinyint unsigned;not null;DEFAULT:'1'" json:"active"`
	ActivatedAt *time.Time `sql:"null;" json:"activated_at" insert:"no" update:"no"`
	Verified    uint8      `sql:"TYPE:tinyint unsigned;not null;DEFAULT:'0'" json:"verified"`
	VerifiedAt  *time.Time `sql:"null;" json:"verified_at" insert:"no" update:"no"`
}

type Timed struct {
	CreatedAt time.Time `sql:"not null;DEFAULT:current_timestamp" json:"created_at" insert:"no" update:"no"`
	UpdatedAt time.Time `sql:"not null;DEFAULT:current_timestamp" json:"updated_at" insert:"no" update:"no"`
}

type Persistent struct {
	Deleted   uint8      `sql:"TYPE:tinyint unsigned;not null;DEFAULT:'0'" json:"deleted" insert:"no"`
	DeletedAt *time.Time `sql:"null;" json:"deleted_at" insert:"no" update:"no"`
}

type WWW struct {
	URLWeb       string   `sql:"TYPE:varchar(256);not null" json:"url_web" unique:"true"`
	URLWebOld    *db.JArr `sql:"TYPE:json;" json:"-" insert:"no" update:"no"`
	MetaTitle    string   `sql:"TYPE:varchar(512);not null;DEFAULT:''" json:"meta_title"`
	MetaDesc     string   `sql:"TYPE:varchar(512);not null;DEFAULT:''" json:"meta_desc"`
	MetaKeywords string   `sql:"TYPE:varchar(512);not null;DEFAULT:''" json:"meta_keywords"`
	Sitemap      uint8    `sql:"TYPE:tinyint unsigned;not null;DEFAULT:'1'" json:"sitemap"`
}

type Tagged struct {
	Tags *db.JSArr `sql:"TYPE:json;" json:"tags"`
}

type Ordered struct {
	Sequence uint16 `sql:"TYPE:smallint unsigned;not null;DEFAULT:'1'" json:"sequence"`
}

type LastAction struct {
	Modifier *db.JDoc `sql:"TYPE:json" json:"modifier"`
}

type AuditTrail struct {
	Action     string    `sql:"TYPE:varchar(6);not null;DEFAULT:'insert'"`
	ActionedAt time.Time `sql:"not null;DEFAULT:current_timestamp"`
}

type DynamicFields struct {
	Info *db.JDoc `sql:"TYPE:json" json:"info"`
}

type Stateful struct {
	MachineState *string `sql:"TYPE:varchar(128);null" json:"machine_state"`
}

type NewDB interface {
	ZeroFill() []interface{}
}

type Triggered interface {
	Triggers() []string
}
