package main

import "net"
import "fmt"
import "bufio"
import "strings"
import "os"
import "unicode"
import "strconv"
import "math"
import "sort"

//The main function
func main() {
  var window [][]string
  window_count := 1
  threshold_lof := 0.0
  file_input := get_input()
  datapoints := get_datapoints(file_input)
  window_size, _ := strconv.Atoi(file_input[0][0])
  k := get_k_val(window_size)
  optimal_k := 0
  fmt.Println("The first value of k is: ", k)
  fmt.Println()
  if len(datapoints) < window_size {
    for _,item := range datapoints {
      window = append(window, item)
    }
    fmt.Println("Window 1 is:")
    fmt.Println(window)
    fmt.Println("Outlier Detection not done for Window number 1.")
    os.Exit(0)
  }
  for _,item := range datapoints {
    if len(window) < window_size {
      window = append(window, item)
      if len(window) == window_size {
        fmt.Printf("Window %d is:\n", window_count)
        fmt.Println(window)
        fmt.Println("Outlier Detection not done for window 1.")
        fmt.Println()
        var lof []float64
        for i := k; i < window_size; i = i + 2 {
          if i == k {
            neighbour_distances, k_distances := get_dist_neighbour(window, i)
            reach_distance, point_neighbors := get_k_neighborhood(window, k_distances, neighbour_distances)
            local_reach_distance := get_local_reach_distance(reach_distance)
            lof = get_lof(local_reach_distance, point_neighbors, i)
            threshold_lof = findMax(lof)
            optimal_k = i
          } else {
            neighbour_distances, k_distances := get_dist_neighbour(window, i)
            reach_distance, point_neighbors := get_k_neighborhood(window, k_distances, neighbour_distances)
            local_reach_distance := get_local_reach_distance(reach_distance)
            lof = get_lof(local_reach_distance, point_neighbors, i)
            threshold_lof_temp := findMax(lof)
            if threshold_lof_temp > threshold_lof {
              threshold_lof = threshold_lof_temp
              optimal_k = i
            } else {
              continue
            }
          }
        }
        fmt.Println("optimal value of k: ", optimal_k)
        fmt.Println("threshold_lof: ", threshold_lof)
        fmt.Println()
      }
    } else {
      window = window[1:]
      window_count++
      window = append(window, item)
      fmt.Printf("Window %d is:\n", window_count)
      fmt.Println(window)
      outlier_count := 0
      neighbour_distances, k_distances := get_dist_neighbour(window, optimal_k)
      reach_distance, point_neighbors := get_k_neighborhood(window, k_distances, neighbour_distances)
      local_reach_distance := get_local_reach_distance(reach_distance)
      lof := get_lof(local_reach_distance, point_neighbors, k)
      fmt.Println("The outliers in this window are:")
      for i := 0; i < len(lof); i++ {
        if lof[i] >= threshold_lof {
          outlier_count++
          fmt.Print(window[i])
        }
      }
      if outlier_count == 0 {
        fmt.Println("None")
      }
      fmt.Println()
      fmt.Println()
    }
  }
}

//the function to get input from the file
func get_input() ([][]string) {
  var file_input [][]string
  count := 0
  scanner_file_input := bufio.NewScanner(os.Stdin)
  scanner_file_input.Split(bufio.ScanLines)
  for scanner_file_input.Scan() {
    str := strings.Fields(scanner_file_input.Text())
    result := check_input(str, count)
    count++
    file_input = append(file_input, result)
  }
  return file_input
}

//the function to check the validity of the input
func check_input(str []string, num int) ([]string) {
  var file_input []string
  flag_wndw_size := 0
  flag_port := 0
  is_ip_or_host := 0
  str_temp := remove_spaces(str)
  if num == 0 {
    flag_wndw_size = chk_num(str_temp)
    if flag_wndw_size == 0 {
      fmt.Println("The window size is: ",str_temp)
      file_input = append(file_input, str_temp)
    } else {
      fmt.Println("Your input for window size is invalid!!")
      os.Exit(0)
    }
  } else {
    index_colon := strings.Index(str_temp, ":")
    host := str_temp[0:index_colon]
    port := str_temp[index_colon + 1:len(str_temp)]
    port_tmp, _ := strconv.Atoi(port)
    flag_port = chk_num(port)
    if port_tmp >= 1024 && port_tmp <= 65535 && flag_port == 0 {
      fmt.Println("The port is: ",port)
      file_input = append(file_input, port)
    } else {
      fmt.Println("The port number is invalid")
      os.Exit(0)
    }
    is_ip_or_host = checkStringForIpOrHostname(host)
    if is_ip_or_host == 1 {
      fmt.Println("The hostname is: ",host)
      file_input = append(file_input, host)
    } else if is_ip_or_host == 2 {
      fmt.Println("The IP Address is: ",host)
      file_input = append(file_input, host)
    } else {
      fmt.Println("The IP Address or Hostname is invalid!!")
      os.Exit(0)
    }
  }
  return file_input
}

//function to remove white spaces from the string
func remove_spaces(str []string) (string) {
  result := ""
  for _, str_temp := range str {
    if str_temp == " " {
      continue
    } else {
      result = result + str_temp
    }
  }
  return result
}

//function to check if input contains character
func chk_num(str string) int {
  flag := 0
  for i := 0; i < len(str); i++ {
    if unicode.IsNumber(rune(str[i])) {
      continue
    } else {
      flag = 1
    }
  }
  temp, _ := strconv.Atoi(str)
  if temp == 0{
    flag = 1
  }
  return flag
}

