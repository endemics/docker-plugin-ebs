package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/endemics/docker-plugin-ebs/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestFindReturnsIdWhenOneMatchingVolume(t *testing.T) {
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
						Value: aws.String("label"),
					},
					&ec2.Tag{
						Key:   aws.String("DockerVolumeName"),
						Value: aws.String("label"),
					},
				},
				VolumeId:   aws.String("vol-681e4aac"),
				VolumeType: aws.String("gp2"),
			},
		},
	}

	m := new(mocks.EC2)
	m.On("DescribeVolumes", mock.AnythingOfType("*ec2.DescribeVolumesInput")).Return(mockOutput, nil)

	wrapper := &EC2Wrapper{m}

	output, err := wrapper.find("label")

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "vol-681e4aac", output, "find should return the volumeId of the volume matching DockerVolumeName=label")
}

func TestFindReturnsEmptyStringWhenNoMatchingVolume(t *testing.T) {
	mockOutput := &ec2.DescribeVolumesOutput{
		Volumes: []*ec2.Volume{},
	}

	m := new(mocks.EC2)
	m.On("DescribeVolumes", mock.AnythingOfType("*ec2.DescribeVolumesInput")).Return(mockOutput, nil)

	wrapper := &EC2Wrapper{m}

	output, err := wrapper.find("nosuchlabel")

	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, "", output, "find should return an empty string when no volume matches DockerVolumeName=label")
}

func TestFindReturnsErrorIfMoreThanOneVolumeMatchesLabel(t *testing.T) {
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
						Value: aws.String("label"),
					},
					&ec2.Tag{
						Key:   aws.String("DockerVolumeName"),
						Value: aws.String("label"),
					},
				},
				VolumeId:   aws.String("vol-681e4aac"),
				VolumeType: aws.String("gp2"),
			},
			&ec2.Volume{
				AvailabilityZone: aws.String("eu-west-1b"),
				CreateTime:       aws.Time(time.Now()),
				Encrypted:        aws.Bool(false),
				Iops:             aws.Int64(3),
				Size:             aws.Int64(1),
				SnapshotId:       aws.String(""),
				State:            aws.String("available"),
				Tags: []*ec2.Tag{
					&ec2.Tag{
						Key:   aws.String("Name"),
						Value: aws.String("label"),
					},
					&ec2.Tag{
						Key:   aws.String("DockerVolumeName"),
						Value: aws.String("label"),
					},
				},
				VolumeId:   aws.String("vol-1234beef"),
				VolumeType: aws.String("gp2"),
			},
		},
	}

	m := new(mocks.EC2)
	m.On("DescribeVolumes", mock.AnythingOfType("*ec2.DescribeVolumesInput")).Return(mockOutput, nil)

	wrapper := &EC2Wrapper{m}

	_, err := wrapper.find("label")

	assert.Error(t, err, "find should return an error when more than one volume matches the label")
}

func TestFindReturnsErrorWhenEC2ReturnsError(t *testing.T) {
	mockOutput := &ec2.DescribeVolumesOutput{
		Volumes: []*ec2.Volume{},
	}

	m := new(mocks.EC2)
	m.On("DescribeVolumes", mock.AnythingOfType("*ec2.DescribeVolumesInput")).Return(mockOutput, fmt.Errorf("this is a mocked AWS error"))

	wrapper := &EC2Wrapper{m}

	_, err := wrapper.find("label")

	assert.Error(t, err, "find should return an error when AWS returns an error")
}

func TestCreateReturnsVolumeIdWhenCreatingVolume(t *testing.T) {
	optsTests := []map[string]string{
		{},
		{
			"size": "10",
		},
		{
			"type": "standard",
		},
		{
			"type": "io1",
			"iops": "1000",
		},
	}

	mockOutput := &ec2.Volume{
		AvailabilityZone: aws.String("eu-west-1a"),
		CreateTime:       aws.Time(time.Now()),
		Encrypted:        aws.Bool(false),
		Iops:             aws.Int64(3),
		Size:             aws.Int64(1),
		SnapshotId:       aws.String(""),
		State:            aws.String("creating"),
		VolumeId:         aws.String("vol-681e4aac"),
		VolumeType:       aws.String("gp2"),
	}

	m := new(mocks.EC2)
	m.On("CreateVolume", mock.AnythingOfType("*ec2.CreateVolumeInput")).Return(mockOutput, nil)

	wrapper := &EC2Wrapper{m}

	for _, tt := range optsTests {
		output, err := wrapper.create("label", tt)

		if err != nil {
			t.Error(err)
		}

		assert.Equal(t, "vol-681e4aac", output, "create should return the volumeId of the volume created")
	}
}

func TestCreateReturnsErrorWhenEC2ReturnsError(t *testing.T) {
	m := new(mocks.EC2)
	m.On("CreateVolume", mock.AnythingOfType("*ec2.CreateVolumeInput")).Return(&ec2.Volume{}, fmt.Errorf("this is a mocked AWS error"))

	wrapper := &EC2Wrapper{m}

	_, err := wrapper.create("label", map[string]string{})

	assert.Error(t, err, "create should return an error when AWS returns an error")
}

func TestCreateReturnsErrorWhenSizeOrIopsAreInvalid(t *testing.T) {
	optsTests := []map[string]string{
		{
			"size": "x",
		},
		{
			"iops": "y",
		},
	}

	mockOutput := &ec2.Volume{
		AvailabilityZone: aws.String("eu-west-1a"),
		CreateTime:       aws.Time(time.Now()),
		Encrypted:        aws.Bool(false),
		Iops:             aws.Int64(3),
		Size:             aws.Int64(1),
		SnapshotId:       aws.String(""),
		State:            aws.String("creating"),
		VolumeId:         aws.String("vol-681e4aac"),
		VolumeType:       aws.String("gp2"),
	}

	m := new(mocks.EC2)
	m.On("CreateVolume", mock.AnythingOfType("*ec2.CreateVolumeInput")).Return(mockOutput, nil)

	wrapper := &EC2Wrapper{m}

	for _, tt := range optsTests {
		_, err := wrapper.create("label", tt)

		assert.Error(t, err, "create should return an error when size or iops cannot be converted in int64")
	}
}

func TestTag(t *testing.T) {
	m := new(mocks.EC2)
	m.On("CreateTags", mock.AnythingOfType("*ec2.CreateTagsInput")).Return(&ec2.CreateTagsOutput{}, nil)

	wrapper := &EC2Wrapper{m}

	err := wrapper.tag("vol-1234beef", "label")

	assert.NoError(t, err, "tag should not return an error when all is fine")
}

func TestTagReturnsErrorWhenEC2ReturnsError(t *testing.T) {
	m := new(mocks.EC2)
	m.On("CreateTags", mock.AnythingOfType("*ec2.CreateTagsInput")).Return(&ec2.CreateTagsOutput{}, fmt.Errorf("this is a mocked AWS error"))

	wrapper := &EC2Wrapper{m}

	err := wrapper.tag("vol-1234beef", "label")

	assert.Error(t, err, "tag should return an error when AWS returns an error")
}
