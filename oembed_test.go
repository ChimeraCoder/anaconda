package anaconda_test

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/ChimeraCoder/anaconda"
)

func TestOEmbed(t *testing.T) {
	// It is the only one that can be tested without auth
	api := anaconda.NewTwitterApi("", "")
	o, err := api.GetOEmbed(url.Values{"id": []string{"99530515043983360"}})
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(o, expectedOEmbed) {
		t.Error("Actual OEmbed differs from expected", o)
	}
}

var expectedOEmbed anaconda.OEmbed = anaconda.OEmbed{
	Cache_age:     "3153600000",
	Url:           "https://twitter.com/twitter/statuses/99530515043983360",
	Height:        0,
	Provider_url:  "https://twitter.com",
	Provider_name: "Twitter",
	Author_name:   "Twitter",
	Version:       "1.0",
	Author_url:    "https://twitter.com/twitter",
	Type:          "rich",
	Html:          "\u003Cblockquote class=\"twitter-tweet\"\u003E\u003Cp\u003ECool! \u201C\u003Ca href=\"https://twitter.com/tw1tt3rart\"\u003E@tw1tt3rart\u003C/a\u003E: \u003Ca href=\"https://twitter.com/hashtag/TWITTERART?src=hash\"\u003E#TWITTERART\u003C/a\u003E \u2571\u2571\u2571\u2571\u2571\u2571\u2571\u2571 \u2571\u2571\u256D\u2501\u2501\u2501\u2501\u256E\u2571\u2571\u256D\u2501\u2501\u2501\u2501\u256E \u2571\u2571\u2503\u2587\u2506\u2506\u2587\u2503\u2571\u256D\u252B\u24E6\u24D4\u24D4\u24DA\u2503 \u2571\u2571\u2503\u25BD\u25BD\u25BD\u25BD\u2503\u2501\u256F\u2503\u2661\u24D4\u24DD\u24D3\u2503 \u2571\u256D\u252B\u25B3\u25B3\u25B3\u25B3\u2523\u256E\u2571\u2570\u2501\u2501\u2501\u2501\u256F \u2571\u2503\u2503\u2506\u2506\u2506\u2506\u2503\u2503\u2571\u2571\u2571\u2571\u2571\u2571 \u2571\u2517\u252B\u2506\u250F\u2513\u2506\u2523\u251B\u2571\u2571\u2571\u2571\u2571\u201D\u003C/p\u003E&mdash; Twitter (@twitter) \u003Ca href=\"https://twitter.com/twitter/status/99530515043983360\"\u003EAugust 5, 2011\u003C/a\u003E\u003C/blockquote\u003E\n\u003Cscript async src=\"//platform.twitter.com/widgets.js\" charset=\"utf-8\"\u003E\u003C/script\u003E",
	Width:         550,
}
