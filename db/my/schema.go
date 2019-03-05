package my

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/nishrit/aero/refl"
	"github.com/nishrit/aero/str"
)

var Dbo *gorm.DB

func Build(history bool, models ...interface{}) {

	if Dbo == nil {
		panic("Dbo referece is null")
	}
	Dbo.LogMode(true)

	// delete tables
	for _, model := range models {
		tbl := NewTable2(model)
		his := tbl.history()

		// delete tables
		if tbl.exists() && tbl.hasData() {
			panic("table has data, can't delete: " + tbl.name)
		}
		if his.exists() && his.hasData() {
			panic("history table has data, can't delete: " + his.name)
		}
		tbl.drop(false)
		his.drop(false)
	}

	// migrate tables
	for _, model := range models {
		Dbo.AutoMigrate(model)
	}

	// create history tables
	if history {
		for _, model := range models {
			setupHistoryAndLogging(model)
		}
	}

	// add foreign keys
	for _, model := range models {
		setupFKeys(model)
	}

	// setup behavior triggers like Timed, Persistent, WWW
	for _, model := range models {
		setupBehaviors(model)
	}

	// add unique indexes
	for _, model := range models {
		setupUniqIndexes(model)
	}

	// add custom triggers
	for _, model := range models {
		populateTriggers(model)
	}

	// initialize records
	for _, model := range models {
		populateRecords(model)
	}
}

func sqlHasRows(sql string) bool {
	rows, err := Dbo.Raw(sql).Rows()
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	return rows.Next()
}

func sqlExec(sql string) {
	err := Dbo.Exec(sql).Error
	if err != nil {
		panic(err)
	}
}

func setupFKeys(model interface{}) {
	tbl := NewTable2(model)

	if tbl.isHistory() {
		return
	}

	// format:
	// fk:"table_name(identity_key)"

	// add foreign keys
	mt := reflect.TypeOf(model).Elem()
	num := mt.NumField()
	for i := 0; i < num; i++ {
		fld := mt.FieldByIndex([]int{i})
		tag := fld.Tag
		if len(tag.Get("fk")) > 0 {
			fk := str.SnakeCase(fld.Name)
			Dbo.Model(model).AddForeignKey(fk, tag.Get("fk"), "RESTRICT", "RESTRICT")
		}
	}
}

func setupUniqIndexes(model interface{}) {
	tbl := NewTable2(model)

	if tbl.isHistory() {
		return
	}

	// formats:
	// unique:"true"
	// unique:"idx_name"
	// unique:"idx_name(field1,field2)"

	allFlds := refl.NestedFields(reflect.ValueOf(model).Elem().Interface())
	for i := 0; i < len(allFlds); i++ {
		fld := allFlds[i]
		if len(fld.Tag.Get("unique")) > 0 {
			name := fld.Tag.Get("unique")
			if name == "true" { // generate index name
				name = "idx_" + str.SnakeCase(fld.Name) + "_unique"
			}
			lbrace := strings.Index(name, "(")
			if lbrace == -1 {
				Dbo.Model(model).AddUniqueIndex(name, str.SnakeCase(fld.Name))
			} else {
				fldCsv := name[lbrace+1 : len(name)-1]
				flds := strings.Split(fldCsv, ",")
				for i := range flds {
					flds[i] = strings.TrimSpace(flds[i])
				}
				Dbo.Model(model).AddUniqueIndex(name[:lbrace], flds...)
			}
		}
	}
}

