package main

import "github.com/aws/aws-sdk-go/service/ec2"

type Ec2er interface {
	DescribeVolumes(*ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error)
	CreateVolume(*ec2.CreateVolumeInput) (*ec2.Volume, error)
}

var _ Ec2er = (*ec2.EC2)(nil)

type Ec2Wrapper struct {
	ec2 Ec2er
}
