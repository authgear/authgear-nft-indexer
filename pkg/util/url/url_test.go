package url_test

import (
	"net/url"
	"testing"

	urlutil "github.com/authgear/authgear-nft-indexer/pkg/util/url"
	. "github.com/smartystreets/goconvey/convey"
)

func TestUrl(t *testing.T) {
	Convey("pagination param normal", t, func() {

		values := url.Values{
			"limit":  []string{"10"},
			"offset": []string{"20"},
		}

		limit, offset, err := urlutil.ParsePaginationParams(values, 100, 1)
		So(err, ShouldBeNil)
		So(limit, ShouldEqual, 10)
		So(offset, ShouldEqual, 20)

	})

	Convey("pagination param defaultValues", t, func() {

		values := url.Values{}

		limit, offset, err := urlutil.ParsePaginationParams(values, 100, 1)
		So(err, ShouldBeNil)
		So(limit, ShouldEqual, 100)
		So(offset, ShouldEqual, 1)
	})

	Convey("pagination param limit failure", t, func() {

		values := url.Values{
			"limit":  []string{"foo"},
			"offset": []string{"20"},
		}

		limit, offset, err := urlutil.ParsePaginationParams(values, 100, 1)
		So(err, ShouldBeError, "strconv.Atoi: parsing \"foo\": invalid syntax")
		So(limit, ShouldEqual, -1)
		So(offset, ShouldEqual, -1)
	})

	Convey("pagination param offset failure", t, func() {

		values := url.Values{
			"limit":  []string{"100"},
			"offset": []string{"foo"},
		}

		limit, offset, err := urlutil.ParsePaginationParams(values, 100, 1)
		So(err, ShouldBeError, "strconv.Atoi: parsing \"foo\": invalid syntax")
		So(limit, ShouldEqual, -1)
		So(offset, ShouldEqual, -1)
	})
}
