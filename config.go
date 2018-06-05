package kotori

var GlobCfg = Config{}

type Config struct {
	PORT         int64    `toml:"port"`
	ADMIN        []Admin  `toml:"admin"`
	ALLOW_ORIGIN []string `toml:"allow_origin"`
}
