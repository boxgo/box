package http

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/boxgo/box/pkg/config/source"
)

type (
	httpSource struct {
		err        error
		namespace  string
		service    string
		version    string
		httpConfig *httpConfig
		client     *http.Client
	}

	httpConfigData struct {
		Namespace string `json:"namespace"`
		Service   string `json:"service"`
		Version   string `json:"version"`
		Format    string `json:"format"`
		Data      string `json:"data"`
	}
)

const (
	timeout = time.Second * 3
)

func NewSource(opts ...source.Option) source.Source {
	var (
		options                     = source.NewOptions(opts...)
		httpCfg                     *httpConfig
		namespace, service, version string
		client                      *http.Client
	)

	if val, ok := options.Context.Value(namespaceKey{}).(string); ok && val != "" {
		namespace = val
	}
	if val, ok := options.Context.Value(serviceKey{}).(string); ok && val != "" {
		service = val
	}
	if val, ok := options.Context.Value(versionKey{}).(string); ok && val != "" {
		version = val
	}

	if config, ok := options.Context.Value(httpConfigKey{}).(httpConfig); !ok {
		log.Panic("config source http is not set.")
	} else {
		httpCfg = &config
		client = http.DefaultClient
		client.Timeout = timeout
	}

	return &httpSource{
		namespace:  namespace,
		service:    service,
		version:    version,
		httpConfig: httpCfg,
		client:     client,
	}
}

func (rs *httpSource) Read() (*source.ChangeSet, error) {
	if rs.err != nil {
		return nil, rs.err
	}

	var (
		configData = httpConfigData{}
		fetchUrl   *url.URL
		header     = http.Header{}
	)

	if fetchUrl, rs.err = url.Parse(rs.httpConfig.Url); rs.err != nil {
		return nil, rs.err
	} else {
		fetchUrl.RawQuery = url.Values{
			"namespace": []string{rs.namespace},
			"service":   []string{rs.service},
			"version":   []string{rs.version},
		}.Encode()
	}

	if len(rs.httpConfig.Authorization.Type) > 0 && len(rs.httpConfig.Authorization.Type) > 0 {
		auth := ""
		switch strings.ToLower(rs.httpConfig.Authorization.Type) {
		case "basic":
			auth = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(rs.httpConfig.Authorization.Credentials)))
		case "bearer":
			auth = fmt.Sprintf("Bearer %s", rs.httpConfig.Authorization.Credentials)
		}
		header.Add("Authorization", auth)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), timeout)
	defer cancel()

	if rsp, err := rs.client.Do((&http.Request{
		Method: "GET",
		URL:    fetchUrl,
		Header: header,
	}).WithContext(ctx)); err != nil {
		log.Printf("config http request error: %#v", err)
		return nil, err
	} else {
		dec := json.NewDecoder(rsp.Body)
		if err = dec.Decode(&configData); err != nil {
			log.Printf("config http decode error: %#v", err)
			return nil, err
		}
	}

	cs := &source.ChangeSet{
		Timestamp: time.Now(),
		Source:    rs.String(),
		Data:      []byte(configData.Data),
		Format:    configData.Format,
	}
	cs.Checksum = cs.Sum()

	return cs, nil

}
func (rs *httpSource) Watch() (source.Watcher, error) {
	if rs.err != nil {
		return nil, rs.err
	}

	return newWatcher(rs)
}

func (rs *httpSource) String() string {
	return "http"
}
