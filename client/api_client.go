package client

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"

    "devicedb/routes"
)

type APIClientConfig struct {
    Servers []string
}

type APIClient struct {
    servers []string
    nextServerIndex int
    httpClient *http.Client
}

func New(config APIClientConfig) *APIClient {
    return &APIClient{
        servers: config.Servers,
        nextServerIndex: 0,
        httpClient: &http.Client{ },
    }
}

func (client *APIClient) nextServer() (server string) {
    if len(client.servers) == 0 {
        return
    }

    server = client.servers[client.nextServerIndex]
    client.nextServerIndex = (client.nextServerIndex + 1) % len(client.servers)

    return
}

func (client *APIClient) AddSite(ctx context.Context, siteID string) error {
    _, err := client.sendRequest(ctx, "PUT", "/sites/" + siteID, nil)

    if err != nil {
        return err
    }

    return nil
}

func (client *APIClient) RemoveSite(ctx context.Context, siteID string) error {
    _, err := client.sendRequest(ctx, "DELETE", "/sites/" + siteID, nil)

    if err != nil {
        return err
    }

    return nil
}

func (client *APIClient) AddRelay(ctx context.Context, relayID string) error {
    _, err := client.sendRequest(ctx, "PUT", "/relays/" + relayID, nil)

    if err != nil {
        return err
    }

    return nil
}

func (client *APIClient) MoveRelay(ctx context.Context, relayID string, siteID string) error {
    var relaySettingsPatch routes.RelaySettingsPatch = routes.RelaySettingsPatch{ Site: siteID }

    body, err := json.Marshal(relaySettingsPatch)

    if err != nil {
        return err
    }

    _, err = client.sendRequest(ctx, "PATCH", "/relays/" + relayID, body)

    if err != nil {
        return err
    }

    return nil
}

func (client *APIClient) RemoveRelay(ctx context.Context, relayID string) error {
    _, err := client.sendRequest(ctx, "DELETE", "/relays/" + relayID, nil)

    if err != nil {
        return err
    }

    return nil
}

func (client *APIClient) Batch(ctx context.Context, siteID string, batch Batch) error {
    return nil
}

func (client *APIClient) Get(ctx context.Context, siteID string, keys []string) ([]Entry, error) {
    return nil, nil
}

func (client *APIClient) GetMatches(ctx context.Context, siteID string, keys []string) (EntryIterator, error) {
    return EntryIterator{}, nil
}

func (client *APIClient) sendRequest(ctx context.Context, httpVerb string, endpointURL string, body []byte) ([]byte, error) {
    u := fmt.Sprintf("http://%s%s", client.nextServer(), endpointURL)
    request, err := http.NewRequest(httpVerb, u, bytes.NewReader(body))

    if err != nil {
        return nil, err
    }

    request = request.WithContext(ctx)

    resp, err := client.httpClient.Do(request)

    if err != nil {
        return nil, err
    }
    
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        errorMessage, err := ioutil.ReadAll(resp.Body)
        
        if err != nil {
            return nil, err
        }
       
        return nil, &ErrorStatusCode{ Message: string(errorMessage), StatusCode: resp.StatusCode }
    }

    responseBody, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        return nil, err
    }

    return responseBody, nil
}