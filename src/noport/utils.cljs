(ns noport.utils)

;; UTILS ;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(defn p! [& args]
  (apply println args)
  (last args))

(defn pjs! [& args]
  (apply js/console.log args)
  (last args))

(defn remove-at [v i]
  (if (< -1 i (count v)) (into (subvec v 0 i) (subvec v (inc i))) v))

;; LOAD/SAVE ;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(defn get! [url]
  (->
    (js/fetch url
      (clj->js {
        :method "GET"
        :headers {"Content-Type" "application/json"
                  "Accept" "application/json"} }))
    (.then
      (fn [response]
        (assert (.-ok response) "doc not found")
        (.json response)))))

(defn post! [url json-data]
  (js/fetch url
    (clj->js {
      :method "POST"
      :headers {"Content-Type" "application/json"
                "Accept" "application/json"}
      :body (js/JSON.stringify json-data nil 2) })))
