package arrays

func Sum(arr []int) int {
    var sum int
	for _, number := range arr {          // index,number
        sum += number             
	}
	return sum
}
  
func SumAll(nums ...[]int) []int {
	lengthOfNumbers := len(nums)
	sums := make([]int,lengthOfNumbers)

	for i , numbers := range nums {
		sums[i] = Sum(numbers)
	}
	
    return sums
}