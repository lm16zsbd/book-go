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
	Language string `json:"language" enums:"en,es,ar,fr,pt,ru,id,zh-hk,zh-cn,zh-tw,English,Spanish,Arabic,French,Portuguese,Russian,Indonesian,Chinese"`
}

// @Summary      AI翻译
// @Tags         Book
// @Accept       json
// @Produce      json
// @Param        body body TransInput true "翻译参数(language: English / Spanish / Arabic / French / Portuguese / Russian / Indonesian / 中文+粤语自动互译)"
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
		"en":         "English",
		"English":    "English",
		"zh-hk":      "Traditional Chinese (zh-HK)",
		"zh-tw":      "Traditional Chinese (zh-HK)",
		"zh-cn":      "Simplified Chinese (zh-CN)",
		"Chinese":    "Simplified Chinese (zh-CN)",
		"es":         "Spanish",
		"Spanish":    "Spanish",
		"ar":         "Arabic",
		"Arabic":     "Arabic",
		"fr":         "French",
		"French":     "French",
		"pt":         "Portuguese",
		"Portuguese": "Portuguese",
		"ru":         "Russian",
		"Russian":    "Russian",
		"id":         "Indonesian",
		"Indonesian": "Indonesian",
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

// escapeXml 转义 XML 特殊字符
func escapeXml(text string) string {
	buffer := new(bytes.Buffer)
	for _, r := range text {
		switch r {
		case '&':
			buffer.WriteString("&amp;")
		case '<':
			buffer.WriteString("&lt;")
		case '>':
			buffer.WriteString("&gt;")
		case '"':
			buffer.WriteString("&quot;")
		case '\'':
			buffer.WriteString("&apos;")
		default:
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

type TTSInput struct {
	Text         string `json:"text"`
	LanguageCode string `json:"languageCode"`
	SSMLGender   string `json:"ssmlGender"`
}

// ttsMaxChars 单次 Polly 合成的文本长度上限（neural 计费字符约 3000，留余量）。
const ttsMaxChars = 2800

// splitTextForTTS 将长文本按字符数切分，尽量在句子/空白边界断开，
// 使每段都在 Polly 单次合成上限内。
func splitTextForTTS(text string, maxChars int) []string {
	runes := []rune(text)
	if len(runes) <= maxChars {
		return []string{text}
	}
	var chunks []string
	start := 0
	for start < len(runes) {
		end := start + maxChars
		if end >= len(runes) {
			chunks = append(chunks, string(runes[start:]))
			break
		}
		// 优先在句子结束符处断开
		split := -1
		for i := end; i > start; i-- {
			switch runes[i-1] {
			case '。', '！', '？', '!', '?', '.', '\n', '；', ';':
				split = i
			}
			if split != -1 {
				break
			}
		}
		// 退而求其次，在空白处断开
		if split == -1 {
			for i := end; i > start; i-- {
				if runes[i-1] == ' ' || runes[i-1] == '\t' {
					split = i
					break
				}
			}
		}
		// 实在没有边界就硬切
		if split <= start {
			split = end
		}
		chunks = append(chunks, string(runes[start:split]))
		start = split
	}
	return chunks
}

// pollyVoiceID 按语言代码选择 Polly 语音，逻辑与 sericamind-service 一致
func pollyVoiceID(languageCode, ssmlGender string) pollytypes.VoiceId {
	switch languageCode {
	case "cmn-CN":
		return pollytypes.VoiceIdZhiyu
	case "yue-HK":
		return pollytypes.VoiceIdHiujin
	case "en-GB":
		return pollytypes.VoiceIdEmma
	case "en-US":
		if ssmlGender == "MALE" {
			return pollytypes.VoiceIdMatthew
		}
		return pollytypes.VoiceIdJoanna
	default:
		return pollytypes.VoiceIdZhiyu
	}
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
		input.LanguageCode = "cmn-CN"
	}

	// 按语言选择对应的 Polly 语音（与 sericamind-service 保持一致）
	voiceID := pollyVoiceID(input.LanguageCode, input.SSMLGender)

	// Polly 的粤语 LanguageCode 为 yue-CN，而对外契约使用 yue-HK
	pollyLang := input.LanguageCode
	if pollyLang == "yue-HK" {
		pollyLang = "yue-CN"
	}

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
		fmt.Printf("TTS config failed: %v\n", err)

		httputil.RespondJSON(w, 500, map[string]string{"error": "TTS service error"})
		return nil, nil
	}

	pollyClient := polly.NewFromConfig(cfg)

	// Polly SynthesizeSpeech 单次有字符上限（neural 约 3000 计费字符），
	// 长文本需分块逐段合成后再拼接 MP3。
	var audioData bytes.Buffer
	for _, chunk := range splitTextForTTS(input.Text, ttsMaxChars) {
		// 对文本进行 XML 转义并包装在 <speak> 标签中
		ssmlText := fmt.Sprintf("<speak>%s</speak>", escapeXml(chunk))

		speechInput := &polly.SynthesizeSpeechInput{
			Text:         aws.String(ssmlText),    // 使用 SSML 文本
			TextType:     pollytypes.TextTypeSsml, // 明确指定文本类型为 SSML
			OutputFormat: pollytypes.OutputFormatMp3,
			VoiceId:      voiceID,
			Engine:       engine,
			LanguageCode: pollytypes.LanguageCode(pollyLang),
		}

		output, err := pollyClient.SynthesizeSpeech(context.Background(), speechInput)
		if err != nil {
			fmt.Printf("TTS synthesis failed: %v\n", err)
			httputil.RespondJSON(w, 500, map[string]string{"error": "TTS synthesis failed"})
			return nil, nil
		}

		if _, err := io.Copy(&audioData, output.AudioStream); err != nil {
			output.AudioStream.Close()
			httputil.RespondJSON(w, 500, map[string]string{"error": "failed to read audio data"})
			return nil, nil
		}
		output.AudioStream.Close()
	}

	w.Header().Set("Content-Type", "audio/mpeg")
	w.Header().Set("Content-Disposition", "inline; filename=speech.mp3")
	w.WriteHeader(http.StatusOK)
	w.Write(audioData.Bytes())
	return nil, nil
}
