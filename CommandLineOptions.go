package main

import (
	"flag"
)

var ClientID = flag.String("id", "00000000-0000-0000-0000-000000000001", "client id for connecting the extension to the server")
var ClientSecret = flag.String("secret", ".", "client secret for connecting the extension to the server")
var User = flag.String("usr", "eruci18@gmail.com", "username for connecting the extension to the server")
var Password = flag.String("pwd", "e12038076", "password for connecting the extension to the server")
var Addr = flag.String("addr", "acc.server.simplified.engineering", "https server service address")
var Port = flag.Int("port", 443, "https server service port")
