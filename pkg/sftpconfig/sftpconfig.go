package sftpconfig

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/ghodss/yaml"
	"github.com/navinds25/styx/pkg/encryption"
)

// RawEncryptionKey is for the AES encryption key, it is base64Encoded.
var RawEncryptionKey = "SlNqVVhZalpzWWg5ZkRWOENvWHdwNzhITnl3RnpnWnFE"

// Config is the Interface for methods on TransferConfig struct
type Config interface {
	EncryptSecureFields(key string) error
	DecryptSecureFields(key string) error
	EncodeGob() error
	DecodeGob() error
}

// TransferConfig is the struct for parsing the configs
type TransferConfig struct {
	TransferID  string `json:"transfer_id"`
	Description string `json:"description"`
	Source      struct {
		Local struct {
			LocalFile string `json:"local_file"`
			LocalPath string `json:"local_path"`
		} `json:"local"`
		Remote struct {
			RemotePath string `json:"remote_path"`
			Host       string `json:"host"`
			Port       int    `json:"port"`
			Auth       struct {
				Username string `json:"username"`
				Password string `json:"password"`
				Key      string `json:"key"`
			} `json:"auth"`
		} `json:"remote"`
	} `json:"source"`
	Destination struct {
		Local struct {
			DestFile string `json:"dest_file"`
			DestPath string `json:"dest_path"`
		} `json:"local"`
		Remote struct {
			RemotePath string `json:"remote_path"`
			Host       string `json:"host"`
			Port       int    `json:"port"`
			Auth       struct {
				Username string `json:"username"`
				Password string `json:"password"`
				Key      string `json:"key"`
			} `json:"auth"`
		} `json:"remote"`
	} `json:"destination"`
}

func aesString(value, op string) (string, error) {
	key, err := base64.StdEncoding.DecodeString(RawEncryptionKey)
	if err != nil {
		return "", err
	}
	if op == "encrypt" {
		ciphertext, err := encryption.AESEncryptCBC(key, []byte(value))
		if err != nil {
			return "", err
		}
		return string(ciphertext), nil
	} else if op == "decrypt" {
		ciphertext, err := encryption.AESDecryptCBC(key, []byte(value))
		if err != nil {
			return "", err
		}
		return string(ciphertext), nil
	} else {
		return "", errors.New("aesString, unknown op")
	}
}

// EncryptSecureFields encrypts the fields of struct instance.
func (config *TransferConfig) EncryptSecureFields() error {
	if err := secureFields(config, "encrypt"); err != nil {
		return err
	}
	return nil
}

// DecryptSecureFields decrypts the fields of the struct instance.
func (config *TransferConfig) DecryptSecureFields() error {
	if err := secureFields(config, "decrypt"); err != nil {
		return err
	}
	return nil
}

func secureFields(config *TransferConfig, op string) error {
	sourceAuthFields := reflect.TypeOf(config.Source.Remote.Auth)
	sourceAuthValues := reflect.ValueOf(config.Source.Remote.Auth)
	numOfSourceAuthFields := sourceAuthFields.NumField()
	for i := 0; i < numOfSourceAuthFields; i++ {
		value := sourceAuthValues.Field(i)
		valueStr := value.String()
		if valueStr != "" {
			if op == "encrypt" {
				encryptedValue, err := aesString(valueStr, op)
				if err != nil {
					return err
				}
				value.SetString(encryptedValue)
			} else if op == "decrypt" {
				decryptedValue, err := aesString(valueStr, op)
				if err != nil {
					return err
				}
				value.SetString(decryptedValue)
			} else {
				return errors.New("secureFields: Unknown op")
			}
		}
	}
	destAuthFields := reflect.TypeOf(config.Destination.Remote.Auth)
	destAuthValues := reflect.ValueOf(config.Destination.Remote.Auth)
	numOfDestAuthFields := destAuthFields.NumField()
	for i := 0; i < numOfDestAuthFields; i++ {
		value := destAuthValues.Field(i)
		valueStr := value.String()
		if valueStr != "" {
			if op == "encrypt" {
				encryptedValue, err := aesString(valueStr, op)
				if err != nil {
					return err
				}
				value.SetString(encryptedValue)
			} else if op == "decrypt" {
				decryptedValue, err := aesString(valueStr, op)
				if err != nil {
					return err
				}
				value.SetString(decryptedValue)
			} else {
				return errors.New("secureFields: Unknown op")
			}
		}
	}
	return nil
}

// EncodeGob Returns Gob encoded byte array of struct
func (config *TransferConfig) EncodeGob() ([]byte, error) {
	value := bytes.Buffer{}
	if err := gob.NewEncoder(&value).Encode(config); err != nil {
		return nil, err
	}
	return value.Bytes(), nil
}

// DecodeGob takes a encoded Gob as a byte array and updates the
// struct instance with the decoded values
func (config *TransferConfig) DecodeGob(value []byte) error {
	valReader := bytes.NewReader(value)
	if err := gob.NewDecoder(valReader).Decode(config); err != nil {
		return err
	}
	return nil
}

// GetConfig returns the parsed config
func GetConfig(confFile string) (map[string][]TransferConfig, error) {
	var config map[string][]TransferConfig
	_, err := os.Stat(confFile)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(confFile)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	jsonData, err := yaml.YAMLToJSON(data)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(jsonData, &config); err != nil {
		return nil, err
	}
	return config, nil
}
