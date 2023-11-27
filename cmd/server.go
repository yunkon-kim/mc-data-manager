/*
Copyright © 2023 cychoi, tykim <dev@zconverter.com>

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
package cmd

import (
	"github.com/cloud-barista/cm-data-mold/internal/log"
	dmsv "github.com/cloud-barista/cm-data-mold/websrc/serve"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var listenPort string
var allowIP []string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start Web Server",
	Long:  `Start Web Server`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.SetFormatter(&log.CustomTextFormatter{CmdName: "server", JobName: "web server"})
		logrus.Info("Start Web Server")
		dmsv.Run(dmsv.InitServer(allowIP...), listenPort)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVarP(&listenPort, "port", "P", "80", "Listen port")
	serverCmd.Flags().StringArrayVarP(&allowIP, "allow-ip", "I", []string{}, "ip to allow")
}
