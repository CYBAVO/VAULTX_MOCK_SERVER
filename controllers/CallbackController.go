// Copyright (c) 2018-2019 The CYBAVO developers
// All Rights Reserved.
// NOTICE: All information contained herein is, and remains
// the property of CYBAVO and its suppliers,
// if any. The intellectual and technical concepts contained
// herein are proprietary to CYBAVO
// Dissemination of this information or reproduction of this materia
// is strictly forbidden unless prior written permission is obtained
// from CYBAVO.

package controllers

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/cybavo/VAULTX_MOCK_SERVER/api"
)

type CallbackController struct {
	beego.Controller
}

func (c *CallbackController) AbortWithError(status int, err error) {
	resp := api.ErrorCodeResponse{
		ErrMsg:  err.Error(),
		ErrCode: status,
	}
	c.Data["json"] = resp
	c.Abort(strconv.Itoa(status))
}

func calcSHA256(data []byte) (calculatedHash []byte, err error) {
	sha := sha256.New()
	_, err = sha.Write(data)
	if err != nil {
		return
	}
	calculatedHash = sha.Sum(nil)
	return
}

func calcChecksum(data []byte, secret string) string {
	payload := string(data) + secret
	sha, _ := calcSHA256([]byte(payload))
	checksum := base64.URLEncoding.EncodeToString(sha)
	return checksum
}

// @router /callback [post]
func (c *CallbackController) Callback() {
	var request api.CallbackStruct
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	if err != nil {
		logs.Error("callback unmarshal error =>", err)
		c.AbortWithError(http.StatusBadRequest, err)
	}

	//
	// get API secret by company ID `request.CompanyID` to replace `api.APISecret` below
	//

	checksum := c.Ctx.Input.Header("X-CHECKSUM")
	checksumVerf := calcChecksum(c.Ctx.Input.RequestBody, api.APISecret)
	if checksum != checksumVerf {
		c.AbortWithError(http.StatusBadRequest, errors.New("Bad checksum"))
	}

	logs.Debug("Callback => %s\n%#v", c.Ctx.Input.RequestBody, request)

	c.Ctx.WriteString("OK")
}

// @router /getnonce [post]
func (c *OuterController) GetNonce() {
	defer c.ServeJSON()

	var request api.GetNonceRequest
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
	}

	//
	// get API secret by company ID `request.CompanyID` to replace `api.APISecret` below
	//

	checksum := c.Ctx.Input.Header("X-CHECKSUM")
	checksumVerf := calcChecksum(c.Ctx.Input.RequestBody, api.APISecret)
	if checksum != checksumVerf {
		c.AbortWithError(http.StatusBadRequest, errors.New("Bad checksum"))
	}

	//
	// should get nonce of `request.Address` from fullnode
	//
	var nonce int64 = 1
	resp := &api.GetNonceResponse{
		Nonce: nonce,
	}
	respb, _ := json.Marshal(resp)
	c.Data["json"] = resp
	c.Ctx.Output.Header("X-CHECKSUM", calcChecksum(respb, api.APISecret))
}
