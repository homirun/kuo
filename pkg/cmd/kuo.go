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

	return nil
}

func (o *KuoOptions) Run() error {
	if o.userSpecifiedFlags == "set" {
		err := o.editKuoConfig(o.userSpecifiedContexts)
		if err != nil {
			return err
		}
	}

	if o.userSpecifiedFlags == "get" {
		o.getRouter()
	}

	return nil
}

func (o *KuoOptions) getRouter() error {
	if o.args[1] == "node" {
		err := o.getNode()
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *KuoOptions) getNode() error {
	contexts, err := o.readKuoConfig()
	if err != nil {
		return err
	}

	for _, context := range contexts {
		fmt.Printf("======== %s ========\n", context)
		o.changeContext(context)
		out, err := exec.Command("kubectl", "get", "node").Output()
		if err != nil {
			return err
		}
		fmt.Printf("%s", out)
	}
	return nil
}

func (o *KuoOptions) changeContext(newContext string) error {
	if o.rawConfig.CurrentContext != newContext {
		o.rawConfig.CurrentContext = newContext

		err := clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), o.rawConfig, true)
		if err != nil {
			return err
		}
	}
	return nil
}

func (o *KuoOptions) readKuoConfig() ([]string, error) {
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

func (o *KuoOptions) editKuoConfig(contexts []string) error {
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

	return nil
}
