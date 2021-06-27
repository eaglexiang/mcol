package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/eaglexiang/mcol/env"
	"github.com/pkg/errors"
)

func getFilename() (fn string, err error) {
	home, err := env.Home()
	if err != nil {
		return
	}

	fn = filepath.Join(home, ".mcol.cache")

	return
}

var cache = make(map[string][]string) // map[db_name] [col_name]

// color
const (
	highlightText = "\033[1;30;47m%s\033[0m"
)

// Matched 被匹配到的结果
type Matched struct {
	Text    string
	Indexes map[int]string // map[index]key
}

// Fmt 格式化后的文本
func (m Matched) Fmt() (fmt string, args []interface{}) {
	// 对 m.Indexes 进行排序
	indexes := make([]int, 0, len(m.Indexes))
	for index := range m.Indexes {
		indexes = append(indexes, index)
	}
	sort.Ints(indexes)

	var lastTail int // 上一个 key 的尾巴
	for _, index := range indexes {
		key := m.Indexes[index]

		// 上一个尾巴到本 key 之前的部分
		fmt += "%s"
		args = append(args, m.Text[lastTail:index])

		// key 使用高亮字体
		fmt += highlightText
		args = append(args, key)

		lastTail = index + len(key)
	}

	fmt += "%s"
	args = append(args, m.Text[lastTail:])

	return
}

// Load 加载缓存
func Load() (err error) {
	filename, err := getFilename()
	if err != nil {
		return
	}

	buf, err := os.ReadFile(filename)
	if err != nil {
		err = errors.WithMessage(err, filename)
		err = errors.WithStack(err)
		return
	}

	err = json.Unmarshal(buf, &cache)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

// Save 保存缓存
func Save(c map[string][]string) (err error) {
	cache = c

	buf, err := json.MarshalIndent(cache, "", "    ")
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	filename, err := getFilename()
	if err != nil {
		return
	}

	err = os.WriteFile(filename, buf, os.FileMode(0666))
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

// Search 搜索
func Search(keys []string) (results []Matched) {
	for dbName, colNames := range cache {
		for _, colName := range colNames {
			dbCol := dbName + "/" + colName

			match := true
			result := Matched{
				Text:    dbCol,
				Indexes: make(map[int]string),
			}
			for _, key := range keys {
				index := strings.Index(dbCol, key)
				if index == -1 {
					match = false
					break
				}

				result.Indexes[index] = key
			}

			if match {
				results = append(results, result)
			}
		}
	}
	return
}
