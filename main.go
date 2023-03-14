package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

type challenge_token struct {
	challenge_domain string
	challenge_txt    string
}

// test for acme.sh
func main() {
	fmt.Println("!... Hello World ...!")

	setting_acme_server("letsencrypt")

	output := issue_ca("letsencrypt", "test.ltpix.link")

	issue_challenge_token := getChallenge(&output)

	output2 := renew_ca("letsencrypt", "test.ltpix.link")
	renew_challenge_token := getChallenge(&output2)

	fmt.Println(issue_challenge_token.challenge_domain)
	fmt.Println(issue_challenge_token.challenge_txt)

	fmt.Println(renew_challenge_token.challenge_domain)
	fmt.Println(renew_challenge_token.challenge_txt)

}
func setting_acme_server(server string) error {

	command := "acme.sh --set-default-ca --server " + server

	result, err := exec.Command("bash", "-c", command).Output()

	if err != nil {
		fmt.Print(err.Error())
	} else {
		fmt.Print(string(result))
	}

	return err
}

func issue_ca(server string, domain string) []string {

	command := fmt.Sprintln("acme.sh --issue -d test.ltpix.link --server letsencrypt --dns --yes-I-know-dns-manual-mode-enough-go-ahead-please --log")

	result, _ := exec.Command("bash", "-c", command).Output()

	data := strings.Split(string(result), "\n")

	return data
}

func renew_ca(server string, domain string) []string {

	command := fmt.Sprintln("acme.sh --issue -d test.ltpix.link --renew --server letsencrypt --dns --yes-I-know-dns-manual-mode-enough-go-ahead-please --log")

	result, _ := exec.Command("bash", "-c", command).Output()

	data := strings.Split(string(result), "\n")

	return data

}
func getChallenge(result *[]string) challenge_token {
	regex := regexp.MustCompile(`_acme-challenge\.[A-Za-z0-9]+\.[A-Za-z0-9]+`)

	re := regexp.MustCompile(`TXT value: '(.+?)'`)

	var txtValue string
	var challenge_domain string
	// 使用正規表達式擷取完整的子域名
	for _, data := range *result {

		if len(regex.FindString(data)) > 1 {
			challenge_domain = regex.FindString(data)
		}

		match := re.FindStringSubmatch(data)
		if len(match) >= 2 {
			fmt.Println(match[1])
			txtValue = match[1]
		} else {
			fmt.Println("找不到 TXT value")
		}
	}

	return challenge_token{challenge_domain, txtValue}

}
