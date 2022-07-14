package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/tidwall/gjson"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	searchDomain := "." + domain

	for i := 0; scanner.Scan(); i++ {
		if scanner.Err() != nil {
			return result, fmt.Errorf("read user data error: %w", scanner.Err())
		}

		email := gjson.GetBytes(scanner.Bytes(), "Email").String()

		if strings.Contains(email, searchDomain) {
			result[strings.ToLower(strings.SplitN(email, "@", 2)[1])]++
		}
	}

	return result, nil
}
