package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
	"os/exec"
)

var (
	// TODO: いい感じに
	kuoExample = `
	# hogehoge
	%[1]s kuo
	# hogehoge
	%[1]s kuo set-ctxs
	# hogehoge
	%[1]s kuo get pod
`
	errNoContext = fmt.Errorf("no context is currently set, use %q to select a new one", "kubectl config use-context <context>")
)

type KuoOptions struct {
	configFlags *genericclioptions.ConfigFlags

	rawConfig api.Config
	args      []string

	userSpecifiedFlags    string
	userSpecifiedContexts []string

	genericclioptions.IOStreams
}

func NewKuoOptions(streams genericclioptions.IOStreams) *KuoOptions {
	return &KuoOptions{
		configFlags: genericclioptions.NewConfigFlags(true),
		IOStreams:   streams,
	}
}

func NewCmdKuo(streams genericclioptions.IOStreams) *cobra.Command {
	o := NewKuoOptions(streams)

	cmd := &cobra.Command{
		Use:          "kuo [flags] [options]",
		Short:        "A kubernetes plugin that operates multiple contexts",
		Example:      fmt.Sprintf(kuoExample, "kubectl"),
		SilenceUsage: true,
		RunE: func(c *cobra.Command, args []string) error {
			if err := o.Complete(c, args); err != nil {
				return err
			}
			if err := o.Validate(); err != nil {
				return err
			}
			if err := o.Run(); err != nil {
				return err
			}

			return nil
		},
	}
	o.configFlags.AddFlags(cmd.Flags())

	return cmd
}

func (o *KuoOptions) Complete(cmd *cobra.Command, args []string) error {
	o.args = args

	if o.args != nil {
		o.userSpecifiedFlags = o.args[0]
		if o.userSpecifiedFlags == "set" {
			if len(o.args) != 3 {
				return fmt.Errorf("'set' requires 2 context arguments")
			}
			o.userSpecifiedContexts = append(o.userSpecifiedContexts, o.args[1])
			o.userSpecifiedContexts = append(o.userSpecifiedContexts, o.args[2])
		}
	}

	var err error
	o.rawConfig, err = o.configFlags.ToRawKubeConfigLoader().RawConfig()
	if err != nil {
		return err
	}

	return nil
}

func (o *KuoOptions) Validate() error {
	if len(o.rawConfig.CurrentContext) == 0 {
		return errNoContext
	}
	if len(o.args) > 4 {
		return fmt.Errorf("invalid arguments")
	}
	if o.args[0] == "apply" {
		for {
			fmt.Printf("It will be applied to multiple contexts.\nDo you want to continue? (y/n)\n")
			scanner := bufio.NewScanner(os.Stdin)
			scanner.Scan()
			result := scanner.Text()
			if result == "y" || result == "yes" {
				break
			} else if result == "n" || result == "no" {
				os.Exit(0)
			}
		}
	}

	return nil
}

func (o *KuoOptions) Run() error {
	if o.userSpecifiedFlags == "set" {
		err := o.EditKuoConfig(o.userSpecifiedContexts)
		if err != nil {
			return err
		}
	} else if o.userSpecifiedFlags != "" {
		err := o.ExecKubectl()
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *KuoOptions) ExecKubectl() error {

	contexts, err := o.ReadKuoConfig()
	if err != nil {
		return err
	}

	for _, context := range contexts {
		fmt.Printf("======== %s ========\n", context)
		err := o.ChangeContext(context)
		if err != nil {
			return err
		}
		out, err := exec.Command("kubectl", o.args...).Output()
		if err != nil {
			return err
		}
		fmt.Printf("%s", out)
	}
	return nil
}

func (o *KuoOptions) ChangeContext(newContext string) error {
	if o.rawConfig.CurrentContext != newContext {
		o.rawConfig.CurrentContext = newContext

		err := clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), o.rawConfig, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *KuoOptions) ReadKuoConfig() ([]string, error) {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	fp, err := os.Open(dirname + "/.kuoconfig")
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	var configContexts []string
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		configContexts = append(configContexts, scanner.Text())
	}
	return configContexts, nil
}

func (o *KuoOptions) EditKuoConfig(contexts []string) error {
	dirname, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	fp, err := os.Create(dirname + "/.kuoconfig")
	if err != nil {
		return err
	}
	defer fp.Close()

	for _, v := range contexts {
		fp.WriteString(v + "\n")
	}
	fmt.Printf("set .kuoconfig: %s\n", contexts)
	return nil
}
