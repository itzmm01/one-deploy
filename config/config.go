package config

import (
	"one-backup/keygen"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"github.com/wonderivan/logger"
)

const Debug = false

// ModelConfig for special case
type ModelConfig struct {
	StoreWith    map[string]string
	CompressType string
	BackupNum    int
	Databases    []map[string]string
	Viper        *viper.Viper
	IsEncrypt    bool
}

// SubConfig sub config info
type SubConfig struct {
	Name  string
	Type  string
	Viper *viper.Viper
}

func passBase(autoEncrypt, pass string, isEncrypt bool) string {
	if autoEncrypt == "no" {
		return pass
	}
	if isEncrypt {
		return pass
	}
	return keygen.AesEncryptCBC(pass, "pass")
}

func interfaceTomap(inter interface{}, autoEncrypt string, isEncrypt bool) map[string]string {
	tmp1 := map[string]string{}
	switch inter.(type) {
	case map[interface{}]interface{}:
		for k, v := range inter.(map[interface{}]interface{}) {
			switch v.(type) {
			case string:
				tmp1[k.(string)] = strings.Replace(v.(string), " ", "", -1)
			case int:
				tmp := v.(int)
				tmp1[k.(string)] = strconv.Itoa(tmp)
			}
		}
	case map[string]string:
		tmp1 = inter.(map[string]string)
	}
	for k, v := range tmp1 {
		if k == "password" {
			tmp1["password"] = passBase(autoEncrypt, v, isEncrypt)
			continue
		}
	}

	return tmp1
}

// init
func Init(autoEncrypt, filename, filepath string) ModelConfig {
	var config ModelConfig
	viper.SetConfigType("yaml")
	viper.SetConfigName(filename)
	viper.AddConfigPath(filepath)
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("read config failed: %v", err)
	}
	if viper.Get("storewith") == nil || viper.Get("compresstype") == nil || viper.Get("backupnum") == nil {
		logger.Error("config file: storewith|compresstype|backupnum error")
		os.Exit(1)
	}
	config.IsEncrypt = viper.Get("isencrypt").(bool)
	switch viper.Get("storewith").(type) {
	case map[string]string:
		tmp1 := map[string]string{}
		for k, v := range viper.Get("storewith").(map[string]string) {
			if k == "password" {
				tmp1["password"] = passBase(autoEncrypt, v, config.IsEncrypt)
				continue
			}
			tmp1[k] = v
		}
		config.StoreWith = tmp1
	case map[string]interface{}:
		tmp1 := map[string]string{}
		for k, v := range viper.Get("storewith").(map[string]interface{}) {
			if k == "password" {
				tmp1["password"] = passBase(autoEncrypt, v.(string), config.IsEncrypt)
				continue
			}
			switch v.(type) {
			case int:
				tmp1[k] = strconv.Itoa(v.(int))
				continue
			default:
				tmp1[k] = v.(string)
			}

		}
		config.StoreWith = tmp1
	default:
		logger.Error("config file: StoreWith error")
		os.Exit(1)
	}
	viper.Set("storewith", config.StoreWith)
	config.CompressType = viper.Get("compresstype").(string)
	config.BackupNum = viper.Get("backupnum").(int)

	switch viper.Get("databases").(type) {
	case []interface{}:
		databasesList := viper.Get("databases").([]interface{})
		for _, database := range databasesList {
			config.Databases = append(config.Databases, interfaceTomap(database, autoEncrypt, config.IsEncrypt))
		}
	case []map[string]string:
		databasesList := viper.Get("databases").([]map[string]string)
		tmp1 := map[string]string{}
		for _, database := range databasesList {
			for k, v := range database {
				if k == "password" {
					tmp1["password"] = passBase(autoEncrypt, v, config.IsEncrypt)
					continue
				}
			}
		}
		config.Databases = append(config.Databases, tmp1)
	default:
		logger.Error("config file: databases error")
		os.Exit(1)
	}

	config.IsEncrypt = true
	viper.Set("databases", config.Databases)
	viper.Set("isencrypt", config.IsEncrypt)
	viper.WriteConfigAs(filepath + "/" + filename + ".yml")
	return config
}
