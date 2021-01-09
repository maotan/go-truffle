/**
* @Author: mo tan
* @Description:
* @Date 2021/1/2 14:39
 */
package yaml_config

//----- active
type ActiveConfig struct {
	Active  string     `mapstructure:"active" json:"active" yaml:"active"`
}

type GinConfig struct{
	ActiveConf ActiveConfig `mapstructure:"gin" json:"gin" yaml:"gin"`
}

//------ log config
type LogConfig struct {
	LogPath  string      `mapstructure:"log-path" json:"logPath" yaml:"log-path"`
	RotationCount uint   `mapstructure:"rotation-count" json:"rotationCount" yaml:"rotation-count"`
}
// ----- server config
type ServerConfig struct {
	Port int
	Name string
}

// ----- consul config
type ConsulConfig struct {
	Host string
	Port int
	Token string
}

// ------ database config
type DatabaseConfig struct {
	Driver string		`mapstructure:"driver" json:"driver" yaml:"driver"`
	Url string		`mapstructure:"url" json:"url" yaml:"url"`
	Username string 	`mapstructure:"username" json:"username" yaml:"username"`
	Password string		`mapstructure:"password" json:"password" yaml:"password"`
}

// ----- total config
type YamlConfig struct {
	LogConf LogConfig `mapstructure:"log" json:"log" yaml:"log"`
	ServerConf ServerConfig `mapstructure:"server" json:"server" yaml:"server"`
	ConsulConf ConsulConfig `mapstructure:"consul" json:"consul" yaml:"consul"`
	DatabaseConf DatabaseConfig `mapstructure:"database" json:"database" yaml:"database"`
}

/*func (this *GinConfig) DefaultGinConfig() {
	activeConf := ActiveConfig{Active: ""}
	this.ActiveConf = activeConf
}*/

/*func (this *GinConfig) InitGinConfig(path string) {
	this.DefaultGinConfig()
	file, _ := ioutil.ReadFile(path)
	if err := yaml.Unmarshal(file, this); err != nil {
		panic(err)
	}
}*/