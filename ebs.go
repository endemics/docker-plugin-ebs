package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (e *EC2Wrapper) find(label string) (string, error) {
	params := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			{
				Name: aws.String("tag:DockerVolumeName"),
				Values: []*string{
					aws.String(label),
				},
			},
		},
	}

	resp, err := e.ec2.DescribeVolumes(params)

	if err != nil {
		return err.Error(), err
	}

	if len(resp.Volumes) == 0 {
		return "", nil
	}

	if len(resp.Volumes) > 1 {
		err := fmt.Errorf("Unable to identify EBS volume with tag DockerVolumeName=%s, more than one volume matches\n", label)
		return "", err
	}

	return *resp.Volumes[0].VolumeId, nil
}

func (e *EC2Wrapper) create(label string) (string, error) {
	params := &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String("eu-west-1a"),
		Size:             aws.Int64(1),
		VolumeType:       aws.String("gp2"),
	}

	resp, err := e.ec2.CreateVolume(params)

	if err != nil {
		return "", err
	}

	return *resp.VolumeId, nil
}

func (e *EC2Wrapper) tag(volumeId string, label string) error {
	params := &ec2.CreateTagsInput{
		Resources: []*string{
			aws.String(volumeId),
		},
		Tags: []*ec2.Tag{
			{
				Key:   aws.String("Name"),
				Value: aws.String(label),
			},
			{
				Key:   aws.String("DockerVolumeName"),
				Value: aws.String(label),
			},
		},
	}
	_, err := e.ec2.CreateTags(params)

	return err
}
