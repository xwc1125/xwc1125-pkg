// Package es
//
// @author: xwc1125
// @date: 2020/9/4 0004
package es

import (
	"fmt"
	"testing"
	"time"

	"github.com/olivere/elastic/v7"
)

type Tweet struct {
	User     string                `json:"user" es:"text"`
	Message  string                `json:"message" es:"text"`
	Retweets int                   `json:"retweets"  es:"text"`
	Image    string                `json:"image,omitempty"  es:"text"`
	Created  time.Time             `json:"created,omitempty"  es:"store"`
	Tags     []string              `json:"tags,omitempty"  es:"text"`
	Location string                `json:"location,omitempty"`
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"  es:"text"`
	A        float32
}

func TestTweet_Mapping(t *testing.T) {
	mapping := Mapping(Tweet{})
	fmt.Println(mapping)
}
