package plugin

import (
	"encoding/json"
	"fmt"
	"gitee.com/jamespi/lecheng-drone/config"
	"net/http"
)

func GetAccessToken() string {
	url := fmt.Sprintf("%s?appId=%s&appKey=%s", config.DjiSettings["url"], config.DjiSettings["appId"], config.DjiSettings["appKey"])

	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	// 解析响应获取token
	var result struct {
		AccessToken string `json:"accessToken"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	return result.AccessToken
}
