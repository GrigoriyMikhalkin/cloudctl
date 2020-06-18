package cmd

import (
	"log"

	"github.com/metal-stack/cloud-go/api/models"

	"git.f-i-ts.de/cloud-native/cloudctl/cmd/output"
	"github.com/metal-stack/cloud-go/api/client/s3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	s3Cmd = &cobra.Command{
		Use:   "s3",
		Short: "manage s3",
		Long:  "manges access to s3 storage located in different partitions",
	}
	s3DescribeCmd = &cobra.Command{
		Use:   "describe",
		Short: "describe an s3 user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s3Describe()
		},
		PreRun: bindPFlags,
	}
	s3CreateCmd = &cobra.Command{
		Use:   "create",
		Short: "create an s3 user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s3Create()
		},
		PreRun: bindPFlags,
	}
	s3DeleteCmd = &cobra.Command{
		Use:     "remove",
		Aliases: []string{"rm", "delete"},
		Short:   "delete an s3 user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s3Delete(args)
		},
		PreRun: bindPFlags,
	}
	s3ListCmd = &cobra.Command{
		Use:     "list",
		Short:   "list s3 users",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return s3List()
		},
		PreRun: bindPFlags,
	}
	s3PartitionListCmd = &cobra.Command{
		Use:     "partitions",
		Short:   "list s3 partitions",
		Aliases: []string{"partition"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return s3ListPartitions()
		},
		PreRun: bindPFlags,
	}
	s3AddKeyCmd = &cobra.Command{
		Use:   "add-key",
		Short: "adds a key for an s3 user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s3AddKey(args)
		},
		PreRun: bindPFlags,
	}
	s3RemoveKeyCmd = &cobra.Command{
		Use:   "remove-key",
		Short: "remove a key for an s3 user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return s3RemoveKey(args)
		},
		PreRun: bindPFlags,
	}
)

