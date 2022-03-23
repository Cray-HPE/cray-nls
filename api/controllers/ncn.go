//
//  MIT License
//
//  (C) Copyright 2022 Hewlett Packard Enterprise Development LP
//
//  Permission is hereby granted, free of charge, to any person obtaining a
//  copy of this software and associated documentation files (the "Software"),
//  to deal in the Software without restriction, including without limitation
//  the rights to use, copy, modify, merge, publish, distribute, sublicense,
//  and/or sell copies of the Software, and to permit persons to whom the
//  Software is furnished to do so, subject to the following conditions:
//
//  The above copyright notice and this permission notice shall be included
//  in all copies or substantial portions of the Software.
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//  IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//  FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL
//  THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
//  OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
//  ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
//  OTHER DEALINGS IN THE SOFTWARE.
//
package controllers

import (
	"github.com/Cray-HPE/cray-nls/services"
	"github.com/Cray-HPE/cray-nls/utils"
	"github.com/gin-gonic/gin"
)

// NcnController data type
type NcnController struct {
	service services.NcnService
	logger  utils.Logger
}

// NewNcnController creates new Ncn controller
func NewNcnController(NcnService services.NcnService, logger utils.Logger) NcnController {
	return NcnController{
		service: NcnService,
		logger:  logger,
	}
}

