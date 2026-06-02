package exception

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Exception struct {
	Group               int
	Code                int
	Name                string
	MultilingualMessage map[string]string
	HTTPStatus          int
}

func NewException(group, code int, name string, messages map[string]string, httpStatus int) *Exception {
	return &Exception{
		Group:               group,
		Code:                code,
		Name:                name,
		MultilingualMessage: messages,
		HTTPStatus:          httpStatus,
	}
}

func (e *Exception) Error() string {
	return fmt.Sprintf("[%d] %s", e.FullCode(), e.DefaultMessage())
}

func (e *Exception) WithMsg(msg string) *ExError {
	return &ExError{Exception: e, Message: msg}
}

type ExError struct {
	*Exception
	Message string
}

func (e *ExError) Error() string {
	return fmt.Sprintf("[%d] %s", e.FullCode(), e.Message)
}

func (e *Exception) FullCode() int {
	return e.Group*100 + e.Code
}

func (e *Exception) DefaultMessage() string {
	if msg, ok := e.MultilingualMessage["en_us"]; ok {
		return msg
	}
	return e.Name
}

func (e *Exception) Write(w http.ResponseWriter, message ...string) {
	msg := e.DefaultMessage()
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(e.HTTPStatus)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"code":    e.FullCode(),
		"message": msg,
		"data": map[string]interface{}{
			"name":                e.Name,
			"multilingualMessage": e.MultilingualMessage,
			"params":              map[string]interface{}{},
		},
	})
}
