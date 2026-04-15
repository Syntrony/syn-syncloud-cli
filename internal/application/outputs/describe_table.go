package outputs

import (
	"encoding/json"
	"fmt"
	dto "synctl/internal/domain/dto"
)

func PrintDescribe(resource dto.ResourceDetailDto, outputFormat string) {
	if outputFormat == "json" {
		PrintDescribeJson(resource)
	} else {
		PrintDescribeTable(resource)
	}
}

func PrintDescribeTable(r dto.ResourceDetailDto) {
	fmt.Println("Name:       ", r.Name)
	fmt.Println("Kind:       ", r.Kind)
	fmt.Println("Runtime:    ", r.Runtime)
	fmt.Println("ID:         ", r.Id)
	fmt.Println("Created:    ", r.CreatedAt)
	fmt.Println("Updated:   ", r.UpdatedAt)

	if len(r.Status) > 0 {
		fmt.Println("\nStatus:")
		for k, v := range r.Status {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}

	if len(r.Spec) > 0 {
		fmt.Println("\nSpec:")
		for k, v := range r.Spec {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}

	if len(r.LiveSpec) > 0 {
		fmt.Println("\nLive Data:")
		for k, v := range r.LiveSpec {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}

	if len(r.Events) > 0 {
		fmt.Println("\nEvents:")
		for _, e := range r.Events {
			fmt.Println("  ", e)
		}
	}
}

func PrintDescribeJson(r dto.ResourceDetailDto) {
	b, _ := json.MarshalIndent(r, "", "  ")
	fmt.Println(string(b))
}
