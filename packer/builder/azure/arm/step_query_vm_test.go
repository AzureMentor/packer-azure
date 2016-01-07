// Copyright (c) Microsoft Open Technologies, Inc.
// All Rights Reserved.
// Licensed under the Apache License, Version 2.0.
// See License.txt in the project root for license information.

package arm

import (
	"fmt"
	"testing"

	"github.com/Azure/azure-sdk-for-go/arm/compute"

	"github.com/Azure/packer-azure/packer/builder/azure/common/constants"

	"github.com/mitchellh/multistep"
)

func TestStepQueryVMShouldFailIfValidateFails(t *testing.T) {
	var testSubject = &StepQueryVM{
		query: func(string, string) (compute.VirtualMachine, error) {
			return createVirtualMachineFromUri("test.vhd"), fmt.Errorf("!! Unit Test FAIL !!")
		},
		say:   func(message string) {},
		error: func(e error) {},
	}

	stateBag := createTestStateBagStepQueryVM()

	var result = testSubject.Run(stateBag)
	if result != multistep.ActionHalt {
		t.Fatalf("Expected the step to return 'ActionHalt', but got '%s'.", result)
	}

	if _, ok := stateBag.GetOk(constants.Error); ok == false {
		t.Fatalf("Expected the step to set stateBag['%s'], but it was not.", constants.Error)
	}
}

func TestStepQueryVMShouldPassIfValidatePasses(t *testing.T) {
	var testSubject = &StepQueryVM{
		query: func(string, string) (compute.VirtualMachine, error) {
			return createVirtualMachineFromUri("test.vhd"), nil
		},
		say:   func(message string) {},
		error: func(e error) {},
	}

	stateBag := createTestStateBagStepQueryVM()

	var result = testSubject.Run(stateBag)
	if result != multistep.ActionContinue {
		t.Fatalf("Expected the step to return 'ActionContinue', but got '%s'.", result)
	}

	if _, ok := stateBag.GetOk(constants.Error); ok == true {
		t.Fatalf("Expected the step to not set stateBag['%s'], but it was.", constants.Error)
	}
}

func TestStepQueryVMShouldTakeValidateArgumentsFromStateBag(t *testing.T) {
	var actualResourceGroupName string
	var actualComputeName string

	var testSubject = &StepQueryVM{
		query: func(resourceGroupName string, computeName string) (compute.VirtualMachine, error) {
			actualResourceGroupName = resourceGroupName
			actualComputeName = computeName

			return createVirtualMachineFromUri("test.vhd"), nil
		},
		say:   func(message string) {},
		error: func(e error) {},
	}

	stateBag := createTestStateBagStepQueryVM()
	var result = testSubject.Run(stateBag)

	if result != multistep.ActionContinue {
		t.Fatalf("Expected the step to return 'ActionContinue', but got '%s'.", result)
	}

	var expectedComputeName = stateBag.Get(constants.ArmComputeName).(string)
	var expectedResourceGroupName = stateBag.Get(constants.ArmResourceGroupName).(string)

	if actualComputeName != expectedComputeName {
		t.Fatalf("Expected the step to source 'constants.ArmResourceGroupName' from the state bag, but it did not.")
	}

	if actualResourceGroupName != expectedResourceGroupName {
		t.Fatalf("Expected the step to source 'constants.ArmResourceGroupName' from the state bag, but it did not.")
	}

	expectedOSDiskVhd, ok := stateBag.GetOk(constants.ArmOSDiskVhd)
	if !ok {
		t.Fatalf("Expected the state bag to have a value for '%s', but it did not.", constants.ArmOSDiskVhd)
	}

	if expectedOSDiskVhd != "test.vhd" {
		t.Fatalf("Expected the value of stateBag[%s] to be '127.0.0.1', but got '%s'.", constants.ArmOSDiskVhd, expectedOSDiskVhd)
	}
}

func createTestStateBagStepQueryVM() multistep.StateBag {
	stateBag := new(multistep.BasicStateBag)

	stateBag.Put(constants.ArmComputeName, "Unit Test: ComputeName")
	stateBag.Put(constants.ArmResourceGroupName, "Unit Test: ResourceGroupName")

	return stateBag
}

func createVirtualMachineFromUri(vhdUri string) compute.VirtualMachine {
	vm := compute.VirtualMachine{
		Properties: &compute.VirtualMachineProperties{
			StorageProfile: &compute.StorageProfile{
				OsDisk: &compute.OSDisk{
					Vhd: &compute.VirtualHardDisk{
						URI: &vhdUri,
					},
				},
			},
		},
	}

	return vm
}
