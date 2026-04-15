package outputs

import (
	"fmt"
	"os"
	dto "synctl/internal/domain/dto"
	"text/tabwriter"
)

func PrintNodes(nodes []dto.NodeDto) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tHOSTNAME\tROLE\tIP\tOS\tARCH")
	for _, n := range nodes {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", n.Id, n.Hostname, n.Role, n.Ip, n.Os, n.Arch)
	}
	w.Flush()
}

func PrintResources(resources []dto.ResourceDto) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tNAME\tKIND\tRUNTIME\tCREATED")
	for _, r := range resources {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n", r.Id, r.Name, r.Kind, r.Runtime, r.CreatedAt)
	}
	w.Flush()
}
