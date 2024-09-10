// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	armPolicy "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/quota/armquota"
)

const (
	location          = "westus2"
	managementGroupId = "E7EC67B3-7657-4966-BFFC-41EFD36BAA09"
	resourceName      = "cores"
	provider          = "Microsoft.Compute"
)

const (
	moduleName    = "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/quota/armquota"
	moduleVersion = "v1.1.0-beta.1"
)

func createGroupQuota(groupName string, managementGroupId string) {
	clientFactory := createClientFactory()
	ctx := context.Background()

	groupQuotaBody := &armquota.GroupQuotasClientBeginCreateOrUpdateOptions{
		GroupQuotaPutRequestBody: &armquota.GroupQuotasEntity{
			Properties: &armquota.GroupQuotasEntityBase{
				AdditionalAttributes: &armquota.AdditionalAttributes{
					GroupID: &armquota.GroupingID{
						GroupingIDType: to.Ptr(armquota.GroupingIDTypeBillingID),
						Value:          to.Ptr("E7EC67B3-7657-4966-BFFC-41EFD36BAA09"),
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

func createGroupQuotaLimitRequest(groupName string, provider string, resourceName string, region string, limitVal int64) {
	ctx := context.Background()
	clientFactory := createClientFactory()

	groupQuotaLimitRequestBody := &armquota.GroupQuotaLimitsRequestClientBeginCreateOrUpdateOptions{
		GroupQuotaRequest: &armquota.SubmittedResourceRequestStatus{
			Properties: &armquota.SubmittedResourceRequestStatusProperties{
				RequestedResource: &armquota.GroupQuotaRequestBase{
					Properties: &armquota.GroupQuotaRequestBaseProperties{
						Limit:  &limitVal,
						Region: to.Ptr(region),
					},
				},
			},
		},
	}
	groupQuotaLimitsRequestClient := clientFactory.NewGroupQuotaLimitsRequestClient()
	poller, err := groupQuotaLimitsRequestClient.BeginCreateOrUpdate(ctx, managementGroupId, groupName, provider, resourceName, groupQuotaLimitRequestBody)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}

	// Get the HTTP response from the quota limit request
	res, err := poller.Poll(ctx)
	if err != nil {
		panic(err)
	}
	opStatus := res.Header["Location"][0]
	opStatusURI, err := url.Parse(opStatus)
	if err != nil {
		panic(err)
	}

	// Get the request ID and the Request URL in case the request doesn't reach a terminal state in 2 minutes
	opStatusURISegments := strings.Split(opStatusURI.Path, "/")
	id := opStatusURISegments[len(opStatusURISegments)-1]

	// Poll for 2 minutes until the group quota limit request has completed
	start := time.Now()
	duration := 2 * time.Minute
	provisioningState := armquota.RequestStateInProgress
	for time.Since(start) < duration {
		group_limit_request_res, err := groupQuotaLimitsRequestClient.Get(ctx, managementGroupId, groupName, id, nil)
		if err != nil {
			panic(err)
		}

		provisioningState = *group_limit_request_res.Properties.ProvisioningState
		if provisioningState == armquota.RequestStateSucceeded {
			groupQuotaLimitClient := clientFactory.NewGroupQuotaLimitsClient()
			filterString := fmt.Sprintf("location eq %s", region)
			limit_res, err := groupQuotaLimitClient.Get(ctx, managementGroupId, groupName, provider, resourceName, filterString, nil)
			if err != nil {
				log.Fatalf("failed to finish the request: %v", err)
			}
			// You could use response here. We use blank identifier for just demo purposes.
			_ = limit_res
			break
		} else if provisioningState == armquota.RequestStateFailed {
			fmt.Println("Group Quota Limit Request Failed")
			log.Fatalf("Request failed")
		}
		fmt.Println("Polling...")
		time.Sleep(30 * time.Second)
	}

	if provisioningState == armquota.RequestStateInProgress || provisioningState == "Escalated" {
		fmt.Println("Did not reach terminal state within 2 minutes. Please perform a get on this URL: %s", opStatus)
	}

	if provisioningState == "Escalated" {
		fmt.Println("Request was escalated please contact your capacity manager. Please perform a get on this URL: %s", opStatus)
	}

	// You could use response here. We use blank identifier for just demo purposes.
	//_ = limit_res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.GroupQuotaLimit = armquota.GroupQuotaLimit{
	// 	Name: to.Ptr("cores"),
	// 	Type: to.Ptr("Microsoft.Quota/groupQuotas/groupQuotaLimits"),
	// 	ID: to.Ptr("/providers/Microsoft.Management/managementGroups/E7EC67B3-7657-4966-BFFC-41EFD36BAA09/providers/Microsoft.Quota/groupQuotas/groupquota1/providers/Microsoft.Compute/locations/westus/groupQuotaLimits/cores"),
	// 	Properties: &armquota.GroupQuotaDetails{
	// 		Name: &armquota.GroupQuotaDetailsName{
	// 			LocalizedValue: to.Ptr("Total vCPUs Regional Cores"),
	// 			Value: to.Ptr("cores"),
	// 		},
	// 		AllocatedToSubscriptions: &armquota.AllocatedQuotaToSubscriptionList{
	// 			Value: []*armquota.AllocatedToSubscription{
	// 				{
	// 					QuotaAllocated: to.Ptr[int64](20),
	// 					SubscriptionID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 			}},
	// 		},
	// 		AvailableLimit: to.Ptr[int64](80),
	// 		Limit: to.Ptr[int64](100),
	// 		Region: to.Ptr("westus"),
	// 		Unit: to.Ptr("count"),
	// 	},
	// }
}

func getGroupQuotaLimit(groupName string, provider string, resourceName string, region string) {
	ctx := context.Background()
	clientFactory := createClientFactory()
	filterString := fmt.Sprintf("location eq %s", region)
	res, err := clientFactory.NewGroupQuotaLimitsClient().Get(ctx, managementGroupId, groupName, provider, resourceName, filterString, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.GroupQuotaLimit = armquota.GroupQuotaLimit{
	// 	Name: to.Ptr("cores"),
	// 	Type: to.Ptr("Microsoft.Quota/groupQuotas/groupQuotaLimits"),
	// 	ID: to.Ptr("/providers/Microsoft.Management/managementGroups/E7EC67B3-7657-4966-BFFC-41EFD36BAA09/providers/Microsoft.Quota/groupQuotas/groupquota1/providers/Microsoft.Compute/locations/westus/groupQuotaLimits/cores"),
	// 	Properties: &armquota.GroupQuotaDetails{
	// 		Name: &armquota.GroupQuotaDetailsName{
	// 			LocalizedValue: to.Ptr("Total vCPUs Regional Cores"),
	// 			Value: to.Ptr("cores"),
	// 		},
	// 		AllocatedToSubscriptions: &armquota.AllocatedQuotaToSubscriptionList{
	// 			Value: []*armquota.AllocatedToSubscription{
	// 				{
	// 					QuotaAllocated: to.Ptr[int64](20),
	// 					SubscriptionID: to.Ptr("00000000-0000-0000-0000-000000000000"),
	// 			}},
	// 		},
	// 		AvailableLimit: to.Ptr[int64](80),
	// 		Limit: to.Ptr[int64](100),
	// 		Region: to.Ptr("westus"),
	// 		Unit: to.Ptr("count"),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/106483d9f698ac3b6c0d481ab0c5fab14152e21f/specification/quota/resource-manager/Microsoft.Quota/preview/2023-06-01-preview/examples/SubscriptionQuotaAllocationRequests/PutSubscriptionQuotaAllocationRequest-Compute.json
func createSubscriptionAllocationRequest(groupName string, provider string, resourceName string, region string, allocVal int64) {
	ctx := context.Background()
	clientFactory := createClientFactory()
	subscriptionQuotaAllocationRequestClient := clientFactory.NewGroupQuotaSubscriptionAllocationRequestClient()
	poller, err := subscriptionQuotaAllocationRequestClient.BeginCreateOrUpdate(ctx, managementGroupId, groupName, provider, resourceName, armquota.AllocationRequestStatus{
		Properties: &armquota.AllocationRequestStatusProperties{
			RequestedResource: &armquota.AllocationRequestBase{
				Properties: &armquota.AllocationRequestBaseProperties{
					Limit:  to.Ptr[int64](allocVal),
					Region: to.Ptr(region),
				},
			},
		},
	}, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	res, err := poller.Poll(ctx)
	if err != nil {
		panic(err)
	}
	opStatus := res.Header["Location"][0]
	opStatusURI, err := url.Parse(opStatus)
	if err != nil {
		panic(err)
	}

	// Get the request ID and the Request URL in case the request doesn't reach a terminal state in 2 minutes
	opStatusURISegments := strings.Split(opStatusURI.Path, "/")
	id := opStatusURISegments[len(opStatusURISegments)-1]

	// Poll for 2 minutes until the group quota limit request has completed
	start := time.Now()
	duration := 2 * time.Minute
	provisioningState := armquota.RequestStateInProgress
	for time.Since(start) < duration {
		allocation_res, err := subscriptionQuotaAllocationRequestClient.Get(ctx, managementGroupId, groupName, id, nil)
		if err != nil {
			panic(err)
		}

		provisioningState = *allocation_res.Properties.ProvisioningState
		if provisioningState == armquota.RequestStateSucceeded {
			subscriptionQuotaAllocationClient := clientFactory.NewGroupQuotaSubscriptionAllocationClient()
			filterString := fmt.Sprintf("location eq %s", region)
			allocation_res, err := subscriptionQuotaAllocationClient.Get(ctx, managementGroupId, provider, resourceName, filterString, nil)
			if err != nil {
				log.Fatalf("failed to finish the request: %v", err)
			}
			// You could use response here. We use blank identifier for just demo purposes.
			_ = allocation_res
			break
		} else if provisioningState == armquota.RequestStateFailed {
			fmt.Println("Group Quota Subscription Allocation Failed")
			log.Fatalf("Request failed")
		}
		fmt.Println("Polling...")
		time.Sleep(30 * time.Second)
	}

	if provisioningState == armquota.RequestStateInProgress || provisioningState == armquota.RequestStateAccepted {
		fmt.Println("Did not reach terminal state within 2 minutes. Please perform a get on this URL: %s", opStatus)
	}

	if provisioningState == "Escalated" {
		fmt.Println("Request was escalated please contact your capacity manager. Please perform a get on this URL: %s", opStatus)
	}

	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SubscriptionQuotaAllocations = armquota.SubscriptionQuotaAllocations{
	// 	Name: to.Ptr("standardav2family"),
	// 	Type: to.Ptr("Microsoft.Quota/groupQuotas/quotaAllocations"),
	// 	ID: to.Ptr("/providers/Microsoft.Management/managementGroups/E7EC67B3-7657-4966-BFFC-41EFD36BAA09/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Quota/groupQuotas/groupquota1/providers/Microsoft.Compute/locations/westus/quotaAllocations/standardav2family"),
	// 	Properties: &armquota.SubscriptionQuotaDetails{
	// 		Name: &armquota.SubscriptionQuotaDetailsName{
	// 			LocalizedValue: to.Ptr("standard Av2 Family vCPUs"),
	// 			Value: to.Ptr("standardav2family"),
	// 		},
	// 		Limit: to.Ptr[int64](100),
	// 		Region: to.Ptr("westus"),
	// 		ShareableQuota: to.Ptr[int64](25),
	// 	},
	// }
}

// Generated from example definition: https://github.com/Azure/azure-rest-api-specs/blob/106483d9f698ac3b6c0d481ab0c5fab14152e21f/specification/quota/resource-manager/Microsoft.Quota/preview/2023-06-01-preview/examples/SubscriptionQuotaAllocation/SubscriptionQuotaAllocation_Get-Compute.json
func getSubscriptionQuotaAllocation(groupName string, resourceName string, region string) {
	ctx := context.Background()
	clientFactory := createClientFactory()
	filterString := fmt.Sprintf("location eq %s", region)
	res, err := clientFactory.NewGroupQuotaSubscriptionAllocationClient().Get(ctx, managementGroupId, groupName, resourceName, filterString, nil)
	if err != nil {
		log.Fatalf("failed to finish the request: %v", err)
	}
	// You could use response here. We use blank identifier for just demo purposes.
	_ = res
	// If the HTTP response code is 200 as defined in example definition, your response structure would look as follows. Please pay attention that all the values in the output are fake values for just demo purposes.
	// res.SubscriptionQuotaAllocations = armquota.SubscriptionQuotaAllocations{
	// 	Name: to.Ptr("standardav2family"),
	// 	Type: to.Ptr("Microsoft.Quota/groupQuotas/quotaAllocations"),
	// 	ID: to.Ptr("/providers/Microsoft.Management/managementGroups/E7EC67B3-7657-4966-BFFC-41EFD36BAA09/subscriptions/00000000-0000-0000-0000-000000000000/providers/Microsoft.Quota/groupQuotas/groupquota1/providers/Microsoft.Compute/locations/westus/quotaAllocations/standardav2family"),
	// 	Properties: &armquota.SubscriptionQuotaDetails{
	// 		Name: &armquota.SubscriptionQuotaDetailsName{
	// 			LocalizedValue: to.Ptr("standard Av2 Family vCPUs"),
	// 			Value: to.Ptr("standardav2family"),
	// 		},
	// 		Limit: to.Ptr[int64](100),
	// 		Region: to.Ptr("westus"),
	// 		ShareableQuota: to.Ptr[int64](25),
	// 	},
	// }
}

func main() {
	//Group Quota Functions
	createGroupQuota("groupquota1", managementGroupId)
	getGroupQuota("groupquota1")
	deleteGroupQuota("groupquota1")

	//Subscription Functions
	addSubscription("groupquota1")
	deleteSubscription("groupquota1")

	//GroupQuotaLimit
	createGroupQuotaLimitRequest("groupquota1", provider, resourceName, location, 64)
	getGroupQuotaLimit("groupquota1", provider, resourceName, location)

	//SubscriptionQuotaAllocation
	createSubscriptionAllocationRequest("groupquota1", provider, resourceName, location, 10)
	getSubscriptionQuotaAllocation("groupquota1", resourceName, location)

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
