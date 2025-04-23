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
	"strings"
)

func (s iufService) getGlobalParams(session iuf.Session, in_product iuf.Product, stages iuf.Stages) map[string]interface{} {
	return map[string]interface{}{
		"product_manifest": s.getGlobalParamsProductManifest(session, in_product),
		"input_params":     s.getGlobalParamsInputParams(session, in_product),
		"site_params":      s.getGlobalParamsSiteParams(session, in_product, stages),
		"stage_params":     s.getGlobalParamsStageParams(session, in_product, stages),
	}
}

func (s iufService) getProductVersionKey(product iuf.Product) string {
	return s.getProductVersionKeyFromNameAndVersion(product.Name, product.Version)
}

func (s iufService) getProductVersionKeyFromNameAndVersion(name string, version string) string {
	return strings.ReplaceAll(name+"-"+s.normalizeProductVersion(version), ".", "-")
}

//normalizeProductVersion normalize the product version so that we force-follow semver format
func (s iufService) normalizeProductVersion(version string) string {
	return strings.ReplaceAll(version, "_", "-")
}

func (s iufService) getGlobalParamsProductManifest(session iuf.Session, in_product iuf.Product) map[string]interface{} {
	resProducts := make(map[string]interface{})
	var currentProductManifest map[string]interface{}
	for _, product := range session.Products {
		manifestBytes := []byte(product.Manifest)
		manifestJsonBytes, err := yaml.YAMLToJSON(manifestBytes)
		if err != nil {
			s.logger.Errorf("getGlobalParamsProductManifest: There was an error converting YAML to JSON for the manifest for the product %s during global param construction for the session %s in activity %s. The YAML contents: %s. Error %v", s.getProductVersionKey(in_product), session.Name, session.ActivityRef, product.Manifest, err)
			continue
		}
		var manifestJson map[string]interface{}
		err = json.Unmarshal(manifestJsonBytes, &manifestJson)
		if err != nil {
			s.logger.Errorf("getGlobalParamsProductManifest: There was an error parsing JSON for the manifest for the product %s during global param construction for the session %s in activity %s. The YAML contents: %s. Error %v", s.getProductVersionKey(in_product), session.Name, session.ActivityRef, product.Manifest, err)
			continue
		}
		if s.getProductVersionKey(product) == s.getProductVersionKey(in_product) {
			currentProductManifest = manifestJson
		}
		if manifestJson["version"] != nil {
			manifestJson["version"] = s.normalizeProductVersion(manifestJson["version"].(string))
		}
		resProducts[s.getProductVersionKey(product)] = map[string]interface{}{
			"manifest":          manifestJson,
			"original_location": product.OriginalLocation,
		}
	}
	return map[string]interface{}{
		// commenting this out as a temporary fix to CASM-3761
		// "products": resProducts,
		"current_product": map[string]interface{}{
			"name":              in_product.Name,
			"version":           in_product.Version,
			"manifest":          currentProductManifest,
			"original_location": in_product.OriginalLocation,
		},
	}
}

func (s iufService) getGlobalParamsInputParams(session iuf.Session, in_product iuf.Product) map[string]interface{} {
	var productsArray []string
	for _, product := range session.Products {
		productsArray = append(productsArray, s.getProductVersionKey(product))
	}

	return map[string]interface{}{
		"products":                                 productsArray,
		"media_dir":                                path.Join(s.env.MediaDirBase, session.InputParameters.MediaDir),
		"bootprep_config_managed":                  session.InputParameters.BootprepConfigManaged,
		"bootprep_config_management":               session.InputParameters.BootprepConfigManagement,
		"limit_management_nodes":                   session.InputParameters.LimitManagementNodes,
		"limit_managed_nodes":                      session.InputParameters.LimitManagedNodes,
		"managed_rollout_strategy":                 session.InputParameters.ManagedRolloutStrategy,
		"management_rollout_strategy":              session.InputParameters.ManagementRolloutStrategy,
		"concurrent_management_rollout_percentage": session.InputParameters.ConcurrentManagementRolloutPercentage,
		"media_host":                               session.InputParameters.MediaHost,
		"concurrency":                              session.InputParameters.Concurrency,
	}
}

func (s iufService) getGlobalParamsStageParams(session iuf.Session, in_product iuf.Product, stages iuf.Stages) map[string]interface{} {
	res := make(map[string]interface{})
	activity, _ := s.GetActivity(session.ActivityRef)
	if activity.OperationOutputs == nil || activity.OperationOutputs["stage_params"] == nil {
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
				mergo.Merge(&currentProduct, value.(map[string]interface{})[s.getProductVersionKey(in_product)])
			}
			// commenting this out as a temporary fix to CASM-3761
			//res[stageName].(map[string]interface{})["products"] = products
			res[stageName].(map[string]interface{})["current_product"] = currentProduct
		} else {
			res[stageName].(map[string]interface{})["global"] = outputValue
		}

	}
	return res
}

func (s iufService) getSiteParams(structSiteParams iuf.SiteParameters, currentStage string, stages iuf.Stages) iuf.SiteParameters {
	var siteParams iuf.SiteParameters

	idx := slices.IndexFunc(stages.Stages, func(stage iuf.Stage) bool { return stage.Name == currentStage })
	var includeDefaultProduct bool
	if idx >= 0 {
		for _, operation := range stages.Stages[idx].Operations {
			if operation.IncludeDefaultProductInSiteParams {
				includeDefaultProduct = true
				break
			}
		}
	}

	if !includeDefaultProduct {
		siteParams.Products = make(map[string]map[string]interface{})
		// remove any "default" product
		for key, value := range structSiteParams.Products {
			if strings.ToLower(key) != "default" {
				siteParams.Products[key] = value
			}
		}
	} else {
		siteParams.Products = structSiteParams.Products
	}

	siteParams.Global = structSiteParams.Global

	return siteParams
}

func (s iufService) getGlobalParamsSiteParams(session iuf.Session, in_product iuf.Product, stages iuf.Stages) iuf.SiteParametersForOperationsAndHooks {
	params := s.getSiteParams(session.SiteParameters, session.CurrentStage, stages)
	return iuf.SiteParametersForOperationsAndHooks{
		SiteParameters: params,
		// Note that we don't key by productName-productVersion here intentionally. There is only one set of configuration
		//  per product being installed.
		CurrentProduct: params.Products[in_product.Name],
	}
}
