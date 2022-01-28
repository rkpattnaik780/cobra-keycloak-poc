package utils

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
)

//go:embed static/success.html
var successPage string

var (
	CloseApp sync.WaitGroup
)

func OpenBrowser(url string) error {
	var browserCommand *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		browserCommand = exec.Command("xdg-open", url)
	case "windows":
		browserCommand = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		browserCommand = exec.Command("open", url)
	default:
		return fmt.Errorf("unsupported operating system: %v", runtime.GOOS)
	}
	err := browserCommand.Run()
	return err
}

type Config struct {
	KeycloakConfig       KeycloakConfig
	EmbeddedServerConfig EmbeddedServerConfig
}

type KeycloakConfig struct {
	KeycloakURL string
	Realm       string
	ClientID    string
}

type EmbeddedServerConfig struct {
	Port         uint32
	CallbackPath string
}

func (c *EmbeddedServerConfig) GetCallbackURL() string {
	return fmt.Sprintf("http://localhost:%v/%v", c.Port, c.CallbackPath)
}

func StartServer(config Config) {

	serverAddress := fmt.Sprintf("localhost:%v", config.EmbeddedServerConfig.Port)

	http.HandleFunc("/sso-callback", func(w http.ResponseWriter, r *http.Request) {

		code := r.URL.Query().Get("code")
		if code != "" {
			request, err := BuildTokenExchangeRequest(config, code)
			if err == nil {
				var resp *http.Response
				var body []byte
				resp, err = http.DefaultClient.Do(request)
				if err == nil {
					body, err = ioutil.ReadAll(io.LimitReader(resp.Body, 1<<20))
					defer resp.Body.Close()
					if resp.StatusCode == 200 {
						content, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))
						switch content {
						case "application/json":
							var f interface{}
							json.Unmarshal(body, &f)
							m := f.(map[string]interface{})

							var cfg *IConfig = &IConfig{
								AccessToken:  m["access_token"].(string),
								RefreshToken: m["refresh_token"].(string),
							}

							SaveConfig(cfg)
							fmt.Println("Successfully Logged In...")
						default:
							fmt.Println(body)
						}
					} else {
						err = fmt.Errorf("invalid Status code (%v)", resp.StatusCode)
					}
					fmt.Fprintf(w, successPage)
					CloseApp.Done()
				}
			}
		}
	})

	go func() {
		log.Print("Booting up the server")
		if err := http.ListenAndServe(serverAddress, nil); err != nil {
			log.Fatalf("Unable to start server: %v\n", err)
			CloseApp.Done()
		}
	}()
}
