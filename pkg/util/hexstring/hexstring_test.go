package hexstring_test

import (
	"math/big"
	"testing"

	"github.com/authgear/authgear-nft-indexer/pkg/util/hexstring"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHexCmp(t *testing.T) {
	Convey("NewFromInt64 normal", t, func() {
		h1, err := hexstring.NewFromInt64(16)
		So(err, ShouldBeNil)
		So(h1, ShouldEqual, "0x10")

		h2, err := hexstring.NewFromInt64(24)
		So(err, ShouldBeNil)
		So(h2, ShouldEqual, "0x18")

		h3, err := hexstring.NewFromInt64(0)
		So(err, ShouldBeNil)
		So(h3, ShouldEqual, "0x0")
	})

	Convey("NewFromInt64 failure", t, func() {
		h1, err := hexstring.NewFromInt64(-20)
		So(err, ShouldBeError, "value must be positive")
		So(h1, ShouldEqual, "")
	})

	Convey("NewFromBigInt normal", t, func() {
		h1, err := hexstring.NewFromBigInt(big.NewInt(5029))
		So(err, ShouldBeNil)
		So(h1, ShouldEqual, "0x13a5")

		h2, err := hexstring.NewFromBigInt(big.NewInt(2738))
		So(err, ShouldBeNil)
		So(h2, ShouldEqual, "0xab2")

		h3, err := hexstring.NewFromBigInt(big.NewInt(0))
		So(err, ShouldBeNil)
		So(h3, ShouldEqual, "0x0")
	})

	Convey("NewFromBigInt failure", t, func() {
		h1, err := hexstring.NewFromBigInt(big.NewInt(-1))
		So(err, ShouldBeError, "value must be positive")
		So(h1, ShouldEqual, "")
	})

	Convey("Parse normal", t, func() {
		h1, err := hexstring.Parse("0x10")
		So(err, ShouldBeNil)
		So(h1, ShouldEqual, "0x10")

		h2, err := hexstring.Parse("0x18")
		So(err, ShouldBeNil)
		So(h2, ShouldEqual, "0x18")

		h3, err := hexstring.Parse("0x0")
		So(err, ShouldBeNil)
		So(h3, ShouldEqual, "0x0")
	})

	Convey("Parse failure", t, func() {
		h, err := hexstring.Parse("10")
		So(err, ShouldBeError, "hex string must start with 0x")
		So(h, ShouldEqual, "")
	})

	Convey("FindSmallest normal", t, func() {
		hexStrings := []hexstring.HexString{
			"0x10",
			"0x18",
			"0x1",
		}
		hex, i, ok := hexstring.FindSmallest(hexStrings)

		So(ok, ShouldBeTrue)
		So(hex, ShouldEqual, "0x1")
		So(i, ShouldEqual, 2)
	})
	Convey("FindSmallest failure", t, func() {
		hexStrings := []hexstring.HexString{}
		hex, i, ok := hexstring.FindSmallest(hexStrings)

		So(ok, ShouldBeFalse)
		So(hex, ShouldEqual, "")
		So(i, ShouldEqual, -1)
	})
}
