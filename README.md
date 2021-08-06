# Outlier-Detection-LOF
 outlier detection in streaming data using local outlier factor (LOF) score
Please implement outlier detection in streaming data using local outlier factor (LOF)
score. You should implement it in C, C++, Go, or Java.

Your program should take input from stdin. The input contains a window size w (32-bit integer), a <host>:<port> pair where you receive an input stream of fixed-ddimensional 
data points in a comma separated values (CSV) format. Your program should output outliers starting from the (w+1)th input data until the end of input
stream. Your program should handle concept drift too. The LOF approach is a normalized distance-based approach. It adjusts for local variations in cluster density by normalizing 
distances with the average point-specific distances in a data locality. For a given data point x, let vk(x) be the distance to its knearest neighbor, and let Lk(x) be the set of 
points within the k-nearest neighbor distance of x. |Lk(x)| ≥ k because of ties in the distance. Then, the asymmetric reachability distance rk(x, y) of object x with respect to y 
is defined as rk(x, y) = max{Dist(x, y), vk(y)}. When y is in a dense region and the distance between x and y is large, rk(x, y) is equal to the true distance Dist(x, y). When the
distance between x and y is small, then rk(x, y) is smoothed out by the k-nearest neighbor distance of y.

The larger the value of k, the greater the smoothing. The average reachability distance ark(x) = MEANy∈Lk(x)rk(x, y), and LOFk(x) = MEANy∈Lk(x)(ark(x) / ark(y)). The maximum 
value of LOFk(x) over a range of different values of k is used as the outlier score to determine the best size of the neighborhood. To handle streaming data with
sliding window, we extend LOF to incremental scenarios: 1) the statistic of the newly inserted data points is computed, 2) only the LOF scores of the affected data points
by the newly inserted data point in the existing data points in the window are updated, and 3) similarly updated the deleted data points.

The data point with LOF score greater than a threshold t will be reported as outlier. 
