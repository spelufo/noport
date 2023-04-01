(ns noport.main
  (:require
    [noport.utils :refer [remove-at p! pjs! get! post!]]
    [reagent.core :as r]
    [reagent.dom :as rd]))

(defonce *state
  (r/atom {:servers [{:domain "noport"  :port 8012} ] }))

(defn next-port []
  (inc (:port (last (:servers @*state)))))

(defn assoc-server [state domain port ssl?]
  (update state :servers conj {:domain domain :port port :ssl? ssl?}))
(defn assoc-server! [domain port ssl?]
  (swap! *state assoc-server domain port ssl?))

(defn update-server [state i f & args]
  (apply update-in state [:servers i] f args))
(defn update-server! [i f & args]
  (apply swap! *state update-server i f args))

(defn remove-server [state i]
  (update state :servers remove-at i))
(defn remove-server! [i]
  (swap! *state remove-server i))

(defn server-by-domain [state domain]
  (filter #(= (:domain %) domain) (:servers state)))

(defn server-by-port [state port]
  (filter #(= (:port %) port) (:servers state)))

(defn distinct-domains? [state]
  (apply distinct? (map :domain (:servers state))))

(defn distinct-ports? [state]
  (apply distinct? (map :port (:servers state))))

(defn validate [state]
  (cond
    (not (distinct-domains? state)) "Duplicate domains."
    (not (distinct-ports? state)) "Duplicate ports."))

(defn load! []
  (.then
    (get! (str js/API_URL "/.noport.json"))
    #(reset! *state %)))

(defn save! []
  (post! (str js/API_URL "/.noportt.json") @*state))

(defn install! []
  (.then
    (save!)
    (post! (str js/API_URL "/install") {})))


;; UI ;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;;

(defn ui-server [i server]
  [:div.server {:key i}
    [:div.server-domain
      [:input {
        :value (:domain server)
        :on-change
          (fn [ev]
            (let [domain (.. ev -target -value)]
              (swap! *state assoc-in [:servers i :domain] domain)))}]
      " " [:a {:href (str "http://" (:domain server) ".localhost")} ".localhost"]]
    [:div "→"]
    [:div.server-port
      [:a {:href (str "http://localhost:" (:port server))} "localhost:"] " "
      [:input {
        :value (:port server)
        :on-change
          (fn [ev]
            (let [port-str (.. ev -target -value) port (js/parseInt port-str)]
              (when (= port-str "")
                (swap! *state assoc-in [:servers i :port] 0))
              (when-not (NaN? port)
                (swap! *state assoc-in [:servers i :port] port))))}]]
    [:div.server-buttons
      [:button.remove-server {:on-click #(remove-server! i)} "×"]]])

(defn ui-servers [servers]
  (into [:div.servers]
      (map-indexed ui-server servers)))

(defn ui-main []
  (let [state @*state invalid-msg (validate state) invalid? (boolean invalid-msg)]
    [:div.main
      [:h1 "Servers"]
      [ui-servers (:servers state)]
      [:div.buttons
        [:button {:on-click #(assoc-server! "" (next-port) true)} "New"]
        [:button {:on-click #(save!) :disabled invalid?} "Save"]
        [:button {:on-click #(install!) :disabled invalid?} "Install"]
        (when invalid? [:span.invalid-message invalid-msg])]]))

(defn mount-app! []
  (rd/render [ui-main] (.getElementById js/document "app")))

(mount-app!)

(defonce _
  (load!))
