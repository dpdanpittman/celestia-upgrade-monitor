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
   git clone https://github.com/your-org/celestia-upgrade-monitor.git
   cd celestia-upgrade-monitor

   ```

2. **Build the binary**:

   ```bash
    go build -o celestia-upgrade-monitor
   ```

3. **Run the service**:

   ```bash
    ./celestia-upgrade-monitor -grpc-addr <GRPC_ENDPOINT> -port <PORT>
   ```

4. **Access the endpoints**:
   - JSON API: `http://<ADDRESS>:<PORT>/upgrade`
   - Prometheus metrics: `http://<ADDRESS>:<PORT>/metrics`

## 📈 Prometheus Metrics

The following metrics are exposed:

- `celestia_upgrade_status`: `1` if an upgrade is scheduled, `0` otherwise
- `celestia_upgrade_version`: The `app_version` of the upgrade
- `celestia_upgrade_height`: The block height of the next upgrade

### Example Output:

```plaintext
# HELP celestia_upgrade_height Height at which the upgrade will take place
# TYPE celestia_upgrade_height gauge
celestia_upgrade_height 6.680339e+06
# HELP celestia_upgrade_status Upgrade status as reported by celestia-app signal service
# TYPE celestia_upgrade_status gauge
celestia_upgrade_status 1
# HELP celestia_upgrade_version Current upgrade version
# TYPE celestia_upgrade_version gauge
celestia_upgrade_version 4
```

## 📋 Requirements

To build and run the Celestia Upgrade Monitor, you'll need:

- **Go 1.20 or later** – for compiling and running the application
- **Access to a Celestia node** (`celestia-appd`) with:
  - The gRPC server enabled and reachable (default port is often `9090`, but your deployment may vary)
  - The `celestia/signal/v1` gRPC service available

## 📝 License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT).
