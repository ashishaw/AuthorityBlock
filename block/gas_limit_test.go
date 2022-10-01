// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package block_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ashishaw/authorityblock/block"
	"github.com/ashishaw/authorityblock/ablock"
)

func TestGasLimit_IsValid(t *testing.T) {

	tests := []struct {
		gl       uint64
		parentGL uint64
		want     bool
	}{
		{ablock.MinGasLimit, ablock.MinGasLimit, true},
		{ablock.MinGasLimit - 1, ablock.MinGasLimit, false},
		{ablock.MinGasLimit, ablock.MinGasLimit * 2, false},
		{ablock.MinGasLimit * 2, ablock.MinGasLimit, false},
		{ablock.MinGasLimit + ablock.MinGasLimit/ablock.GasLimitBoundDivisor, ablock.MinGasLimit, true},
		{ablock.MinGasLimit*2 + ablock.MinGasLimit/ablock.GasLimitBoundDivisor, ablock.MinGasLimit * 2, true},
		{ablock.MinGasLimit*2 - ablock.MinGasLimit/ablock.GasLimitBoundDivisor, ablock.MinGasLimit * 2, true},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, block.GasLimit(tt.gl).IsValid(tt.parentGL))
	}
}

func TestGasLimit_Adjust(t *testing.T) {

	tests := []struct {
		gl    uint64
		delta int64
		want  uint64
	}{
		{ablock.MinGasLimit, 1, ablock.MinGasLimit + 1},
		{ablock.MinGasLimit, -1, ablock.MinGasLimit},
		{math.MaxUint64, 1, math.MaxUint64},
		{ablock.MinGasLimit, int64(ablock.MinGasLimit), ablock.MinGasLimit + ablock.MinGasLimit/ablock.GasLimitBoundDivisor},
		{ablock.MinGasLimit * 2, -int64(ablock.MinGasLimit), ablock.MinGasLimit*2 - (ablock.MinGasLimit*2)/ablock.GasLimitBoundDivisor},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, block.GasLimit(tt.gl).Adjust(tt.delta))
	}
}

func TestGasLimit_Qualify(t *testing.T) {
	tests := []struct {
		gl       uint64
		parentGL uint64
		want     uint64
	}{
		{ablock.MinGasLimit, ablock.MinGasLimit, ablock.MinGasLimit},
		{ablock.MinGasLimit - 1, ablock.MinGasLimit, ablock.MinGasLimit},
		{ablock.MinGasLimit, ablock.MinGasLimit * 2, ablock.MinGasLimit*2 - (ablock.MinGasLimit*2)/ablock.GasLimitBoundDivisor},
		{ablock.MinGasLimit * 2, ablock.MinGasLimit, ablock.MinGasLimit + ablock.MinGasLimit/ablock.GasLimitBoundDivisor},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, block.GasLimit(tt.gl).Qualify(tt.parentGL))
	}
}
