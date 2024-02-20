package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//go:embed app
var embededFiles embed.FS

func getFileSystem(useOS bool) http.FileSystem {
	if useOS {
		log.Print("using live mode")

		return http.FS(os.DirFS("app"))
	}

	log.Print("using embed mode")

	fsys, err := fs.Sub(embededFiles, "app")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}

type UrlModel struct {
	Url string `json:"url"`
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s : %s \n", r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

func main() {
	// reverseProxyURL, ok := os.LookupEnv("REVERSE_PROXY_URL")
	// if !ok || reverseProxyURL == "" {
	// 	log.Fatalf("web: environment variable not declared: reverseProxyURL")
	// }

	// webPort, ok := os.LookupEnv("WEB_PORT")
	// if !ok || webPort == "" {
	// 	log.Fatalf("web: environment variable not declared: webPort")
	// }

	//router := gin.Default()

	useOS := len(os.Args) > 1 && os.Args[1] == "live"
	assetHandler := http.FileServer(getFileSystem(useOS))

	mux := http.NewServeMux()

	mux.Handle("/", requestLogger(assetHandler))

	//mux.Handle("/static/*", http.StripPrefix("/static/", assetHandler))

	server := &http.Server{
		Addr:    "0.0.0.0:8999",
		Handler: mux,
	}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal("Err while running server", err)
		}
	}()

	log.Println("Http server running successfully!")

	<-quit

	// e.Logger.Fatal(e.Start(fmt.Sprintf(":%v", webPort)))
}
