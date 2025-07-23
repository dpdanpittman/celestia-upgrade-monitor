package main

import "github.com/prometheus/client_golang/prometheus"

var (
	GrpcServerAddress      string
	HttpServerPort         string
	RequiredThresholdPower float64 = 0.80
)

var (
	upgradeStatus = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "celestia_upgrade_status",
			Help: "Upgrade status as reported by celestia-app signal service, this is 1 if signal quorom is reached and upgrade is happening, 0 otherwise",
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
	tallyThresholdPower = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "celestia_tally_threshold_power",
			Help: "Threshold power signalled for the upgrade",
		},
	)
	tallyTotalVotingPower = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "celestia_tally_total_voting_power",
			Help: "Total voting power in the network",
		},
	)
	tallyThresholdPercent = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "celestia_tally_threshold_percent",
			Help: "Threshold percent signalled for the upgrade",
		},
	)
)

type UpgradeData struct {
	UpgradeData UpgradeResponse `json:"upgrade_data"`
	TallyData   TallyResponse   `json:"tally_data"`
}

type UpgradeResponse struct {
	Upgrade Upgrade `json:"upgrade"`
}

type Upgrade struct {
	AppVersion    int   `json:"app_version"`
	UpgradeHeight int64 `json:"upgrade_height"`
}

type TallyResponse struct {
	TotalVotingPower int64   `json:"total_voting_power"`
	ThresholdPower   int64   `json:"threshold_power"`
	ThresholdPercent float64 `json:"threshold_percent"`
}
