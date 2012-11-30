// Package assetcache provides a caching asset handler for static assets.
// If the asset is not in memory, assetcache will cache it for faster
// access.
package assetcache

import (
        "fmt"
        "github.com/gokyle/filecache"
        "github.com/gokyle/webshell"
        "net/http"
        "regexp"
)

// The cache is pre-configured to store up to 128 2MB sized items for up
// to 30 days, checking every hour for expired items.
var (
        MaxItems = 128
        MaxExpire = 2592000     // 30 days in seconds
        MaxSize int64 = 2 * filecache.Megabyte
        Every = 3600
)

// Type AssetCache represents a file cache for static assets.
type AssetCache struct {
        cache   *filecache.FileCache
        path    string
        stripre *regexp.Regexp
}

// CreateAssetCache creates a basic static asset cache.
func CreateAssetCache(route, path string) *AssetCache {
        cache := new(filecache.FileCache)
        cache.MaxItems = MaxItems
        cache.ExpireItem = MaxExpire
        cache.MaxSize = MaxSize
        cache.Every = Every

        regex_string := fmt.Sprintf("^%s(.*)$", route)
        stripre := regexp.MustCompile(regex_string)
        return &AssetCache{cache, path, stripre}
}

// IsRunning indicates whether the cache is active.
func (ac *AssetCache) IsRunning() bool {
        return ac.cache.Active()
}

// Start performs the necessary cache startup initialisation, and starts it.
// The cache must have been created already.
func (ac *AssetCache) Start() error {
        return ac.cache.Start()
}

// AssetHandler returns an HTTP handler for the given AssetCache.
// It is suitable for use in the WebApp.AddRoute method.
func AssetHandler(ac *AssetCache) webshell.RouteHandler {
        return func(w http.ResponseWriter, r *http.Request) {
                r.URL.Path = ac.stripre.ReplaceAllString(r.URL.Path,
                        "/" + ac.path + "$1")
                if ac.cache.InCache(r.URL.Path) {
                        fmt.Println("assetcache <- cache")
                }
                ac.cache.HttpWriteFile(w, r)
        }
}

// BackgroundAttachAssetCache will transparently set up an asset cache for a WebApp.
func BackgroundAttachAssetCache(app *webshell.WebApp, route, path string) (err error) {
        ac := CreateAssetCache(route, path)
        err = ac.cache.Start()
        if err != nil || ! ac.IsRunning() {
                return err
        }
        app.AddRoute(route, AssetHandler(ac))
        return
}