//function to check if the input has hostname or IP
func checkStringForIpOrHostname(host string) (int) { 
  flag := 0
  addr := net.ParseIP(host)
  addr_host, _ := net.LookupIP(host) 
  if addr  != nil {  
    flag = 2 
  } else if addr_host != nil {
    flag = 1
  } else {
    flag = -1
  }
  return flag 
}  

//function to get the datapoints
func get_datapoints(file_input [][]string) ([][]string) {
  var datapoints [][]string
  conn_string := file_input[1][1] + ":" + file_input[1][0]
  conn, _ := net.Dial("tcp", conn_string)
  scanner_server_data := bufio.NewScanner(conn)
  scanner_server_data.Split(bufio.ScanLines)
  for scanner_server_data.Scan() {
    str := strings.Fields(scanner_server_data.Text())
    datapoints = append(datapoints, str)
  }
  return datapoints
}

//function to choose the right value of k manually
func get_k_val(window_size int) int {
  k := 0
  if window_size == 1 {
    k = 1
  } else if window_size <= 5 {
    k = window_size - 1
  } else if window_size <= 10 {
    k = 5
  } else if window_size <= 20 {
    k = 15
  } else if window_size <= 100 {
    k = 25
  } else if window_size <= 500 {
    k = 155
  } else if window_size <= 1500 {
    k = 555
  } else {
    k = 1055
  }
  return k
}

//function to calculate distance between 2 points
func calc_distance(point_1 []string, point_2 []string) (float64) {
  sum := 0
  distance := 0.0
  for i := 0; i < len(point_1); i++ {
    item1, _ := strconv.Atoi(point_1[i])
    item2, _ := strconv.Atoi(point_2[i])
    sum = sum + (item1 - item2) * (item1 - item2)
  }
  distance = math.Sqrt(float64(sum))
  return distance
}

//function to calculate the distance from and to all the points in a window
func get_dist_neighbour(window [][]string, k int) ([][]float64, []float64) {
  k_distance := make([]float64, len(window))
  neighbour_distance := make([][]float64, len(window))
  for i := 0; i < len(window); i++ {
    neighbour_distance[i] = make([]float64, len(window) - 1)
  }
  for i := 0; i < len(window); i++ {
    iterator := 0
    temp := make([]float64, len(window) - 1)
    point_1 := get_int_points(window[i][0])
    for j := 0; j < len(window); j++ {
      if i == j {
        continue
      } else {
        point_2 := get_int_points(window[j][0])
        neighbour_distance[i][iterator] = calc_distance(point_1, point_2)
        temp[iterator] = neighbour_distance[i][iterator]
        iterator++
      }
    }
    sort.Float64s(temp)
    k_distance[i] = temp[k - 1]
    iterator = 0
    temp = nil
  }
  return neighbour_distance, k_distance
}

//function to split the string points into array 
func get_int_points(point string) []string {
  return strings.Split(point, ",")
}

//function to get k-neighborhood of each point in a window
func get_k_neighborhood(window [][]string, k_distance []float64, neighbour_distances [][]float64) ([]float64, [][]int) {
  result_mtrx_1 := make([][]int, len(window))
  for i := 0; i < len(result_mtrx_1); i++ {
    result_mtrx_1[i] = make([]int, len(window))
  }
  result_mtrx_2 := make([]float64, 0.0)  // to store the distance to that point from k neighbors for each element
  result_mtrx_3 := make([]float64, 0.0)  //to store the k distances of the selected neighbors
  avg_reach_distance := make([]float64, 0.0)
  count := 0
  for i := 0; i < len(neighbour_distances); i++ {
    for j := 0; j < len(neighbour_distances[0]); j++ {
      if count == i {
        count++
      }
      if neighbour_distances[i][j] <= k_distance[i] {
        result_mtrx_1[i][count] = 1
        result_mtrx_2 = append(result_mtrx_2, neighbour_distances[i][j])
        result_mtrx_3 = append(result_mtrx_3, k_distance[count])
        count++
      } else {
        count++
      }
    }
    avg_reach_distance = append(avg_reach_distance, get_avg_reach_distance(result_mtrx_2, result_mtrx_3))
    result_mtrx_2 = nil
    result_mtrx_3 = nil
    count = 0
  }
  return avg_reach_distance, result_mtrx_1
}

//function to calculate the average reachability distance
func get_avg_reach_distance(k_neighbor_distance []float64, k_distance []float64) float64{
  sum := 0.0
  for i := 0; i < len(k_distance); i++ {
    if k_distance[i] > k_neighbor_distance[i] {
      sum = sum + k_distance[i]
    } else {
      sum = sum + k_neighbor_distance[i]
    }
  }
  return sum / float64(len(k_distance))
}

//function to calculate the local reachability distance
func get_local_reach_distance(reach_distance []float64) []float64 {
  result_mtrx := make([]float64, 0.0)
  for i := 0; i < len(reach_distance); i++ {
    temp := 1 / reach_distance[i]
    result_mtrx = append(result_mtrx, temp)
  }
  return result_mtrx
}

//function to calculate LOF score
func get_lof(local_reach_distance []float64, neighbors [][]int, k int) []float64{
  result_mtrx_1 := make([]float64, 0.0)
  sum := 0.0
  for i := 0; i < len(neighbors); i++ {
    for j := 0; j < len(neighbors[0]); j++ {
      if neighbors[i][j] == 1 {
        sum += local_reach_distance[j]
      }
    }
    sum = sum / float64(k)
    result_mtrx_1 = append(result_mtrx_1, sum / local_reach_distance[i])
  }
  return result_mtrx_1
}

//function to find max value in te array
func findMax(arr []float64) (float64) {
  max := arr[0]
  for _, value := range arr {
    if value > max {
      max = value
    }
  }
  return max
}