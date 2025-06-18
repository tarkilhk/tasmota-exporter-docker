package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"tailscale.com/envknob"
)

var overrideListenAddr = envknob.String("TASMOTA_EXPORTER_LISTEN_ADDR")

var (
	onGauge,
	voltageGauge,
	currentGauge,
	powerGauge,
	apparentPowerGauge,
	reactivePowerGauge,
	factorGauge,
	todayGauge,
	yesterdayGauge,
	totalGauge,
	dailyLastGauge prometheus.Gauge

	registry *prometheus.Registry

	// Track the last day we sent the daily last metric
	lastDailyMetricSent time.Time
)

func init() {
	onGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tasmota_on",
		Help: "Indicates if the tasmota plug is on/off",
	})
	voltageGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tasmota_voltage_volts",
		Help: "voltage of tasmota plug in volt (V)",
	})
	currentGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tasmota_current_amperes",
		Help: "current of tasmota plug in ampere (A)",
	})
	powerGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tasmota_power_watts",
		Help: "current power of tasmota plug in watts (W)",
	})
	apparentPowerGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tasmota_apparent_power_voltamperes",
		Help: "apparent power of tasmota plug in volt-amperes (VA)",
	})
	reactivePowerGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tasmota_reactive_power_voltamperesreactive",
		Help: "reactive power of tasmota plug in volt-amperes reactive (VAr)",
	})
	factorGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tasmota_power_factor",
		Help: "current power factor of tasmota plug",
	})
	todayGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tasmota_today_kwh_total",
		Help: "todays energy usage total in kilowatts hours (kWh)",
	})
	yesterdayGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tasmota_yesterday_kwh_total",
		Help: "yesterdays energy usage total in kilowatts hours (kWh)",
	})
	totalGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tasmota_kwh_total",
		Help: "total energy usage in kilowatts hours (kWh)",
	})
	dailyLastGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tasmota_daily_last_kwh_total",
		Help: "The last kWh reading of the day, sent once per day between 23:58:00 and 23:59:59",
	})

	registry = prometheus.NewRegistry()
	registry.MustRegister(onGauge)
	registry.MustRegister(voltageGauge)
	registry.MustRegister(currentGauge)
	registry.MustRegister(powerGauge)
	registry.MustRegister(apparentPowerGauge)
	registry.MustRegister(reactivePowerGauge)
	registry.MustRegister(factorGauge)
	registry.MustRegister(todayGauge)
	registry.MustRegister(yesterdayGauge)
	registry.MustRegister(totalGauge)
	registry.MustRegister(dailyLastGauge)
}

func main() {
	http.HandleFunc("/probe", tasmotaHandler)

	listenAddr := ":9090"
	if overrideListenAddr != "" {
		listenAddr = overrideListenAddr
	}

	log.Printf("starting tasmota exporter on %s", listenAddr)
	err := http.ListenAndServe(listenAddr, nil)
	if errors.Is(err, http.ErrServerClosed) {
		log.Printf("server closed")
	} else if err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}

