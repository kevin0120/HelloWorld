// This file is auto-generated, don't edit it. Thanks.
package main

import (
	"fmt"
	openapi "github.com/alibabacloud-go/darabonba-openapi/client"
	dysmsapi20170525 "github.com/alibabacloud-go/dysmsapi-20170525/v2/client"
	console "github.com/alibabacloud-go/tea-console/client"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"os"
)

func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		value = fallback
	}
	return value
}

var AccessKeyID = GetEnv("AccessKeyID", "LTAI5tMAhS8PwCyN9CPHaxdY")
var AccessKeySecret = GetEnv("AccessKeySecret", "X4jz25XAgNwcBBC9OaRq2MnO8Ojxvm")
var PhoneNumbers = GetEnv("PhoneNumbers", "15221753062")
var SignName = GetEnv("SignName", "阿里云短信测试")
var TemplateCode = GetEnv("TemplateCode", "SMS_154950909")

var Template = GetEnv("Template", "12345")
var TemplateParam = GetEnv("TemplateParam", fmt.Sprintf(`{"code":%s}`, Template))

/**
 * 使用AK&SK初始化账号Client
 * @param accessKeyId
 * @param accessKeySecret
 * @return Client
 * @throws Exception
 */
func CreateClient(accessKeyId *string, accessKeySecret *string) (_result *dysmsapi20170525.Client, _err error) {
	config := &openapi.Config{
		// 您的AccessKey ID
		AccessKeyId: accessKeyId,
		// 您的AccessKey Secret
		AccessKeySecret: accessKeySecret,
	}
	// 访问的域名
	config.Endpoint = tea.String("dysmsapi.aliyuncs.com")
	_result = &dysmsapi20170525.Client{}
	_result, _err = dysmsapi20170525.NewClient(config)
	return _result, _err
}

func _main() (_err error) {
	client, _err := CreateClient(&AccessKeyID, &AccessKeySecret)
	if _err != nil {
		return _err
	}

	sendSmsRequest := &dysmsapi20170525.SendSmsRequest{
		PhoneNumbers:  &PhoneNumbers,
		SignName:      &SignName,
		TemplateCode:  &TemplateCode,
		TemplateParam: &TemplateParam,
	}
	resp, _err := client.SendSms(sendSmsRequest)
	if _err != nil {
		return _err
	}

	console.Log(util.ToJSONString(tea.ToMap(resp)))
	return _err
}

func main() {
	err := _main()
	if err != nil {
		panic(err)
	}
}
