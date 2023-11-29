package adstxt

import (
	"fmt"
	"source/pkg/utility"
	"strings"
)

type LineSchema struct {
	Domain        string
	AccountId     string
	Type          string
	Certification string
	Comment       string
}

type LineAdsTxt string

type lines []string

func (output lines) ToString() string {
	return strings.Join(output, "\n")
}

func StandardizedWithText(linesText string) (outputs lines) {
	return Standardized(strings.Split(linesText, "\n"))
}

func Standardized(linesInput []string) lines {
	var outputs []string
	for _, line := range linesInput {
		var lineOutput string
		if strings.Contains(line, ",") {
			lineOutput = FormatAdsTxt(line)
			if utility.ValidateString(lineOutput) == "" {
				continue
			}
		} else {
			lineOutput = line
		}
		if ok, _ := utility.InStringArray(lineOutput, outputs); !ok {
			outputs = append(outputs, lineOutput)
		}
	}
	return outputs
}

func FormatAdsTxt(lineString string) string {
	line := Parse(lineString)

	if certification, ok := Certifications[line.Domain]; ok {
		if line.Certification != certification {
			line.Certification = certification
		}
	}

	var output string
	if line.Domain != "" && line.AccountId != "" && line.Type != "" {
		output = fmt.Sprintf("%s, %s, %s",
			line.Domain,
			line.AccountId,
			line.Type,
		)
	}

	if line.Certification != "" {
		output += ", " + line.Certification
	}

	if line.Comment != "" {
		output += " #" + line.Comment
	}

	return output
}

func Parse(lineString string) (line LineSchema) {
	lineWithComment := strings.Split(lineString, "#")
	if len(lineWithComment) > 1 {
		line.Comment = strings.TrimSpace(lineWithComment[1])
	}
	lineArray := strings.Split(lineWithComment[0], ",")
	for key, value := range lineArray {
		value = strings.TrimSpace(value)
		switch key {
		case 0:
			line.Domain = formatDomain(value)
		case 1:
			line.AccountId = formatAccountId(value)
		case 2:
			line.Type = formatType(value)
		case 3:
			line.Certification = formatCertification(value)
		}
	}
	return
}

func formatDomain(domain string) string {
	domain = strings.ToLower(domain)
	return domain
}

func formatAccountId(accountId string) string {
	return accountId
}

func formatType(typeString string) string {
	typeString = strings.ToUpper(typeString)
	return typeString
}

func formatCertification(certification string) string {
	return certification
}
