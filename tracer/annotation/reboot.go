package annotation

import "C"

// from reboot(2) manpages
var rebootMagics = map[uintptr]string{
	0xfee1dead: "LINUX_REBOOT_MAGIC1",
	672274793:  "LINUX_REBOOT_MAGIC2",
	85072278:   "LINUX_REBOOT_MAGIC2A",
	369367448:  "LINUX_REBOOT_MAGIC2B",
	537993216:  "LINUX_REBOOT_MAGIC2C",
}

func AnnotateRebootMagic(arg Arg, _ int) {
	if magic, ok := rebootMagics[arg.Raw()]; ok {
		arg.SetAnnotation(magic, true)
	}
}

var rebootCmds = map[uintptr]string{
	0x00000000: "LINUX_REBOOT_CMD_CAD_OFF",
	0x89ABCDEF: "LINUX_REBOOT_CMD_CAD_ON",
	0xCDEF0123: "LINUX_REBOOT_CMD_HALT",
	0x45584543: "LINUX_REBOOT_CMD_KEXEC",
	0x4321FEDC: "LINUX_REBOOT_CMD_POWER_OFF",
	0x01234567: "LINUX_REBOOT_CMD_RESTART",
	0xA1B2C3D4: "LINUX_REBOOT_CMD_RESTART2",
	0xD000FCE2: "LINUX_REBOOT_CMD_SW_SUSPEND",
}

func AnnotateRebootCmd(arg Arg, _ int) {
	if cmd, ok := rebootCmds[arg.Raw()]; ok {
		arg.SetAnnotation(cmd, true)
	}
}
