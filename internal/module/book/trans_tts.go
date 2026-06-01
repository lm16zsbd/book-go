package book

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	u "serica-go/internal/module/user"
	"serica-go/internal/pkg/httputil"
)

type TransInput struct {
	Data     string `json:"data"`
	Language string `json:"language"`
}

func (h *Handler) Trans(w http.ResponseWriter, r *http.Request) {
	userID := u.GetUserID(r)

	var input TransInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.RespondJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}
	if input.Data == "" {
		httputil.RespondJSON(w, 200, "")
		return
	}

	_ = userID

	langMap := map[string]string{
		"en":    "English",
		"英文":    "English",
		"zh-hk": "Traditional Chinese (zh-HK)",
		"繁体中文":  "Traditional Chinese (zh-HK)",
		"zh-tw": "Traditional Chinese (zh-HK)",
		"zh-cn": "Simplified Chinese (zh-CN)",
		"简体中文":  "Simplified Chinese (zh-CN)",
	}

	target := ""
	if input.Language != "" {
		target = langMap[input.Language]
		if target == "" {
			target = input.Language
		}
	}

	instructions := ""
	if target != "" {
		instructions = fmt.Sprintf("You are a professional translator. Translate the input into %s. Reply ONLY with the translated text. Do not include any explanations, notes, quotes, or additional content.", target)
	} else {
		instructions = "You are a professional translator. If the input contains any Chinese characters (simplified or traditional), translate it to English. If the input is in English, translate it to Traditional Chinese (zh-HK). Reply ONLY with the translated text. Do not include any explanations, notes, quotes, or additional content."
	}

	payload := map[string]interface{}{
		"model":        "gpt-4o",
		"input":        input.Data,
		"instructions": instructions,
		"temperature":  0.3,
	}
	bodyBytes, _ := json.Marshal(payload)

	apiKey := h.getOpenAIKey()
	if apiKey == "" {
		httputil.RespondJSON(w, 500, map[string]string{"error": "OPENAI_API_KEY not configured"})
		return
	}

	req, _ := http.NewRequestWithContext(r.Context(), "POST", "https://api.openai.com/v1/responses", bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		httputil.RespondJSON(w, 502, map[string]string{"error": "Translation service error"})
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		httputil.RespondJSON(w, resp.StatusCode, map[string]string{"error": "Translation service error"})
		return
	}

	var result map[string]interface{}
	json.Unmarshal(respBody, &result)

	output, _ := result["output"].([]interface{})
	if len(output) > 0 {
		first, _ := output[0].(map[string]interface{})
		content, _ := first["content"].([]interface{})
		for _, c := range content {
			cMap, _ := c.(map[string]interface{})
			if cMap["type"] == "output_text" {
				text, _ := cMap["text"].(string)
				httputil.RespondJSON(w, 200, text)
				return
			}
		}
	}

	httputil.RespondJSON(w, 200, "")
}

func (h *Handler) getOpenAIKey() string {
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		return key
	}
	if h.cfg != nil && h.cfg.OpenAI.APIKey != "" {
		return h.cfg.OpenAI.APIKey
	}
	return ""
}

type TTSInput struct {
	Text         string `json:"text"`
	LanguageCode string `json:"languageCode"`
	SSMLGender   string `json:"ssmlGender"`
}

func (h *Handler) TextToSpeech(w http.ResponseWriter, r *http.Request) {
	var input TTSInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.RespondJSON(w, 400, map[string]string{"error": "invalid body"})
		return
	}

	if input.LanguageCode == "" {
		input.LanguageCode = "cmn-CN"
	}
	if input.SSMLGender == "" {
		input.SSMLGender = "FEMALE"
	}

	voiceMap := map[string]string{
		"cmn-CN": "Zhiyu",
		"yue-HK": "Hiujin",
		"en-GB":  "Emma",
		"en-US":  "Joanna",
	}
	if input.SSMLGender == "MALE" && input.LanguageCode == "en-US" {
		voiceMap["en-US"] = "Matthew"
	}

	voiceID := voiceMap[input.LanguageCode]
	if voiceID == "" {
		voiceID = "Zhiyu"
	}

	ssml := fmt.Sprintf(`<speak>%s</speak>`, escapeXML(input.Text))

	payload := map[string]interface{}{
		"OutputFormat": "mp3",
		"Text":         ssml,
		"TextType":     "ssml",
		"VoiceId":      voiceID,
		"Engine":       "neural",
		"LanguageCode": normalizePollyLang(input.LanguageCode),
	}

	_ = payload

	httputil.RespondJSON(w, 500, map[string]string{"error": "TTS requires AWS Polly SDK integration"})
}

func normalizePollyLang(lang string) string {
	if lang == "yue-HK" {
		return "yue-CN"
	}
	return lang
}

func escapeXML(s string) string {
	s = replaceAll(s, "&", "&amp;")
	s = replaceAll(s, "<", "&lt;")
	s = replaceAll(s, ">", "&gt;")
	s = replaceAll(s, "\"", "&quot;")
	s = replaceAll(s, "'", "&apos;")
	return s
}

func replaceAll(s, old, new string) string {
	result := ""
	for i := 0; i < len(s); i++ {
		if i+len(old) <= len(s) && s[i:i+len(old)] == old {
			result += new
			i += len(old) - 1
		} else {
			result += string(s[i])
		}
	}
	return result
}
