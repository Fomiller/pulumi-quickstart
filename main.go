package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create state bucket
		// Bucket was created manually and then imported using pulumi cli
		// imported resources are marked as protected, remove pulumi.Protect(true) to destroy resource
		_, err := s3.NewBucket(ctx, "fomiller-pulumi-state", &s3.BucketArgs{
			Arn:          pulumi.String("arn:aws:s3:::fomiller-pulumi-state"),
			Bucket:       pulumi.String("fomiller-pulumi-state"),
			HostedZoneId: pulumi.String("Z3AQBSTGFYJSTF"),
			RequestPayer: pulumi.String("BucketOwner"),
		}, pulumi.Protect(true))
		if err != nil {
			return err
		}

		// Create an AWS resource (S3 Bucket)
		bucket, err := s3.NewBucket(ctx, "my-pulumi-bucket", &s3.BucketArgs{
			Bucket: pulumi.String("my-pulumi-bucket"),
			Website: s3.BucketWebsiteArgs{
				IndexDocument: pulumi.String("index.html"),
			},
		})

		// Create bucket object and serve as static web page
		_, err = s3.NewBucketObject(ctx, "index.html", &s3.BucketObjectArgs{
			Acl:         pulumi.String("public-read"),
			ContentType: pulumi.String("text/html"),
			Bucket:      bucket.ID(),
			Source:      pulumi.NewFileAsset("index.html"),
		})

		if err != nil {
			return err
		}

		// Export the name of the bucket
		ctx.Export("bucketName", bucket.ID())
		ctx.Export("bucketEndpoint", pulumi.Sprintf("http://%s", bucket.WebsiteEndpoint))
		return nil
	})
}
