package banner

import "fmt"

func PrintBanner() {
	const reset = "\033[0m"
	const blue = "\033[38;2;30;80;160m"
	const blueL = "\033[38;2;90;150;220m"
	const blueW = "\033[38;2;180;210;240m"
	const white = "\033[38;2;220;235;255m"
	const gray = "\033[38;2;140;170;210m"

	banner := []string{
		blue + "      ██████       " + reset,
		blue + "    ██      ██     " + reset,
		blueL + "   ██  " + blueW + "██████" + blueL + "  ██    " + reset,
		blueL + "   ██  " + blueW + "██    " + blueL + "████    " + reset,
		blue + "    ██  " + blueW + "██████" + blue + "  ██     " + reset,
		blue + "      ██████       " + reset,
		"",
		white + "  ███████╗██╗   ██╗███╗   ██╗████████╗██████╗  ██████╗ ███╗   ██╗██╗   ██╗" + reset,
		white + "  ██╔════╝╚██╗ ██╔╝████╗  ██║╚══██╔══╝██╔══██╗██╔═══██╗████╗  ██║╚██╗ ██╔╝" + reset,
		white + "  ███████╗ ╚████╔╝ ██╔██╗ ██║   ██║   ██████╔╝██║   ██║██╔██╗ ██║ ╚████╔╝ " + reset,
		white + "  ╚════██║  ╚██╔╝  ██║╚██╗██║   ██║   ██╔══██╗██║   ██║██║╚██╗██║  ╚██╔╝  " + reset,
		white + "  ███████║   ██║   ██║ ╚████║   ██║   ██║  ██║╚██████╔╝██║ ╚████║   ██║   " + reset,
		white + "  ╚══════╝   ╚═╝   ╚═╝  ╚═══╝   ╚═╝   ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝   ╚═╝   " + reset,
		"",
		gray + "                         S Y N C L O U D  P L A T F O R M" + reset,
		"",
	}

	for _, line := range banner {
		fmt.Println(line)
	}
}
