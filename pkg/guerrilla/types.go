package guerrilla

import (
	"encoding/json"
	"strconv"
	"time"
)

type EmailSummary struct {
	Att       int
	Date      string
	Excerpt   string
	From      string
	ID        string
	Read      bool
	Subject   string
	Timestamp time.Time
}

type Email struct {
	EmailSummary
	ReplyTo      string
	ContentType  string
	Recipient    string
	SourceID     string
	SourceMailID string
	Body         string
	Size         int
	RefMid       string
}

type apiEmailSummary struct {
	Att       vagueType `json:"att"`
	Date      string    `json:"mail_date"`
	Excerpt   string    `json:"mail_excerpt"`
	From      string    `json:"mail_from"`
	ID        vagueType `json:"mail_id"`
	Read      vagueType `json:"mail_read"`
	Subject   string    `json:"mail_subject"`
	Timestamp vagueType `json:"mail_timestamp"`
}

type apiEmail struct {
	apiEmailSummary
	ReplyTo      string    `json:"reply_to"`
	ContentType  string    `json:"content_type"`
	Recipient    string    `json:"mail_recipient"`
	SourceID     vagueType `json:"source_id"`
	SourceMailID vagueType `json:"source_mail_id"`
	Body         string    `json:"mail_body"`
	Size         int       `json:"size"`
	RefMid       string    `json:"ref_mid"`
}

type vagueType struct {
	int    int
	bool   bool
	string string
}

func (v vagueType) String() string {
	return v.string
}

func (v vagueType) Int() int {
	return v.int
}

func (v vagueType) Bool() bool {
	return v.bool
}

func (v *vagueType) UnmarshalJSON(raw []byte) error {

	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		v.string = s
		v.bool = s != ""
		v.int, _ = strconv.Atoi(s)
	}
	var i int
	if err := json.Unmarshal(raw, &i); err == nil {
		v.int = i
		v.string = strconv.Itoa(i)
		v.bool = i > 0
	}
	var b bool
	if err := json.Unmarshal(raw, &b); err == nil {
		v.bool = b
		if b {
			v.int = 1
			v.string = "1"
		} else {
			v.string = "0"
			v.int = 0
		}
	}
	return nil
}

func (a apiEmail) Email() Email {
	return Email{
		EmailSummary: a.apiEmailSummary.Summary(),
		ReplyTo:      a.ReplyTo,
		ContentType:  a.ContentType,
		Recipient:    a.Recipient,
		SourceID:     a.SourceID.String(),
		SourceMailID: a.SourceMailID.String(),
		Body:         a.Body,
		Size:         a.Size,
		RefMid:       a.RefMid,
	}
}

func (a apiEmailSummary) Summary() EmailSummary {

	// for some reason the api sends 0 timestamps most of the time
	var t time.Time
	if a.Timestamp.Int() == 0 {
		t = time.Now()
	} else {
		t = time.Unix(int64(a.Timestamp.Int()), 0)
	}

	return EmailSummary{
		Att:       a.Att.Int(),
		Date:      a.Date,
		Excerpt:   a.Excerpt,
		From:      a.From,
		ID:        a.ID.String(),
		Read:      a.Read.Bool(),
		Subject:   a.Subject,
		Timestamp: t,
	}
}
