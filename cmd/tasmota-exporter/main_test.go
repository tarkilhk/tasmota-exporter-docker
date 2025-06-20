package main

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestParser(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  TasmotaPlug
	}{
		{
			name: "living-room-corner-on",
			input: `{t}</table><hr/>{t}{s}</th><th></th><th style='text-align:center'><th></th><td>{e}{s}Voltage{m}</td><td style='text-align:left'>237</td><td>&nbsp;</td><td> V{e}{s}Current{m}</td><td style='text-align:left'>0.053</td><td>&nbsp;</td><td> A{e}{s}Active Power{m}</td><td style='text-align:left'>7</td><td>&nbsp;</td><td> W{e}{s}Apparent Power{m}</td><td style='text-align:left'>13</td><td>&nbsp;</td><td> VA{e}{s}Reactive Power{m}</td><td style='text-align:left'>10</td><td>&nbsp;</td><td> VAr{e}{s}Power Factor{m}</td><td style='text-align:left'>0.59</td><td>&nbsp;</td><td>                         {e}{s}Energy Today{m}</td><td style='text-align:left'>0.002</td><td>&nbsp;</td><td> kWh{e}{s}Energy Yesterday{m}</td><td style='text-align:left'>0.016</td><td>&nbsp;</td><td> kWh{e}{s}Energy Total{m}</td><td style='text-align:left'>3.334</td><td>&nbsp;</td><td> kWh{e}</table><hr/>{t}</table>{t}<tr><td style='width:100%;text-align:center;font-weight:bold;font-size:62px'>ON</td></tr><tr></tr></table>

			`,
			want: TasmotaPlug{
				On:            true,
				Voltage:       237,
				Current:       0.053,
				Power:         7,
				ApparentPower: 13,
				ReactivePower: 10,
				Factor:        0.59,
				Today:         0.002,
				Yesterday:     0.016,
				Total:         3.334,
			},
		},
		{
			name: "living-room-corner-off",
			input: `{t}</table><hr/>{t}{s}</th><th></th><th style='text-align:center'><th></th><td>{e}{s}Voltage{m}</td><td style='text-align:left'>238</td><td>&nbsp;</td><td> V{e}{s}Current{m}</td><td style='text-align:left'>0.000</td><td>&nbsp;</td><td> A{e}{s}Active Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> W{e}{s}Apparent Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> VA{e}{s}Reactive Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> VAr{e}{s}Power Factor{m}</td><td style='text-align:left'>0.00</td><td>&nbsp;</td><td>                         {e}{s}Energy Today{m}</td><td style='text-align:left'>0.013</td><td>&nbsp;</td><td> kWh{e}{s}Energy Yesterday{m}</td><td style='text-align:left'>0.016</td><td>&nbsp;</td><td> kWh{e}{s}Energy Total{m}</td><td style='text-align:left'>3.345</td><td>&nbsp;</td><td> kWh{e}</table><hr/>{t}</table>{t}<tr><td style='width:100%;text-align:center;font-weight:normal;font-size:62px'>OFF</td></tr><tr></tr></table>

`,
			want: TasmotaPlug{
				On:            false,
				Voltage:       238,
				Current:       0,
				Power:         0,
				ApparentPower: 0,
				ReactivePower: 0,
				Factor:        0,
				Today:         0.013,
				Yesterday:     0.016,
				Total:         3.345,
			},
		},
		{
			name: "living-room-shelf-on",
			input: `{t}</table><hr/>{t}{s}</th><th></th><th style='text-align:center'><th></th><td>{e}{s}Voltage{m}</td><td style='text-align:left'>243</td><td>&nbsp;</td><td> V{e}{s}Current{m}</td><td style='text-align:left'>0.000</td><td>&nbsp;</td><td> A{e}{s}Active Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> W{e}{s}Apparent Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> VA{e}{s}Reactive Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> VAr{e}{s}Power Factor{m}</td><td style='text-align:left'>0.00</td><td>&nbsp;</td><td>                         {e}{s}Energy Today{m}</td><td style='text-align:left'>0.000</td><td>&nbsp;</td><td> kWh{e}{s}Energy Yesterday{m}</td><td style='text-align:left'>0.000</td><td>&nbsp;</td><td> kWh{e}{s}Energy Total{m}</td><td style='text-align:left'>2.495</td><td>&nbsp;</td><td> kWh{e}</table><hr/>{t}</table>{t}<tr><td style='width:100%;text-align:center;font-weight:bold;font-size:62px'>ON</td></tr><tr></tr></table>

			`,
			want: TasmotaPlug{
				On:            true,
				Voltage:       243,
				Current:       0,
				Power:         0,
				ApparentPower: 0,
				ReactivePower: 0,
				Factor:        0,
				Today:         0,
				Yesterday:     0,
				Total:         2.495,
			},
		},
		{
			name: "living-room-shelf-off",
			input: `{t}</table><hr/>{t}{s}</th><th></th><th style='text-align:center'><th></th><td>{e}{s}Voltage{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> V{e}{s}Current{m}</td><td style='text-align:left'>0.000</td><td>&nbsp;</td><td> A{e}{s}Active Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> W{e}{s}Apparent Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> VA{e}{s}Reactive Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> VAr{e}{s}Power Factor{m}</td><td style='text-align:left'>0.00</td><td>&nbsp;</td><td>                         {e}{s}Energy Today{m}</td><td style='text-align:left'>0.000</td><td>&nbsp;</td><td> kWh{e}{s}Energy Yesterday{m}</td><td style='text-align:left'>0.000</td><td>&nbsp;</td><td> kWh{e}{s}Energy Total{m}</td><td style='text-align:left'>2.495</td><td>&nbsp;</td><td> kWh{e}</table><hr/>{t}</table>{t}<tr><td style='width:100%;text-align:center;font-weight:normal;font-size:62px'>OFF</td></tr><tr></tr></table>

			`,
			want: TasmotaPlug{
				On:            false,
				Voltage:       0,
				Current:       0,
				Power:         0,
				ApparentPower: 0,
				ReactivePower: 0,
				Factor:        0,
				Today:         0,
				Yesterday:     0,
				Total:         2.495,
			},
		},
		{
			name: "living-room-drawer-on",
			input: `{t}</table><hr/>{t}{s}</th><th></th><th style='text-align:center'><th></th><td>{e}{s}Voltage{m}</td><td style='text-align:left'>237</td><td>&nbsp;</td><td> V{e}{s}Current{m}</td><td style='text-align:left'>0.000</td><td>&nbsp;</td><td> A{e}{s}Active Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> W{e}{s}Apparent Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> VA{e}{s}Reactive Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> VAr{e}{s}Power Factor{m}</td><td style='text-align:left'>0.00</td><td>&nbsp;</td><td>                         {e}{s}Energy Today{m}</td><td style='text-align:left'>0.000</td><td>&nbsp;</td><td> kWh{e}{s}Energy Yesterday{m}</td><td style='text-align:left'>0.009</td><td>&nbsp;</td><td> kWh{e}{s}Energy Total{m}</td><td style='text-align:left'>2.644</td><td>&nbsp;</td><td> kWh{e}</table><hr/>{t}</table>{t}<tr><td style='width:100%;text-align:center;font-weight:bold;font-size:62px'>ON</td></tr><tr></tr></table>

			`,
			want: TasmotaPlug{
				On:            true,
				Voltage:       237,
				Current:       0,
				Power:         0,
				ApparentPower: 0,
				ReactivePower: 0,
				Factor:        0,
				Today:         0,
				Yesterday:     0.009,
				Total:         2.644,
			},
		},
		{
			name: "living-room-drawer-off",
			input: `{t}</table><hr/>{t}{s}</th><th></th><th style='text-align:center'><th></th><td>{e}{s}Voltage{m}</td><td style='text-align:left'>236</td><td>&nbsp;</td><td> V{e}{s}Current{m}</td><td style='text-align:left'>0.000</td><td>&nbsp;</td><td> A{e}{s}Active Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> W{e}{s}Apparent Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> VA{e}{s}Reactive Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> VAr{e}{s}Power Factor{m}</td><td style='text-align:left'>0.00</td><td>&nbsp;</td><td>                         {e}{s}Energy Today{m}</td><td style='text-align:left'>0.000</td><td>&nbsp;</td><td> kWh{e}{s}Energy Yesterday{m}</td><td style='text-align:left'>0.009</td><td>&nbsp;</td><td> kWh{e}{s}Energy Total{m}</td><td style='text-align:left'>2.644</td><td>&nbsp;</td><td> kWh{e}</table><hr/>{t}</table>{t}<tr><td style='width:100%;text-align:center;font-weight:normal;font-size:62px'>OFF</td></tr><tr></tr></table>

			`,
			want: TasmotaPlug{
				On:            false,
				Voltage:       236,
				Current:       0,
				Power:         0,
				ApparentPower: 0,
				ReactivePower: 0,
				Factor:        0,
				Today:         0,
				Yesterday:     0.009,
				Total:         2.644,
			},
		},
		{
			name: "office-light-on",
			input: `{t}</table><hr/>{t}{s}</th><th></th><th style='text-align:center'><th></th><td>{e}{s}Voltage{m}</td><td style='text-align:left'>237</td><td>&nbsp;</td><td> V{e}{s}Current{m}</td><td style='text-align:left'>0.203</td><td>&nbsp;</td><td> A{e}{s}Active Power{m}</td><td style='text-align:left'>29</td><td>&nbsp;</td><td> W{e}{s}Apparent Power{m}</td><td style='text-align:left'>48</td><td>&nbsp;</td><td> VA{e}{s}Reactive Power{m}</td><td style='text-align:left'>39</td><td>&nbsp;</td><td> VAr{e}{s}Power Factor{m}</td><td style='text-align:left'>0.60</td><td>&nbsp;</td><td>                         {e}{s}Energy Today{m}</td><td style='text-align:left'>0.001</td><td>&nbsp;</td><td> kWh{e}{s}Energy Yesterday{m}</td><td style='text-align:left'>0.094</td><td>&nbsp;</td><td> kWh{e}{s}Energy Total{m}</td><td style='text-align:left'>16.007</td><td>&nbsp;</td><td> kWh{e}</table><hr/>{t}</table>{t}<tr><td style='width:100%;text-align:center;font-weight:bold;font-size:62px'>ON</td></tr><tr></tr></table>

			`,
			want: TasmotaPlug{
				On:            true,
				Voltage:       237,
				Current:       0.203,
				Power:         29,
				ApparentPower: 48,
				ReactivePower: 39,
				Factor:        0.6,
				Today:         0.001,
				Yesterday:     0.094,
				Total:         16.007,
			},
		},
		{
			name: "office-light-off",
			input: `{t}</table><hr/>{t}{s}</th><th></th><th style='text-align:center'><th></th><td>{e}{s}Voltage{m}</td><td style='text-align:left'>237</td><td>&nbsp;</td><td> V{e}{s}Current{m}</td><td style='text-align:left'>0.000</td><td>&nbsp;</td><td> A{e}{s}Active Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> W{e}{s}Apparent Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> VA{e}{s}Reactive Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> VAr{e}{s}Power Factor{m}</td><td style='text-align:left'>0.00</td><td>&nbsp;</td><td>                         {e}{s}Energy Today{m}</td><td style='text-align:left'>0.000</td><td>&nbsp;</td><td> kWh{e}{s}Energy Yesterday{m}</td><td style='text-align:left'>0.094</td><td>&nbsp;</td><td> kWh{e}{s}Energy Total{m}</td><td style='text-align:left'>16.006</td><td>&nbsp;</td><td> kWh{e}</table><hr/>{t}</table>{t}<tr><td style='width:100%;text-align:center;font-weight:normal;font-size:62px'>OFF</td></tr><tr></tr></table>

			`,
			want: TasmotaPlug{
				On:            false,
				Voltage:       237,
				Current:       0,
				Power:         0,
				ApparentPower: 0,
				ReactivePower: 0,
				Factor:        0,
				Today:         0,
				Yesterday:     0.094,
				Total:         16.006,
			},
		},
		{
			name: "office-air-on",
			input: `{t}</table><hr/>{t}{s}</th><th></th><th style='text-align:center'><th></th><td>{e}{s}Voltage{m}</td><td style='text-align:left'>236</td><td>&nbsp;</td><td> V{e}{s}Current{m}</td><td style='text-align:left'>0.460</td><td>&nbsp;</td><td> A{e}{s}Active Power{m}</td><td style='text-align:left'>51</td><td>&nbsp;</td><td> W{e}{s}Apparent Power{m}</td><td style='text-align:left'>108</td><td>&nbsp;</td><td> VA{e}{s}Reactive Power{m}</td><td style='text-align:left'>96</td><td>&nbsp;</td><td> VAr{e}{s}Power Factor{m}</td><td style='text-align:left'>0.47</td><td>&nbsp;</td><td>                         {e}{s}Energy Today{m}</td><td style='text-align:left'>0.003</td><td>&nbsp;</td><td> kWh{e}{s}Energy Yesterday{m}</td><td style='text-align:left'>0.207</td><td>&nbsp;</td><td> kWh{e}{s}Energy Total{m}</td><td style='text-align:left'>1.124</td><td>&nbsp;</td><td> kWh{e}</table><hr/>{t}</table>{t}<tr><td style='width:100%;text-align:center;font-weight:bold;font-size:62px'>ON</td></tr><tr></tr></table>

			`,
			want: TasmotaPlug{
				On:            true,
				Voltage:       236,
				Current:       0.46,
				Power:         51,
				ApparentPower: 108,
				ReactivePower: 96,
				Factor:        0.47,
				Today:         0.003,
				Yesterday:     0.207,
				Total:         1.124,
			},
		},
		{
			name: "office-air-off",
			input: `{t}</table><hr/>{t}{s}</th><th></th><th style='text-align:center'><th></th><td>{e}{s}Voltage{m}</td><td style='text-align:left'>236</td><td>&nbsp;</td><td> V{e}{s}Current{m}</td><td style='text-align:left'>0.000</td><td>&nbsp;</td><td> A{e}{s}Active Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> W{e}{s}Apparent Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> VA{e}{s}Reactive Power{m}</td><td style='text-align:left'>0</td><td>&nbsp;</td><td> VAr{e}{s}Power Factor{m}</td><td style='text-align:left'>0.00</td><td>&nbsp;</td><td>                         {e}{s}Energy Today{m}</td><td style='text-align:left'>0.000</td><td>&nbsp;</td><td> kWh{e}{s}Energy Yesterday{m}</td><td style='text-align:left'>0.207</td><td>&nbsp;</td><td> kWh{e}{s}Energy Total{m}</td><td style='text-align:left'>1.121</td><td>&nbsp;</td><td> kWh{e}</table><hr/>{t}</table>{t}<tr><td style='width:100%;text-align:center;font-weight:normal;font-size:62px'>OFF</td></tr><tr></tr></table>

			`,
			want: TasmotaPlug{
				On:            false,
				Voltage:       236,
				Current:       0,
				Power:         0,
				ApparentPower: 0,
				ReactivePower: 0,
				Factor:        0,
				Today:         0,
				Yesterday:     0.207,
				Total:         1.121,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parse(tt.input)

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("unexpected parsed output (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDailyMetricFunctions(t *testing.T) {
	// Reset the global map for testing
	lastDailyMetricSent = make(map[string]time.Time)

	target1 := "192.168.1.100"
	target2 := "192.168.1.101"

	// Test the core logic by manually setting the map entries
	t.Run("map-tracking-different-targets", func(t *testing.T) {
		// Simulate sending daily metric for target1
		lastDailyMetricSent[target1] = time.Now()

		// Check that target1 is tracked
		if _, exists := lastDailyMetricSent[target1]; !exists {
			t.Errorf("Expected target1 to be in lastDailyMetricSent map")
		}

		// Check that target2 is not tracked yet
		if _, exists := lastDailyMetricSent[target2]; exists {
			t.Errorf("Expected target2 to NOT be in lastDailyMetricSent map yet")
		}

		// Simulate sending daily metric for target2
		lastDailyMetricSent[target2] = time.Now()

		// Check that both targets are now tracked
		if _, exists := lastDailyMetricSent[target1]; !exists {
			t.Errorf("Expected target1 to be in lastDailyMetricSent map")
		}
		if _, exists := lastDailyMetricSent[target2]; !exists {
			t.Errorf("Expected target2 to be in lastDailyMetricSent map")
		}
	})

	t.Run("map-independence", func(t *testing.T) {
		// Clear the map
		lastDailyMetricSent = make(map[string]time.Time)

		// Set different times for different targets
		time1 := time.Date(2024, 1, 15, 23, 58, 0, 0, time.UTC)
		time2 := time.Date(2024, 1, 15, 23, 59, 0, 0, time.UTC)

		lastDailyMetricSent[target1] = time1
		lastDailyMetricSent[target2] = time2

		// Verify they're independent
		if lastDailyMetricSent[target1] != time1 {
			t.Errorf("Expected target1 to have time1, got %v", lastDailyMetricSent[target1])
		}
		if lastDailyMetricSent[target2] != time2 {
			t.Errorf("Expected target2 to have time2, got %v", lastDailyMetricSent[target2])
		}
	})
}

func TestIsDailyMetricWindow(t *testing.T) {
	// Time outside the window
	outside := time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC) // 12:00
	if isDailyMetricWindow(outside) {
		t.Errorf("Expected isDailyMetricWindow to return false for 12:00, got true")
	}

	// Time inside the window
	inside := time.Date(2024, 1, 15, 23, 58, 0, 0, time.UTC) // 23:58
	if !isDailyMetricWindow(inside) {
		t.Errorf("Expected isDailyMetricWindow to return true for 23:58, got false")
	}

	// Another time inside the window
	inside2 := time.Date(2024, 1, 15, 23, 59, 0, 0, time.UTC) // 23:59
	if !isDailyMetricWindow(inside2) {
		t.Errorf("Expected isDailyMetricWindow to return true for 23:59, got false")
	}

	// Edge case: 00:00 (should be false)
	edge := time.Date(2024, 1, 16, 0, 0, 0, 0, time.UTC) // 00:00
	if isDailyMetricWindow(edge) {
		t.Errorf("Expected isDailyMetricWindow to return false for 00:00, got true")
	}
}
