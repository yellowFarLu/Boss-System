package models

import (
	"testing"

	"fmt"
)

func Test_NewTagsServer(t *testing.T) {
	InitTest()
	tagsCache := NewTagsServer()
	for _, v := range tagsCache.tags {
		fmt.Println(v.TagId, v.Name, v.Type, v.Note)
	}
}
