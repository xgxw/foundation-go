package utils

import (
	"encoding/json"
	"strings"
)

// ParseOssLsPaths is 处理oss获取的文件列表
func ParseOssLsPaths(paths []string, delimiter string) ([]byte, error) {
	var result = make(map[string]interface{})
	for _, path := range paths {
		parsePath(path, "/", result)
	}
	return json.Marshal(result)
}

func parsePath(path string, delimiter string, result map[string]interface{}) map[string]interface{} {
	if path == "" {
		return nil
	}
	dirs := strings.Split(path, delimiter)
	parseFile(dirs, result)
	return result
}

func parseFile(dirs []string, result map[string]interface{}) map[string]interface{} {
	var key = dirs[0]
	switch len(dirs) {
	case 1:
		if key == "" {
			return nil
		}
		result[key] = key
	case 0:
		return nil
	default:
		// 检测是否添加过该key,且当key为文件时复写为map
		var temp map[string]interface{}
		if v, ok := result[key]; ok {
			var isMap bool
			temp, isMap = v.(map[string]interface{})
			if !isMap {
				temp = make(map[string]interface{})
				result[key] = temp
			}
		} else {
			temp = make(map[string]interface{})
			result[key] = temp
		}
		parseFile(dirs[1:], temp)
	}
	return result
}
