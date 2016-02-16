package refl

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestIsAddress(t *testing.T) {

	type Strukt struct {
		AField string
	}
	str := "abc"
	var direct interface{} = Strukt{}
	var reference interface{} = &Strukt{}

	Convey("Given pointers, fn:IsAddress must return true", t, func() {
		So(IsAddress(&str), ShouldBeTrue)
		So(IsAddress(&Strukt{}), ShouldBeTrue)
		So(IsAddress(reference), ShouldBeTrue)
	})

	Convey("Given literals, fn:IsAddress must return false", t, func() {
		So(IsAddress(str), ShouldBeFalse)
		So(IsAddress(Strukt{}), ShouldBeFalse)
		So(IsAddress(direct), ShouldBeFalse)
	})
}

func TestComposedOf(t *testing.T) {

	type GrandFather struct {
		AField string
	}
	type Father struct {
		BField string
		GrandFather
	}
	type Son struct {
		CField string
		Father
	}
	type Uncle struct {
		DField string
	}

	Convey("Directly composed struct must give true", t, func() {
		So(ComposedOf(Son{}, Father{}), ShouldBeTrue)
	})

	Convey("Inirectly composed struct must give true", t, func() {
		So(ComposedOf(Son{}, GrandFather{}), ShouldBeTrue)
	})

	Convey("Disconnected struct must give false", t, func() {
		So(ComposedOf(Son{}, Uncle{}), ShouldBeFalse)
	})

}
