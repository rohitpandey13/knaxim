package config

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/go-test/deep"
)

var testconfig = Configuration{
	Address:    "test.test",
	StaticPath: "/static/path",
	IndexPath:  "/index/path",
	Server: &http.Server{
		Addr:              "0.0.0.0:8765",
		ReadHeaderTimeout: 123456789,
	},
	Cert: &Ssl{
		CertFile: "/cert/file",
		KeyFile:  "/key/file",
		HTTPport: "5678",
	},
	GracefulTimeout: Duration{time.Second * 5},
	BasicTimeout:    Duration{time.Minute},
	FileTimeoutRate: 5000000,
	MaxFileTimeout:  Duration{time.Second * 7},
	MinFileTimeout:  Duration{time.Second * 4},
	DatabaseType:    "memory",
	Database:        json.RawMessage(`{}`),
	DatabaseReset:   true,
	Tika: Tika{
		Type: "external",
	},
	GotenPath: "gotenburg.tst",
	FileLimit: 12341,
	FreeSpace: 54321,
	AdminKey:  "helloboys",
	GuestUser: &Guest{
		Name:  "guestusername",
		Pass:  "guestpassword",
		Email: "guest@guest.guest",
	},
	SetupTimeout: Duration{time.Minute * 9},
	UserTimeouts: struct {
		Inactivity Duration
		Total      Duration
	}{
		Inactivity: Duration{time.Hour * 4},
		Total:      Duration{time.Hour * 6},
	},
	Email: SMTP{
		From:   "test@test.test",
		Server: "emailServer",
	},
}

var testconfigjson = []byte(`{
  "Address": "test.test",
	"static": "/static/path",
	"index": "/index/path",
	"server": {
		"Addr": "0.0.0.0:8765",
		"ReadHeaderTimeout": 123456789
	},
	"cert": {
		"cert": "/cert/file",
		"key":  "/key/file",
		"http_port": "5678"
	},
	"close_time": "5s",
	"basic_timeout": "1m",
	"file_timeout_rate": 5000000,
	"max_file_timeout":  "7s",
	"min_file_timeout":  "4s",
	"db_type": "memory",
  "db": {},
	"db_clear": true,
	"tika": {
		"type": "external"
	},
	"gotenpath": "gotenburg.tst",
	"filelimit": 12341,
	"total_free_space": 54321,
	"AdminKey":  "helloboys",
	"GuestUser": {
		"Name":  "guestusername",
		"Pass":  "guestpassword",
		"Email": "guest@guest.guest"
	},
	"SetupTimeout": "9m",
	"UserTimeouts": {
		"Inactivity": "4h",
		"Total": "6h"
	},
	"Email": {
		"from":   "test@test.test",
		"server": "emailServer"
	}
}`)

func TestConfigDecode(t *testing.T) {
	var decodedConfig = &Configuration{}
	if err := json.Unmarshal(testconfigjson, decodedConfig); err != nil {
		t.Fatal("unable to decode configuration: ", err.Error())
	}
	if diff := deep.Equal(*decodedConfig, testconfig); diff != nil {
		t.Log(diff)
		t.Fatal("decoded Configuration does not equal config")
	}
}