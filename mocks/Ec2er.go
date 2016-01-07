package mocks

import "github.com/stretchr/testify/mock"

import "github.com/aws/aws-sdk-go/service/ec2"

type EC2er struct {
	mock.Mock
}

func (_m *EC2er) DescribeVolumes(_a0 *ec2.DescribeVolumesInput) (*ec2.DescribeVolumesOutput, error) {
	ret := _m.Called(_a0)

	var r0 *ec2.DescribeVolumesOutput
	if rf, ok := ret.Get(0).(func(*ec2.DescribeVolumesInput) *ec2.DescribeVolumesOutput); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ec2.DescribeVolumesOutput)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ec2.DescribeVolumesInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *EC2er) CreateVolume(_a0 *ec2.CreateVolumeInput) (*ec2.Volume, error) {
	ret := _m.Called(_a0)

	var r0 *ec2.Volume
	if rf, ok := ret.Get(0).(func(*ec2.CreateVolumeInput) *ec2.Volume); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ec2.Volume)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*ec2.CreateVolumeInput) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
