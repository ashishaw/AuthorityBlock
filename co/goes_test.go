// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>

package co_test

import (
	"testing"

	"github.com/ashishaw/authorityblock/co"
)

func TestGoes(t *testing.T) {
	var g co.Goes
	g.Go(func() {})
	g.Go(func() {})
	g.Wait()

	<-g.Done()
}
