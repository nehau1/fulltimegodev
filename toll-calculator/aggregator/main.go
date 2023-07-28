package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/Stiffjobs/toll-calculator/types"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("httpAddr", ":3000", "the listen address of the HTTP server")
	grpcListenAddr := flag.String("grpcAddr", ":3001", "the listen address of the GRPC server")
	flag.Parse()
	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLogMiddleware(svc)

	go func() {
		log.Fatal(makeGRPCTansport(*grpcListenAddr, svc))
	}()
	log.Fatal(makeHTTPTransport(*httpListenAddr, svc))
}

func makeGRPCTansport(listenAddr string, svc Aggregator) error {
	fmt.Println("HTTP transport running on port ", listenAddr)
	// make a tcp listener
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	// make a new grpc native server with (options)
	server := grpc.NewServer([]grpc.ServerOption{}...)
	//Register (our) grpc server implementation to the grpc package
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(ln)
}

func makeHTTPTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("HTTP transport running on port ", listenAddr)
	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))

	return http.ListenAndServe(listenAddr, nil)
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		obuID := r.URL.Query().Get("obu")
		if obuID == "" {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing obuID"})
			return
		}
		num, err := strconv.Atoi(obuID)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "failed to convert obuID, please make sure it is an integer"})
			return
		}
		invoice, err := svc.CalculateInvoice(num)

		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		writeJSON(w, http.StatusOK, invoice)
	}
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance

		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}

		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(v)
}
