package actions

func init() {
	// Register uninstall action
	Register(&Action{
		ID:           ActionUninstall,
		Use:          "uninstall",
		Short:        "Completely uninstall dnstm",
		Long:         "Remove all dnstm components from the system.\n\nThis will:\n  - Stop and remove all instance services\n  - Stop and remove DNS router service\n  - Stop and remove microsocks service\n  - Remove all configuration in /etc/dnstm\n  - Remove dnstm user\n  - Remove transport binaries (dnstt-server, slipstream-server, ssserver, microsocks)\n  - Remove firewall rules\n\nNote: The dnstm binary itself is kept for easy reinstallation.",
		MenuLabel:    "Uninstall",
		RequiresRoot: true,
		Confirm: &ConfirmConfig{
			Message:     "Are you sure you want to uninstall everything?",
			Description: "This will remove all dnstm components from your system.",
			DefaultNo:   true,
			ForceFlag:   "force",
		},
	})

	// Register install action
	Register(&Action{
		ID:           ActionInstall,
		Use:          "install",
		Short:        "Verify transport binaries are present and configure system",
		Long:         "Verify transport binaries and configure the system for DNS tunneling.\n\nThis will:\n  - Create dnstm system user\n  - Initialize router configuration and directories\n  - Set operating mode (defaults to single)\n  - Create DNS router service\n  - Verify transport binaries are present\n  - Configure firewall rules (port 53 UDP/TCP)\n\nTransport binaries must be copied to /usr/local/bin manually (e.g. via SCP).\n\nOptionally use --mode to set the operating mode:\n  single  Single-tunnel mode (default) - one tunnel at a time\n  multi   Multi-tunnel mode - multiple tunnels with DNS router",
		MenuLabel:    "Install",
		RequiresRoot: true,
		Inputs: []InputField{
			{
				Name:  "force",
				Label: "Force reinstall if already installed",
				Type:  InputTypeBool,
			},
			{
				Name:      "mode",
				Label:     "Operating Mode",
				ShortFlag: 'm',
				Type:      InputTypeSelect,
				Options:   OperatingModeOptions(),
				Default:   "single",
				// Skip mode selection in interactive mode - defaults to single,
				// user will be prompted to switch to multi when adding second tunnel
				ShowIf: func(ctx *Context) bool { return !ctx.IsInteractive },
			},
		},
	})

	// Register ssh-users action (TUI-only, hidden from CLI help)
	Register(&Action{
		ID:                ActionSSHUsers,
		Use:               "ssh-users",
		Short:             "Manage SSH tunnel users",
		Long:              "Launch sshtun-user for managing SSH tunnel users and hardening",
		MenuLabel:         "SSH Users",
		Hidden:            true,
		RequiresRoot:      true,
		RequiresInstalled: true,
	})

}

// SetSystemHandler sets the handler for a system action.
func SetSystemHandler(actionID string, handler Handler) {
	SetHandler(actionID, handler)
}
