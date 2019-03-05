package orm

import (
	"testing"

	"github.com/nishrit/aero/db"
	. "github.com/smartystreets/goconvey/convey"
)

func TestModelMustInsert(t *testing.T) {

	type Abc struct {
		Alphabet string `insert:"must"`
		Numbers  string `insert:"no"`
	}

	Convey("Missing must-insert field should give an error", t, func() {
		ok, errs := Insertable(Abc{}, map[string]string{"some": "thing"})
		So(ok, ShouldBeFalse)
		So(len(errs), ShouldEqual, 1)

	})
	Convey("Present must-insert and missing no-insert field should be ok", t, func() {
		ok, errs := Insertable(Abc{}, map[string]string{"alphabet": "thing"})
		So(ok, ShouldBeTrue)
		So(errs, ShouldBeNil)
	})
	Convey("Present no-insert field should give an error", t, func() {
		ok, errs := Insertable(Abc{}, map[string]string{"alphabet": "thing", "numbers": "1234"})
		So(ok, ShouldBeFalse)
		So(len(errs), ShouldEqual, 1)
	})
}

func TestModelMustUpdate(t *testing.T) {

	type Def struct {
		Alphabet string `update:"no"`
		Numbers  string `update:"must"`
	}

	Convey("Missing must-update field should give an error", t, func() {
		ok, errs := Updatable(Def{}, map[string]string{"some": "thing"})
		So(ok, ShouldBeFalse)
		So(len(errs), ShouldEqual, 1)
	})
	Convey("Present must-update and missing no-update field should be ok", t, func() {
		ok, errs := Updatable(Def{}, map[string]string{"numbers": "thing"})
		So(ok, ShouldBeTrue)
		So(errs, ShouldBeNil)
	})
	Convey("Present no-update field should give an error", t, func() {
		ok, errs := Updatable(Def{}, map[string]string{"alphabet": "thing", "numbers": "1234"})
		So(ok, ShouldBeFalse)
		So(len(errs), ShouldEqual, 1)
	})

}

