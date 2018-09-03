package consensus

import (
	"crypto/ecdsa"
	"errors"
	"math/big"
)

var VRFNIL = VRFValue{}

type VRFValue struct {
	VrfAddress    string
	VrfPublickKey *ecdsa.PublicKey
	VrfValue      *big.Int
}

type VRFList []VRFValue

func (vrf VRFList) findThePositionOfVRFValue(v VRFValue) int {
	for index, value := range vrf {
		if value.VrfValue == v.VrfValue && value.VrfAddress == v.VrfAddress {
			return index
		}
	}
	return -1
}

func (vrf VRFList) IsInVRFList(v VRFValue) bool {
	for _, vv := range vrf {
		if vv.VrfValue == v.VrfValue && vv.VrfAddress == v.VrfAddress && vv.VrfPublickKey == v.VrfPublickKey {
			return true
		}
	}
	return false
}

func (vrf VRFList) AddVRFValue(v VRFValue) error {
	if vrf.IsInVRFList(v) {
		return errors.New("The vrf value is exist")
	}
	vrf = append(vrf, v)
}

func (vrf VRFList) DeleteVRFValue(v VRFValue) (VRFValue, error) {
	if !vrf.IsInVRFList(v) {
		return VRFNIL, errors.New("The vrf value is not exist")
	}
	position := vrf.findThePositionOfVRFValue(v)
	delValue := vrf[position+1]
	vrf = append(vrf[:position], vrf[position+1:]...)
	return delValue, nil
}
