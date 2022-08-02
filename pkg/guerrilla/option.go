package guerrilla

import "net/http"

type Option func(*client)

func OptionWithHTTPClient(c *http.Client) Option {
	return func(client *client) {
		client.inner = c
	}
}

func OptionWithEndpoint(endpoint string) Option {
	return func(client *client) {
		client.endpoint = endpoint
	}
}

func OptionWithAgent(agent string) Option {
	return func(client *client) {
		client.agent = agent
	}
}

func OptionWithLanguage(language string) Option {
	return func(client *client) {
		client.language = language
	}
}