func TestModelJsonArray(t *testing.T) {

	var ok bool
	var errs []error

	type Ghi struct {
		Array    db.JArr
		ArrayPtr *db.JArr
	}

	Convey("Given a struct that expects Json array field", t, func() {

		Convey("Then passing a string value should give error", func() {
			ok, errs = Insertable(Ghi{}, map[string]string{"array": "thing"})
			So(ok, ShouldBeFalse)
			So(len(errs), ShouldEqual, 1)

			ok, errs = Updatable(Ghi{}, map[string]string{"array": "thing"})
			So(ok, ShouldBeFalse)
			So(len(errs), ShouldEqual, 1)

			ok, errs = Insertable(Ghi{}, map[string]string{"array_ptr": "thinger"})
			So(ok, ShouldBeFalse)
			So(len(errs), ShouldEqual, 1)

			ok, errs = Updatable(Ghi{}, map[string]string{"array_ptr": "thinger"})
			So(ok, ShouldBeFalse)
			So(len(errs), ShouldEqual, 1)
		})

		Convey("And passing an Array of int should be Ok", func() {
			ok, errs = Insertable(Ghi{}, map[string]string{"array": "[1,2,3]"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)

			ok, errs = Updatable(Ghi{}, map[string]string{"array": "[1,2,3]"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)

		})

		Convey("And passing an Array of string,int should be Ok", func() {
			ok, errs = Insertable(Ghi{}, map[string]string{"array_ptr": "[\"abc\", 1]"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)

			ok, errs = Updatable(Ghi{}, map[string]string{"array_ptr": "[\"abc\", 1]"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)
		})
	})
}

func TestModelJsonDoc(t *testing.T) {

	var ok bool
	var errs []error

	type Jkl struct {
		Mapper    db.JDoc
		MapperPtr *db.JDoc
	}

	Convey("Given a struct that expects Json Doc field", t, func() {

		Convey("Then passing a String should give error", func() {
			ok, errs = Insertable(Jkl{}, map[string]string{"mapper": "thing"})
			So(ok, ShouldBeFalse)
			So(len(errs), ShouldEqual, 1)

			ok, errs = Updatable(Jkl{}, map[string]string{"mapper": "thing"})
			So(ok, ShouldBeFalse)
			So(len(errs), ShouldEqual, 1)

			ok, errs = Insertable(Jkl{}, map[string]string{"mapper_ptr": "thinger"})
			So(ok, ShouldBeFalse)
			So(len(errs), ShouldEqual, 1)

			ok, errs = Updatable(Jkl{}, map[string]string{"mapper_ptr": "thinger"})
			So(ok, ShouldBeFalse)
			So(len(errs), ShouldEqual, 1)
		})

		Convey("Then passing an Array should give error", func() {
			ok, errs = Insertable(Jkl{}, map[string]string{"mapper": "[1,2,3]"})
			So(ok, ShouldBeFalse)
			So(len(errs), ShouldEqual, 1)

			ok, errs = Updatable(Jkl{}, map[string]string{"mapper": "[1,2,3]"})
			So(ok, ShouldBeFalse)
			So(len(errs), ShouldEqual, 1)

			ok, errs = Insertable(Jkl{}, map[string]string{"mapper": "[1,2,3]"})
			So(ok, ShouldBeFalse)
			So(len(errs), ShouldEqual, 1)

			ok, errs = Updatable(Jkl{}, map[string]string{"mapper_ptr": "[1,2,3]"})
			So(ok, ShouldBeFalse)
			So(len(errs), ShouldEqual, 1)
		})

		Convey("Then passing a Document should be Ok", func() {
			ok, errs = Insertable(Jkl{}, map[string]string{"mapper": "{\"key\":\"value\"}"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)

			ok, errs = Updatable(Jkl{}, map[string]string{"mapper": "{\"key\":\"value\"}"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)

			ok, errs = Insertable(Jkl{}, map[string]string{"mapper_ptr": "{\"key\":\"value\"}"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)

			ok, errs = Updatable(Jkl{}, map[string]string{"mapper_ptr": "{\"key\":\"value\"}"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)
		})

	})

}

func TestModelJson(t *testing.T) {
	var ok bool
	var errs []error

	type Mno struct {
		Any    db.JRaw
		AnyPtr *db.JRaw
	}

	Convey("Given a struct that expects any valid Json field", t, func() {

		Convey("Then passing a String literal should be Ok", func() {
			ok, errs = Insertable(Mno{}, map[string]string{"any": "\"thing\""})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)

			ok, errs = Updatable(Mno{}, map[string]string{"any": "\"thing\""})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)

			ok, errs = Insertable(Mno{}, map[string]string{"any_ptr": "\"thinger\""})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)

			ok, errs = Updatable(Mno{}, map[string]string{"any_ptr": "\"thinger\""})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)
		})

		Convey("Then passing an Array should give error", func() {
			ok, errs = Insertable(Mno{}, map[string]string{"mapper": "[1,2,3]"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)

			ok, errs = Updatable(Mno{}, map[string]string{"mapper": "[1,2,3]"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)

			ok, errs = Insertable(Mno{}, map[string]string{"mapper": "[1,2,3]"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)

			ok, errs = Updatable(Mno{}, map[string]string{"mapper_ptr": "[1,2,3]"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)
		})

		Convey("Then passing a Document should be Ok", func() {
			ok, errs = Insertable(Mno{}, map[string]string{"mapper": "{\"key\":\"value\"}"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)

			ok, errs = Updatable(Mno{}, map[string]string{"mapper": "{\"key\":\"value\"}"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)

			ok, errs = Insertable(Mno{}, map[string]string{"mapper_ptr": "{\"key\":\"value\"}"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)

			ok, errs = Updatable(Mno{}, map[string]string{"mapper_ptr": "{\"key\":\"value\"}"})
			So(ok, ShouldBeTrue)
			So(len(errs), ShouldEqual, 0)
		})

	})

}
