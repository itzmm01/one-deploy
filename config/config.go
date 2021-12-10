package config

import (
	"one-backup/keygen"
	"strconv"
	"strings"

	"github.com/spf13/viper"
	"github.com/wonderivan/logger"
)

// ModelConfig for special case
type ModelConfig struct {
	StoreWith    map[string]interface{}
	CompressType string
	BackupNum    int
	Databases    []map[string]string
	Viper        *viper.Viper
}

// SubConfig sub config info
type SubConfig struct {
	Name  string
	Type  string
	Viper *viper.Viper
}

func passBase(autoEncrypt, pass string) string {
	if autoEncrypt == "no" {
		return pass
	}

	passStr := keygen.AesDecryptCBC(strings.Replace(pass, " ", "", -1), "pass")
	if passStr == "base64 error" {
		return keygen.AesEncryptCBC(strings.Replace(pass, " ", "", -1), "pass")
	} else {
		return pass
	}

}
func Init(autoEncrypt, filename, filepath string) ModelConfig {
	var config ModelConfig
	viper.SetConfigType("yaml")
	viper.SetConfigName(filename)
	viper.AddConfigPath(filepath)
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("read config failed: %v", err)
	}

	config.StoreWith = viper.Get("storewith").(map[string]interface{})
	config.CompressType = viper.Get("compresstype").(string)
	config.BackupNum = viper.Get("backupnum").(int)

	switch viper.Get("databases").(type) {
	case []interface{}:
		databasesList := viper.Get("databases").([]interface{})
		for _, database := range databasesList {
			tmp1 := map[string]string{}
			switch database.(type) {
			case map[interface{}]interface{}:
				for k, v := range database.(map[interface{}]interface{}) {
					switch v.(type) {
					case string:
						tmp1[k.(string)] = strings.Replace(v.(string), " ", "", -1)
					case int:
						tmp := v.(int)
						tmp1[k.(string)] = strconv.Itoa(tmp)
					}
				}
			case map[string]string:
				tmp1 = database.(map[string]string)
			}
			for k, v := range tmp1 {
				if k == "password" {
					tmp1["password"] = passBase(autoEncrypt, v)
					continue
				}
			}
			config.Databases = append(config.Databases, tmp1)
		}
	case []map[string]string:
		databasesList := viper.Get("databases").([]map[string]string)
		tmp1 := map[string]string{}
		for _, database := range databasesList {
			for k, v := range database {
				if k == "password" {
					tmp1["password"] = passBase(autoEncrypt, v)
					continue
				}
			}
		}
		config.Databases = append(config.Databases, tmp1)
	}
	viper.Set("databases", config.Databases)
	viper.WriteConfigAs(filepath + "/" + filename + ".yml")
	return config
}
