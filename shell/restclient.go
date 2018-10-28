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

// RestClient -- An object that makes REST calls
type RestClient struct {
	Debug   bool
	Verbose bool
	Headers string
	Client  *http.Client
}

// RestResponse -- The response structure returned by a REST interface
type RestResponse struct {
	Text     string
	httpResp *http.Response
}

func NewRestClient() RestClient {
	return RestClient{
		Debug:   false,
		Verbose: false,
		Headers: "",
		Client:  &http.Client{Timeout: time.Duration(30 * time.Second)},
	}
}

func NewRestClientFromOptions() RestClient {

	client := RestClient{
		Debug:   IsCmdDebugEnabled(),
		Verbose: IsCmdVerboseEnabled() && !IsCmdSilentEnabled(),
		Client: &http.Client{
			Timeout: time.Duration(GetCmdTimeoutValueMs()) * time.Millisecond,
		},
	}

	// Create non-default transport to break connecxtion pooling
	if IsCmdReconnectEnabled() {
		client.Client.Transport = &http.Transport{
			MaxIdleConns:          1000,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		}
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

func (r *RestClient) DisableRedirect() {
	r.Client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error { return http.ErrUseLastResponse }
}

func (r *RestClient) EnableRedirect() {
	r.Client.CheckRedirect = nil
}

// DisableCertValidation -- A function to disable cert validation in the rest client
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

func (r *RestClient) DoGet(authContext Auth, url string) (resultResponse *RestResponse, resultError error) {
	return r.DoMethod(http.MethodGet, authContext, url)
}

func (r *RestClient) DoMethod(method string, authContext Auth, url string) (resultResponse *RestResponse, resultError error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, errors.New("Building request: " + err.Error())
	}
	if authContext != nil {
		authContext.AddAuth(req)
	}

	// Add headers from command parsing/client configuration
	if err := addHeaders(req, r.Headers); err != nil {
		fmt.Fprintf(OutputWriter(), "Warning: %s\n", err.Error())
	}

	// TODO: What is the best content handling; was hardcoded to json
	contentType := "application/json"
	addDefaultContentType(req, contentType)

	if r.Debug {
		fmt.Fprintf(OutputWriter(), "Executing: (GET) %s\n", req.URL.String())
		fmt.Fprintln(OutputWriter(), "Sending Headers:")
		dumpHeaders(OutputWriter(), req)

		fmt.Fprintln(OutputWriter(), "Sending Cookies:")
		for _, c := range req.Cookies() {
			fmt.Fprintf(OutputWriter(), "%s=%s\n", c.Name, c.Value)
		}
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
	JsonBodyValidate(r.Debug, data)
	return r.DoMethodWithBody(method, authContext, url, "application/json", data)
}

// DoWithXml -- Perform an HTTP method request with an XML body
func (r *RestClient) DoWithXml(method string, authContext Auth, url string, data string) (resultResponse *RestResponse, resultError error) {

	XmlBodyValidate(r.Debug, data)
	return r.DoMethodWithBody(method, authContext, url, "application/xml", data)
}

func (r *RestClient) DoWithForm(method string, authContext Auth, url string, data string) (resultResponse *RestResponse, resultError error) {

	FormBodyValidate(r.Debug, data)
	return r.DoMethodWithBody(method, authContext, url, "application/x-www-form-urlencoded", data)
}

// DoMethodWithBody - Perform a HTTP request for the given method type and content provided
func (r *RestClient) DoMethodWithBody(method string, authContext Auth, url string, contentType string, data string) (resultResponse *RestResponse, resultError error) {
	if r.Debug {
		fmt.Fprintf(OutputWriter(), "Body:\n%s\n", data)
	}

	if method == http.MethodGet {
		fmt.Fprintf(OutputWriter(), "Warning: using a HTTP body with a Get method is not best practice")
	}

	req, err := http.NewRequest(method, url, strings.NewReader(data))
	if err != nil {
		return nil, errors.New("Building request: " + err.Error())
	}
	if authContext != nil {
		authContext.AddAuth(req)
	}

	// Add headers from command parsing/client configuration
	if err := addHeaders(req, r.Headers); err != nil {
		fmt.Fprintf(OutputWriter(), "Warning: %s\n", err.Error())
	}

	addDefaultContentType(req, contentType)

	if r.Debug {
		fmt.Fprintf(OutputWriter(), "Executing: (%s) %s\n", method, req.URL.String())
		fmt.Fprintln(OutputWriter(), "Sending Headers:")
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

// GetX509Pool - Get the X509 pool to use; the system is used by default and
// a cert.pem file is appended if it can be found in the init directory,
// or the exe directory.
func GetX509Pool() *x509.CertPool {
	// server cert is self signed -> server_cert == ca_cert
	CAPool, err := x509.SystemCertPool()
	if err != nil {
		fmt.Fprintf(ConsoleWriter(), "Warning: SystemCertPool returned error: %s", err.Error())
		CAPool = x509.NewCertPool()
	}
	certpath := filepath.Join(GetInitDirectory(), "cert.pem")
	severCert, err := ioutil.ReadFile(certpath)
	if err != nil {
		certpath := filepath.Join(GetExeDirectory(), "cert.pem")
		severCert, err = ioutil.ReadFile(certpath)
		if err != nil {
			fmt.Fprintln(ConsoleWriter(), "Could not load certificate:", certpath)
			return CAPool
		}
	}
	CAPool.AppendCertsFromPEM(severCert)
	return CAPool
}

// GetStatus - Get the status return code in the response; if the response is
// invalid or not initialized it will return -1
func (resp *RestResponse) GetStatus() int {
	if resp != nil && resp.httpResp != nil {
		return resp.httpResp.StatusCode
	}
	return -1
}

// GetStatusString - Get the status code in string format for the response; if the response is
// invalid or not initialized it will return "Unknown Status".
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

func (resp *RestResponse) GetContentType() string {
	contentType := "application/octet-stream"
	for k, v := range resp.httpResp.Header {
		if strings.ToLower(k) == "content-type" {
			if len(v) > 0 {
				contentType = strings.ToLower(v[0])
			}
		}
	}
	return contentType
}

// JsonBodyValidate -  Validate form data is ok; only display a warning if not
func JsonBodyValidate(debug bool, body string) {
	raw := make(map[string]interface{}, 0)
	bytes := []byte(body)
	err := json.Unmarshal(bytes, &raw)
	if err != nil {
		fmt.Fprintf(ErrorWriter(), "Warning: Failed to decode JSON body: %s\n", err.Error())
	}
}

// XmlBodyValidate -  Validate form data is ok; only display a warning if not
func XmlBodyValidate(debug bool, body string) {
	// TODO: Maybe parse the xml into a dom
}

// FormBodyValidate -  Validate form data is ok; only display a warning if not
func FormBodyValidate(debug bool, body string) {
	// TODO
}

// PerformHealthCheck - Valiate a default /health endpoint
// TODO: Evaluate moving this to specialized commands; it is not generic enough to be heres
func PerformHealthCheck(client RestClient, url string) error {
	url = url + "/health"

	resp, err := client.DoGet(nil, url)
	if err != nil || resp.GetStatus() != http.StatusOK {
		return errors.New("Failed Health Check")
	}
	fmt.Fprintln(OutputWriter(), "Healthy")
	return nil
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

func addHeaders(req *http.Request, headerParam string) error {
	if len(headerParam) == 0 {
		return nil
	}

	parts := strings.Split(headerParam, ",")
	headers := make(map[string]string, 0)

	for _, pair := range parts {
		kv := strings.SplitN(pair, "=", 2)
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

// addDefaultContentType -- adds the content type unless header exists
func addDefaultContentType(req *http.Request, contentType string) {
	if len(strings.TrimSpace(contentType)) > 0 {
		contentTypeSet := false
		for k := range req.Header {
			if strings.ToLower(k) == "content-type" {
				contentTypeSet = true
			}
		}
		if !contentTypeSet {
			req.Header.Add("Content-Type", contentType)
		}
	}
}

func dumpHeaders(w io.Writer, req *http.Request) {
	for k, h := range req.Header {
		value := strings.Join(h, ",")
		fmt.Fprintf(w, "%s=%s\n", k, value)
	}
}
