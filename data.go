/*
 * Copyright (c) 2013 Dario Castañé.
 * This file is part of Zas.
 *
 * Zas is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Zas is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Zas.  If not, see <http://www.gnu.org/licenses/>.
 */
package main

import (
	"errors"
	"fmt"
	thtml "html/template"
	"path"
	"strings"
)

/*
 * Context data store used in templates.
 */
type ZasData struct {
	// Template used as body from current file.
	Body thtml.HTML
	// Current title, from first level header (H1).
	Title string
	// Current path (usable in URLs).
	Path string
	// Site configuration, as found in ZAS_CONF_FILE.
	Site ZasSiteData
	// Current configuration, from first HTML comment (expected as YAML map).
	Page map[interface{}]interface{}
	// Config loaded from ZAS_CONF_FILE.
	config ConfigSection
}

/*
 * Site configuration.
 *
 * They are required fields in order to complete social/semantic meta tags.
 */
type ZasSiteData struct {
	BaseURL string
	Image   string
}

/*
 * Builds URL from current configuration.
 */
func (zd *ZasData) URL() string {
	return fmt.Sprintf("%s%s", zd.Site.BaseURL, zd.Path)
}

/*
 * Helper template method to get any value from ZasData.config using pathes.
 */
func (zd *ZasData) Extra(keypath string) (value string, err error) {
	keypath = path.Clean(keypath)
	if path.IsAbs(keypath) {
		keypath = keypath[1:]
	}
	steps := strings.Split(keypath, "/")
	last := len(steps) - 1
	key, steps := steps[last], steps[:last]
	section := zd.config
	for _, step := range steps {
		section = section.GetSection(step)
		if section == nil {
			err = errors.New("not found")
			return
		}
	}
	value = section.GetString(key)
	return
}

func NewZasData(filepath string, config ConfigSection) (data ZasData) {
	// Any path must finish in ".html".
	if strings.HasSuffix(filepath, ".md") {
		filepath = strings.Replace(filepath, ".md", ".html", -1)
	}
	data.Path = fmt.Sprintf("/%s", filepath)
	data.config = config
	data.Site.BaseURL = config.GetSection("site").GetString("baseurl")
	data.Site.Image = config.GetSection("site").GetString("image")
	return
}
