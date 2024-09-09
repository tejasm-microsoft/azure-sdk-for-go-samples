// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	armPolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/quota/armquota"
)

var subscriptionId string

const (
	resourceGroupName = "sample-resource-group"
	vmName            = "sample-vm"
	vnetName          = "sample-vnet"
	subnetName        = "sample-subnet"
	nsgName           = "sample-nsg"
	nicName           = "sample-nic"
	diskName          = "sample-disk"
	publicIPName      = "sample-public-ip"
	location          = "westus2"
	managementGroupId = "testMgIdRoot"
)

const (
	moduleName    = "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/quota/armquota"
	moduleVersion = "v1.1.0-beta.1"
)

func createGroupQuota(groupName string) {
	clientFactory := createClientFactory()
	ctx := context.Background()

	groupQuotaBody := &armquota.GroupQuotasClientBeginCreateOrUpdateOptions{
		GroupQuotaPutRequestBody: &armquota.GroupQuotasEntity{
			Properties: &armquota.GroupQuotasEntityBase{
				AdditionalAttributes: &armquota.AdditionalAttributes{
					GroupID: &armquota.GroupingID{
						GroupingIDType: to.Ptr(armquota.GroupingIDTypeBillingID),
						Value:          to.Ptr("f1d1800e-d38e-41f2-b63c-72d59ecaf9c0"),
					},
				},
				DisplayName: to.Ptr(groupName),
			},
		},
	}

	poller, err := clientFactory.NewGroupQuotasClient().BeginCreateOrUpdate(ctx, managementGroupId, groupName, groupQuotaBody)

	if err != nil {
		log.Fatalf("failed to create GroupQuota: %v", err)
	}

	fmt.Println("waiting for GroupQuota creation to complete")

	res, err := poller.PollUntilDone(ctx, nil)

	if err != nil {
		log.Fatalf("failed to poll the result of the request: %v", err)
	}

	_ = res
	//return res
}

