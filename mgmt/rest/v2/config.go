/*
http://www.apache.org/licenses/LICENSE-2.0.txt


Copyright 2017 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v2

import (
	"net/http"
	"strconv"

	"github.com/intelsdi-x/snap/control/plugin/cpolicy"
	"github.com/intelsdi-x/snap/core"
	"github.com/intelsdi-x/snap/core/cdata"
	"github.com/julienschmidt/httprouter"
)

type PolicyTable cpolicy.RuleTable

type PolicyTableSlice []cpolicy.RuleTable

// cdata.ConfigDataNode implements it's own UnmarshalJSON
type PluginConfigItem struct {
	cdata.ConfigDataNode
}

func (s *apiV2) getPluginConfigItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	styp := p.ByName("type")
	if styp == "" {
		cdn := s.configManager.GetPluginConfigDataNodeAll()
		item := &PluginConfigItem{ConfigDataNode: cdn}
		Write(200, item, w)
		return
	}

	typ, err := getPluginType(styp)
	if err != nil {
		Write(400, FromError(err), w)
		return
	}

	name := p.ByName("name")
	sver := p.ByName("version")
	iver := -2
	if sver != "" {
		if iver, err = strconv.Atoi(sver); err != nil {
			Write(400, FromError(err), w)
			return
		}
	}

	cdn := s.configManager.GetPluginConfigDataNode(typ, name, iver)
	item := &PluginConfigItem{ConfigDataNode: cdn}
	Write(200, item, w)
}

func (s *apiV2) deletePluginConfigItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var typ core.PluginType
	styp := p.ByName("type")
	if styp != "" {
		typ, err = getPluginType(styp)
		if err != nil {
			Write(400, FromError(err), w)
			return
		}
	}

	name := p.ByName("name")
	sver := p.ByName("version")
	iver := -2
	if sver != "" {
		if iver, err = strconv.Atoi(sver); err != nil {
			Write(400, FromError(err), w)
			return
		}
	}

	src := []string{}
	errCode, err := core.UnmarshalBody(&src, r.Body)
	if errCode != 0 && err != nil {
		Write(400, FromError(err), w)
		return
	}

	var res cdata.ConfigDataNode
	if styp == "" {
		res = s.configManager.DeletePluginConfigDataNodeFieldAll(src...)
	} else {
		res = s.configManager.DeletePluginConfigDataNodeField(typ, name, iver, src...)
	}

	item := &PluginConfigItem{ConfigDataNode: res}
	Write(200, item, w)
}

func (s *apiV2) setPluginConfigItem(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var err error
	var typ core.PluginType
	styp := p.ByName("type")
	if styp != "" {
		typ, err = getPluginType(styp)
		if err != nil {
			Write(400, FromError(err), w)
			return
		}
	}

	name := p.ByName("name")
	sver := p.ByName("version")
	iver := -2
	if sver != "" {
		if iver, err = strconv.Atoi(sver); err != nil {
			Write(400, FromError(err), w)
			return
		}
	}

	src := cdata.NewNode()
	errCode, err := core.UnmarshalBody(src, r.Body)
	if errCode != 0 && err != nil {
		Write(400, FromError(err), w)
		return
	}

	var res cdata.ConfigDataNode
	if styp == "" {
		res = s.configManager.MergePluginConfigDataNodeAll(src)
	} else {
		res = s.configManager.MergePluginConfigDataNode(typ, name, iver, src)
	}

	item := &PluginConfigItem{ConfigDataNode: res}
	Write(200, item, w)
}

func getPluginType(t string) (core.PluginType, error) {
	if ityp, err := strconv.Atoi(t); err == nil {
		return core.PluginType(ityp), nil
	}
	ityp, err := core.ToPluginType(t)
	if err != nil {
		return core.PluginType(-1), err
	}
	return ityp, nil
}
