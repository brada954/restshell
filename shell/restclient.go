package shell

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
	"time"
)

type RestClient struct {
	Debug   bool
	Verbose bool
	History bool
	Headers string
	Client  *http.Client
}

type RestResponse struct {
	Text     string
	httpResp *http.Response
}

func NewRestClient() RestClient {
	return RestClient{
		Debug:   false,
		Verbose: false,
		History: true,
		Headers: "",
		Client:  &http.Client{Timeout: time.Duration(30 * time.Second)},
	}
}

func NewRestClientFromOptions() RestClient {
	client := RestClient{
		Debug:   IsCmdDebugEnabled(),
		Verbose: IsCmdVerboseEnabled() && !IsCmdSilentEnabled(),
		History: true,
		Client:  &http.Client{Timeout: time.Duration(GetCmdTimeoutValueMs()) * time.Millisecond},
	}

	if IsCmdLocalCertsEnabled() {
		newTlsWithLocalCertificates(client.Client)
	}
	if IsCmdSkipCertValidationEnabled() {
		client.DisableCertValidation()
	}
	if IsCmdNoRedirectEnabled() {
		client.DisableRedirect()
	}

	if transport, ok := client.Client.Transport.(*http.Transport); ok {
		transport.MaxIdleConnsPerHost = 1000
	}

	client.Headers = GetCmdHeaderValues("")
	return client
}

func (r *RestClient) DisableHistory() {
	r.History = false
}

func (r *RestClient) DisableRedirect() {
	r.Client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }
}

func (r *RestClient) EnableRedirect() {
	r.Client.CheckRedirect = nil
}

// TODO: cleanup relationship be between skipping verification and including certificates
func (r *RestClient) DisableCertValidation() {
	if t, ok := r.Client.Transport.(*http.Transport); ok {
		if t.TLSClientConfig != nil {
			t.TLSClientConfig.InsecureSkipVerify = true
			return
		}
	}

	r.Client.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}

}

// Warning: this does not work well on windows as system
// certs may not be available
// TODO: cleanup relationship with disabling verification as this needs to be done first
func newTlsWithLocalCertificates(client *http.Client) {
	var tlsconfig = &tls.Config{RootCAs: GetX509Pool()}
	if t, ok := client.Transport.(*http.Transport); ok {
		t.TLSClientConfig = tlsconfig
	} else {
		client.Transport = &http.Transport{TLSClientConfig: tlsconfig}
	}
}

func (r *RestClient) DoGet(authContext Auth, url string) (resultResponse *RestResponse, resultError error) {
	return r.DoMethod(http.MethodGet, authContext, url)
}

func (r *RestClient) DoMethod(method string, authContext Auth, url string) (resultResponse *RestResponse, resultError error) {
	if r.History {
		defer func() {
			if resultError != nil {
				PushError(resultError)
			} else {
				PushResponse(resultResponse, resultError)
			}
		}()
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.New("Building request: " + err.Error())
	}
	if authContext != nil {
		authContext.AddAuth(req)
	}
	req.Header.Add("Content-Type", "application/json")
	if r.Debug {
		for _, c := range req.Cookies() {
			fmt.Fprintf(OutputWriter(), "Cookie:\n%s=%s\n", c.Name, c.Value)
		}
	}

	// Add headers from command parsing/client configuration
	if err := addHeaders(req, r.Headers); err != nil {
		fmt.Fprintf(OutputWriter(), "Warning: %s\n", err.Error())
	}

	if r.Debug {
		fmt.Fprintf(OutputWriter(), "Executing: (GET) %s\n", req.URL.String())
		fmt.Fprintln(OutputWriter(), "Headers:")
		dumpHeaders(OutputWriter(), req)
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		errMsg := "Response Error: " + err.Error()
		if r.Debug {
			fmt.Fprintln(OutputWriter(), errMsg)
		}
		return nil, errors.New(errMsg)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Reading Response: " + err.Error())
	}

	result := &RestResponse{string(body), resp}
	return result, nil
}

func (r *RestClient) DoWithJsonMarshal(method string, authContext Auth, url string, data interface{}) (*RestResponse, error) {
	body, error := json.Marshal(data)
	if error == nil {
		return r.DoWithJson(method, authContext, url, string(body))
	}
	return nil, errors.New("Bad Request")
}

func (r *RestClient) DoWithJson(method string, authContext Auth, url string, data string) (resultResponse *RestResponse, resultError error) {
	if r.History {
		defer func() {
			if resultError != nil {
				PushError(resultError)
			} else {
				PushResponse(resultResponse, resultError)
			}
		}()
	}

	if method == http.MethodGet {
		fmt.Fprintf(OutputWriter(), "Warning: using a JSON body with a Get method is not best practice")
	}

	JsonBodyValidate(r.Debug, data)

	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return nil, errors.New("Building request: " + err.Error())
	}
	if authContext != nil {
		authContext.AddAuth(req)
	}
	req.Header.Add("Content-Type", "application/json")

	// Add headers from command parsing/client configuration
	if err := addHeaders(req, r.Headers); err != nil {
		fmt.Fprintf(OutputWriter(), "Warning: %s\n", err.Error())
	}

	if r.Debug {
		fmt.Fprintf(OutputWriter(), "Executing: (%s) %s\n", method, req.URL.String())
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		errMsg := "Response Error: " + err.Error()
		if r.Debug {
			fmt.Fprintln(OutputWriter(), errMsg)
		}
		return nil, errors.New(errMsg)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Reading Response: " + err.Error())
	}

	result := &RestResponse{string(body), resp}
	return result, nil
}

