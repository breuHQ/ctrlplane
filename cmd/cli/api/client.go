// Copyright © 2023, Breu, Inc. <info@breu.io>. All rights reserved.
//
// This software is made available by Breu, Inc., under the terms of the BREU COMMUNITY LICENSE AGREEMENT, Version 1.0,
// found at https://www.breu.io/license/community. BY INSTALLING, DOWNLOADING, ACCESSING, USING OR DISTRIBUTING ANY OF
// THE SOFTWARE, YOU AGREE TO THE TERMS OF THE LICENSE AGREEMENT.
//
// The above copyright notice and the subsequent license agreement shall be included in all copies or substantial
// portions of the software.
//
// Breu, Inc. HEREBY DISCLAIMS ANY AND ALL WARRANTIES AND CONDITIONS, EXPRESS, IMPLIED, STATUTORY, OR OTHERWISE, AND
// SPECIFICALLY DISCLAIMS ANY WARRANTY OF MERCHANTABILITY OR FITNESS FOR A PARTICULAR PURPOSE, WITH RESPECT TO THE
// SOFTWARE.
//
// Breu, Inc. SHALL NOT BE LIABLE FOR ANY DAMAGES OF ANY KIND, INCLUDING BUT NOT LIMITED TO, LOST PROFITS OR ANY
// CONSEQUENTIAL, SPECIAL, INCIDENTAL, INDIRECT, OR DIRECT DAMAGES, HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY,
// ARISING OUT OF THIS AGREEMENT. THE FOREGOING SHALL APPLY TO THE EXTENT PERMITTED BY APPLICABLE LAW.

package api

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"go.breu.io/quantm/internal/auth"
	"go.breu.io/quantm/internal/core"
	"go.breu.io/quantm/internal/providers/github"
)

var Client client

type client struct {
	AuthClient   *auth.Client
	CoreClient   *core.Client
	GithubClient *github.Client
}

// CheckStatus returns false if the status code is other than provided in parameters.
func (c *client) CheckStatus(r *http.Response, successCodes ...int) {
	pass := false
	for _, c := range successCodes {
		if r.StatusCode == c {
			pass = true
		}
	}
	if pass == false {
		fmt.Printf("Command failed with status code: %d\r\n", r.StatusCode)
	}
}

func (c *client) CheckError(err error) {
	if err != nil {
		if strings.Contains(err.Error(), "No connection") {
			fmt.Print("Quantum server is not running\n")
		} else {
			fmt.Printf("Command failed: %v", err.Error())
		}
		os.Exit(1)
	}
}

// init initializes the auth, github and core clients to connect with quantm.
func (c *client) Init(url string) {
	c.AuthClient, _ = auth.NewClient(url)
	c.CoreClient, _ = core.NewClient(url)
	c.GithubClient, _ = github.NewClient(url)
}