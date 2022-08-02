package guerrilla

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Client interface {
	GetAllEmails() ([]EmailSummary, error)
	GetNewEmails() ([]EmailSummary, error)
	GetEmail(id string) (*Email, error)
	GetAddress() string
}

var _ Client = (*client)(nil)

type client struct {
	inner    *http.Client
	endpoint string
	agent    string
	language string
	session  session
}

type session struct {
	token   string
	email   string
	lastSeq string
}

var DefaultClient = client{
	inner: &http.Client{
		Timeout: time.Second * 10,
	},
	endpoint: "https://api.guerrillamail.com/ajax.php",
	agent:    "https://github.com/liamg/guerrilla",
	language: "en",
}

func Init(options ...Option) (Client, error) {
	client := DefaultClient
	for _, option := range options {
		option(&client)
	}
	if err := client.init(); err != nil {
		return nil, err
	}
	return &client, nil
}

const (
	methodGetEmailAddress = "get_email_address"
	methodGetEmailList    = "get_email_list"
	methodCheckEmail      = "check_email"
)

func (c *client) sendRequest(function string, params map[string]string, target interface{}) error {

	req, err := http.NewRequest(http.MethodGet, c.endpoint, nil)
	if err != nil {
		return err
	}

	query := req.URL.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	query.Set("f", function)
	query.Set("agent", c.agent)
	if c.session.token != "" {
		query.Set("sid_token", c.session.token)
	}
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

func (c *client) init() error {
	var resp getEmailAddressResponse
	if err := c.sendRequest(methodGetEmailAddress, map[string]string{
		"lang": c.language,
	}, &resp); err != nil {
		return err
	}
	c.session = session{
		token: resp.SidToken,
		email: strings.TrimSpace(resp.EmailAddr),
	}
	return nil
}

type getEmailListResponse struct {
	Alias    string            `json:"alias"`
	Count    string            `json:"count"`
	Email    string            `json:"email"`
	List     []apiEmailSummary `json:"list"`
	SidToken string            `json:"sid_token"`
	Stats    struct {
		CreatedAddresses int    `json:"created_addresses"`
		ReceivedEmails   string `json:"received_emails"`
		SequenceMail     string `json:"sequence_mail"`
		Total            string `json:"total"`
		TotalPerHour     string `json:"total_per_hour"`
	} `json:"stats"`
}

func (c *client) GetAllEmails() ([]EmailSummary, error) {
	var offset int
	var emails []EmailSummary
	for {
		var resp getEmailListResponse
		if err := c.sendRequest("get_email_list", map[string]string{
			"offset": strconv.Itoa(offset),
		}, &resp); err != nil {
			return nil, err
		}
		count, err := strconv.Atoi(resp.Count)
		if err != nil {
			return nil, err
		}
		for _, email := range resp.List {
			emails = append(emails, email.Summary())
		}
		if len(emails) >= count || len(resp.List) == 0 {
			break
		}
	}
	if len(emails) > 0 {
		c.session.lastSeq = emails[len(emails)-1].ID
	}
	return emails, nil
}

type checkEmailResponse struct {
	List     []apiEmailSummary `json:"list"`
	Count    vagueType         `json:"count"`
	Email    string            `json:"email"`
	Alias    string            `json:"alias"`
	Ts       vagueType         `json:"ts"`
	SidToken string            `json:"sid_token"`
	Stats    struct {
		SequenceMail     vagueType `json:"sequence_mail"`
		CreatedAddresses vagueType `json:"created_addresses"`
		ReceivedEmails   vagueType `json:"received_emails"`
		Total            vagueType `json:"total"`
		TotalPerHour     vagueType `json:"total_per_hour"`
	} `json:"stats"`
}

func (c *client) GetNewEmails() ([]EmailSummary, error) {
	var emails []EmailSummary
	for {
		var resp checkEmailResponse
		if err := c.sendRequest("check_email", map[string]string{
			"seq": c.session.lastSeq,
		}, &resp); err != nil {
			return nil, err
		}
		count := resp.Count.Int()
		for _, email := range resp.List {
			emails = append(emails, email.Summary())
		}
		if len(resp.List) > 0 {
			c.session.lastSeq = resp.List[len(resp.List)-1].ID.String()
		}
		if len(emails) >= count || len(resp.List) == 0 {
			break
		}
	}
	return emails, nil
}

func (c *client) GetEmail(id string) (*Email, error) {
	var resp apiEmail
	if err := c.sendRequest("fetch_email", map[string]string{
		"email_id": id,
	}, &resp); err != nil {
		return nil, err
	}
	email := resp.Email()
	return &email, nil
}

func (c *client) GetAddress() string {
	return c.session.email
}
