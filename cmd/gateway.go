package cmd

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/fi-ts/cloud-go/api/client/gateway"
	"github.com/fi-ts/cloud-go/api/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/utils/pointer"
)

var (
	gatewayCmd = &cobra.Command{
		Use:   "gateway",
		Short: "Manage gateways",
		Long:  "Manage gateways, which enable access to services in another cluster",
	}
	gatewayCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create a gateway",
		RunE: func(cmd *cobra.Command, args []string) error {
			return gatewayCreate()
		},
		PreRun: bindPFlags,
	}
)

func init() {
	gatewayCreateCmd.Flags().String("name", "", "Name of the gateway")
	gatewayCreateCmd.Flags().String("clients", "", "Name of the gateway")
	gatewayCreateCmd.Flags().String("project", "", "Project-UID which the gateway belongs to")
	gatewayCreateCmd.Flags().String("pipes", "", "Comma-separated-list of pipes (e.g. PIPE_1,PIPE_2). Pipe has format SVC_NAME_1:CLIENT_POD_PORT_1:REMOTE_SVC_ENDPOINT_1.")
	if err := gatewayCreateCmd.MarkFlagRequired("name"); err != nil {
		log.Fatal(err)
	}
	if err := gatewayCreateCmd.MarkFlagRequired("project"); err != nil {
		log.Fatal(err)
	}
	if err := gatewayCreateCmd.MarkFlagRequired("pipes"); err != nil {
		log.Fatal(err)
	}
	gatewayCmd.AddCommand(gatewayCreateCmd)
}

func gatewayCreate() error {
	pipes := viper.GetString("pipes")
	parsed, err := parseFlagPipes(pipes)
	if err != nil {
		return fmt.Errorf("failed to parse pipes flag %s: %w", pipes, err)
	}

	params := gateway.NewCreateGatewayParams()
	params.SetBody(&models.V1GatewayCreateRequest{
		ProjectUID: ptr(viper.GetString("project")),
		Name:       ptr(viper.GetString("name")),
		Pipes:      parsed,
	})

	resp, err := cloud.Gateway.CreateGateway(params, nil)
	if err != nil {
		return fmt.Errorf("failed to create a gateway with params %v: %w", params, err)
	}
	return printer.Print(resp.Payload)
}

func parseFlagPipes(flag string) ([]*models.V1PipeSpec, error) {
	ss := strings.Split(flag, ",")
	pipes := []*models.V1PipeSpec{}
	for i := range ss {
		pipe, err := parsePipe(ss[i])
		if err != nil {
			return nil, fmt.Errorf("failed to parse flag `pipes`: %w", err)
		}
		pipes = append(pipes, pipe)
	}

	return pipes, nil
}

func parsePipe(unparsed string) (*models.V1PipeSpec, error) {
	ss := strings.Split(unparsed, ":")
	if len(ss) < 3 {
		return nil, errors.New("pipe incomplete: it should be a colon-separated-list `SVC_NAME_1:CLIENT_POD_PORT_1:REMOTE_SVC_ENDPOINT_1`")
	}

	pipe := &models.V1PipeSpec{}
	pipe.Name = ptr(ss[0])
	port, err := u16StrToI64Ptr(ss[1])
	if err != nil {
		return nil, fmt.Errorf("failed to convert `%s` to pointer to int64: %w", ss[2], err)
	}
	pipe.Port = port
	pipe.Remote = ptr(ss[2])
	return pipe, nil
}

func ptr(s string) *string {
	return pointer.StringPtr(s)
}

// Convert an uint16 as string to a pointer to int64
func u16StrToI64Ptr(s string) (*int64, error) {
	u16AsU64, err := strconv.ParseUint(s, 10, 16) // uint16 in gateway k8s api
	if err != nil {
		return nil, fmt.Errorf("failed to convert the port in pipe %s to uint16: %w", s, err)
	}

	return pointer.Int64Ptr(int64(u16AsU64)), nil
}