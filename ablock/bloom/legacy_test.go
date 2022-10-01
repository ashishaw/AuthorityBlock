// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package bloom_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ashishaw/authorityblock/ablock/bloom"
)

func TestLegacyBloom(t *testing.T) {
	itemCount := 100

	bloom := bloom.NewLegacyBloom(bloom.LegacyEstimateBloomK(itemCount))

	for i := 0; i < itemCount; i++ {
		bloom.Add([]byte(fmt.Sprintf("%v", i)))
	}

	for i := 0; i < itemCount; i++ {
		assert.Equal(t, true, bloom.Test([]byte(fmt.Sprintf("%v", i))))
	}
}
