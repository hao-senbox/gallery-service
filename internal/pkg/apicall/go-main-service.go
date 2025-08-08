package apicall

import (
	"context"
	"encoding/json"
	"gallery-service/internal/application/dto/responses"
	"gallery-service/internal/pkg/apicall/dto"
	"gallery-service/pkg/consul"
	"net/http"

	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

var (
	goMainService = "go-main-service"
	basePath      = "/v1"
)

type CallAPI struct {
	client        consul.ServiceDiscovery
	clientService *api.CatalogService
}

func NewGoMainServiceAPI(client *api.Client) (*CallAPI, error) {
	// Create ServiceDiscovery instance with Consul address and service name
	sd, err := consul.NewServiceDiscovery(client, goMainService)
	if err != nil {
		return nil, err
	}

	// Discover service
	service, err := sd.DiscoverService()
	if err != nil {
		return nil, err
	}

	return &CallAPI{
		client:        sd,
		clientService: service,
	}, nil
}

func (c *CallAPI) GetUserByID(ctx context.Context, userID string) (*dto.UserEntityResponse, error) {
	currentToken, ok := ctx.Value("current_token").(string)
	if !ok {
		return nil, errors.New("unauthorized")
	}

	res, err := c.client.CallAPI(c.clientService, basePath+"/user/"+userID, http.MethodGet, nil, map[string]string{
		"Authorization": currentToken,
	})
	if err != nil {
		return nil, err
	}
	var resDTO dto.UserEntityResponse
	err = json.Unmarshal(res, &resDTO)

	return &resDTO, err
}

func (c *CallAPI) GetUserByToken(_ context.Context, token string) (*dto.UserEntityResponse, error) {
	res, err := c.client.CallAPI(
		c.clientService,
		basePath+"/user/current-user",
		http.MethodGet,
		nil,
		map[string]string{
			"Authorization": token,
		},
	)
	if err != nil {
		return nil, err
	}

	var apiRes responses.APIResponse[dto.UserEntityResponse]
	err = json.Unmarshal(res, &apiRes)
	if err != nil {
		return nil, err
	}

	return &apiRes.Data, nil
}
