package main

var GlobCfg = Config{}

type Config struct {
	PORT	   int64         `toml:"port"`
	ADMIN	   []Admin	     `toml:"admin"`
}