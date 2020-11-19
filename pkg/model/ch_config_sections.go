// Copyright 2019 Altinity Ltd and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	chi "github.com/altinity/clickhouse-operator/pkg/apis/clickhouse.altinity.com/v1"
	"github.com/altinity/clickhouse-operator/pkg/util"
)

// configSectionsGenerator
type configSectionsGenerator struct {
	// ClickHouse config generator
	chConfigGenerator *ClickHouseConfigGenerator
	// clickhouse-operator configuration
	chopConfig *chi.OperatorConfig
}

// NewConfigSections
func NewConfigSectionsGenerator(
	chConfigGenerator *ClickHouseConfigGenerator,
	chopConfig *chi.OperatorConfig,
) *configSectionsGenerator {
	return &configSectionsGenerator{
		chConfigGenerator: chConfigGenerator,
		chopConfig:        chopConfig,
	}
}

// CreateConfigsCommon
func (c *configSectionsGenerator) CreateConfigsCommon() map[string]string {
	commonConfigSections := make(map[string]string)
	// commonConfigSections maps section name to section XML chopConfig of the following sections:
	// 1. remote servers
	// 2. common settings
	// 3. common files
	util.IncludeNonEmpty(commonConfigSections, createConfigSectionFilename(configRemoteServers), c.chConfigGenerator.GetRemoteServers())
	util.IncludeNonEmpty(commonConfigSections, createConfigSectionFilename(configSettings), c.chConfigGenerator.GetSettings(nil))
	util.MergeStringMaps(commonConfigSections, c.chConfigGenerator.GetFiles(chi.SectionCommon, true, nil))
	// Extra user-specified config files
	util.MergeStringMaps(commonConfigSections, c.chopConfig.CHCommonConfigs)

	return commonConfigSections
}

// CreateConfigsUsers
func (c *configSectionsGenerator) CreateConfigsUsers() map[string]string {
	commonUsersConfigSections := make(map[string]string)
	// commonUsersConfigSections maps section name to section XML chopConfig of the following sections:
	// 1. users
	// 2. quotas
	// 3. profiles
	// 4. user files
	util.IncludeNonEmpty(commonUsersConfigSections, createConfigSectionFilename(configUsers), c.chConfigGenerator.GetUsers())
	util.IncludeNonEmpty(commonUsersConfigSections, createConfigSectionFilename(configQuotas), c.chConfigGenerator.GetQuotas())
	util.IncludeNonEmpty(commonUsersConfigSections, createConfigSectionFilename(configProfiles), c.chConfigGenerator.GetProfiles())
	util.MergeStringMaps(commonUsersConfigSections, c.chConfigGenerator.GetFiles(chi.SectionUsers, false, nil))
	// Extra user-specified config files
	util.MergeStringMaps(commonUsersConfigSections, c.chopConfig.CHUsersConfigs)

	return commonUsersConfigSections
}

// CreateConfigsHost
func (c *configSectionsGenerator) CreateConfigsHost(host *chi.ChiHost) map[string]string {
	// Prepare for this replica deployment chopConfig files map as filename->content
	hostConfigSections := make(map[string]string)
	util.IncludeNonEmpty(hostConfigSections, createConfigSectionFilename(configMacros), c.chConfigGenerator.GetHostMacros(host))
	util.IncludeNonEmpty(hostConfigSections, createConfigSectionFilename(configPorts), c.chConfigGenerator.GetHostPorts(host))
	util.IncludeNonEmpty(hostConfigSections, createConfigSectionFilename(configZookeeper), c.chConfigGenerator.GetHostZookeeper(host))
	util.IncludeNonEmpty(hostConfigSections, createConfigSectionFilename(configSettings), c.chConfigGenerator.GetSettings(host))
	util.MergeStringMaps(hostConfigSections, c.chConfigGenerator.GetFiles(chi.SectionHost, true, host))
	// Extra user-specified config files
	util.MergeStringMaps(hostConfigSections, c.chopConfig.CHHostConfigs)

	return hostConfigSections
}

// createConfigSectionFilename
func createConfigSectionFilename(section string) string {
	return "chop-generated-" + section + ".xml"
}
