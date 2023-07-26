package client

import openapi "fybe.com/apiclient"

func NewClient(
	apiUrl string,
	authUrl string,
	clientId string,
	clientSecret *string,
	username string,
	password *string,
) (*openapi.APIClient, error) {
	configuration := openapi.NewConfiguration()
	configuration.AddDefaultHeader("x-trace-id", "fybe_terraform_provider")

	httpClient, err := BearerHttpClient(
		authUrl,
		clientId,
		*clientSecret,
		username,
		*password,
	)

	if err != nil {
		return nil, err
	}

	configuration.HTTPClient = httpClient

	var server openapi.ServerConfiguration
	server.URL = apiUrl

	var serverConfigurations []openapi.ServerConfiguration
	serverConfigurations = append(serverConfigurations, server)
	configuration.Servers = serverConfigurations

	return openapi.NewAPIClient(configuration), nil
}
