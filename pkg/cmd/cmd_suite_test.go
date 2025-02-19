package cmd_test

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/redhat-et/copilot-ops/pkg/openai"
	gogpt "github.com/sashabaranov/go-gpt3"
)

func TestCmd(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cmd Suite")
}

// OpenAITestServer Creates a mocked OpenAI server which can pretend to handle requests during testing.
func OpenAITestServer() *httptest.Server {
	// FIXME: replace these patterns once this PR is merged:
	// https://github.com/sashabaranov/go-gpt3/pull/28
	var completionPath = regexp.MustCompile(`/v1/engines/[a-zA-Z0-9\-]+/completions`)

	return httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resBytes []byte
		log.Println("received request at path '", r.URL.Path, "'")
		switch {
		case r.URL.Path == "/v1/edits":
			res := gogpt.EditsResponse{
				Object:  "test-object",
				Created: uint64(time.Now().Unix()),
				Usage: gogpt.EditsUsage{
					PromptTokens:     23,
					CompletionTokens: 24,
					TotalTokens:      47,
				},
				Choices: []gogpt.EditsChoice{
					{
						Text: `
# @path/to/kubernetes.yaml
apiVersion: v1
kind: Pod
metadata:
	name: cute-cats
spec: 
	priority: high
`,
						Index: 0,
					},
				},
			}
			resBytes, _ = json.Marshal(res)
			fmt.Fprint(w, string(resBytes))
			return
		case completionPath.MatchString(r.URL.Path):
			res := gogpt.CompletionResponse{
				Object:  "test-object",
				Created: uint64(time.Now().Unix()),
				Model:   openai.OpenAICodeDavinciV2,
				Usage: gogpt.CompletionUsage{
					PromptTokens:     23,
					CompletionTokens: 24,
					TotalTokens:      47,
				},
				ID: "test-id",
				Choices: []gogpt.CompletionChoice{
					{
						Text:  "choice 1",
						Index: 0,
					},
				},
			}
			resBytes, _ = json.Marshal(res)
			fmt.Fprintln(w, string(resBytes))
			return
		default:
			// the endpoint doesn't exist
			log.Println("test server was accessed, but no endpoint was found")
			http.Error(w, "the resource path doesn't exist", http.StatusNotFound)
			return
		}
	}))
}