func getGroupQuota(groupName string) {
	clientFactory := createClientFactory()
	ctx := context.Background()
	res, err := clientFactory.NewGroupQuotasClient().Get(ctx, managementGroupId, groupName, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.GroupQuotasEntity = armquota.GroupQuotasEntity{
	// 	Name: to.Ptr("groupquota1"),
	// 	Type: to.Ptr("Microsoft.Quota/groupQuotas"),
	// 	ID: to.Ptr("/providers/Microsoft.Management/managementGroups/E7EC67B3-7657-4966-BFFC-41EFD36BAA09/providers/Microsoft.Quota/groupQuotas/groupquota1"),
	// 	Properties: &armquota.GroupQuotasEntityBase{
	// 		AdditionalAttributes: &armquota.AdditionalAttributes{
	// 			Environment: to.Ptr(armquota.EnvironmentTypeProduction),
	// 			GroupID: &armquota.GroupingID{
	// 				GroupingIDType: to.Ptr(armquota.GroupingIDTypeServiceTreeID),
	// 				Value: to.Ptr("yourServiceTreeIdHere"),
	// 			},
	// 		},
	// 		DisplayName: to.Ptr("GroupQuota1"),
	// 		ProvisioningState: to.Ptr(armquota.RequestStateSucceeded),
	// 	},
	// }
}

func deleteGroupQuota(groupName string) {
	ctx := context.Background()
	clientFactory := createClientFactory()

	poller, err := clientFactory.NewGroupQuotasClient().BeginDelete(ctx, managementGroupId, groupName, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

func addSubscription(groupName string) {
	ctx := context.Background()
	clientFactory := createClientFactory()

	poller, err := clientFactory.NewGroupQuotaSubscriptionsClient().BeginCreateOrUpdate(ctx, managementGroupId, groupName, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.GroupQuotaSubscriptionID = armquota.GroupQuotaSubscriptionID{
	// 	Name: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	Type: to.Ptr("Microsoft.Quota/groupQuotas/subscriptions"),
	// 	ID: to.Ptr("/providers/Microsoft.Management/managementGroups/E7EC67B3-7657-4966-BFFC-41EFD36BAA09/providers/Microsoft.Quota/groupQuotas/groupquota1/subscriptions/00000000-0000-0000-0000-000000000000"),
	// 	Properties: &armquota.GroupQuotaSubscriptionIDProperties{
	// 		ProvisioningState: to.Ptr(armquota.RequestStateSucceeded),
	// 		SubscriptionID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 	},
	// }
}

func deleteSubscription(groupName string) {
	ctx := context.Background()
	clientFactory := createClientFactory()

	poller, err := clientFactory.NewGroupQuotaSubscriptionsClient().BeginDelete(ctx, managementGroupId, groupName, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	_, err = poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
}

func createGroupQuotaLimitRequest(groupName string) {
	ctx := context.Background()
	clientFactory := createClientFactory()

	var limitVal int64 = 64

	groupQuotaLimitRequestBody := &armquota.GroupQuotaLimitsRequestClientBeginCreateOrUpdateOptions{
		GroupQuotaRequest: &armquota.SubmittedResourceRequestStatus{
			Properties: &armquota.SubmittedResourceRequestStatusProperties{
				RequestedResource: &armquota.GroupQuotaRequestBase{
					Properties: &armquota.GroupQuotaRequestBaseProperties{
						Limit:  &limitVal,
						Region: to.Ptr("westus"),
					},
				},
			},
		},
	}

	poller, err := clientFactory.NewGroupQuotaLimitsRequestClient().BeginCreateOrUpdate(ctx, managementGroupId, groupName, "Microsoft.Compute", "cores", groupQuotaLimitRequestBody)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.PollUntilDone(ctx, nil)
	if err != nil {
		log.Fatalf("failed to pull the result: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SubmittedResourceRequestStatus = armquota.SubmittedResourceRequestStatus{
	// 	Name: to.Ptr("requestId1"),
	// 	Type: to.Ptr("Microsoft.Quota/groupQuotas/groupQuotaLimitsRequests"),
	// 	ID: to.Ptr("/providers/Microsoft.Management/managementGroups/E7EC67B3-7657-4966-BFFC-41EFD36BAA09/providers/Microsoft.Quota/groupQuotas/groupquota1/resourceProviders/Microsoft.Compute/groupQuotaLimitsRequests/requestId1"),
	// 	Properties: &armquota.SubmittedResourceRequestStatusProperties{
	// 		ProvisioningState: to.Ptr(armquota.RequestStateSucceeded),
	// 		RequestSubmitTime: to.Ptr(func() time.Time { t, _ := time.Parse(time.RFC3339Nano, "2024-03-08T12:09:27.978Z"); return t}()),
	// 		RequestedResource: &armquota.GroupQuotaRequestBase{
	// 			Properties: &armquota.GroupQuotaRequestBaseProperties{
	// 				Name: &armquota.GroupQuotaRequestBasePropertiesName{
	// 					LocalizedValue: to.Ptr("Standard AV2 Family vCPUs"),
	// 					Value: to.Ptr("standardav2family"),
	// 				},
	// 				Comments: to.Ptr("Contoso requires more quota."),
	// 				Limit: to.Ptr[int64](100),
	// 				Region: to.Ptr("westus"),
	// 			},
	// 		},
	// 	},
	// }
}

func main() {
	//createGroupQuota("test-sdk-tejas-go", "testMgIdRoot")
	//getGroupQuota("test-sdk-tejas-go")
	//deleteGroupQuota("test-sdk-tejas-go-0")

	//Subscription Functions
	addSubscription("test-sdk-tejas-go")
	//deleteSubscription("test-sdk-tejas-go")

	//GroupQuotaLimit Request
	//createGroupQuotaLimitRequest("test-sdk-tejas-go")
}

func createClientFactory() *armquota.ClientFactory {
	//TODO: wrap in if statement for only corp tenant and only CLI
	credOptions := azidentity.DefaultAzureCredentialOptions{
		AdditionallyAllowedTenants: []string{"7be17286-b3ad-482d-ab84-4c772eebd529"},
	}

	cred, err := azidentity.NewDefaultAzureCredential(&credOptions)

	if err != nil {
		log.Fatalf("failed to obtain a credential: %v\n", err)
	}

	clientFactory, err := armquota.NewClientFactory("65a85478-2333-4bbd-981b-1a818c944faf", cred, getArmClientOptions())

	if err != nil {
		log.Fatalf("failed to create client factory: %v", err)
	}

	return clientFactory
}

func getArmClientOptions() *armPolicy.ClientOptions {
	return &armPolicy.ClientOptions{
		ClientOptions: azcore.ClientOptions{
			Cloud: cloud.Configuration{
				ActiveDirectoryAuthorityHost: "https://login.microsoft.com/",
				Services: map[cloud.ServiceName]cloud.ServiceConfiguration{
					cloud.ResourceManager: {
						Audience: "https://management.azure.com/",
						Endpoint: "https://management.azure.com/",
					},
				},
			},
		},
		AuxiliaryTenants:      []string{},
		DisableRPRegistration: false,
	}
}
