// Copyright (c) 2018-2019 The CYBAVO developers
// All Rights Reserved.
// NOTICE: All information contained herein is, and remains
// the property of CYBAVO and its suppliers,
// if any. The intellectual and technical concepts contained
// herein are proprietary to CYBAVO
// Dissemination of this information or reproduction of this materia
// is strictly forbidden unless prior written permission is obtained
// from CYBAVO.

package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
	_ "github.com/cybavo/VAULTX_MOCK_SERVER/preinit"
	_ "github.com/cybavo/VAULTX_MOCK_SERVER/routers"
)

func main() {
	enableCORS()
	beego.Run()
}

func enableCORS() {
	v := beego.AppConfig.DefaultBool("enable_cors", true)
	if v {
		beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
			AllowAllOrigins: true,
			AllowMethods:    []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowHeaders:    []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Content-Type"},
			ExposeHeaders:   []string{"Content-Length", "Access-Control-Allow-Origin"},
		}))
	}
}
