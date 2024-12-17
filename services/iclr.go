package services

import (
	"crawler/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type OpenReviewNote struct {
	Id      string `json:"id"`
	Content struct {
		Title string `json:"title"`
		Pdf   string `json:"pdf"`
	} `json:"content"`
}

type OpenReviewResponse struct {
	Notes []OpenReviewNote `json:"notes"`
}

func ICLR(year int) ([]models.Paper, error) {
	// 构建API URL
	baseURL := "https://api.openreview.net/notes"
	params := url.Values{}
	params.Add("invitation", fmt.Sprintf("ICLR.cc/%d/Conference/-/Blind_Submission", year))
	params.Add("details", "replyCount")
	params.Add("limit", "1000")

	// 发送请求
	apiURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())
	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("API请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}

	// 解析JSON
	var result OpenReviewResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("JSON解析失败: %v", err)
	}

	// 转换为Paper模型
	papers := make([]models.Paper, 0)
	for _, note := range result.Notes {
		pdf := "https://openreview.net" + note.Content.Pdf
		paper := models.Paper{
			PaperName: note.Content.Title,
			URL:       pdf,
		}
		papers = append(papers, paper)
	}

	if len(papers) == 0 {
		return nil, fmt.Errorf("未找到论文")
	}

	return papers, nil
}
