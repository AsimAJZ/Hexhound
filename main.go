package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func scanner(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("reading directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		filePath := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			err := scanner(filePath)
			if err != nil {
				fmt.Println("Directory warning:", err)
			}
			continue
		}

		pack, err := fileSniff(filePath)
		if err != nil {
			fmt.Println("Detection warning: ", err)
			continue
		}

		err = detectFormat(&pack)
		if err != nil {
			fmt.Println("Detection failed: ", err)
		}

		status := packVerifier(&pack)

		fmt.Printf(
			"File: %-35s Size: %-10d Extension: %-8s Detected: %-8s Result: %s\n",
			filePath,
			pack.FileSize,
			pack.ext,
			pack.detectedFormat,
			status,
		)
		ihdrReport, err := IhdrCheck(&pack)
		if err != nil {
			return err
		}

		fmt.Println("IHDR Report")
		for i, r := range ihdrReport.Checks {
			fmt.Println(i, " - ", r)
		}

		idatReport, err := IDATcheck(&pack, &ihdrReport)
		if err != nil {
			return err
		}

		fmt.Println("IDAT Report")
		for i, r := range idatReport {
			fmt.Println(i, " - ", r.Checks)
		}

	}
	return nil
}

func main() {
	err := scanner("malformed_pngs")
	if err != nil {
		fmt.Println("Scanner error:", err)
	}
}
