package http

import (
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
		err     error
		service string
		version string
		config  *httpConfig
		client  *http.Client
	}

	httpConfigData struct {
		Format string `json:"format"`
		Data   string `json:"data"`
	}
)

func NewSource(opts ...source.Option) source.Source {
	var (
		options          = source.NewOptions(opts...)
		service, version string
		client           *http.Client
	)

	if val, ok := options.Context.Value(serviceKey{}).(string); ok && val != "" {
		service = val
	}
	if val, ok := options.Context.Value(versionKey{}).(string); ok && val != "" {
		version = val
	}

	config, ok := options.Context.Value(httpConfigKey{}).(httpConfig)
	if !ok {
		log.Panic("config source http is not set.")
	} else {
		client = http.DefaultClient

		if config.Config == "" {
			config.Config = service
		}
	}

	return &httpSource{
		service: service,
		version: version,
		config:  &config,
		client:  client,
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

	if fetchUrl, rs.err = url.Parse(rs.config.Url); rs.err != nil {
		return nil, rs.err
	} else {
		fetchUrl.RawQuery = url.Values{
			"service": []string{rs.service},
			"version": []string{rs.version},
			"config":  []string{rs.config.Config},
		}.Encode()
	}

	if rs.config.Authorization != nil {
		auth := ""
		switch strings.ToLower(rs.config.Authorization.Type) {
		case "basic":
			auth = fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(rs.config.Authorization.Credentials)))
		case "bearer":
			auth = fmt.Sprintf("Bearer %s", rs.config.Authorization.Credentials)
		}
		header.Add("Authorization", auth)
	}

	if rsp, err := rs.client.Do(&http.Request{
		Method: "GET",
		URL:    fetchUrl,
		Header: header,
	}); err != nil {
		return nil, err
	} else {
		dec := json.NewDecoder(rsp.Body)
		if err = dec.Decode(&configData); err != nil {
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