// NcnCreateBakcup 	perform backup action on a NCN
// @Summary               Create a NCN backup
// @description.markdown  ncn-create-backup
// @Tags                  NCN
// @Param                 hostname        path  string                 true  "Hostname"
// @Accept    json
// @Produce   json
// @Failure               400  {object}  utils.ResponseError
// @Failure               401  {object}  utils.ResponseError
// @Failure               403  {object}  utils.ResponseError
// @Failure               404  {object}  utils.ResponseError
// @Failure               500  {object}  utils.ResponseError
// @Router                /ncn/{hostname}/backup [post]
// @Security  OAuth2Application[admin]
func (u NcnController) NcnCreateBakcup(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnRestoreBakcup 	restore a NCN backup
// @Summary               Restore a NCN backup
// @description.markdown  ncn-restore-backup
// @Tags                  NCN
// @Param                 hostname  path  string  true  "Hostname"
// @Accept    json
// @Produce   json
// @Failure               400  {object}  utils.ResponseError
// @Failure               401  {object}  utils.ResponseError
// @Failure               403  {object}  utils.ResponseError
// @Failure               404  {object}  utils.ResponseError
// @Failure               500  {object}  utils.ResponseError
// @Router                /ncn/{hostname}/restore [post]
// @Security  OAuth2Application[admin]
func (u NcnController) NcnRestoreBakcup(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnWipe 		perform disk wipe on a NCN
// @Summary               Perform disk wipe on a NCN
// @description.markdown  ncn-wipe-disk
// @Tags                  NCN
// @Param                 hostname  path  string  true  "Hostname"
// @Accept    json
// @Produce   json
// @Failure               400  {object}  utils.ResponseError
// @Failure               401  {object}  utils.ResponseError
// @Failure               403  {object}  utils.ResponseError
// @Failure               404  {object}  utils.ResponseError
// @Failure               500  {object}  utils.ResponseError
// @Router                /ncn/{hostname}/wipe [post]
// @Security              OAuth2Application[admin]
func (u NcnController) NcnWipe(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnSetBootParam 		set boot parameters before reboot a NCN
// @Summary               Set boot parameters before reboot a NCN
// @description.markdown  ncn-set-boot-parameters
// @Tags                  NCN
// @Param                 hostname  path  string  true  "Hostname"
// @Param                 bootParameters  body  models.BootParameters  true  "TODO: use data model from `csi/bss`"
// @Accept    json
// @Produce   json
// @Failure               400  {object}  utils.ResponseError
// @Failure               401  {object}  utils.ResponseError
// @Failure               403  {object}  utils.ResponseError
// @Failure               404  {object}  utils.ResponseError
// @Failure               500  {object}  utils.ResponseError
// @Router                /ncn/{hostname}/boot-parameters [put]
// @Security              OAuth2Application[admin]
func (u NcnController) NcnSetBootParam(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnReboot 		perform reboot on a NCN
// @Summary               Perform reboot on a NCN
// @description.markdown  ncn-reboot
// @Tags                  NCN
// @Param                 hostname  path  string  true  "Hostname"
// @Accept    json
// @Produce   json
// @Failure               400  {object}  utils.ResponseError
// @Failure               401  {object}  utils.ResponseError
// @Failure               403  {object}  utils.ResponseError
// @Failure               404  {object}  utils.ResponseError
// @Failure               500  {object}  utils.ResponseError
// @Router                /ncn/{hostname}/reboot [post]
// @Security              OAuth2Application[admin]
func (u NcnController) NcnReboot(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnPostRebuild 		perform post rebuild action on a NCN
// @Summary               Perform post rebuild action on a NCN
// @description.markdown  ncn-post-rebuild
// @Tags                  NCN
// @Param                 hostname  path  string  true  "Hostname"
// @Accept    json
// @Produce   json
// @Failure               400  {object}  utils.ResponseError
// @Failure               401  {object}  utils.ResponseError
// @Failure               403  {object}  utils.ResponseError
// @Failure               404  {object}  utils.ResponseError
// @Failure               500  {object}  utils.ResponseError
// @Router                /ncn/{hostname}/post-rebuild [post]
// @Security              OAuth2Application[admin]
func (u NcnController) NcnPostRebuild(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnValidate 		perform validation on a NCN
// @Summary               Perform validation on a NCN
// @description.markdown  ncn-validate
// @Tags                  NCN
// @Param                 hostname  path  string  true  "Hostname"
// @Accept                json
// @Produce               json
// @Failure               400  {object}  utils.ResponseError
// @Failure               401  {object}  utils.ResponseError
// @Failure               403  {object}  utils.ResponseError
// @Failure               404  {object}  utils.ResponseError
// @Failure               500  {object}  utils.ResponseError
// @Router                /ncn/{hostname}/validate [post]
// @Security  OAuth2Application[admin,read]
func (u NcnController) NcnValidate(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnPostUpgrade 		perform post upgrade actions
// @Summary               Perform post upgrade actions
// @description.markdown  ncn-post-upgrade
// @Tags                  NCN
// @Param                 type  path  string  true  "Type of ncn"
// @Accept                json
// @Produce               json
// @Failure               400  {object}  utils.ResponseError
// @Failure               401  {object}  utils.ResponseError
// @Failure               403  {object}  utils.ResponseError
// @Failure               404  {object}  utils.ResponseError
// @Failure               500  {object}  utils.ResponseError
// @Router                /ncn/{type}/post-upgrade [post]
// @Security              OAuth2Application[admin]
func (u NcnController) NcnPostUpgrade(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnCreateRebuildJob		create a rebuild request
// @Summary   Perform post upgrade actions
// @Tags      V2 APIs (draft)
// @Accept                json
// @Produce               json
// @Failure   501  "Not Implemented"
// @Router    /ncn/rebuild [post]
// @Security              OAuth2Application[admin]
func (u NcnController) NcnCreateRebuildRequest(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnGetRebuildJob
// @Summary   Get status of a rebuild job
// @Param     rebuild_job_id  path  string  true  "job id"
// @Tags      V2 APIs (draft)
// @Accept                json
// @Produce               json
// @Failure   501  "Not Implemented"
// @Router    /ncn/rebuild/{rebuild_job_id} [get]
// @Security  OAuth2Application[admin,read]
func (u NcnController) NcnGetRebuildRequest(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnDeleteRebuildJob
// @Summary   Delete a rebuild job
// @Param     rebuild_job_id  path  string  true  "job id"
// @Tags      V2 APIs (draft)
// @Accept                json
// @Produce               json
// @Failure   501  "Not Implemented"
// @Router    /ncn/rebuild/{rebuild_job_id} [delete]
// @Security  OAuth2Application[admin,read]
func (u NcnController) NcnDeleteRebuildRequest(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnCreateRebootJob		create a reboot request
// @Summary   Perform post upgrade actions
// @Tags      V2 APIs (draft)
// @Accept                json
// @Produce               json
// @Failure   501  "Not Implemented"
// @Router    /ncn/reboot [post]
// @Security              OAuth2Application[admin]
func (u NcnController) NcnCreateRebootRequest(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnGetRebootJob
// @Summary   Get status of a reboot job
// @Param     reboot_job_id  path  string  true  "job id"
// @Tags      V2 APIs (draft)
// @Accept                json
// @Produce               json
// @Failure   501  "Not Implemented"
// @Router    /ncn/reboot/{reboot_job_id} [get]
// @Security  OAuth2Application[admin,read]
func (u NcnController) NcnGetRebootRequest(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}

// NcnDeleteRebootJob
// @Summary   Delete a reboot job
// @Param     reboot_job_id  path  string  true  "job id"
// @Tags      V2 APIs (draft)
// @Accept                json
// @Produce               json
// @Failure   501  "Not Implemented"
// @Router    /ncn/reboot/{reboot_job_id} [delete]
// @Security              OAuth2Application[admin,read]
func (u NcnController) NcnDeleteRebootRequest(c *gin.Context) {
	c.JSON(200, gin.H{"data": "Ncn updated"})
}
