package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

var (
	blockedIPs = make(map[string]time.Time)
	mutex      = sync.Mutex{}
)

// Adiciona o IP ao iptables para bloqueio
func blockIP(ip string) {
	cmd := exec.Command("iptables", "-A", "INPUT", "-s", ip, "-j", "DROP")
	cmd.Run()
}

// Middleware para filtrar requisições
func requestFilter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		clientIP := r.RemoteAddr[:len(r.RemoteAddr)-6] // Remove a porta do IP
		userAgent := r.Header.Get("User-Agent")

		log.Printf("Requisição de %s - User-Agent: %s", clientIP, userAgent)

		// Verifica se o IP está bloqueado
		mutex.Lock()
		if unblockTime, exists := blockedIPs[clientIP]; exists {
			if time.Now().Before(unblockTime) {
				http.Error(w, "403 Forbidden", http.StatusForbidden)
				mutex.Unlock()
				return
			} else {
				delete(blockedIPs, clientIP)
			}
		}
		mutex.Unlock()

		// Permite apenas requisições de wget e curl
		if !strings.HasPrefix(userAgent, "curl") && !strings.HasPrefix(userAgent, "Wget") {
			mutex.Lock()
			blockedIPs[clientIP] = time.Now().Add(10 * time.Minute) // Bloqueia por 10 minutos
			mutex.Unlock()
			blockIP(clientIP) // Adiciona ao iptables
			log.Printf("IP bloqueado: %s", clientIP)
			http.Error(w, "403 Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", requestFilter(http.FileServer(http.Dir(cwd))))

	httpServer := &http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: mux,
	}

	log.Println("Servidor rodando na porta 8080, servindo arquivos de", cwd)
	log.Fatal(httpServer.ListenAndServe())
}
