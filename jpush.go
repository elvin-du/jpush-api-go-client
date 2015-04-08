package jpush

import (
	"fmt"

	"github.com/elvin-du/jpush-api-go-client/common"
	"github.com/elvin-du/jpush-api-go-client/device"
	"github.com/elvin-du/jpush-api-go-client/httplib"
	"github.com/elvin-du/jpush-api-go-client/push"
)

// JPush 的 Golang 推送客户端
// 详情： http://docs.jpush.io/server/rest_api_v3_push/
type JPushClient struct {
	appKey       string
	masterSecret string
	headers      map[string]string
	http         *httplib.HTTPClient
}

func NewJPushClient(appKey string, masterSecret string) *JPushClient {
	client := JPushClient{
		appKey:       appKey,
		masterSecret: masterSecret,
	}
	headers := make(map[string]string)
	headers["User-Agent"] = "jpush-api-go-client"
	headers["Connection"] = "keep-alive"
	headers["Authorization"] = "Basic " + common.BasicAuth(appKey, masterSecret)
	client.headers = headers

	client.http = httplib.NewClient()

	return &client
}

// 设置调试模式，调试模式下，会输出日志
func (jpc *JPushClient) SetDebug(debug bool) {
	jpc.http.Debug(debug)
}

// 推送 API
func (jpc *JPushClient) Push(payload *push.PushObject) (*push.PushResult, error) {
	return jpc.doPush(common.PUSH_URL, payload)
}

// 推送校验 API， 只用于验证推送调用是否能够成功，与推送 API 的区别在于：不向用户发送任何消息。
func (jpc *JPushClient) PushValidate(payload *push.PushObject) (*push.PushResult, error) {
	return jpc.doPush(common.PUSH_VALIDATE_URL, payload)
}

func (jpc *JPushClient) doPush(url string, payload *push.PushObject) (*push.PushResult, error) {
	resp, err := jpc.http.PostJson(url, payload, jpc.headers)
	if err != nil {
		return nil, err
	}

	result := &push.PushResult{}
	err = result.FromResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 查询设备(设备的别名与标签)
func (jpc *JPushClient) QueryDevice(registrationId string) (*device.DeviceInfoResult, error) {
	url := fmt.Sprintf(common.DEVICE_URL, registrationId)
	resp, err := jpc.http.Get(url, nil, jpc.headers)
	if err != nil {
		return nil, err
	}

	result := &device.DeviceInfoResult{}
	err = result.FromResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 更新设备 （设置的别名与标签）
func (jpc *JPushClient) UpdateDevice(registrationId string, payload *device.DeviceUpdate) (*common.ResponseBase, error) {
	url := fmt.Sprintf(common.DEVICE_URL, registrationId)
	resp, err := jpc.http.PostJson(url, payload, jpc.headers)
	if err != nil {
		return nil, err
	}

	result := common.NewResponseBase(resp)
	return &result, nil
}
