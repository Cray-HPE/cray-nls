/*
 *
 *  MIT License
 *
 *  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
 *
 *  Permission is hereby granted, free of charge, to any person obtaining a
 *  copy of this software and associated documentation files (the "Software"),
 *  to deal in the Software without restriction, including without limitation
 *  the rights to use, copy, modify, merge, publish, distribute, sublicense,
 *  and/or sell copies of the Software, and to permit persons to whom the
 *  Software is furnished to do so, subject to the following conditions:
 *
 *  The above copyright notice and this permission notice shall be included
 *  in all copies or substantial portions of the Software.
 *
 *  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
 *  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
 *  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
 *  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
 *  OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package services_iuf

import (
	_ "embed"
	"encoding/json"
	iuf "github.com/Cray-HPE/cray-nls/src/api/models/iuf"
	"github.com/imdario/mergo"
	"golang.org/x/exp/slices"
	"path"
	"sigs.k8s.io/yaml"
)

func (s iufService) getGlobalParams(session iuf.Session, in_product iuf.Product, stages iuf.Stages) map[string]interface{} {
	return map[string]interface{}{
		"product_manifest": s.getGlobalParamsProductManifest(session, in_product),
		"input_params":     s.getGlobalParamsInputParams(session, in_product),
		"site_params":      s.getGlobalParamsSiteParams(session, in_product),
		"stage_params":     s.getGlobalParamsStageParams(session, in_product, stages),
	}
}

func (s iufService) getGlobalParamsProductManifest(session iuf.Session, in_product iuf.Product) map[string]interface{} {
	resProducts := make(map[string]interface{})
	var currentProductManifest map[string]interface{}
	for _, product := range session.Products {
		manifestBytes := []byte(product.Manifest)
		manifestJsonBytes, _ := yaml.YAMLToJSON(manifestBytes)
		var manifestJson map[string]interface{}
		json.Unmarshal(manifestJsonBytes, &manifestJson)
		if product.Name == in_product.Name {
			currentProductManifest = manifestJson
		}
		resProducts[product.Name] = map[string]interface{}{
			"manifest":          manifestJson,
			"original_location": product.OriginalLocation,
		}
	}
	return map[string]interface{}{
		"products": resProducts,
		"current_product": map[string]interface{}{
			"name":              in_product.Name,
			"manifest":          currentProductManifest,
			"original_location": in_product.OriginalLocation,
		},
	}
}

func (s iufService) getGlobalParamsInputParams(session iuf.Session, in_product iuf.Product) map[string]interface{} {
	var productsArray []string
	for _, product := range session.Products {
		productsArray = append(productsArray, product.Name)
	}

	var bootPrepManagedContent []map[string]string
	for _, bootPrepManagedItem := range session.InputParameters.BootprepConfigManaged {
		bootPrepManagedContent = append(bootPrepManagedContent, map[string]string{
			"content": bootPrepManagedItem,
		})
	}

	var bootPrepManagementContent []map[string]string
	for _, bootPrepManagementItem := range session.InputParameters.BootprepConfigManagement {
		bootPrepManagementContent = append(bootPrepManagementContent, map[string]string{
			"content": bootPrepManagementItem,
		})
	}

	return map[string]interface{}{
		"products":                   productsArray,
		"media_dir":                  path.Join(s.env.MediaDirBase, session.InputParameters.MediaDir),
		"bootprep_config_managed":    bootPrepManagedContent,
		"bootprep_config_management": bootPrepManagementContent,
		"limit_nodes":                session.InputParameters.LimitNodes,
	}
}

func (s iufService) getGlobalParamsStageParams(session iuf.Session, in_product iuf.Product, stages iuf.Stages) map[string]interface{} {
	res := make(map[string]interface{})
	activity, _ := s.GetActivity(session.ActivityRef)
	if activity.OperationOutputs == nil {
		return map[string]interface{}{}
	}
	stageParams := activity.OperationOutputs["stage_params"].(map[string]interface{})
	// loop through each stage's output
	for stageName, v := range stageParams {
		idx := slices.IndexFunc(stages.Stages, func(stage iuf.Stage) bool { return stage.Name == stageName })
		stageType := stages.Stages[idx].Type
		outputValue := v.(map[string]interface{})
		res[stageName] = make(map[string]interface{})
		s.logger.Debugf("stage: %s, type: %s, outputs: %v", stageName, stageType, v)
		if stageType == "product" || stageName == "process-media" {
			var currentProduct map[string]interface{}
			var products map[string]interface{}
			for _, value := range outputValue {
				mergo.Merge(&products, value.(map[string]interface{}))
				mergo.Merge(&currentProduct, value.(map[string]interface{})[in_product.Name])
			}
			res[stageName].(map[string]interface{})["products"] = products
			res[stageName].(map[string]interface{})["current_product"] = currentProduct
		} else {
			res[stageName].(map[string]interface{})["global"] = outputValue
		}

	}
	return res
}

// Gets the correct version of site parameters based on either session.InputParameters.SiteParameters (deprecated) or
//  session.SiteParameters
func (s iufService) getSiteParams(deprecatedSiteParameters string, structSiteParams map[string]interface{}) map[string]interface{} {
	// check which site parameters we are using.
	var siteParams map[string]interface{}
	if deprecatedSiteParameters != "" {
		// preference is given to the deprecated site params saved as a string to not break existing functionality.
		err := json.Unmarshal([]byte(deprecatedSiteParameters), siteParams)
		if err != nil {
			// fallback
			siteParams = structSiteParams
		}
	} else {
		siteParams = structSiteParams
	}

	return siteParams
}

func (s iufService) getGlobalParamsSiteParams(session iuf.Session, in_product iuf.Product) map[string]interface{} {
	siteParams := s.getSiteParams(session.InputParameters.SiteParameters, session.SiteParameters)

	if siteParams != nil && siteParams["products"] != nil {
		// need to add "current_product" in there
		products := (siteParams["products"]).(map[string]interface{})
		siteParams["current_product"] = products[in_product.Name]
	}

	return siteParams
}