func tasmotaHandler(w http.ResponseWriter, r *http.Request) {
	probeSuccessGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "probe_success",
		Help: "Displays whether or not the probe was a success",
	})
	probeDurationGauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "probe_duration_seconds",
		Help: "Returns how long the probe took to complete in seconds",
	})

	params := r.URL.Query()

	target := params.Get("target")
	if target == "" {
		http.Error(w, "Target parameter is missing", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	r = r.WithContext(ctx)

	start := time.Now()
	success := probeTasmota(ctx, target, registry)
	duration := time.Since(start).Seconds()
	probeDurationGauge.Set(duration)
	if success {
		probeSuccessGauge.Set(1)
		log.Printf("%s: probe succeeded, duration: %fs", target, duration)
	} else {
		log.Printf("%s: probe failed, duration: %fs", target, duration)
	}

	h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
	h.ServeHTTP(w, r)
}

// isMidnightTransition checks if we're in the window around midnight (23:59:00 to 00:01:00).
// This is necessary because the tasmota_today_kwh_total metric from the device carries over
// to the next day until the next scrape happens. By forcing it to 0 during this transition
// period, we ensure that each day's maximum consumption is correctly attributed to its
// proper day in Grafana, preventing the previous day's value from being incorrectly
// associated with the new day.
func isMidnightTransition() bool {
	now := time.Now()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()

	// Check if time is between 23:59:00 and 00:01:00
	// time.Hour() method in Go's standard library always returns in 24-hour format (0-23), independent of the system's locale or time format settings
	// This is why we can use 23 and 0 for the hour check
	if hour == 23 && minute == 59 && second >= 00 {
		return true
	}
	if hour == 0 && minute == 1 && second <= 00 {
		return true
	}
	return false
}

// isDailyMetricWindow checks if we're in the window to send the daily last metric (23:58:00 to 23:59:59)
func isDailyMetricWindow() bool {
	now := time.Now()
	hour := now.Hour()
	minute := now.Minute()
	second := now.Second()

	return hour == 23 && (minute == 58 || minute == 59)
}

// shouldSendDailyMetric checks if we should send the daily metric
// Returns true if we're in the time window and haven't sent it today
func shouldSendDailyMetric() bool {
	now := time.Now()
	
	// If we're not in the time window, don't send
	if !isDailyMetricWindow() {
		return false
	}

	// If we haven't sent it today, we should send it
	if lastDailyMetricSent.Day() != now.Day() || 
	   lastDailyMetricSent.Month() != now.Month() || 
	   lastDailyMetricSent.Year() != now.Year() {
		return true
	}

	return false
}

func probeTasmota(ctx context.Context, target string, registry *prometheus.Registry) (success bool) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(fmt.Sprintf("http://%s?m", target))
	if err != nil {
		log.Printf("failed to query tasmota target (%s): %s", target, err)
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("failed to read data from tasmota target (%s): %s", target, err)
		return false
	}

	tp := parse(string(body))

	if tp.On {
		onGauge.Set(1)
	} else {
		onGauge.Set(0)
	}
	voltageGauge.Set(tp.Voltage)
	currentGauge.Set(tp.Current)
	powerGauge.Set(tp.Power)
	apparentPowerGauge.Set(tp.ApparentPower)
	reactivePowerGauge.Set(tp.ReactivePower)
	factorGauge.Set(tp.Factor)
	
	// Set todayGauge to 0 during midnight transition period
	if isMidnightTransition() {
		todayGauge.Set(0)
	} else {
		todayGauge.Set(tp.Today)
	}
	
	yesterdayGauge.Set(tp.Yesterday)
	totalGauge.Set(tp.Total)

	// Handle daily last metric
	if shouldSendDailyMetric() {
		dailyLastGauge.Set(tp.Today)
		lastDailyMetricSent = time.Now()
	}

	return true
}

type TasmotaPlug struct {
	// On indicates if the plug is on or off.
	On bool `json:"On"`

	// Voltage describes the voltage used of the appliance
	// denoted in V.
	Voltage float64 `json:"Voltage"`

	// Current describes the amount of amperes used, denoted
	// in A.
	Current float64 `json:"Current"`

	// Power describes the current power used, denoted in W (watt)
	Power float64 `json:"Power"`

	// ApparentPower describes the volt-ampere (VA)
	ApparentPower float64 `json:"ApparentPower"`

	// ReactivePower describes Volt-Amps Reactive (VAr)
	ReactivePower float64 `json:"ReactivePower"`

	// Factor describes the power factor
	Factor float64 `json:"Factor"`

	// Today is the total usage of energy in kilowatts hours (kWh)
	// meassured by the internal clock of the plug for today.
	Today float64 `json:"Today"`

	// Yesterday is the total usage of energy in kilowatts hours (kWh)
	// meassured by the internal clock of the plug for yesterday.
	Yesterday float64 `json:"Yesterday"`

	// Total is the total usage of energy in kilowatts hours (kWh)
	// since the plug was last factory reset.
	Total float64 `json:"Total"`
}

func parse(input string) TasmotaPlug {
	ret := TasmotaPlug{
		On: strings.Contains(input, "ON"),
	}

	rows := strings.Split(input, "{s}")
	for _, row := range rows {
		rowRaw := strings.Split(row, "{m}")

		if len(rowRaw) < 2 {
			continue
		}

		label := rowRaw[0]
		valueRaw := rowRaw[1]

		valueSplit := strings.Split(valueRaw, "{e}")

		if len(valueSplit) == 0 {
			continue
		}

		valueStrWithUnit := valueSplit[0]
		if strings.Contains(valueStrWithUnit, "<td") {
			valueStrWithUnit = strings.ReplaceAll(valueStrWithUnit, "</td><td style='text-align:left'>", "")
			valueStrWithUnit = strings.ReplaceAll(valueStrWithUnit, "</td><td>&nbsp;</td><td>", "")
		}

		valueSplitWithUnit := strings.Split(valueStrWithUnit, " ")
		if len(valueSplitWithUnit) == 0 {
			continue
		}

		value, err := strconv.ParseFloat(valueSplitWithUnit[0], 64)
		if err != nil {
			continue
		}

		switch label {
		case "Voltage":
			ret.Voltage = value
		case "Current":
			ret.Current = value
		case "Active Power":
			ret.Power = value
		case "Apparent Power":
			ret.ApparentPower = value
		case "Reactive Power":
			ret.ReactivePower = value
		case "Power Factor":
			ret.Factor = value
		case "Energy Today":
			ret.Today = value
		case "Energy Yesterday":
			ret.Yesterday = value
		case "Energy Total":
			ret.Total = value
		default:
			log.Printf("unable to match label, got: %s, value: %f", label, value)

		}
	}

	return ret
}
