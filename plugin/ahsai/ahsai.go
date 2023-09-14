// Package ahsai AH Soft フリーテキスト音声合成 demo API
package ahsai

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/storeHWID", storeHWIDHandler)
	http.ListenAndServe(":8080", nil)
}

func storeHWIDHandler(w http.ResponseWriter, r *http.Request) {
	// 解析请求参数
	content := r.URL.Query().Get("content")

	// 检查是否收到特定命令消息
	if strings.HasPrefix(content, "#验证") {
		// 获取发送者的HWID
		hwid := strings.TrimSpace(strings.TrimPrefix(content, "#验证"))

		// 存储HWID到Gitee库的HWID.txt文件中
		err := storeHWID(hwid)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "#验证失败")
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "#验证成功")
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Invalid request.")
	}
}

func storeHWID(hwid string) error {
	// Gitee库的API URL
	apiURL := "https://gitee.com/api/v5/repos/wanf_1_0/breeze/contents/HWID.txt"

	// Gitee库的信息
	owner := "wanf_1_0"  // 替换为你的Gitee用户名
	repo := "breeze"      // 替换为你的Gitee仓库名
	path := "HWID.txt"

	// 构建请求URL
	url := strings.ReplaceAll(apiURL, "{owner}", owner)
	url = strings.ReplaceAll(url, "{repo}", repo)
	url = strings.ReplaceAll(url, "{path}", path)

	// 读取现有的HWID.txt文件内容
	existingContent, err := getExistingContent(url)
	if err != nil {
		return err
	}

	// 将新的HWID添加到现有内容中，每行一个HWID
	newContent := existingContent + hwid + "\n"

	// 构建请求体
	body := fmt.Sprintf(`{
		"message": "Update HWID",
		"content": "%s"
	}`, base64.StdEncoding.EncodeToString([]byte(newContent)))

	// 发送请求
	resp, err := http.Put(url, "application/json", strings.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		respBody, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("Failed to store HWID: %s", respBody)
	}

	return nil
}

func getExistingContent(url string) (string, error) {
	// 发送请求获取现有的HWID.txt文件内容
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 解析响应JSON
	var content struct {
		Content string `json:"content"`
	}
	err = json.Unmarshal(respBody, &content)
	if err != nil {
		return "", err
	}

	// 解码Base64内容
	decodedContent, err := base64.StdEncoding.DecodeString(content.Content)
	if err != nil {
		return "", err
	}

	return string(decodedContent), nil
}
