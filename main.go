package main

import (
	"github.com/pulumi/pulumi-gcp/sdk/v6/go/gcp/compute"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Allocate static IP address
		static, err := compute.NewAddress(ctx, "static", nil)
		if err != nil {
			return err
		}

		// Create a compute VM with the SA to access DynamoDB
		instance, err := compute.NewInstance(ctx, "test", &compute.InstanceArgs{
			AllowStoppingForUpdate: pulumi.Bool(true),
			MachineType: pulumi.String("e2-micro"),
			Zone:        pulumi.String("us-east5-a"),
			Tags: pulumi.StringArray{
				pulumi.String("foo"),
				pulumi.String("bar"),
			},
			BootDisk: &compute.InstanceBootDiskArgs{
				InitializeParams: &compute.InstanceBootDiskInitializeParamsArgs{
					Image: pulumi.String("debian-cloud/debian-11"),
				},
			},

			NetworkInterfaces: compute.InstanceNetworkInterfaceArray{
				&compute.InstanceNetworkInterfaceArgs{
					Network: pulumi.String("default"),
					AccessConfigs: compute.InstanceNetworkInterfaceAccessConfigArray{
						&compute.InstanceNetworkInterfaceAccessConfigArgs{
							NatIp: static.Address,
						},
					},
				},
			},

			MetadataStartupScript: pulumi.String("echo hi > /test.txt"),
			
			ServiceAccount: &compute.InstanceServiceAccountArgs{
				Email: pulumi.String("test-aws@tokyo-ring-351304.iam.gserviceaccount.com"),
				Scopes: pulumi.StringArray{
					pulumi.String("cloud-platform"),
				},
			},
		})
		if err != nil {
			return err
		}
		ctx.Export("Instance Name", instance.Name)

		return nil
	})
}
