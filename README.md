Pipes is a simple lib to do stuff in the pipline.

Code example
```
ctx := context.Background()
	data := []interface{}{1, 2, 3, 4, 5, 6, 7, 8}
	gen := Source(ctx, data)

	data1 := []interface{}{10, 42, 63, 84, 95, 6, 57, 48}
	gen1 := Source(ctx, data1)

	data2 := []interface{}{1, -2, 3, -4, 5, -6, 7, 8}
	gen2 := Source(ctx, data2)

	data3 := []interface{}{1.2, -2.6, 3.6, -4.7, 5.4, -6.8, 7.2, 8.0}
	gen3 := Source(ctx, data3)

	p1 := func(data []interface{}) interface{} {
		total := 0
		for _, elem := range data {
			switch elemTyped := elem.(type) {
			case int:
				return elemTyped + 2
			}
		}
		return total
	}
	pipe1 := Pipe(ctx, gen, p1)

	p2 := func(data []interface{}) interface{} {
		total := 0
		for _, elem := range data {
			switch elemTyped := elem.(type) {
			case int:
				return elemTyped * 10
			}
		}
		return total
	}
	pipe2 := Pipe(ctx, pipe1, p2)
	sink := Sink(ctx, []chan interface{}{gen1, pipe2, gen3})

	fsf := func(data []interface{}) interface{} {
		total := 0
		for _, elem := range data {
			switch elemTyped := elem.(type) {
			case int:
				if elemTyped%2 == 0 {
					return elemTyped
				}
			case float64:
				if elemTyped > 0 {
					return elemTyped
				}
			}
		}
		return total
	}
	filtersink := FilterSink(ctx, []chan interface{}{sink, gen2}, fsf)
	for val := range filtersink {
		fmt.Println(val)
	}
```
