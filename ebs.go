package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (e *Ec2Wrapper) find(label string) (string, error) {
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
