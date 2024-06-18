package inst

import "strings"

func ToLines(data []byte) []string {
	s := string(data)
	body := strings.Split(s, "bits 16")[1]
	return strings.Split(strings.TrimSpace(body), "\n")
}

type Tokens struct{}
type SimulationResult struct{}

// TODO: Implement
func Tokenize(lines []string) (Tokens, error) {
	return Tokens{}, nil
}

// TODO: Implement
func Simulate(t Tokens) (SimulationResult, error) {
	return SimulationResult{}, nil
}

// TODO: Implement
func PrintResult(result SimulationResult) {
}

// TODO: Implement
func (sr *SimulationResult) String() string {
	return ""
}
