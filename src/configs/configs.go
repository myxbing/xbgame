package configs

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

const _DEF_CONFIG = "game.config"

var (
	_map         map[string]string
	_lock        sync.RWMutex
	_config_file = flag.String("config", _DEF_CONFIG, "specify absolute path for game.config")
)

// 获取数值配置
func Int64(key string) (int64, error) {
	_lock.RLock()
	defer _lock.RUnlock()
	return strconv.ParseInt(_map[key], 0, 0)
}

// 获取数值配置
func UInt64(key string) (uint64, error) {
	_lock.RLock()
	defer _lock.RUnlock()
	return strconv.ParseUint(_map[key], 0, 0)
}

// 获取数值配置
func Float(key string) (float64, error) {
	_lock.RLock()
	defer _lock.RUnlock()
	return strconv.ParseFloat(_map[key], 0)
}

// 获取布尔类型配置
func Bool(key string) (bool, error) {
	_lock.RLock()
	defer _lock.RUnlock()
	return strconv.ParseBool(_map[key])
}

// 获取数值配置
func Int(key string) (int, error) {
	_lock.RLock()
	defer _lock.RUnlock()
	return strconv.Atoi(_map[key])
}

// 获取字符串配置
func String(key string) string {
	_lock.RLock()
	defer _lock.RUnlock()
	return _map[key]
}

// 判断配置是否正确
func Is(key, value string) bool {
	_lock.RLock()
	defer _lock.RUnlock()
	val, exists := _map[key]
	if !exists {
		return false
	}
	return strings.EqualFold(value, val)
}

func init() {
	Reload()
}

//获取加载配置信息
func Load() map[string]string {
	_lock.RLock()
	defer _lock.RUnlock()
	return _map
}

//重新载入配置
func Reload() {
	path := *_config_file
	log.Println("Loading config from file", path, "...")
	defer log.Println("Config Loaded.")
	_lock.Lock()
	_map = _load_config(path)
	_lock.Unlock()
}

func _load_config(path string) (ret map[string]string) {
	ret = make(map[string]string)
	file, err := os.Open(path)

	if err != nil {
		fmt.Println("error opening file %v\n", err)
		os.Exit(-1)
	}
	defer file.Close()

	// using scanner to read config file
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	multiLine := false
	var value, key string
	for scanner.Scan() {
		line := scanner.Text()
		if multiLine {
			temp := strings.TrimRightFunc(line, unicode.IsSpace)
			if strings.HasSuffix(temp, "\"") {
				value += strings.TrimRight(temp, "\"")
				ret[key] = strings.Replace(value, "\\\"", "\"", -1)
				multiLine = false
			} else {
				value += line
				value += "\n"
			}
			continue
		}

		unspaces := strings.TrimSpace(line)
		if strings.HasPrefix(unspaces, "#") || len(unspaces) == 0 {
			continue
		}

		slice := strings.SplitN(unspaces, "=", 2)
		if len(slice) == 2 {
			key = strings.TrimSpace(slice[0])
			temp := strings.TrimSpace(slice[1])
			if strings.HasPrefix(temp, "\"") {
				temp = strings.TrimLeft(temp, "\"")
				if strings.HasSuffix(temp, "\"") {
					value = strings.TrimRight(temp, "\"")
				} else {
					multiLine = true
					temp = strings.SplitN(line, "=", 2)[1]
					value = strings.TrimLeftFunc(temp, unicode.IsSpace)
					value = strings.TrimLeft(value, "\"")
					value += "\n"
					continue
				}
			} else {
				value = temp
			}

			ret[key] = value
		}
	}

	return
}
