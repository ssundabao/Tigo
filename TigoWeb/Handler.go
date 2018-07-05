// Copyright 2018 The Tigo Authors. All rights reserved.
package WebFramework

import (
	"net/http"
	"encoding/json"
	"fmt"
	"strings"
)

// Handler的基础类，开发者开发的handler继承此类
type BaseHandler struct {
	ResponseWriter  http.ResponseWriter
	Request        *http.Request
}

// 初始化Handler的方法
func (baseHandler *BaseHandler)InitHandler(responseWriter http.ResponseWriter, request *http.Request) {
	baseHandler.Request = request
	baseHandler.ResponseWriter = responseWriter
	baseHandler.Request.ParseForm()
}

// 将对象转化为Json字符串，转换失败则返回空字符串。
// 传入参数Response为一个interface，必须有成员函数Print。
func (baseHandler *BaseHandler)ToJson(response Response) (result string) {
	// 将该对象转换为byte字节数组
	jsonResult, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		return ""
	}
	// 将byte数组转换为string
	return string(jsonResult)
}

// 向客户端响应一个Json结果
func (baseHandler *BaseHandler)ResponseAsJson(response Response) {
	// 将对象转换为Json字符串
	jsonResult := baseHandler.ToJson(response)
	// 设置http报文头内的Content-Type
	baseHandler.ResponseWriter.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(baseHandler.ResponseWriter, jsonResult)
}

// 向客户端响应一个Text结果
func (baseHandler *BaseHandler)ResponseAsText(result string) {
	fmt.Fprintf(baseHandler.ResponseWriter, result)
}

// 向客户端响应一个html结果
func (baseHandler *BaseHandler)ResponseAsHtml(result string) {
	baseHandler.ResponseWriter.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(baseHandler.ResponseWriter, result)
}

// 检查请求是否被允许
func (baseHandler *BaseHandler)CheckRequestMethodValid(methods ...string) (result bool) {
	// 获取请求方式
	requestMethod := baseHandler.Request.Method
	// 遍历被允许的请求方式，判断是否合法
	for _, value := range methods {
		if requestMethod == value || strings.ToLower(requestMethod) == value {
			return true
		}
	}
	// 如果不合法返回405
	baseHandler.ResponseWriter.WriteHeader(405)
	return false
}

// 设置cookie
func (baseHandler *BaseHandler)SetCookie(name string, value string) {
	cookie := http.Cookie{Name:name, Value:value}
	http.SetCookie(baseHandler.ResponseWriter, &cookie)
}

// 设置高级cookie选项
func (baseHandler *BaseHandler)SetCookieObject(cookie Cookie) {
	responseCookie := cookie.ToHttpCookie()
	http.SetCookie(baseHandler.ResponseWriter, &responseCookie)
}

// 设置加密cookie
func (baseHandler *BaseHandler)SetSecureCookie(name string, value string, key string) {
	cookie := Cookie{
		Name:        name,
		Value:       value,
		IsSecurity:  true,
		SecurityKey: key,
	}
	baseHandler.SetCookieObject(cookie)
}

// 获取cookie值，如果获取失败则返回空字符串
func (baseHandler *BaseHandler)GetCookie(name string) (value string) {
	cookie, err := baseHandler.Request.Cookie(name)
	if err != nil {
		return ""
	}
	value = cookie.Value
	return value
}

// 获取加密cookie值，如果获取失败则返回空
func (baseHandler *BaseHandler)GetSecureCookie(name string, key string) (value string) {
	httpCookie, err := baseHandler.Request.Cookie(name)
	if err != nil {
		return ""
	}
	cookie := Cookie{}
	cookie.ConvertFromHttpCookie(*httpCookie)
	cookie.IsSecurity = true
	cookie.SecurityKey = key
	value = cookie.GetCookieDecodeValue()
	return value
}

// 获取cookie对象，多参数输入，参数如下：
//   - 无参数：默认cookieName为空字符串
//   - 一个参数：传入的参数为cookieName
//   - 多个参数：传入的第一个参数为cookieName，第二个参数为加密/解密cookie所用的Key，此时认为cookie是需要进行加密/解密处理的
func (baseHandler *BaseHandler)GetCookieObject(name ...string) (Cookie, error) {
	cookie := Cookie{}
	var cookieName, key string
	length := len(name)
	switch {
	case length < 1:
		cookieName = ""
	case length == 1:
		cookieName = name[0]
	case length > 1:
		cookieName = name[0]
		key = name[1]
	}
	httpCookie, err := baseHandler.Request.Cookie(cookieName)
	if err != nil {
		return cookie, nil
	}
	cookie.ConvertFromHttpCookie(*httpCookie)
	if len(key) > 0 {
		cookie.SetSecurityKey(key)
	}
	return cookie, nil
}

// 获取header

// 设置header
