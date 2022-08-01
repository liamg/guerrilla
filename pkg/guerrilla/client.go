package guerrilla

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type client struct {
    inner    *http.Client
    endpoint string
    agent    string
    language string
    token    string
}

var DefaultClient = client{
    inner: &http.Client{
        Timeout: time.Second * 10,
    },
    endpoint: "https://api.guerrillamail.com/ajax.php",
    agent:    "https://github.com/liamg/guerrilla",
    language: "en",
}

func New(options ...Option) *client {
    client := DefaultClient
    for _, option := range options {
        option(&client)
    }
    return &client
}

const (
    methodGetEmailAddress = "get_email_address"
    methodGetEmailList    = "get_email_list"
    methodCheckEmail      = "check_email"
)

func (c *client) sendRequest(function string, params map[string]string, target interface{}) error {

    req, err := http.NewRequest(http.MethodGet, c.endpoaint, nil)
    if err != nil {
        return err
    }

    query := req.URL.Query()
    for key, value := range params {
        query.Set(key, value)
    }
    query.Set("f", function)
    query.Set("agent", c.agent)
    req.URL.RawQuery = query.Encode()

    req.Header.Set("Accept", "application/json")
    resp, err := c.inner.Do(req)
    if err != nil {
        return err
    }
    defer func() { _ = resp.Body.Close() }()

    if resp.StatusCode >= 400 {
        return fmt.Errorf("request failed with status code %d", resp.StatusCode)
    }
    if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
        return err
    }

    return nil
}

type getEmailAddressResponse struct {
    EmailAddr      string `json:"email_addr"`
    EmailTimestamp int    `json:"email_timestamp"`
    Alias          string `json:"alias"`
    AliasError     string `json:"alias_error"`
    SidToken       string `json:"sid_token"`
}

func (c *client) GetEmailAddress() (string, error) {
    var resp getEmailAddressResponse
    if err := c.sendRequest(methodGetEmailAddress, map[string]string{
        "lang":      c.language,
        "sid_token": c.token,
    }, &resp); err != nil {
        return "", err
    }
    c.token = resp.SidToken
    return resp.EmailAddr, nil
}
