(ns noport.main
  (:require
    [noport.utils :refer [p! pjs!]]
    [reagent.core :as r]
    [reagent.dom :as rd]))


;;; UI

(defonce *state (r/atom {
  :servers [
    {:domain "noport.dev" :port 9500}
    {:domain "spelufo.com" :port 1313}
  ]
  }))

(defn ui-server [i server]
  [:div.server {:key i}
    [:div.server-domain
      [:input {
        :value (:domain server)
        :on-change
          (fn [ev]
            (let [domain (.. ev -target -value)]
              (swap! *state assoc-in [:servers i :domain] domain)))}]
      ".localhost"]
    [:div.server-port
      [:input {
        :value (:port server)
        :on-change
          (fn [ev]
            (let [port-str (.. ev -target -value) port (js/parseInt port-str)]
              (when (= port-str "")
                (swap! *state assoc-in [:servers i :port] 0))
              (when-not (NaN? port)
                (swap! *state assoc-in [:servers i :port] port))))}]]])

(defn ui-servers [servers]
  (into [:div.servers]
      (map-indexed ui-server servers)))

(defn ui-main []
  (let [state @*state]
    [:div.main
      [:h1 "Servers"]
      [ui-servers (:servers state)]
      [:button {} "Save"]
      [:pre.nginx-config ]]))

(defn mount-app! []
  (rd/render [ui-main] (.getElementById js/document "app")))

(mount-app!)

(defonce _
  nil)
