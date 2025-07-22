package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	signaltypes "celestia-upgrade-monitor/celestia/signal/v1"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	// TODO: change to UpgradeData struct and add tally data
	prometheus.MustRegister(upgradeStatus, upgradeVersion, upgradeHeight)
}

func main() {
	log.Println("Starting gRPC client...")
	addr := flag.String("grpc-addr", "string", "gRPC server address")
	port := flag.String("server-port", "string", "HTTP server port, used to serve JSON data from this HTTP server")
	flag.Parse()
	GrpcServerAddress = *addr
	HttpServerPort = *port

	if GrpcServerAddress == "" {
		log.Fatal("gRPC server address must be provided using -grpc-addr flag")
	}
	if HttpServerPort == "" {
		log.Println("HTTP server port not provided, defaulting to :8080")
		HttpServerPort = "8088"
	}

	go func() {
		for {
			log.Println("Querying upgrade status for Prometheus...")
			updatePromMetrics()
			time.Sleep(30 * time.Minute)
		}
	}()

	go httpServer()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
	log.Println("HTTP server started on :8080")
	log.Println("gRPC client and HTTP server are running...")

	select {}
}

func grpcClient(addr string) (*grpc.ClientConn, error) {
	clientOptions := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.NewClient(
		addr,
		clientOptions,
	)
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	return conn, nil
}

func getUpgrade(client signaltypes.QueryClient) (*signaltypes.QueryGetUpgradeResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := client.GetUpgrade(ctx, &signaltypes.QueryGetUpgradeRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get upgrade: %w", err)
	}

	tally, err := client.VersionTally(ctx, &signaltypes.QueryVersionTallyRequest{})
	if err != nil {
		return nil, fmt.Errorf("failed to get version tally: %w", err)
	}
	// TODO: change return to UpgradeData struct and add tally data
	x := float64(tally.ThresholdPower) / float64(tally.TotalVotingPower)
	log.Printf("Version tally: %v, Total Voting Power: %d, Threshold Power: %d, Percentage: %.2f%%",
		tally, tally.TotalVotingPower, tally.ThresholdPower, x*100)

	return resp, nil
}

// HTTP server to responsed with JSON data from gRPC response
func httpServer() {
	http.HandleFunc("/upgrade", func(w http.ResponseWriter, r *http.Request) {
		// Handle the request and respond with JSON data
		conn, err := grpcClient(GrpcServerAddress)
		if err != nil {
			log.Printf("Failed to connect to gRPC server: %v", err)
		}
		defer conn.Close()
		client := signaltypes.NewQueryClient(conn)
		resp, err := getUpgrade(client)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get upgrade: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		log.Println("HTTP request handled successfully: /upgrade")
	})

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Starting HTTP server on :%s", HttpServerPort)
	log.Fatal(http.ListenAndServe(":"+HttpServerPort, nil))

}

func updatePromMetrics() {
	conn, err := grpcClient(GrpcServerAddress)
	if err != nil {
		log.Printf("Failed to connect to gRPC server: %v", err)
		return
	}
	defer conn.Close()

	client := signaltypes.NewQueryClient(conn)
	resp, err := getUpgrade(client)
	if err != nil {
		log.Printf("Failed to get upgrade: %v", err)
		return
	}

	// Step 1: Marshal the gRPC response to JSON
	jsonBytes, err := json.Marshal(resp) // marshal just the Upgrade message
	if err != nil {
		log.Fatalf("marshal failed: %v", err)
	}

	// Step 2: Unmarshal JSON into your custom struct
	var upgrade UpgradeResponse
	if err := json.Unmarshal(jsonBytes, &upgrade); err != nil {
		log.Fatalf("unmarshal failed: %v", err)
	}

	// TODO: change to UpgradeData struct and add tally data
	upgradeHeight.Set(float64(upgrade.Upgrade.UpgradeHeight))
	upgradeVersion.Set(float64(resp.Upgrade.AppVersion))
	if upgrade.Upgrade.UpgradeHeight > 0 {
		upgradeStatus.Set(1)
	} else {
		upgradeStatus.Set(0)
	}
}