func setupBehaviors(model interface{}) {
	tbl := NewTable2(model)

	if tbl.isHistory() {
		return
	}

	// setup on update trigger (updated_at field)
	if refl.ComposedOf(model, Timed{}) {
		updAt := tbl.field("updated_at")
		if !strings.Contains(strings.ToLower(updAt.Extra), "on update current_timestamp") {
			sql := fmt.Sprintf("ALTER TABLE %s MODIFY COLUMN %s ON UPDATE CURRENT_TIMESTAMP", tbl.name, updAt.info())
			sqlExec(sql)
		}
	}

	// setup user triggers
	if refl.ComposedOf(model, User{}) {
		sql := fmt.Sprintf(`CREATE TRIGGER %s_user_insert BEFORE INSERT ON %s FOR EACH ROW
        BEGIN
            IF (NEW.active = 1) THEN
                SET NEW.activated_at = NOW();
            END IF;
            IF (NEW.verified = 1) THEN
                SET NEW.verified_at = NOW();
            END IF;
        END`, tbl.name, tbl.name)
		sqlExec(sql)

		sql = fmt.Sprintf(`CREATE TRIGGER %s_user_update BEFORE UPDATE ON %s FOR EACH ROW
        BEGIN
            IF (OLD.active = 0) AND (NEW.active = 1) THEN
                SET NEW.activated_at = NOW();
            END IF;
            IF (OLD.verified = 0) AND (NEW.verified = 1) THEN
                SET NEW.verified_at = NOW();
            END IF;
        END`, tbl.name, tbl.name)
		sqlExec(sql)

	}

	// do not allow deletes
	if refl.ComposedOf(model, Persistent{}) {
		sql := fmt.Sprintf(`CREATE TRIGGER %s_persistent_delete BEFORE DELETE ON %s FOR EACH ROW
            IF TRUE THEN 
                SIGNAL SQLSTATE '45000'
                SET MESSAGE_TEXT = 'Cannot delete records from table. Instead set deleted=1';
            END IF;`, tbl.name, tbl.name)
		sqlExec(sql)

		sql = fmt.Sprintf(`CREATE TRIGGER %s_persistent_update BEFORE UPDATE ON %s FOR EACH ROW
        BEGIN
            IF (OLD.deleted = 0) AND (NEW.deleted = 1) THEN
                SET NEW.deleted_at = NOW();
            END IF;
            IF (OLD.deleted = 1) AND (NEW.deleted = 0) THEN
                SET NEW.deleted_at = NULL;
            END IF;
        END`, tbl.name, tbl.name)
		sqlExec(sql)
	}

	// url_past : push changes of url_web into url_past
	if refl.ComposedOf(model, WWW{}) {
		sql := fmt.Sprintf(`CREATE TRIGGER %s_www_url_update BEFORE UPDATE ON %s FOR EACH ROW
        BEGIN
            IF STRCMP(LEFT(NEW.url_web,1),'/') <> 0 THEN
                SET NEW.url_web = CONCAT('/', NEW.url_web);
            END IF;

            IF (OLD.url_web <> "") AND (NEW.url_web <> OLD.url_web) THEN
                IF NEW.url_web_old IS NULL THEN
                    SET NEW.url_web_old = JSON_ARRAY();
                END IF;
                IF JSON_CONTAINS(NEW.url_web_old, JSON_ARRAY(OLD.url_web)) = 0 THEN
                    SET NEW.url_web_old = JSON_ARRAY_APPEND(NEW.url_web_old, "$", OLD.url_web);
                END IF;
            END IF;
        END`, tbl.name, tbl.name)
		sqlExec(sql)

		sql = fmt.Sprintf(`CREATE TRIGGER %s_www_url_insert BEFORE INSERT ON %s FOR EACH ROW
        BEGIN
            IF STRCMP(LEFT(NEW.url_web,1),'/') <> 0 THEN
                SET NEW.url_web = CONCAT('/', NEW.url_web);
            END IF;
        END`, tbl.name, tbl.name)
		sqlExec(sql)
	}
}

func setupHistoryAndLogging(model interface{}) *table {

	tbl := NewTable2(model)
	his := tbl.history()
	if his.exists() {
		panic("history table already exists")
	}

	// create table (alike)
	sql := fmt.Sprintf("create table %s like %s;", his.name, tbl.name)
	sqlExec(sql)

	// add action and actioned_at columns
	sql = fmt.Sprintf("alter table %s add column action varchar(6) not null default 'insert' first, add column actioned_at TIMESTAMP default current_timestamp after action", his.name)
	sqlExec(sql)

	// remove auto_increment
	autoInc := his.autoIncrField()
	if autoInc != nil {
		noAuto := strings.Replace(autoInc.info(), "auto_increment", "", -1)
		sql = fmt.Sprintf("ALTER TABLE %s MODIFY %s", his.name, noAuto)
		sqlExec(sql)
	}

	// drop primary key
	pkeys := his.primaryKeys()
	if len(pkeys) > 0 {
		sql := fmt.Sprintf("alter table %s drop primary key", his.name)
		sqlExec(sql)
	}

	// Setup audit triggers on original table

	// insert trigger
	sqlExec(fmt.Sprintf("drop trigger if exists %s_audit_trail_insert", tbl.name))
	sql = fmt.Sprintf(`CREATE TRIGGER %s_audit_trail_insert AFTER INSERT ON %s FOR EACH ROW
        INSERT INTO %s SELECT 'insert',null, src.* 
        FROM %s as src WHERE src.id = NEW.id;`, tbl.name, tbl.name, his.name, tbl.name)
	sqlExec(sql)

	// update trigger
	sqlExec(fmt.Sprintf("drop trigger if exists %s_audit_trail_update", tbl.name))
	sql = fmt.Sprintf(`CREATE TRIGGER %s_audit_trail_update AFTER UPDATE ON %s FOR EACH ROW
    	INSERT INTO %s SELECT 'update',null, src.*
        FROM %s as src WHERE src.id = NEW.id;`, tbl.name, tbl.name, his.name, tbl.name)
	sqlExec(sql)

	// delete trigger
	sqlExec(fmt.Sprintf("drop trigger if exists %s_audit_trail_delete", tbl.name))
	sql = fmt.Sprintf(`CREATE TRIGGER %s_audit_trail_delete BEFORE DELETE ON %s FOR EACH ROW
        INSERT INTO %s SELECT 'delete',null, src.*
        FROM %s as src WHERE src.id = OLD.id;`, tbl.name, tbl.name, his.name, tbl.name)
	sqlExec(sql)

	return his
}

func populateTriggers(model interface{}) {
	if m, ok := model.(Triggered); ok {
		trigs := m.Triggers()
		for _, trg := range trigs {
			err := Dbo.Exec(trg).Error
			if err != nil {
				panic(err)
			}
		}
	}
}

func populateRecords(model interface{}) {
	if m, ok := model.(NewDB); ok {
		tx := Dbo.Begin()
		{
			recs := m.ZeroFill()
			for _, rec := range recs {
				err := Dbo.Model(model).Create(rec).Error
				if err != nil {
					tx.Rollback()
					panic(err)
				}
			}
		}
		tx.Commit()
	}
}
