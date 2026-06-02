package book

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/polly"
	pollytypes "github.com/aws/aws-sdk-go-v2/service/polly/types"

	u "serica-go/internal/module/user"
	"serica-go/internal/pkg/httputil"
)

type TransInput struct {
	Data     string `json:"data"`
	Language string `json:"language"`
}

// @Summary      AI翻译
// @Tags         Book
// @Accept       json
// @Produce      json
// @Param        body body TransInput true "翻译参数"
// @Success      200  {string} string
// @Security     BearerAuth
// @Router       /v1/client/trans [post]
func (h *Handler) Trans(w http.ResponseWriter, r *http.Request) (any, error) {
	userID := u.GetUserID(r)

	var input TransInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.RespondJSON(w, 400, map[string]string{"error": "invalid body"})
		return nil, nil
	}
	if input.Data == "" {
		return "", nil
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
		return nil, nil
	}

	req, _ := http.NewRequestWithContext(r.Context(), "POST", "https://api.openai.com/v1/responses", bytes.NewReader(bodyBytes))
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := httputil.NewHTTPClient(120 * time.Second)
	resp, err := client.Do(req)
	if err != nil {
		httputil.RespondJSON(w, 502, map[string]string{"error": "Translation service error"})
		return nil, nil
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		httputil.RespondJSON(w, resp.StatusCode, map[string]string{"error": "Translation service error"})
		return nil, nil
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
				return text, nil
			}
		}
	}

	return "", nil
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

// @Summary      AWS Polly 文字转语音
// @Tags         Book
// @Accept       json
// @Produce      audio/mpeg
// @Param        body body TTSInput true "TTS参数"
// @Success      200  {file}   audio/mpeg
// @Router       /v1/client/text-to-speech [post]
func (h *Handler) TextToSpeech(w http.ResponseWriter, r *http.Request) (any, error) {
	var input TTSInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		httputil.RespondJSON(w, 400, map[string]string{"error": "invalid body"})
		return nil, nil
	}

	if input.Text == "" {
		httputil.RespondJSON(w, 400, map[string]string{"error": "text is required"})
		return nil, nil
	}
	if input.LanguageCode == "" {
		input.LanguageCode = "yue-HK"
	}

	voiceID := "Hiujin"
	engine := pollytypes.EngineNeural

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(h.cfg.AWS.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			h.cfg.AWS.AccessKeyID,
			h.cfg.AWS.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		httputil.RespondJSON(w, 500, map[string]string{"error": "TTS service error"})
		return nil, nil
	}

	pollyClient := polly.NewFromConfig(cfg)
	speechInput := &polly.SynthesizeSpeechInput{
		Text:         aws.String(input.Text),
		OutputFormat: pollytypes.OutputFormatMp3,
		VoiceId:      pollytypes.VoiceId(voiceID),
		Engine:       engine,
	}

	if input.LanguageCode != "" {
		langCode := pollytypes.LanguageCode(input.LanguageCode)
		speechInput.LanguageCode = langCode
		if voiceID == "Hiujin" && input.LanguageCode == "cmn-CN" {
			speechInput.VoiceId = pollytypes.VoiceIdZhiyu
		}
	}

	output, err := pollyClient.SynthesizeSpeech(context.Background(), speechInput)
	if err != nil {
		httputil.RespondJSON(w, 500, map[string]string{"error": "TTS synthesis failed"})
		return nil, nil
	}

	audioData, err := io.ReadAll(output.AudioStream)
	if err != nil {
		httputil.RespondJSON(w, 500, map[string]string{"error": "failed to read audio data"})
		return nil, nil
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Disposition", "inline; filename=speech.mp3")
	w.WriteHeader(http.StatusOK)
	w.Write(audioData)
	return nil, nil
}
