package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// ルート証明書の読み込み
	caCert, err := os.ReadFile("cert_files/rootCA.pem")
	if err != nil {
		log.Fatal("ルート証明書の読み込みエラー:", err)
	}
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		log.Fatal("ルート証明書のプールへの追加に失敗しました")
	}

	// TLS設定
	tlsConfig := &tls.Config{
		ClientCAs:  caCertPool,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	// 複数のエンドポイントを設定
	mux := http.NewServeMux()

	// ルートページ
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		logCertificateInfo(r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
            <h1>クライアント証明書テスト - メインページ</h1>
            <h2>認証成功！クライアント証明書のCN: %s</h2>
            <h2>アクセス時刻: %s</h2>
            <h2><a href="/page1">ページ1へ移動</a></h2>
            <h2><a href="/page2">ページ2へ移動</a></h2>
        `, getCertificateCN(r), time.Now().Format("2006-01-02 15:04:05"))
	})

	// ページ1
	mux.HandleFunc("/page1", func(w http.ResponseWriter, r *http.Request) {
		logCertificateInfo(r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
            <h1>クライアント証明書テスト - ページ1</h1>
            <h2>認証成功！クライアント証明書のCN: %s</h2>
            <h2>アクセス時刻: %s</h2>
            <h2><a href="/">メインページへ戻る</a></h2>
            <h2><a href="/page2">ページ2へ移動</a></h2>
        `, getCertificateCN(r), time.Now().Format("2006-01-02 15:04:05"))
	})

	// ページ2
	mux.HandleFunc("/page2", func(w http.ResponseWriter, r *http.Request) {
		logCertificateInfo(r)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprintf(w, `
            <h1>クライアント証明書テスト - ページ2</h1>
            <h2>認証成功！クライアント証明書のCN: %s</h2>
            <h2>アクセス時刻: %s</h2>
            <h2><a href="/">メインページへ戻る</a></h2>
            <h2><a href="/page1">ページ1へ移動</a></h2>
        `, getCertificateCN(r), time.Now().Format("2006-01-02 15:04:05"))
	})

	// HTTPSサーバーの設定
	server := &http.Server{
		Addr:      ":8443",
		Handler:   mux,
		TLSConfig: tlsConfig,
	}

	// サーバーの起動
	log.Printf("サーバーを起動します: https://localhost:8443")
	log.Fatal(server.ListenAndServeTLS("cert_files/server.pem", "cert_files/server.key"))
}

// 証明書情報をログに出力する関数
func logCertificateInfo(r *http.Request) {
	log.Printf("アクセスパス: %s", r.URL.Path)
	if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
		log.Printf("クライアント証明書検出 - CN: %s, 発行者: %s",
			r.TLS.PeerCertificates[0].Subject.CommonName,
			r.TLS.PeerCertificates[0].Issuer.CommonName)
	} else {
		log.Printf("クライアント証明書が見つかりません")
	}
}

// 証明書のCNを取得する関数
func getCertificateCN(r *http.Request) string {
	if r.TLS != nil && len(r.TLS.PeerCertificates) > 0 {
		return r.TLS.PeerCertificates[0].Subject.CommonName
	}
	return "不明"
}
