(ns day1.hw1)
(defn loop-sum [nums]
	(loop [x 0 sum 0]
		(if (nil? (get nums x))
					sum
					(recur
						(+ x 1)
						(+ sum (get nums x))))))

(println (loop-sum [1, 2, 3, 6]))