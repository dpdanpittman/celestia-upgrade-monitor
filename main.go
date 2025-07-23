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
	prometheus.MustRegister(upgradeStatus,
		upgradeVersion,
		upgradeHeight,
		tallyThresholdPower,
		tallyTotalVotingPower,
		tallyThresholdPercent,
	)
}

func main() {
	log.Println("Starting gRPC client...")

	// Define flags for gRPC server address and HTTP server port
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
		HttpServerPort = "8080"
	}

	// Start Prometheus metrics update func
	go func() {
		for {
			log.Println("Querying upgrade status for Prometheus /metrics...")
			updatePromMetrics()
			time.Sleep(30 * time.Minute)
		}
	}()

	// Start the HTTP server
	go httpServer()
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
	log.Println("HTTP server started on :8080")
	log.Println("gRPC client and HTTP server are running...")

	select {}
}

func grpcClient(addr string) (*grpc.ClientConn, error) {
	// Create a gRPC client connection to the specified address
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

func getUpgrade(client signaltypes.QueryClient) (UpgradeData, error) {
	// Create a context with a timeout for the gRPC request
	// Used for the Prometheus /metrics endpoint
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Get the upgrade information from the gRPC client
	resp, err := client.GetUpgrade(ctx, &signaltypes.QueryGetUpgradeRequest{})
	if err != nil {
		return UpgradeData{}, fmt.Errorf("failed to get upgrade: %w", err)
	}
	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("marshal failed: %v", err)
	}
	var upgrade UpgradeResponse
	if err := json.Unmarshal(jsonBytes, &upgrade); err != nil {
		log.Fatalf("unmarshal failed: %v", err)
	}

	// Get the version tally information from the gRPC client
	tally, err := client.VersionTally(ctx, &signaltypes.QueryVersionTallyRequest{})
	if err != nil {
		return UpgradeData{}, fmt.Errorf("failed to get version tally: %w", err)
	}

	// Prepare the return data
	percent := float64(tally.ThresholdPower) / float64(tally.TotalVotingPower)
	returnData := UpgradeData{
		UpgradeData: UpgradeResponse{
			Upgrade: Upgrade{
				AppVersion:    upgrade.Upgrade.AppVersion,
				UpgradeHeight: upgrade.Upgrade.UpgradeHeight,
			},
		},
		TallyData: TallyResponse{
			TotalVotingPower: int64(tally.TotalVotingPower),
			ThresholdPower:   int64(tally.ThresholdPower),
			ThresholdPercent: percent,
		},
	}

	return returnData, nil
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

	// Handle Prometheus metrics endpoint
	http.Handle("/metrics", promhttp.Handler())

	log.Printf("Starting HTTP server on :%s", HttpServerPort)
	log.Fatal(http.ListenAndServe(":"+HttpServerPort, nil))
}

func updatePromMetrics() {
	// Create a gRPC connection
	conn, err := grpcClient(GrpcServerAddress)
	if err != nil {
		log.Printf("Failed to connect to gRPC server: %v", err)
		return
	}
	defer conn.Close()

	// Create a gRPC client
	client := signaltypes.NewQueryClient(conn)
	resp, err := getUpgrade(client)
	if err != nil {
		log.Printf("Failed to get upgrade: %v", err)
		return
	}

	// Marshal the response to JSON
	jsonBytes, err := json.Marshal(resp) // marshal just the Upgrade message
	if err != nil {
		log.Fatalf("marshal failed: %v", err)
	}
	var upgrade UpgradeResponse
	if err := json.Unmarshal(jsonBytes, &upgrade); err != nil {
		log.Fatalf("unmarshal failed: %v", err)
	}

	// Collect tally data
	tally, err := client.VersionTally(context.Background(), &signaltypes.QueryVersionTallyRequest{})
	if err != nil {
		log.Printf("Failed to get version tally: %v", err)
		return
	}

	// Update Prometheus metrics
	tallyThresholdPower.Set(float64(tally.ThresholdPower))
	tallyTotalVotingPower.Set(float64(tally.TotalVotingPower))
	percent := float64(tally.ThresholdPower) / float64(tally.TotalVotingPower)
	tallyThresholdPercent.Set(percent)
	upgradeHeight.Set(float64(resp.UpgradeData.Upgrade.UpgradeHeight))
	upgradeVersion.Set(float64(resp.UpgradeData.Upgrade.AppVersion))
	if resp.UpgradeData.Upgrade.UpgradeHeight > 0 {
		upgradeStatus.Set(1)
	} else {
		upgradeStatus.Set(0)
	}
}