func init() {
	s3CreateCmd.Flags().StringP("id", "i", "", "id of the s3 user [required]")
	s3CreateCmd.Flags().StringP("partition", "p", "", "name of s3 partition to create the s3 user in [required]")
	s3CreateCmd.Flags().String("project", "", "id of the project that the s3 user belongs to [required]")
	s3CreateCmd.Flags().StringP("tenant", "t", "", "create s3 for given tenant, defaults to logged in tenant")
	s3CreateCmd.Flags().StringP("name", "n", "", "name of s3 user, only for display")
	s3CreateCmd.Flags().Int64("max-buckets", 0, "maximum number of buckets for the s3 user")
	s3CreateCmd.Flags().StringP("access-key", "", "", "specify the access key, otherwise will be generated")
	s3CreateCmd.Flags().StringP("secret-key", "", "", "specify the secret key, otherwise will be generated")
	err := s3CreateCmd.MarkFlagRequired("id")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = s3CreateCmd.MarkFlagRequired("partition")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = s3CreateCmd.MarkFlagRequired("project")
	if err != nil {
		log.Fatal(err.Error())
	}
	s3CreateCmd.RegisterFlagCompletionFunc("partition", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return s3ListPartitionsCompletion()
	})
	s3CreateCmd.RegisterFlagCompletionFunc("project", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return projectListCompletion()
	})

	s3ListCmd.Flags().StringP("partition", "p", "", "name of s3 partition.")
	s3ListCmd.Flags().String("project", "", "id of the project that the s3 user belongs to")
	s3ListCmd.RegisterFlagCompletionFunc("partition", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return s3ListPartitionsCompletion()
	})
	s3ListCmd.RegisterFlagCompletionFunc("project", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return projectListCompletion()
	})

	s3DescribeCmd.Flags().StringP("id", "i", "", "id of the s3 user [required]")
	s3DescribeCmd.Flags().StringP("partition", "p", "", "name of s3 partition where this user is in [required]")
	s3DescribeCmd.Flags().String("project", "", "id of the project that the s3 user belongs to [required]")
	s3DescribeCmd.Flags().StringP("tenant", "t", "", "tenant of the s3 user, defaults to logged in tenant")
	err = s3DescribeCmd.MarkFlagRequired("id")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = s3DescribeCmd.MarkFlagRequired("partition")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = s3DescribeCmd.MarkFlagRequired("project")
	if err != nil {
		log.Fatal(err.Error())
	}
	s3DescribeCmd.RegisterFlagCompletionFunc("partition", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return s3ListPartitionsCompletion()
	})
	s3DescribeCmd.RegisterFlagCompletionFunc("project", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return projectListCompletion()
	})

	s3DeleteCmd.Flags().StringP("id", "i", "", "id of the s3 user [required]")
	s3DeleteCmd.Flags().StringP("partition", "p", "", "name of s3 partition where this user is in [required]")
	s3DeleteCmd.Flags().String("project", "", "id of the project that the s3 user belongs to [required]")
	s3DeleteCmd.Flags().StringP("tenant", "t", "", "tenant of the s3 user, defaults to logged in tenant")
	s3DeleteCmd.Flags().Bool("force", false, "forces s3 user deletion along with buckets and bucket objects even if those still exist (dangerous!)")
	err = s3DeleteCmd.MarkFlagRequired("id")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = s3DeleteCmd.MarkFlagRequired("partition")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = s3DeleteCmd.MarkFlagRequired("project")
	if err != nil {
		log.Fatal(err.Error())
	}
	s3DeleteCmd.RegisterFlagCompletionFunc("partition", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return s3ListPartitionsCompletion()
	})
	s3DeleteCmd.RegisterFlagCompletionFunc("project", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return projectListCompletion()
	})

	s3AddKeyCmd.Flags().StringP("id", "i", "", "id of the s3 user [required]")
	s3AddKeyCmd.Flags().StringP("partition", "p", "", "name of s3 partition where this user is in [required]")
	s3AddKeyCmd.Flags().String("project", "", "id of the project that the s3 user belongs to [required]")
	s3AddKeyCmd.Flags().StringP("tenant", "t", "", "tenant of the s3 user, defaults to logged in tenant")
	s3AddKeyCmd.Flags().StringP("access-key", "", "", "specify the access key, otherwise will be generated")
	s3AddKeyCmd.Flags().StringP("secret-key", "", "", "specify the secret key, otherwise will be generated")
	err = s3AddKeyCmd.MarkFlagRequired("id")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = s3AddKeyCmd.MarkFlagRequired("partition")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = s3AddKeyCmd.MarkFlagRequired("project")
	if err != nil {
		log.Fatal(err.Error())
	}
	s3AddKeyCmd.RegisterFlagCompletionFunc("partition", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return s3ListPartitionsCompletion()
	})
	s3AddKeyCmd.RegisterFlagCompletionFunc("project", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return projectListCompletion()
	})

	s3RemoveKeyCmd.Flags().StringP("id", "i", "", "id of the s3 user [required]")
	s3RemoveKeyCmd.Flags().StringP("partition", "p", "", "name of s3 partition where this user is in [required]")
	s3RemoveKeyCmd.Flags().String("project", "", "id of the project that the s3 user belongs to [required]")
	s3RemoveKeyCmd.Flags().StringP("tenant", "t", "", "tenant of the s3 user, defaults to logged in tenant")
	s3RemoveKeyCmd.Flags().StringP("access-key", "", "", "specify the access key to delete the access / secret key pair")
	err = s3RemoveKeyCmd.MarkFlagRequired("id")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = s3RemoveKeyCmd.MarkFlagRequired("partition")
	if err != nil {
		log.Fatal(err.Error())
	}
	err = s3RemoveKeyCmd.MarkFlagRequired("project")
	if err != nil {
		log.Fatal(err.Error())
	}
	s3RemoveKeyCmd.RegisterFlagCompletionFunc("partition", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return s3ListPartitionsCompletion()
	})
	s3RemoveKeyCmd.RegisterFlagCompletionFunc("project", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return projectListCompletion()
	})

	s3Cmd.AddCommand(s3CreateCmd)
	s3Cmd.AddCommand(s3DescribeCmd)
	s3Cmd.AddCommand(s3DeleteCmd)
	s3Cmd.AddCommand(s3ListCmd)
	s3Cmd.AddCommand(s3PartitionListCmd)
	s3Cmd.AddCommand(s3AddKeyCmd)
	s3Cmd.AddCommand(s3RemoveKeyCmd)
}

func s3Describe() error {
	tenant := viper.GetString("tenant")
	id := viper.GetString("id")
	partition := viper.GetString("partition")
	project := viper.GetString("project")

	p := &models.V1S3GetRequest{
		ID:        &id,
		Partition: &partition,
		Tenant:    &tenant,
		Project:   &project,
	}

	request := s3.NewGets3Params()
	request.SetBody(p)

	response, err := cloud.S3.Gets3(request, cloud.Auth)
	if err != nil {
		switch e := err.(type) {
		case *s3.Gets3Default:
			return output.HTTPError(e.Payload)
		default:
			return output.UnconventionalError(err)
		}
	}

	return output.YAMLPrinter{}.Print(response.Payload)
}

