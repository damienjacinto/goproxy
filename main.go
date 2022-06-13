package main

import (
    "fmt"
    "github.com/rs/zerolog/log"
    "github.com/damienjacinto/internal/utils"
    "net/url"
	"net/http"
    "net/http/httputil"
	"github.com/gorilla/mux"
    "os"
    "github.com/steveyen/gkvlite"
)

func openFile(path string) *os.File {
    file := fmt.Sprintf("%s/db.gkvlite", path)
    log.Info().Msg(file)
    f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Fatal().Err(err)
    }
    return f
}

func createKeyValue(f *os.File) *gkvlite.Store {
    s, err := gkvlite.NewStore(f)
    if err != nil {
        log.Fatal().Err(err)
    }
    return s
}

func frontEndsHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}

func catchAllHandler(p *httputil.ReverseProxy, remote *url.URL) http.HandlerFunc {
    log.Info().Msg("Proxifing")
    return func(w http.ResponseWriter, r *http.Request) {
        log.Info().Msg(r.URL.String())
        r.Host = remote.Host
        p.ServeHTTP(w, r)
    }
}

func main() {
    utils.InitLog()
    config := utils.GetFlag()

    host := utils.GetEnv("HOST", "search-service-api.eu.sandbox.finalcad.cloud")
    scheme := utils.GetEnv("SCHEME", "https")
    path := utils.GetEnv("STORAGEPATH", "/home/damien/project/goproxy")
    f := openFile(path)
    store := createKeyValue(f)
    collection := store.SetCollection("frontends", nil)
    collection.Set([]byte("mercedes"), []byte("$$"))
    store.Flush()
    f.Sync()

    urlRedirect := fmt.Sprintf("%s://%s", scheme, host)
    remote, err := url.Parse(urlRedirect)
    if err != nil {
        log.Fatal().Err(err)
    }
    proxy := httputil.NewSingleHostReverseProxy(remote)

    log.Info().Msg("Starting sever...")
    log.Info().Msgf("%t", config.Debug)
    log.Info().Msgf("%s %s", host, scheme)
    router := mux.NewRouter().StrictSlash(true)
    //router.Host("{subdomain:[a-z]+}.example.com")
    router.HandleFunc("/admin/frontends", frontEndsHandler).Methods("GET")
    router.PathPrefix("/").Handler(http.HandlerFunc(catchAllHandler(proxy, remote)))

    srv := &http.Server{
        Addr:    ":8080",
        Handler: router,
    }
	log.Print(srv.ListenAndServe())
}
