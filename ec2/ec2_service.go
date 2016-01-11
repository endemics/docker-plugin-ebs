package ec2

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"strconv"
)

const (
	defaultAZ               = "ap-southeast-2a"
	defaultVolumeType       = "gp2" // "standard", "io1", "gp2"
	defaultIops       int64 = 1     // needed only when VolumeType is io1
	defaultSize       int64 = 1     // in GB
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

func (e *EC2Wrapper) create(label string, opts map[string]string) (string, error) {
	var err error
	size := defaultSize
	volumeType := defaultVolumeType
	iops := defaultIops

	if opts["size"] != "" {
		if size, err = strconv.ParseInt(opts["size"], 10, 64); err != nil {
			return "", err
		}
	}

	if opts["type"] != "" {
		volumeType = opts["type"]
	}

	if opts["iops"] != "" {
		if iops, err = strconv.ParseInt(opts["iops"], 10, 64); err != nil {
			return "", err
		}
	}

	params := &ec2.CreateVolumeInput{
		AvailabilityZone: aws.String(defaultAZ),
		Size:             aws.Int64(size),
		VolumeType:       aws.String(volumeType),
	}

	if volumeType == "io1" {
		params.Iops = aws.Int64(iops)
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