func s3Create() error {
	tenant := viper.GetString("tenant")
	id := viper.GetString("id")
	partition := viper.GetString("partition")
	project := viper.GetString("project")
	name := viper.GetString("name")
	maxBuckets := viper.GetInt64("max-buckets")
	accessKey := viper.GetString("access-key")
	secretKey := viper.GetString("secret-key")

	p := &models.V1S3CreateRequest{
		ID:         &id,
		Partition:  &partition,
		Tenant:     &tenant,
		Project:    &project,
		Name:       &name,
		MaxBuckets: &maxBuckets,
	}

	if accessKey != "" {
		p.Key = &models.V1S3Key{
			AccessKey: &accessKey,
			SecretKey: &secretKey,
		}
	}

	request := s3.NewCreates3Params()
	request.SetBody(p)

	response, err := cloud.S3.Creates3(request, cloud.Auth)
	if err != nil {
		switch e := err.(type) {
		case *s3.Creates3Default:
			return output.HTTPError(e.Payload)
		default:
			return output.UnconventionalError(err)
		}
	}

	return output.YAMLPrinter{}.Print(response.Payload)
}

func s3Delete(args []string) error {
	tenant := viper.GetString("tenant")
	id := viper.GetString("id")
	partition := viper.GetString("partition")
	project := viper.GetString("project")
	force := viper.GetBool("force")

	p := &models.V1S3DeleteRequest{
		ID:        &id,
		Partition: &partition,
		Tenant:    &tenant,
		Project:   &project,
		Force:     &force,
	}

	request := s3.NewDeletes3Params()
	request.SetBody(p)

	response, err := cloud.S3.Deletes3(request, cloud.Auth)
	if err != nil {
		switch e := err.(type) {
		case *s3.Deletes3Default:
			return output.HTTPError(e.Payload)
		default:
			return output.UnconventionalError(err)
		}
	}

	return output.YAMLPrinter{}.Print(response.Payload)
}

func s3AddKey(args []string) error {
	tenant := viper.GetString("tenant")
	id := viper.GetString("id")
	partition := viper.GetString("partition")
	project := viper.GetString("project")
	accessKey := viper.GetString("access-key")
	secretKey := viper.GetString("secret-key")

	p := &models.V1S3UpdateRequest{
		ID:        &id,
		Partition: &partition,
		Tenant:    &tenant,
		Project:   &project,
		AddKeys: []*models.V1S3Key{
			{
				AccessKey: &accessKey,
				SecretKey: &secretKey,
			},
		},
	}

	request := s3.NewUpdates3Params()
	request.SetBody(p)

	response, err := cloud.S3.Updates3(request, cloud.Auth)
	if err != nil {
		switch e := err.(type) {
		case *s3.Updates3Default:
			return output.HTTPError(e.Payload)
		default:
			return output.UnconventionalError(err)
		}
	}

	return output.YAMLPrinter{}.Print(response.Payload)
}

func s3RemoveKey(args []string) error {
	tenant := viper.GetString("tenant")
	id := viper.GetString("id")
	partition := viper.GetString("partition")
	project := viper.GetString("project")
	accessKey := viper.GetString("access-key")

	p := &models.V1S3UpdateRequest{
		ID:        &id,
		Partition: &partition,
		Tenant:    &tenant,
		Project:   &project,
		RemoveAccessKeys: []string{
			accessKey,
		},
	}

	request := s3.NewUpdates3Params()
	request.SetBody(p)

	response, err := cloud.S3.Updates3(request, cloud.Auth)
	if err != nil {
		switch e := err.(type) {
		case *s3.Updates3Default:
			return output.HTTPError(e.Payload)
		default:
			return output.UnconventionalError(err)
		}
	}

	return output.YAMLPrinter{}.Print(response.Payload)
}

func s3List() error {
	partition := viper.GetString("partition")
	project := viper.GetString("project")

	p := &models.V1S3ListRequest{
		Partition: &partition,
	}

	request := s3.NewLists3Params()
	request.SetBody(p)

	response, err := cloud.S3.Lists3(request, cloud.Auth)
	if err != nil {
		switch e := err.(type) {
		case *s3.Lists3Default:
			return output.HTTPError(e.Payload)
		default:
			return output.UnconventionalError(err)
		}
	}

	if project == "" {
		return printer.Print(response.Payload)
	}

	var result []*models.V1S3Response
	for _, s3 := range response.Payload {
		if *s3.Project == project {
			result = append(result, s3)
		}
	}
	return printer.Print(result)
}

func s3ListPartitions() error {
	request := s3.NewLists3partitionsParams()

	response, err := cloud.S3.Lists3partitions(request, cloud.Auth)
	if err != nil {
		switch e := err.(type) {
		case *s3.Lists3partitionsDefault:
			return output.HTTPError(e.Payload)
		default:
			return output.UnconventionalError(err)
		}
	}
	return printer.Print(response.Payload)
}
