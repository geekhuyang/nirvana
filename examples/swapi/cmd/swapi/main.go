/*
Copyright 2017 Caicloud Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/caicloud/nirvana/cli"
	"github.com/caicloud/nirvana/examples/swapi/pkg/api"
	"github.com/caicloud/nirvana/examples/swapi/pkg/loader"
	"github.com/spf13/cobra"
)

const ErrCode = 1

func main() {
	var port uint
	var dp string
	c := cli.NewCommand(&cobra.Command{
		RunE: func(cmd *cobra.Command, _ []string) error {
			ml, err := loader.New(toAbs(dp))
			if err != nil {
				return err
			}
			s := api.CreateWebServer(ml)
			hostAndPort := fmt.Sprintf(":%d", port)
			fmt.Printf("start listening on %v\n", hostAndPort)
			http.ListenAndServe(hostAndPort, s)
			return nil
		},
	})

	c.AddFlag(
		cli.UintFlag{
			Name:        "port",
			Shorthand:   "p",
			Usage:       "port that the server listens to",
			Destination: &port,
			DefValue:    8000,
		},
		cli.StringFlag{
			Name:        "data-path",
			Shorthand:   "d",
			Usage:       "supply the data path",
			Destination: &dp,
		},
	)

	if err := c.Execute(); err != nil {
		os.Exit(ErrCode)
	}
}

func toAbs(dataPath string) string {
	if path.IsAbs(dataPath) {
		return dataPath
	} else if p, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		return path.Join(p, dataPath)
	}
}