// Copyright (c) Microsoft Open Technologies, Inc.
// All Rights Reserved.
// Licensed under the Apache License, Version 2.0.
// See License.txt in the project root for license information.

package request_tests

import (
	"testing"
	"github.com/MSOpenTech/packer-azure/packer/builder/azure/driver_restapi/response"
	"fmt"
)

func TestGetOsImages(t *testing.T) {

	errMassage := "GetOsImages: %s\n"

	reqManager, err := getRequestManager()
	if err != nil {
		t.Errorf(errMassage, err.Error())
	}

	requestData := reqManager.GetOsImages()
	resp, err := reqManager.Execute(requestData)

	if err != nil {
		t.Errorf(errMassage, err.Error())
	}

	list, err := response.ParseOsImageList(resp.Body)

	if err != nil {
		t.Errorf(errMassage, err.Error())
	}

	fmt.Printf("osImageList:\n\n")

//	for _, image := range(list.OSImages){
//		fmt.Printf("label: '%s'\nfamily: '%s'\nlocations: '%s'\nPublishedDate: '%s'\n\n", image.Label, image.ImageFamily, image.Location, image.PublishedDate)
//	}

//	label := "Ubuntu Server 12.04 LTS"
	label := "Windows Server 2012 R2 Datacenter"
//	label := "SUSE Linux Enterprise Server 11 SP3 (Premium Image)"
	location := "West US"

	filteredImageList := list.Filter(label, location)
	list.SortByDateDesc(filteredImageList)
	fmt.Printf("Filtered and Sorted ----------------------------------:\n\n")

	for _, image := range(filteredImageList){
		fmt.Printf("label: '%s'\nfamily: '%s'\nlocations: '%s'\nPublishedDate: '%s'\n\n", image.Label, image.ImageFamily, image.Location, image.PublishedDate)
	}

	fmt.Println("Os Image Label: " + filteredImageList[0].Label)
	fmt.Println("Os Image Family: " + filteredImageList[0].ImageFamily)



	t.Error("eom")
}
