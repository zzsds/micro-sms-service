package conf

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"flag"

	"github.com/BurntSushi/toml"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var (
	Conf       Config
	once       sync.Once
	ConfigPath string

	// ******* acm config ***********************
	CFG_ENDPOINT    = GetEnv("CFG_ENDPOINT", "acm.aliyun.com")
	CFG_NAMESPACEID = GetEnv("CFG_NAMESPACEID", "a23c93cd-491c-44dd-be30-fb1df6e6ddaf")
	CFG_ACCESSKEY   = GetEnv("CFG_ACCESSKEY", "")
	CFG_SECRETKEY   = GetEnv("CFG_SECRETKEY", "")
	CFG_CLUSTER     = GetEnv("CFG_CLUSTER", "test")
)

const (
	DataIdDefault = "srv.sms"
)

type Config struct {
	Db  *Db `json:"mysql"`
	Sms struct {
		YunPian *Sms `json:"yunpian"`
	}
}

type Db struct {
	Host     string
	User     string
	Password string
	Name     string
	Charset  string
	Debug    bool
}

type Sms struct {
	Name   string
	APIKey string
	Debug  bool
}

// 初始化配置
func InitConfig() {
	flag.StringVar(&ConfigPath, "c", GetFileConfFile(), "this default local conf.toml")
	flag.Parse()

	if isExists(ConfigPath) {
		local()
	} else {
		load()
	}
}

func CheckErr(err error) {
	if err != nil {
		panic(err)
	}
}

// 加载本地配置
func local() {
	once.Do(func() {
		_, err := toml.DecodeFile(ConfigPath, &Conf)
		CheckErr(err)
	})
}

// 加载线上配置
func load() {
	clientConfig := constant.ClientConfig{
		//
		Endpoint:       CFG_ENDPOINT + ":8080",
		NamespaceId:    CFG_NAMESPACEID,
		AccessKey:      CFG_ACCESSKEY,
		SecretKey:      CFG_SECRETKEY,
		TimeoutMs:      5 * 1000,
		ListenInterval: 30 * 1000,
	}
	configClient, err := clients.CreateConfigClient(map[string]interface{}{
		"clientConfig": clientConfig,
	})

	CheckErr(err)

	content, err := configClient.GetConfig(vo.ConfigParam{
		DataId: DataIdDefault,
		Group:  CFG_CLUSTER})
	CheckErr(err)
	CheckErr(json.Unmarshal([]byte(content), &Conf))
}

func GetEnv(key, value string) string {
	newValue := os.Getenv(key)
	if newValue == "" {
		return value
	}
	return newValue
}

func isExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil && !os.IsNotExist(err)
}

func GetFileConfFile() string {
	_, f, _, _ := runtime.Caller(1)
	return filepath.Join(filepath.Dir(f), "/config.toml")
}

func ValidateMobile(mobile string) bool {
	ok, _ := regexp.MatchString(`^((\+[0-9]\d{10,12})|1[1-9]\d{9})$`, mobile)
	return ok
}

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}
