(ns noport.server
  (:require
    [clojure.java.io :as io]
    [compojure.core :refer :all]
    [compojure.route :as route]
    ))

;; The figwheel dev server calls this handler if it doesn't find a static file to
;; serve in resources/public. We can handle /api, but not /.
;; see: https://github.com/ring-clojure/ring/wiki/Concepts
(defroutes handler

  ; (POST "/documents/:file" {body :body {file :file} :params}
  ;   (with-open [o (io/output-stream (document-path file))]
  ;     (io/copy body o)
  ;     {:status 200}))

  ; (DELETE "/documents/:file" {{file :file} :params}
  ;   (io/delete-file (document-path file) true)
  ;     {:status 200})

  ; (GET "/documents/:file" {{file :file} :params}
  ;   (try
  ;     { :status 200
  ;       :headers {"Content-Type" "application/json"}
  ;       :body (slurp (document-path file))}
  ;     (catch java.io.FileNotFoundException e
  ;       {:status 404})))

  (route/not-found "Not found"))
