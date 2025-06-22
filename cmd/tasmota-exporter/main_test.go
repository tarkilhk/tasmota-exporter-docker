package main

import (
	"math"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	promtest "github.com/prometheus/client_golang/prometheus/testutil"
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

func TestGetTodayValue(t *testing.T) {
	// We need to be able to control the time for this test
	originalNowFunc := getNow
	defer func() { getNow = originalNowFunc }()

	mockTasmotaTodayValue := 42.42

	tests := []struct {
		name        string
		mockTime    time.Time
		expectToday float64
	}{
		{
			name:        "10:00 - not in transition window",
			mockTime:    time.Date(2024, 7, 26, 10, 0, 0, 0, time.UTC),
			expectToday: 42.42,
		},
		{
			name:        "23:59 - in transition window",
			mockTime:    time.Date(2024, 7, 26, 23, 59, 15, 0, time.UTC),
			expectToday: 0,
		},
		{
			name:        "00:00 - in transition window",
			mockTime:    time.Date(2024, 7, 27, 0, 0, 30, 0, time.UTC),
			expectToday: 0,
		},
		{
			name:        "00:01 - not in transition window anymore",
			mockTime:    time.Date(2024, 7, 27, 0, 1, 0, 0, time.UTC),
			expectToday: 42.42,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getNow = func() time.Time { return tt.mockTime }
			got := getTodayValue(mockTasmotaTodayValue)
			if got != tt.expectToday {
				t.Errorf("incorrect today value: got %v, want %v", got, tt.expectToday)
			}
		})
	}
}

func TestHandleDailyLastMetric(t *testing.T) {
	// We need to be able to control the time for this test
	originalNowFunc := getNow
	defer func() { getNow = originalNowFunc }()

	// We need to reset the state for each test run
	originalLastDailyMetricSent := lastDailyMetricSent
	defer func() { lastDailyMetricSent = originalLastDailyMetricSent }()

	target := "test-target"
	mockPlug := TasmotaPlug{Today: 1.234}

	tests := []struct {
		name             string
		mockTime         time.Time
		setupSentMap     func()
		expectGaugeValue float64
	}{
		{
			name:             "not in window - should set gauge to NaN",
			mockTime:         time.Date(2024, 7, 26, 10, 0, 0, 0, time.UTC),
			expectGaugeValue: math.NaN(),
		},
		{
			name:             "in window, not sent yet - should set gauge to value",
			mockTime:         time.Date(2024, 7, 26, 23, 58, 0, 0, time.UTC),
			expectGaugeValue: mockPlug.Today,
		},
		{
			name:     "in window, but already sent today - should set gauge to NaN",
			mockTime: time.Date(2024, 7, 26, 23, 59, 0, 0, time.UTC),
			setupSentMap: func() {
				lastDailyMetricSent[target] = time.Date(2024, 7, 26, 23, 58, 0, 0, time.UTC)
			},
			expectGaugeValue: math.NaN(),
		},
		{
			name: "next day, in window - should set gauge to value again",
			setupSentMap: func() {
				lastDailyMetricSent[target] = time.Date(2024, 7, 26, 23, 58, 0, 0, time.UTC)
			},
			mockTime:         time.Date(2024, 7, 27, 23, 58, 0, 0, time.UTC),
			expectGaugeValue: mockPlug.Today,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset the state for a clean test run
			lastDailyMetricSent = make(map[string]time.Time)
			dailyLastGauge.Set(0) // Reset gauge to a known state

			if tt.setupSentMap != nil {
				tt.setupSentMap()
			}

			getNow = func() time.Time { return tt.mockTime }

			// Call the function that contains the logic we are testing
			handleDailyLastMetric(target, mockPlug)

			// Get the resulting metric value
			metricValue := promtest.ToFloat64(dailyLastGauge)

			if math.IsNaN(tt.expectGaugeValue) {
				if !math.IsNaN(metricValue) {
					t.Errorf("dailyLastGauge value = %v, want NaN", metricValue)
				}
			} else {
				if metricValue != tt.expectGaugeValue {
					t.Errorf("dailyLastGauge value = %v, want %v", metricValue, tt.expectGaugeValue)
				}
			}
		})
	}
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

func TestTodayValue_MidnightTransitionLogic(t *testing.T) {
	mockTasmotaData := `{t}</table><hr/>{t}{s}</th><th></th><th style='text-align:center'><th></th><td>{e}{s}Voltage{m}</td><td style='text-align:left'>237</td><td>&nbsp;</td><td> V{e}{s}Current{m}</td><td style='text-align:left'>0.053</td><td>&nbsp;</td><td> A{e}{s}Active Power{m}</td><td style='text-align:left'>7</td><td>&nbsp;</td><td> W{e}{s}Apparent Power{m}</td><td style='text-align:left'>13</td><td>&nbsp;</td><td> VA{e}{s}Reactive Power{m}</td><td style='text-align:left'>10</td><td>&nbsp;</td><td> VAr{e}{s}Power Factor{m}</td><td style='text-align:left'>0.59</td><td>&nbsp;</td><td>                         {e}{s}Energy Today{m}</td><td style='text-align:left'>42.42</td><td>&nbsp;</td><td> kWh{e}{s}Energy Yesterday{m}</td><td style='text-align:left'>0.016</td><td>&nbsp;</td><td> kWh{e}{s}Energy Total{m}</td><td style='text-align:left'>3.334</td><td>&nbsp;</td><td> kWh{e}</table><hr/>{t}</table>{t}<tr><td style='width:100%;text-align:center;font-weight:bold;font-size:62px'>ON</td></tr><tr></tr></table>`

	tp := parse(mockTasmotaData)

	tests := []struct {
		name     string
		time     time.Time
		expected float64
	}{
		{
			name:     "outside midnight window",
			time:     time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC),
			expected: tp.Today,
		},
		{
			name:     "inside midnight window (23:59:30)",
			time:     time.Date(2024, 1, 15, 23, 59, 30, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "inside midnight window (00:00:30)",
			time:     time.Date(2024, 1, 16, 0, 0, 30, 0, time.UTC),
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var today float64
			if isMidnightTransition(tt.time) {
				today = 0
			} else {
				today = tp.Today
			}
			if today != tt.expected {
				t.Errorf("For time %v, expected Today value %v, got %v", tt.time, tt.expected, today)
			}
		})
	}
}
