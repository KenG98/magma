/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config

import (
	"encoding/json"
	"fbc/cwf/radius/monitoring/census"
	"io/ioutil"
)

// LiveTier name
const LiveTier = "live"

type (
	// ModuleDescriptor a descriptor for loading a single module
	ModuleDescriptor struct {
		Name   string                 `json:"name"`
		Config map[string]interface{} `json:"config"`
	}

	// ListenerConfig for a single listener (server has a listerner per each port)
	ListenerConfig struct {
		Name    string                 `json:"name"`
		Type    string                 `json:"type"`
		Modules []ModuleDescriptor     `json:"modules"`
		Extra   map[string]interface{} `json:"extra"` // Extra config, per listener
	}

	// ServiceTier represents a uniquely identifiable named set of upstream hosts
	ServiceTier struct {
		Name          string   `json:"name"`
		UpstreamHosts []string `json:"upstreamHosts"`
	}

	// ListenerRoute maps a listener to a ServiceTier name
	ListenerRoute struct {
		Listener    string `json:"listener"`
		ServiceTier string `json:"serviceTier"`
	}

	// TierRouting a set of ListenerRoute representing a map from Listener names to their respective ServiceTier
	TierRouting struct {
		Routes []ListenerRoute `json:"tierRoutes"`
	}

	// Canary represents a definition of a canary
	Canary struct {
		Name                string      `json:"name"`
		TrafficSlicePercent int         `json:"trafficSlicePercent"`
		Routing             TierRouting `json:"routing"`
	}

	// LoadBalanceConfig holds the complete configuration for a server load balancer
	LoadBalanceConfig struct {
		ServiceTiers []ServiceTier `json:"serviceTiers"`
		LiveTier     TierRouting   `json:"liveTier"`
		Canaries     []Canary      `json:"canaries"`
		DefaultTier  string        `json:"defaultTier"`
	}

	// RedisConfig the configuration of redus server
	RedisConfig struct {
		Addr     string `json:"addr"`
		Password string `json:"password"`
		DB       int    `json:"db"`
	}

	//SessionStorageConfig ...
	SessionStorageConfig struct {
		StorageType string      `json:"storageType"`
		Redis       RedisConfig `json:"redis"`
	}

	// ServerConfig Encapsulates the configuration of a radius server
	ServerConfig struct {
		Secret         string                `json:"secret"`
		DedupWindow    Duration              `json:"dedupWindow"`
		LoadBalance    LoadBalanceConfig     `json:"loadBalance"`
		Listeners      []ListenerConfig      `json:"listeners"`
		Filters        []string              `json:"filters"`
		SessionStorage *SessionStorageConfig `json:"sessionStorage"`
	}

	// MonitoringConfig ...
	MonitoringConfig struct {
		Census *census.Config `json:"census"`
		Ods    *Ods           `json:"ods"`
		Scuba  *Scuba         `json:"scuba"`
	}

	// RadiusConfig the configuration file format
	RadiusConfig struct {
		Monitoring *MonitoringConfig `json:"monitoring"`
		Server     ServerConfig      `json:"server"`
	}
)

// Read reads and parses a configuration file into a RadiusConfig
func Read(filename string) (*RadiusConfig, error) {
	configBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config RadiusConfig
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
