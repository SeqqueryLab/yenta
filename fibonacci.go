package yenta

import "fmt"

func fibonacci(n int) ([]int, error) {
	if n < 0 {
		return nil, fmt.Errorf("got %d, want n >= 0", n)
	}
	switch n {
	case 0:
		return []int{}, nil
	case 1:
		return []int{0}, nil
	case 2:
		return []int{0, 1}, nil
	default:
		cur, err := fibonacci(n - 1)
		if err != nil {
			return nil, err
		}

		cur = append(cur, cur[len(cur)-2]+cur[len(cur)-1])
		return cur, nil
	}
}
