// Copyright (c) 2022 Ashish Waingankar

// Distributed under the GNU Lesser General Public License v3.0 software license, see the accompanying
// file LICENSE or <https://www.gnu.org/licenses/lgpl-3.0.html>
package node_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/ashishaw/authorityblock/api/node"
	"github.com/ashishaw/authorityblock/chain"
	"github.com/ashishaw/authorityblock/comm"
	"github.com/ashishaw/authorityblock/genesis"
	"github.com/ashishaw/authorityblock/muxdb"
	"github.com/ashishaw/authorityblock/state"
	"github.com/ashishaw/authorityblock/txpool"
)

var ts *httptest.Server

func TestNode(t *testing.T) {
	initCommServer(t)
	res := httpGet(t, ts.URL+"/node/network/peers")
	var peersStats map[string]string
	if err := json.Unmarshal(res, &peersStats); err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 0, len(peersStats), "count should be zero")
}

func initCommServer(t *testing.T) {
	db := muxdb.NewMem()
	stater := state.NewStater(db)
	gene := genesis.NewDevnet()

	b, _, _, err := gene.Build(stater)
	if err != nil {
		t.Fatal(err)
	}
	repo, _ := chain.NewRepository(db, b)
	comm := comm.New(repo, txpool.New(repo, stater, txpool.Options{
		Limit:           10000,
		LimitPerAccount: 16,
		MaxLifetime:     10 * time.Minute,
	}))
	router := mux.NewRouter()
	node.New(comm).Mount(router, "/node")
	ts = httptest.NewServer(router)
}

func httpGet(t *testing.T, url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	r, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	return r
}
