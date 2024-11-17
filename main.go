package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
//	"os"
	"os/exec"
)

// SatelliteData represents the output from the Python script
type SatelliteData struct {
	Position [3]float64 `json:"position"`
	Velocity [3]float64 `json:"velocity"`
	Error    string     `json:"error,omitempty"`
}

// WriteTLEsToFile writes TLE data to a temporary file and returns the file path
func WriteTLEsToFile(tles []map[string]string) (string, error) {
	// Create a temporary file
	tmpFile, err := ioutil.TempFile("", "tle_data_*.json")
	if err != nil {
		return "", fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer tmpFile.Close()

	// Marshal TLE data into JSON
	tleData, err := json.Marshal(tles)
	if err != nil {
		return "", fmt.Errorf("failed to serialize TLEs: %w", err)
	}

	// Write the TLE data to the file
	if _, err := tmpFile.Write(tleData); err != nil {
		return "", fmt.Errorf("failed to write TLEs to file: %w", err)
	}

	// Return the file path for use in the Python script
	return tmpFile.Name(), nil
}

// CalculatePositions calls the Python script to calculate satellite positions based on TLE data
func CalculatePositions(tles []map[string]string) ([]SatelliteData, error) {
	// Write TLEs to a temporary file
	tleFile, err := WriteTLEsToFile(tles)
	if err != nil {
		return nil, err
	}

	// Pass the file path to the Python script
	cmd := exec.Command("python3", "tle_processor.py", tleFile)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Run the command and capture the output
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to process TLEs: %v, stderr: %s", err, stderr.String())
	}

	// Unmarshal the Python script output into SatelliteData
	var positions []SatelliteData
	err = json.Unmarshal(out.Bytes(), &positions)
	if err != nil {
		return nil, fmt.Errorf("failed to parse satellite positions: %w", err)
	}

	return positions, nil
}

// Example TLEs
func getExampleTLEs() []map[string]string {
	return []map[string]string{
		{
			"TLE_LINE1": "1 25544U 98067A   23274.70850694  .00002431  00000-0  51909-4 0  9995",
			"TLE_LINE2": "2 25544  51.6415  23.2393 0004089  52.7238  18.2595 15.53187687632727",
		},
		{
			"TLE_LINE1": "1 43212U 99012B   23275.56350988  .00003141  00000-0  63272-4 0  9991",
			"TLE_LINE2": "2 43212  72.9781  15.6894 0003971  47.2410  66.5753 15.72026516322512",
		},
	}
}

func main() {
	// Fetch the example TLE data
	tles := getExampleTLEs()

	// Calculate satellite positions using the TLE data
	positions, err := CalculatePositions(tles)
	if err != nil {
		fmt.Printf("Error calculating positions: %v\n", err)
		return
	}

	// Print the satellite positions
	fmt.Println("Satellite Positions:")
	for _, pos := range positions {
		fmt.Printf("Position: %v, Velocity: %v\n", pos.Position, pos.Velocity)
	}
}
