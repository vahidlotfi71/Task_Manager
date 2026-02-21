
package Controller

func getString(body map[string]interface{}, key string) string {
	if val, ok := body[key]; ok && val != nil {
		return val.(string)
	}
	return ""
}
