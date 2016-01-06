package main

import (
	"docker-ebs-plugin/mocks"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestFindVolumeReturnsIdWhenOneVolume(t *testing.T) {
	mockOutput := &ec2.DescribeVolumesOutput{
		Volumes: []*ec2.Volume{
			&ec2.Volume{
				AvailabilityZone: aws.String("eu-west-1a"),
				CreateTime:       aws.Time(time.Now()),
				Encrypted:        aws.Bool(false),
				Iops:             aws.Int64(3),
				Size:             aws.Int64(1),
				SnapshotId:       aws.String(""),
				State:            aws.String("available"),
				Tags: []*ec2.Tag{
					&ec2.Tag{
						Key:   aws.String("Name"),
						Value: aws.String("docker1"),
					},
					&ec2.Tag{
						Key:   aws.String("DockerVolumeName"),
						Value: aws.String("docker1"),
					},
				},
				VolumeId:   aws.String("vol-681e4aac"),
				VolumeType: aws.String("gp2"),
			},
		},
	}

	m := new(mocks.Ec2er)
	m.On("DescribeVolumes", mock.AnythingOfType("*ec2.DescribeVolumesInput")).Return(mockOutput, nil)

	wrapper := &Ec2Wrapper{m}

	output, err := wrapper.findVolume("docker1")

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "vol-681e4aac", output)
}
