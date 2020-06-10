// Copyright 2020. Akamai Technologies, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
        "io/ioutil"

	dnsv2 "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

func cmdCreateRecordsets(c *cli.Context) error {

	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}
	dnsv2.Init(config)

	var (
		zonename string
		//outputPath string
		inputPath string
	)

	if c.NArg() == 0 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("zonename is required"), 1)
	}

	zonename = c.Args().First()
	
	if c.IsSet("file") {
		inputPath = c.String("file")
		inputPath = filepath.FromSlash(inputPath)
	} else {
                return cli.NewExitError(color.RedString("Input file is required"), 1)
	}
        akamai.StartSpinner("Fetching Recordset data ", "")
        // Read in json file
        data, err := ioutil.ReadFile(inputPath)
        if err != nil {
                akamai.StopSpinnerFail()
                return cli.NewExitError(color.RedString("Failed to read input file"), 1)
        }
	recordsets := &dnsv2.Recordsets{}
        err = json.Unmarshal(data, recordsets)
        if err != nil {
                akamai.StopSpinnerFail()
                return cli.NewExitError(color.RedString("Failed to parse json file content"), 1)
	}
	akamai.StopSpinnerOk()
       	akamai.StartSpinner("Creating Recordsets ", "")
	err = recordsets.Save(zonename, true)
	if err != nil {
		akamai.StopSpinnerFail()
		return cli.NewExitError(color.RedString(fmt.Sprintf("Recordset creation failed. Error: %s", err.Error())), 1)
	}
	akamai.StopSpinnerOk()

	return nil
}

