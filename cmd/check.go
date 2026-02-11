package cmd

import (
	"fmt"
	"time"

	checker "github.com/Arnel-Rah/uptimego/internal"
)

func CheckAndLogEndpoint(endpoint map[string]interface{}) {
	name, _ := endpoint["name"].(string)
	url, _ := endpoint["url"].(string)
	timeoutStr, _ := endpoint["timeout"].(string)

	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		fmt.Printf("Timeout invalide pour %s: %v\n", name, err)
		return
	}
	result := checker.CheckEndpoint(url, timeout)
	fmt.Println(checker.FormatResult(name, url, result))
}
