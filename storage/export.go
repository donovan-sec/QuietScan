package storage

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"quietscan/types"
)

// sanitizeCSVCell neutralizes CSV formula injection. Spreadsheet apps (Excel,
// LibreOffice, Sheets) execute any cell value beginning with =, +, -, @, or a
// leading tab/carriage return as a formula. Vendor and Hostname values come from
// network responses, so a hostile device on the LAN could set a hostname like
// `=cmd|'/c calc'!A1` and have it run when the exported CSV is opened. The
// encoding/csv writer escapes for CSV parsing but does nothing about formula
// interpretation, so we prefix any dangerous value with a single quote, which
// every major spreadsheet treats as "render as literal text".
func sanitizeCSVCell(value string) string {
	if value == "" {
		return value
	}
	switch value[0] {
	case '=', '+', '-', '@', '\t', '\r':
		return "'" + value
	}
	return value
}

// sanitizeCSVRecord applies sanitizeCSVCell to every field in a record.
func sanitizeCSVRecord(record []string) []string {
	out := make([]string, len(record))
	for i, field := range record {
		out[i] = sanitizeCSVCell(field)
	}
	return out
}

var allowOverwrite = false // Set via SetAllowOverwrite

// SetAllowOverwrite sets whether existing files can be overwritten
func SetAllowOverwrite(allow bool) {
	allowOverwrite = allow
}

// ExportToCSV exports the scan result to a CSV file in the current working directory
func ExportToCSV(result *types.ScanResult) (string, error) {
	if result == nil || len(result.Devices) == 0 {
		return "", fmt.Errorf("no scan results to export")
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	// Generate filename with timestamp
	timestamp := result.Timestamp.Format("20060102-150405")
	filename := fmt.Sprintf("quietscan-export-%s.csv", timestamp)
	filepath := filepath.Join(cwd, filename)

	// Validate file path before creating (consistent with ExportToCSVWithMetadata)
	if err := ValidateFilePath(filepath, allowOverwrite); err != nil {
		return "", fmt.Errorf("cannot write to output path: %v", err)
	}

	// Create CSV file
	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"IP Address", "MAC Address", "Vendor", "Hostname"}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %v", err)
	}

	// Write device data
	for _, device := range result.Devices {
		record := []string{
			device.IP,
			device.MAC,
			device.Vendor,
			device.Hostname,
		}
		if err := writer.Write(sanitizeCSVRecord(record)); err != nil {
			return "", fmt.Errorf("failed to write CSV record: %v", err)
		}
	}

	return filepath, nil
}

// ExportToCSVWithMetadata exports the scan result to a CSV file with metadata header
func ExportToCSVWithMetadata(result *types.ScanResult) (string, error) {
	if result == nil || len(result.Devices) == 0 {
		return "", fmt.Errorf("no scan results to export")
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	// Generate filename with timestamp
	timestamp := result.Timestamp.Format("20060102-150405")
	filename := fmt.Sprintf("quietscan-export-%s.csv", timestamp)
	filepath := filepath.Join(cwd, filename)

	// Validate file path before creating
	if err := ValidateFilePath(filepath, allowOverwrite); err != nil {
		return "", fmt.Errorf("cannot write to output path: %v", err)
	}

	// Check if file exists and log overwrite warning
	if CheckFileExists(filepath) && allowOverwrite {
		fmt.Fprintf(os.Stderr, "Warning: Overwriting existing file: %s\n", filepath)
	}

	// Create CSV file
	file, err := os.Create(filepath)
	if err != nil {
		return "", fmt.Errorf("failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write metadata as header rows (using a comment-like format)
	// Note: Standard CSV doesn't support comments, so we'll add metadata as regular rows
	// that can be easily identified
	metadata := []string{
		fmt.Sprintf("Scan Date: %s", result.Timestamp.Format(time.RFC1123)),
		fmt.Sprintf("Adapter: %s", result.Adapter),
		fmt.Sprintf("Subnet: %s", result.Subnet),
		fmt.Sprintf("Device Count: %d", len(result.Devices)),
		"", // Empty row separator
	}

	for _, meta := range metadata {
		if err := writer.Write([]string{meta}); err != nil {
			return "", fmt.Errorf("failed to write metadata: %v", err)
		}
	}

	// Write header
	header := []string{"IP Address", "MAC Address", "Vendor", "Hostname"}
	if err := writer.Write(header); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %v", err)
	}

	// Write device data
	for _, device := range result.Devices {
		record := []string{
			device.IP,
			device.MAC,
			device.Vendor,
			device.Hostname,
		}
		if err := writer.Write(sanitizeCSVRecord(record)); err != nil {
			return "", fmt.Errorf("failed to write CSV record: %v", err)
		}
	}

	return filepath, nil
}
