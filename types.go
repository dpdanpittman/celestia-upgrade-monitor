package main

import "github.com/prometheus/client_golang/prometheus"

var (
	GrpcServerAddress      string
	HttpServerPort         string
	RequiredThresholdPower float64 = 0.80
	upgradeStatus                  = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "celestia_upgrade_status",
			Help: "Upgrade status as reported by celestia-app signal service, this is 1 if upgrade if signal threshold and upgrade is happening, 0 otherwise",
		},
	)
	upgradeVersion = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "celestia_upgrade_version",
			Help: "Current upgrade version",
		},
	)
	upgradeHeight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "celestia_upgrade_height",
			Help: "Height at which the upgrade will take place",
		},
	)
)

type UpgradeResponse struct {
	Upgrade struct {
		AppVersion    int   `json:"app_version"`
		UpgradeHeight int64 `json:"upgrade_height"`
	} `json:"upgrade"`
}

type TallyResponse struct {
	TotalVotingPower int64   `json:"total_voting_power"`
	ThresholdPower   int64   `json:"threshold_power"`
	ThresholdPercent float64 `json:"threshold_percent"`
}

type UpgradeData struct {
	UpgradeData UpgradeResponse `json:"upgrade_data"`
	TallyData   TallyResponse   `json:"tally_data"`
}
