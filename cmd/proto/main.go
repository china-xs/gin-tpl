package main

import (
	"github.com/china-xs/gin-tpl/cmd/proto/add"
	"github.com/china-xs/gin-tpl/cmd/proto/server"
	"github.com/spf13/cobra"
	"log"
)

const release = "v1.0.1"

func init() {
	rootCmd.AddCommand(add.CmdAdd)
	rootCmd.AddCommand(server.CmdServer)
	//rootCmd.AddCommand(client.CmdClient)
	//rootCmd.AddCommand(server.CmdServer)
}

var rootCmd = &cobra.Command{
	Use:     "proto",
	Short:   "proto 快速生成proto文件",
	Long:    `proto 快速生成proto文件`,
	Version: release,
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
