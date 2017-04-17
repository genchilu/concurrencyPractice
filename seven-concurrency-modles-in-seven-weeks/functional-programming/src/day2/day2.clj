(ns day2.day2
	(:require [clojure.core.protocols :refer [CollReduce coll-reduce]]
						[clojure.core.reducers :refer [CollFold coll-fold]]))
(defn get-words [text] (re-seq #"\w+" text))
(def pages ["one potato two potato three potato four" "five potato six potato seven potato more"])
(def merge-counts (partial merge-with +))
(defn count-words-parallel [pages]
	(reduce (partial merge-with +)
		(pmap #(frequencies (get-words %)) pages)))

; (defn count-words [pages]
; (reduce (partial merge-with +)
;	(pmap count-words-sequential (partition-all 100 pages))) )

(defn my-reduce
	([f coll] (coll-reduce coll f))
	([f init coll] (coll-reduce coll f init)))

(defn make-reducer [reducible transformf]
	(reify CollReduce
		(coll-reduce [_ f1]
			(coll-reduce reducible
				(transformf f1) (f1)))
		(coll-reduce [_ f1 init]
			(coll-reduce reducible (transformf f1) init))))

(defn my-map [mapf reducible]
	(make-reducer reducible
		(fn [reducef]
			(fn [acc v] (reducef acc (mapf v))))))

(println (into [] (my-map (partial * 2) [1 2 3 4])))
;(println (pmap #(frequencies (get-words %)) pages))
;(println (merge-counts {:x 1 :y 2} {:y 1 :z 1}))
;(println (count-words-parallel pages))
;(println (reduce conj [] (r/map (partial * 2) [1 2 3 4])))
;(println (into [] (r/map (partial * 2) [1 2 3 4])))
;(println (into [] (r/map (partial + 1) (r/filter even? [1 2 3 4]))))