func (r *RestClient) DoWithForm(method string, authContext Auth, url string, data string) (resultResponse *RestResponse, resultError error) {
	if r.History {
		defer func() {
			if resultError != nil {
				PushError(resultError)
			} else {
				PushResponse(resultResponse, resultError)
			}
		}()
	}

	if method == http.MethodGet {
		fmt.Fprintf(OutputWriter(), "Warning: using a body with a Get method is not best practice")
	}
	if r.Debug {
		fmt.Fprintf(OutputWriter(), "Executing: (%s) %s\n", method, url)
	}

	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return nil, errors.New("Building request: " + err.Error())
	}
	if authContext != nil {
		authContext.AddAuth(req)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Add headers from command parsing/client configuration
	if err := addHeaders(req, r.Headers); err != nil {
		fmt.Fprintf(OutputWriter(), "Warning: %s\n", err.Error())
	}

	resp, err := r.Client.Do(req)
	if err != nil {
		errMsg := "Response Error: " + err.Error()
		if r.Debug {
			fmt.Fprintln(OutputWriter(), errMsg)
		}
		return nil, errors.New(errMsg)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("Reading Response: " + err.Error())
	}

	result := &RestResponse{string(body), resp}
	return result, nil
}

func GetX509Pool() *x509.CertPool {
	// server cert is self signed -> server_cert == ca_cert
	CA_Pool, err := x509.SystemCertPool()
	if err != nil {
		fmt.Fprintf(ConsoleWriter(), "Warning: SystemCertPool returned error: %s", err.Error())
		CA_Pool = x509.NewCertPool()
	}
	certpath := filepath.Join(GetInitDirectory(), "cert.pem")
	severCert, err := ioutil.ReadFile(certpath)
	if err != nil {
		certpath := filepath.Join(GetExeDirectory(), "cert.pem")
		severCert, err = ioutil.ReadFile(certpath)
		if err != nil {
			fmt.Fprintln(ConsoleWriter(), "Could not load certificate:", certpath)
			return CA_Pool
		}
	}
	CA_Pool.AppendCertsFromPEM(severCert)
	return CA_Pool
}

func (resp *RestResponse) GetStatus() int {
	if resp != nil && resp.httpResp != nil {
		return resp.httpResp.StatusCode
	}
	return -1
}

func (resp *RestResponse) GetStatusString() string {
	if resp != nil && resp.httpResp != nil {
		return fmt.Sprintf("%s (%d)", resp.httpResp.Status, resp.httpResp.StatusCode)
	}
	return "Unknown Status"
}

func (resp *RestResponse) GetCookies() []*http.Cookie {
	return resp.httpResp.Cookies()
}

func (resp *RestResponse) GetHeader() http.Header {
	return resp.httpResp.Header
}

func JsonBodyValidate(debug bool, body string) {
	raw := make(map[string]interface{}, 0)
	bytes := []byte(body)
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		fmt.Fprintf(ErrorWriter(), "Warning: Failed to decode JSON body: %s\n", err.Error())
	}
	if debug {
		fmt.Fprintf(OutputWriter(), "Body:\n%s\n", body)
	}
}

// Global Helpers
func PerformHealthCheck(client RestClient, url string) error {
	url = url + "/health"

	resp, err := client.DoGet(nil, url)
	if err != nil || resp.GetStatus() != http.StatusOK {
		// if IsCmdVerboseEnabled() {
		// 	if err != nil {
		// 		fmt.Fprintf(OutputWriter(), "Health Check Error: %s\n", err.Error())
		// 	} else {
		// 		fmt.Fprintf(OutputWriter(), "Health Check Status: %d\n", resp.GetStatus())
		// 	}
		// }
		return errors.New("Failed Health Check")
	}
	fmt.Fprintln(OutputWriter(), "Healthy")
	return nil
}

func addHeaders(req *http.Request, headerParam string) error {
	if len(headerParam) == 0 {
		return nil
	}

	parts := strings.Split(headerParam, ",")
	headers := make(map[string]string, 0)

	for _, pair := range parts {
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			return errors.New("Error parsing header parameter; no headers used")
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		if strings.ToLower(key) == "host" {
			req.Host = value
		} else {
			headers[key] = value
		}
	}
	for key, value := range headers {
		fmt.Printf("Adding header: %s=%s\n", key, value)
		req.Header.Add(key, value)
	}
	return nil
}

func dumpHeaders(w io.Writer, req *http.Request) {
	for k, h := range req.Header {
		value := strings.Join(h, ",")
		fmt.Fprintf(w, "%s=%s\n", k, value)
	}
}
