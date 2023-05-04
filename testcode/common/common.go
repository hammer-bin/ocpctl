package common

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type AppConfigProperties map[string]string

var ConfInfo AppConfigProperties

// init() 함수는 패키지 내에서 가장먼저 실행되는 함수
// main() 함수 내에 포함된 패키지가 있을 경우 패키지내에 포함된 init() 함수가 먼저 실행된다.
func init() {

	profile := "local"
	if len(os.Getenv("PROFILE")) > 0 {
		profile = os.Getenv("PROFILE")
	}

	if profile == "local" {
		_, err := ReadPropertiesFile("ocp.properties")
		if err != nil {
			return
		}
	}

}

func ReadPropertiesFile(filename string) (AppConfigProperties, error) {
	ConfInfo = AppConfigProperties{}

	if len(filename) == 0 {
		return ConfInfo, nil
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				ConfInfo[key] = value
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return ConfInfo, nil
}
