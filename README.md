# Celestia Upgrade Monitor

A lightweight Go service that connects to a Celestia `celestia-appd` gRPC node to monitor upcoming upgrades and expose the data:

- As a JSON API (`/upgrade`)
- As Prometheus metrics (`/metrics`)

Useful for dashboards, alerts, or integrations that track scheduled upgrades in Celestia.

---

## 🚀 Features

- Periodic polling of the Celestia `signal` gRPC service (every 30 minutes)
- JSON endpoint at `/upgrade`
- Prometheus metrics at `/metrics`
- Runs a single HTTP server with both endpoints

---

## 🛠 Installation

1. **Clone the repository**:

   ```bash
   git clone https://github.com/sryps/celestia-upgrade-monitor.git
   cd celestia-upgrade-monitor

   ```

2. **Build the binary**:

   ```bash
    go build -o celestia-upgrade-monitor
   ```

3. **Run the service**:

   ```bash
    ./celestia-upgrade-monitor -grpc-addr <GRPC_ENDPOINT> -server-port <PORT>
   ```

4. **Access the endpoints**:
   - JSON API: `http://<ADDRESS>:<PORT>/upgrade`
   - Prometheus metrics: `http://<ADDRESS>:<PORT>/metrics`

## 📈 Prometheus Metrics

The following metrics are exposed:

```plaintext
# HELP celestia_tally_threshold_percent Threshold percent signalled for the upgrade
# TYPE celestia_tally_threshold_percent gauge
celestia_tally_threshold_percent 83.33333344006462
# HELP celestia_tally_threshold_power Threshold power signalled for the upgrade
# TYPE celestia_tally_threshold_power gauge
celestia_tally_threshold_power 5.20518014e+08
# HELP celestia_tally_total_voting_power Total voting power in the network
# TYPE celestia_tally_total_voting_power gauge
celestia_tally_total_voting_power 6.24621616e+08
# HELP celestia_upgrade_height Height at which the upgrade will take place
# TYPE celestia_upgrade_height gauge
celestia_upgrade_height 6.680339e+06
# HELP celestia_upgrade_status Upgrade status as reported by celestia-app signal service, this is 1 if signal quorom is reached and upgrade is happening, 0 otherwise
# TYPE celestia_upgrade_status gauge
celestia_upgrade_status 1
# HELP celestia_upgrade_version Current upgrade version
# TYPE celestia_upgrade_version gauge
celestia_upgrade_version 4
```

## 🛠 RPC JSON API

// TODO: Add RPC JSON API details for tally data
The JSON API at `/upgrade` provides the following structure:

```json
{
  "upgrade_data": {
    "upgrade": {
      "app_version": 4,
      "upgrade_height": 6680339
    }
  },
  "tally_data": {
    "total_voting_power": 624621492,
    "threshold_power": 520517910,
    "threshold_percent": 0.8333333333333334
  }
}
```

## 📋 Requirements

To build and run the Celestia Upgrade Monitor, you'll need:

- **Go 1.20 or later** – for compiling and running the application
- **Access to a Celestia node** (`celestia-appd`) with:
  - The gRPC server enabled and reachable (default port is often `9090`, but your deployment may vary)
  - The `celestia/signal/v1` gRPC service available

## 📝 License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT).
