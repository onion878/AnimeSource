package utils

import (
	"../structs"
	"fmt"
	"time"
)

func SaveOrUpdateIndex(name string, chapter string) structs.Index {
	var index structs.Index
	var indexs []structs.Index
	engine := GetCon()
	engine.Where("name = ?", name).Find(&indexs)
	if len(indexs) == 0 {
		index.Id = NewKeyId()
		index.Update = time.Now()
		index.Total = 0
		index.Name = name
		index.Chapter = chapter
		engine.Insert(&index)
		index.Flag = true
	} else {
		index = indexs[0]
		index.Chapter = chapter
		index.Update = time.Now()
		engine.Update(&index)
		index.Flag = false
	}
	fmt.Printf("%+v\n", index)
	return index
}

func SaveChapter(name string, pid string, url string, num int) structs.Chapter {
	var chapter structs.Chapter
	var chapters []structs.Chapter
	engine := GetCon()
	engine.Where("name = ?", name).Find(&chapters)
	if len(chapters) == 0 {
		chapter.Id = NewKeyId()
		chapter.Pid = pid
		chapter.Name = name
		chapter.Path = url
		chapter.Num = num
		engine.Insert(&chapter)
	} else {
		chapter = chapters[0]
		chapter.Pid = pid
		chapter.Name = name
		chapter.Path = url
		chapter.Num = num
	}
	return chapter
}